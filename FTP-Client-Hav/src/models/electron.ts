import { fileType } from "@/const";
import { getExtension } from "@/utils/files";

export type fileDetail = {
    name: string;
    isDirectory: boolean;
    size: number;
    mtime: Date;
}

export class FileDetail {
  name: string;
  isDirectory: boolean;
  size: number;
  time: String;
  mtime: Date | string;
  type: number;

  constructor(
    name: string,
    size: number,
    isDirectory?: boolean,
    mtime?: Date | string
  ) {
    this.name = name;
    this.isDirectory = isDirectory ?? false; // todo 暂时未设计文件夹转移
    this.size = size;
    this.mtime = mtime ?? new Date();
    this.type = this.determineFileType();
    this.time = this.filterTime();
  }

  private determineFileType(): number {
    if (this.isDirectory) {
      return fileType["dir"];
    } else {
      const extension: string = getExtension(this.name);
      // console.debug("extension:", extension, fileType[extension]);
      return fileType[extension] !== undefined
        ? fileType[extension]
        : Object.keys(fileType).length - 1;
    }
  }

  private filterTime(): String {
    if (typeof this.mtime === "string") return this.mtime;
    const formatter = new Intl.DateTimeFormat("en-CA", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    });
    return formatter.format(this.mtime).replace(/-/g, "/");
  }
}
