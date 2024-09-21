using System.Net;

namespace FTP_server.consts
{
    public static class ftpConfig
    {
        public const int PORT = 8000;
        public static readonly IPAddress LISTENER_IP_ADDRESS = IPAddress.Any; // 监听所有网络接口
        public const char MESSAGE_EOF = '\n';
    }
}
