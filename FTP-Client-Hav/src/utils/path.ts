
export function splitPath(path: string) {
  // console.info("path.ts", path);
  // 首先替换所有的 `/` 为 `\`
  const normalizedPath = path.replace(/\\/g, "/");
  const lastSlashIndex = normalizedPath.lastIndexOf("/");
  if (lastSlashIndex === -1) {
    return [path, ""]; // 如果没有找到`\`，返回整个路径和空字符串
  }
  // console.info("path.ts-2", normalizedPath, lastSlashIndex);
  const directory = normalizedPath.substring(0, lastSlashIndex+1); // 目录部分
  const filename = normalizedPath.substring(lastSlashIndex + 1); // 文件名部分
  return [directory, filename];
}