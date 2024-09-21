using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

using FTP_server;
using FTP_server.consts;
using FTP_server.utils;
using FTP_server.file;

public class FtpSocketServer
{
    private const int port = ftpConfig.PORT;
    private Socket serverSocket;
    private const char EOF = ftpConfig.MESSAGE_EOF;
    public FtpSocketServer()
    {
        serverSocket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    }
    public void Start()
    {
        // 设置服务器的 IP 地址和端口
        IPAddress ipAddress = ftpConfig.LISTENER_IP_ADDRESS;
        IPEndPoint localEndPoint = new IPEndPoint(ipAddress, port);
        serverSocket.Bind(localEndPoint);
        serverSocket.Listen(10);
        Log.info($"FTP Server started on port {{{port}}}");
        Log.info("The root is: " + storeConfig.STORE_ROOT_PATH);

        try
        {
            while (true) // info FTP 业务入口
            {
                Socket clientSocket = serverSocket.Accept();
                Log.info("", "LINK", clientSocket);

                Task.Run(() => HandleClient(clientSocket)); // info 开一个线程处理
            }
        }
        catch (SocketException ex) // Socket相关操作中可能抛出的异常
        { 
            Log.error($"Socket error: {ex.Message}");
            serverSocket?.Close();
        }
        catch (FormatException ex) // 解析格式不正确的字符串时可能抛出的异常
        {
            Log.error($"Invalid IP address format: {ex.Message}");
            serverSocket?.Close();
        }
    }

    private void HandleClient(Socket clientSocket)
    {
        try
        {
            byte[] buffer = new byte[1024];

            // info 登录建立连接
            int bytesReceived = clientSocket.Receive(buffer);
            string receivedMessage = Encoding.UTF8.GetString(buffer, 0, bytesReceived);
            var (command, argument) = stringHandler.SplitCommand(receivedMessage);
            Log.info(argument, command, clientSocket);

            if (command.ToUpper() != "LOGIN")
            {
                SendConfirmation(clientSocket, "Please link FTP server by [LOGIN]"); // UPDATE
                return;
            }
            else
            {
                var (username, password) = stringHandler.SplitCommand(argument);
                password = AES256.DecryptString(password, AES256Config.KEY, AES256Config.IV);
                Log.debug($"{username} {password}", "PASSWORD-DECRYPT", clientSocket);
                SendConfirmation(clientSocket, "Login successfully.");
            }
            // TODO 增加鉴权 

            while (true)
            {
                bytesReceived = clientSocket.Receive(buffer);
                if(bytesReceived > 0)
                {
                    receivedMessage = Encoding.UTF8.GetString(buffer, 0, bytesReceived);

                    (command, argument) = stringHandler.SplitCommand(receivedMessage);
                    Log.info(argument, command, clientSocket);

                    if (command.ToUpper() == "EXIT")    // info 退出 TCP 连接
                        break;

                    switch (command.ToUpper())
                    {
                        case "ECHO":
                            SendConfirmation(clientSocket, argument);
                            break;
                        case "LIST":
                            //DirectoryLister.ListDirectoryContents(storeConfig.STORE_PATH); //! info 暂时弃用。
                            //if (argument == ".")    // info 特化处理
                            //    SendDirectoryList(clientSocket, storeConfig.STORE_ROOT_PATH);
                            //else
                            SendDirectoryList(clientSocket, storeConfig.STORE_ROOT_PATH + argument);
                            break;
                        case "DOWN":    // info 下载文件
                            Task.Run(() => SendFile(clientSocket, storeConfig.STORE_ROOT_PATH + argument)); // INFO 再开一个子任务
                            break;
                        case "UPLOAD":
                            string[] parts = stringHandler.SplitCommandArray(argument);
                            long fileSize = long.Parse(parts[1]);
                            Task.Run(() => receiveFile(clientSocket, parts[0], parts[2], fileSize)); // 同上
                            break;
                        default:
                            SendConfirmation(clientSocket, "500 Command not understood");
                            break;
                    }
                    // clear buffer
                    Array.Clear(buffer, 0, buffer.Length);
                }
            }
        }
        catch (SocketException ex)
        {
            Log.error($"Socket error: {ex.Message}");
        }
        finally
        {
            Log.info("", "UNLINK", clientSocket);
            clientSocket.Close();
        }
        //Console.WriteLine("Tag: end.");
    }

    private void SendConfirmation(Socket clientSocket, string message)
    {
        // 返回一般的 string 作为基本的消息传递方式。
        byte[] responseBytes = Encoding.UTF8.GetBytes(message + EOF);
        clientSocket.Send(responseBytes);
        Log.info(message, "BACK-MESSAGE", clientSocket);
        return;
    }

    private void SendStopCommand(Socket clientSocket)
    {
        // 用于 Send Directory 中停止发送列表的方式。
        string stopCommand = "STOP";
        byte[] byteData = Encoding.ASCII.GetBytes(stopCommand + '\n');

        clientSocket.Send(byteData);
    }

    private Socket buildNewLink()
    {
        var transferListener = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
        transferListener.Bind(new IPEndPoint(IPAddress.Any, 0));
        transferListener.Listen(1);
        return transferListener;
    }

    private void SendDirectoryList(Socket clientSocket, string path)
    {
        const string DATE_FORMATTING = "MM dd yyyy";

        var dirs = Directory.GetDirectories(path);
        var files = Directory.GetFiles(path);
        StringBuilder stringBuilder = new StringBuilder();

        foreach ( var dir in dirs)
        {
            DirectoryInfo directoryInfo = new DirectoryInfo(dir);
            string date = directoryInfo.LastWriteTime.ToString(DATE_FORMATTING);
            stringBuilder.AppendLine($"drwxr-xr-x 2 owner group {date} {Path.GetFileName(dir)}");
        }

        foreach ( var file in files)
        {
            FileInfo fileInfo = new FileInfo(file);
            string date = fileInfo.LastWriteTime.ToString(DATE_FORMATTING);
            stringBuilder.AppendLine($"-rw-r--r-- 1 owner group {date} {fileInfo.Length} {Path.GetFileName(file)}");
        }

        byte[] byteData = Encoding.UTF8.GetBytes(stringBuilder.ToString());
        clientSocket.Send(byteData);

        SendStopCommand(clientSocket); // 

        return;
    }

    async private void SendFile(Socket clientSocket, string fullFilePath)
    {
        var transferListener = buildNewLink();
        int transferPort = ((IPEndPoint)transferListener.LocalEndPoint).Port;
        Log.debug($"Transfer server started on port {{{transferPort}}} for ...");

        SendConfirmation(clientSocket, $"DOWNLOAD_PORT {transferPort}");
        var transferClient = await transferListener.AcceptAsync();
        Log.info($"Connection established with {transferClient.RemoteEndPoint} for file transfer", "SUB-LINK", transferClient);

        var receiveTask = new FileTransferController(transferClient, fullFilePath);
        await receiveTask.StartAsync();

        SendConfirmation(clientSocket, "Download file successfully.");
    }

    //private static string ReceiveString(Socket socket)
    //{
    //    byte[] buffer = new byte[1024];
    //    int bytesRec = socket.Receive(buffer);
    //    string data = Encoding.ASCII.GetString(buffer, 0, bytesRec);
    //    return data.Trim();
    //}

    async private void receiveFile(Socket clientSocket, string saveFilePath, string fileName, long fileSize) 
    {
        var transferListener = buildNewLink();
        int transferPort = ((IPEndPoint)transferListener.LocalEndPoint).Port;
        Log.debug($"Transfer server started on port {{{transferPort}}} for ...");

        //string response = $"UPLOAD_PORT {transferPort}\n";
        //clientSocket.Send(Encoding.UTF8.GetBytes(response));
        SendConfirmation(clientSocket, $"UPLOAD_PORT {transferPort}");
        var transferClient = await transferListener.AcceptAsync();
        Log.info($"Connection established with {transferClient.RemoteEndPoint} for file transfer", "SUB-LINK", transferClient);

        string fullFilePath = Path.Combine(storeConfig.STORE_ROOT_PATH, saveFilePath, fileName);
        var receiveTask = new FileTransferController(transferClient, fullFilePath, fileSize);
        await receiveTask.StartAsync();

        SendConfirmation(clientSocket, "Upload file successfully.");
    }
}

