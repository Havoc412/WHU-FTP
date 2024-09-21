export type errRes = {
    err_code: number;
    err_msg: string;
}

export type fileDetail = {
    Permissions: string;
    Owner: string;
    Group: string;
    Size: number; // Bytes
    Modified: string;
    Name: string;
}

export type listRes = {
    msg: string;
    dir_list: fileDetail[];
    file_list: fileDetail[];
}

export type uploadApiType = {
    targetpath: string, // server 目标文件夹
    localfilepath: string // 本地目标文件
}

export type downloadApiType = {
  targetpath: string; // server 目标文件
  savepath: string; // 本地文件
};