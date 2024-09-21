public static class stringHandler
{
    public static (string, string) SplitCommand(string str)
    {
        int spaceIndex = str.IndexOf(' ');
        if (spaceIndex == -1)
        {
            return (str, string.Empty);
        }
        else
        {
            string firstPart = str.Substring(0, spaceIndex);
            string secondPart = str.Substring(spaceIndex + 1);
            return (firstPart, secondPart);
        }
    }

    public static string[] SplitCommandArray(string argument)
    {
        // INFO 目前就是用于 UPLOAD 指令的三段
        // 找到第一个空格的位置
        int firstSpaceIndex = argument.IndexOf(' ');
        if (firstSpaceIndex == -1)
        {
            throw new ArgumentException("The input string does not contain enough parts.");
        }

        // 找到第二个空格的位置
        int secondSpaceIndex = argument.IndexOf(' ', firstSpaceIndex + 1);
        if (secondSpaceIndex == -1)
        {
            throw new ArgumentException("The input string does not contain enough parts.");
        }

        string[] res = { argument[..firstSpaceIndex], argument.Substring(firstSpaceIndex + 1, secondSpaceIndex - firstSpaceIndex - 1), argument[(secondSpaceIndex + 1)..] };

        Console.WriteLine(res[0]+"|");
        Console.WriteLine(res[1]+"|");
        Console.WriteLine(res[2] + "|");
        return res;
    }
}