// 导入所需的 Node.js 模块
import { ipcRenderer } from "electron";
import { FileDetail, type fileDetail } from "@/models/electron";
// import { fsFilterTypr } from "@/const/electron.js";

// 创建并导出一个函数用于获取文件
export async function fetchFiles(filterType: number, path: string): Promise<FileDetail[]> {
  return new Promise((resolve, reject) => {
    ipcRenderer.send("get-files", filterType, path);

    ipcRenderer.once(
      "files-list",
      (event, filesList: fileDetail[] | string) => {
        if (typeof filesList === "string") {
          console.error(filesList); // 处理错误
          reject(new Error(filesList));
        } else {
          const files = filesList.map(
            (file) =>
              new FileDetail(file.name, file.size, file.isDirectory, file.mtime)
          );
        //   console.info("Received files list:", files);
          resolve(files); // 在这里解析带有文件详细信息的数组
        }
      }
    );
  });
}

export async function selectDirectory(): Promise<string>  {
  return new Promise((resolve, reject) => {
    ipcRenderer.send("dialog:openFile");

    ipcRenderer.once("dialog:fileSelected", (event, path) => {
      if(typeof path === 'number') {
        const error = "Error select directory.";
        console.error(error);
        reject(new Error(error));
      }
      // console.info("ts", path);
      resolve(path);
    });
  })
}

// tag Common

export function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";

  const k = 1024; // 或者使用 1000，取决于你的转换需求（1024是常用于计算机存储）
  const sizes = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${Math.ceil(bytes / Math.pow(k, i))} ${sizes[i]}`;
}

export function getExtension(fileName: string): string {
  const lastIndex = fileName.lastIndexOf(".");
  return lastIndex !== -1 ? fileName.substring(lastIndex + 1) : "";
}
