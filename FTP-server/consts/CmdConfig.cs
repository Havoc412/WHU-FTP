using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FTP_server.consts
{
    //public struct cmdInfo
    //{
    //    public string Name;
    //    public ConsoleColor Color;

    //    public cmdInfo(string name, ConsoleColor color = ConsoleColor.Black)
    //    {
    //        Name = name;
    //        Color = color;
    //    }
    //}
    public static class CmdConfig
    {
        // 基本枚举
        //public const int LINK = 0;
        //public const int LOGIN = 1;

        //public static readonly cmdInfo[] CmdInfo =
        //{
        //    new cmdInfo("LINK", ConsoleColor.Yellow),
        //    new cmdInfo("LOGIN")
        //};
        public static Dictionary<string, ConsoleColor> CmdColor = new Dictionary<string, ConsoleColor>();

        // INFO 逐个添加
        private static void add(string key, ConsoleColor color = ConsoleColor.Black)
        {
            CmdColor[key] = color;
        }

        public static void init()
        {
            add("LINK", ConsoleColor.Yellow);
            add("SUB-LINK", ConsoleColor.Yellow);
            add("LOGIN");
        }
    }
}
