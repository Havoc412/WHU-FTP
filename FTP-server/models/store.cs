namespace FTP_server.models
{
    public struct fileItem
    {
        public string name { get; }
        public bool type { get; } // INFO 0: dir; 1: file

        public fileItem(string name, bool type)
        {
            this.name = name;
            this.type = type;
        }
    }
}
