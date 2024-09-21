using System.Net.Sockets;
using FTP_server.consts;

namespace FTP_server
{
    public static class Log
    {
        private enum LogLevel
        {
            DBG = 0,
            INF = 1,
            WAR = 2,
            ERR = 3,
            FTL = 4,
            UNF = 5
        }

        private static readonly string[] Rank = new string[] {
            "DBG", "INF", "WAR", "ERR", "FTL", "UDF"
        };

        private static void WriteTimeWithRank(LogLevel rankNum)
        {
            Console.ForegroundColor = ConsoleColor.DarkGreen;

            Console.Write("[");

            DateTime now = DateTime.Now;
            // 格式化时间为 HH:mm:ss
            string timeFormatted = now.ToString("HH:mm:ss");
            Console.ForegroundColor = ConsoleColor.Black;
            Console.Write(timeFormatted + " ");

            switch(rankNum)
            {
                case LogLevel.DBG: 
                    Console.ForegroundColor = ConsoleColor.Blue; 
                    break;
                case LogLevel.INF:
                    Console.ForegroundColor = ConsoleColor.White;
                    break;
                case LogLevel.WAR:
                    Console.ForegroundColor = ConsoleColor.Yellow;
                    break;
                case LogLevel.ERR:
                    Console.ForegroundColor = ConsoleColor.Red;
                    break;
                case LogLevel.FTL:
                    Console.ForegroundColor = ConsoleColor.DarkRed;
                    break;
                default:
                    Console.ForegroundColor = ConsoleColor.Gray;
                    break;
            }
            Console.Write(Rank[(int)rankNum]);
            Console.ForegroundColor = ConsoleColor.DarkGreen;
            Console.Write("] ");
            Console.ResetColor();
        }

        private static void WriteCmd(string cmd)
        {
            Console.ForegroundColor = CmdConfig.CmdColor.ContainsKey(cmd) ? CmdConfig.CmdColor[cmd] : ConsoleColor.Black;
            Console.Write($"[{cmd}] ");
            Console.ResetColor();
        }

        private static void WriteClient(Socket client)
        {
            Console.ForegroundColor = ConsoleColor.Green;
            Console.WriteLine($" {{{client.RemoteEndPoint}}}");

            Console.ResetColor();
        }

        private static void LogMessage(LogLevel rankNum, string msg, string cmd, Socket client)
        {
            WriteTimeWithRank(rankNum);

            if(cmd != null && client != null)
            {
                WriteCmd(cmd);
                Console.Write(msg);
                WriteClient(client);
            } else
            {
                Console.WriteLine(msg);
            }
        }
        // TIP 通用日志方法，不包括 Socket 参数
        private static void LogMessage(LogLevel rankNum, string msg)
        {
            LogMessage(rankNum, msg, null, null);
        }

        // INFO 对外函数接口
        static public void debug(string msg)
        {
            LogMessage(LogLevel.DBG, msg);
        }
        static public void debug(string msg, string cmd, Socket client)
        {
            LogMessage(LogLevel.DBG, msg, cmd, client);
        }

        static public void info(string msg)
        {
            LogMessage(LogLevel.INF, msg);
        }
        static public void info(string msg, string cmd, Socket client)
        {
            LogMessage(LogLevel.INF, msg, cmd, client);
        }
        static public void warning(string msg)
        {
            LogMessage(LogLevel.WAR, msg);
        }
        static public void warning(string msg, string cmd, Socket client)
        {
            LogMessage(LogLevel.WAR, msg, cmd, client);
        }
        static public void error(string msg)
        {
            LogMessage(LogLevel.ERR, msg);
        }
        static public void error(string msg, string cmd, Socket client)
        {
            LogMessage(LogLevel.ERR, msg, cmd, client);
        }

        static public void fatal(string msg)
        {
            LogMessage(LogLevel.FTL, msg);
        }
        static public void fatal(string msg, string cmd, Socket client)
        {
            LogMessage(LogLevel.FTL, msg, cmd, client);
        }
    }
}
