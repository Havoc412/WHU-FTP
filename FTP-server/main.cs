using FTP_server.consts;

namespace Server // tip 然后 链接器 才能顺利启动
{
    class Program
    {
        static void Main()
        {
            CmdConfig.init();
            // CORE
            FtpSocketServer ftpSocketServer = new FtpSocketServer();
            ftpSocketServer.Start();
        }
    }
}