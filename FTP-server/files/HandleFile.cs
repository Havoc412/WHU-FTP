using System.Net.Sockets;
using System.Text;

namespace FTP_server.file
{
    public class FileTransferController
    {
        // 异步处理文件的具体操作
        private readonly Socket transferSocket;
        private readonly string filePath;
        private readonly long fileSize; // 这也就是 upload 文件的上限。
        private volatile bool paused;
        private const string ControlPrefix = "CTRL:";
        private const string DATAPrefix = "DATA:";

        public FileTransferController(Socket transferSocket, string filePath, long fileSize = -1)
        {
            this.transferSocket = transferSocket;
            this.filePath = filePath;
            this.fileSize = fileSize;
            this.paused = false;

            // 自启动
            //Task.Run(() => StartAsync());
        }

        public async Task StartAsync()
        {
            if (fileSize == -1)
                await sendFile();
            else 
                await receiveFile();

            transferSocket.Close(); // 任务完成，关闭连接。
        }

        private async Task receiveFile()
        {
            string fileName = Path.GetFileName(filePath);
            Log.debug($"Receiving file: {{{fileName}}} of size {{{fileSize}}} bytes to {{{filePath}}}", "", transferSocket);
            // TIP 打开文件流
            using (var fileStream = new FileStream(filePath, FileMode.Create))
            {
                byte[] buffer = new byte[1024 * 1024 + ControlPrefix.Length]; // 1MB buffer
                long totalReceived = 0;
                int bytesRead;

                while (totalReceived < fileSize)
                {
                    //if(totalReceived > 250000000)
                    //    Console.WriteLine(totalReceived);
                    if(paused)
                    {
                        await Task.Delay(500);
                        bytesRead = transferSocket.Receive(buffer);
                        string command = Encoding.UTF8.GetString(buffer, 0, Math.Min(ControlPrefix.Length, bytesRead));
                        if (command.StartsWith(ControlPrefix))
                        {
                            string receiveText = Encoding.UTF8.GetString(buffer, 0, bytesRead);
                            HandleCtrlMessage(receiveText);
                        }
                        else
                        {
                            Log.error("Task is PAUSEed, waiting for other CTRL, not " + command);
                        }
                        continue;
                    }
                    bytesRead = transferSocket.Receive(buffer);
                    string receivePrefix = Encoding.UTF8.GetString(buffer, 0, Math.Min(ControlPrefix.Length, bytesRead));

                    //string testDATA = Encoding.UTF8.GetString((byte[])buffer, 0, bytesRead);
                    //Log.debug(testDATA, "test", transferSocket);
                    if(receivePrefix.StartsWith(ControlPrefix)) 
                    {
                        string receiveText = Encoding.UTF8.GetString(buffer, 0, bytesRead);
                        HandleCtrlMessage(receiveText);
                    } 
                    else
                    {
                        int dataLength = bytesRead - DATAPrefix.Length;
                        fileStream.Write(buffer, DATAPrefix.Length,dataLength);
                        totalReceived += dataLength;
                    }
                }
            }
            Log.debug($"File {fileName} received and saved successfully.", "", transferSocket);
        }

        private async Task sendFile()
        {
            // INFO 1, send file size
            FileInfo fileInfo = new FileInfo(filePath);
            byte[] fileLength = BitConverter.GetBytes(fileInfo.Length);
            transferSocket.Send(fileLength);

            using (FileStream fs = new FileStream(filePath, FileMode.Open, FileAccess.Read))
            {
                byte[] buffer = new byte[1024 * 1024]; //  INFO 1MB 块大小
                int bytesRead;
                while (true) // 
                {
                    if (!paused)
                    {
                        bytesRead = await fs.ReadAsync(buffer, 0, buffer.Length);
                        if (bytesRead == 0)
                        {
                            break; // 文件读取完毕，退出循环
                        }

                        transferSocket.Send(buffer, bytesRead, SocketFlags.None);
                    } 
                    else 
                    {
                        await Task.Delay(500); // Wait for resume command
                    }

                    // Check for control commands
                    if (transferSocket.Available > 0) // 轮询的方式
                    {
                        byte[] commandBuffer = new byte[1024];
                        int commandBytes = transferSocket.Receive(commandBuffer);
                        string command = Encoding.UTF8.GetString(commandBuffer, 0, commandBytes).Trim();

                        HandleCtrlMessage(command);
                    }
                }
            }
            Log.debug($"File {filePath} received and saved successfully.", "", transferSocket);
        }

        private void HandleCtrlMessage(string message)
        {
            if (message.Contains("pause"))
            {
                paused = true;
            }
            else if (message.Contains("resume"))
            {
                paused = false;
            }
            else if (message.Contains("stop"))
            {
                transferSocket.Close(); // TODO 待定
            }
            Log.info(filePath, message, transferSocket);
        }
    }
}