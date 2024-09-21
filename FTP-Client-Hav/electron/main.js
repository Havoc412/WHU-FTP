import { app, BrowserWindow, ipcMain, dialog } from "electron";

import { entryConfig, mainConfig } from "./config.js";

// import reloader from "electron-reloader"
// reloader(module);   // test hot reloade

let mainWindow = null;

function createWindow(config) {
  mainWindow = new BrowserWindow(config.windowConfig);
  // info 使用 loadFile 加载 electron/index.html 文件
  // mainWindow.loadFile(path.join(process.cwd(), "electron/index.html")); // start
  mainWindow.loadURL("http://localhost:3000" + config.startRoute); // develope
  // mainWindow.loadURL(`file://${path.join(__dirname, "../dist/index.html")}`); // packager

  // 设置窗口是否可以由用户手动最大化。
  mainWindow.setMaximizable(config.resizeable);
  // 设置用户是否可以调节窗口尺寸
  if(config.resizeable) {

  }
};

// 在应用准备就绪时调用函数
// app.whenReady().then(createWindow);
app.whenReady().then(() => createWindow(entryConfig));

ipcMain.on("minimize-window", () => {
  mainWindow.minimize();
});
ipcMain.on("toggle-maximize", () => {
  if (mainWindow.isMaximized()) {
    mainWindow.unmaximize();
    mainWindow.webContents.send('window-unmaximized');
  } else {
    mainWindow.maximize();
    mainWindow.webContents.send('window-maximized');
  }
})
ipcMain.on("close-window", () => {
  mainWindow.close();
});

// info 进入正式的窗口
ipcMain.on("entry_mainWindow", () => {
  if (mainWindow) {
    mainWindow.on("closed", () => {
      mainWindow = null; // 清除引用
      app.whenReady().then(() => createWindow(mainConfig));
    });
    mainWindow.close();
  } else {
    app.whenReady().then(() => createWindow(mainConfig));
  }
})

// info file elections
import fs from "fs";
import path from "path";
import { fsFilterType } from "../src/const/electron.js";

ipcMain.on("get-files", (event, filterType, dirPath) => {
  fs.readdir(dirPath, (err, files) => {
    if(err) {
      console.error(err)
      event.sender.send('files-list', "Error reading directory");
      return;
    }

    let targetDetails = [];
    let completedFiles = 0;
    files.forEach(file => {
      const filePath = path.join(dirPath, file);
      fs.stat(filePath, (err, stats) => {
        if(err)
          console.error(`Error getting stats of ${file}`, err);
        else {
          let isDir = stats.isDirectory();
          if (
            filterType === fsFilterType.ALL ||
            filterType === fsFilterType.Directory && isDir  ||
            filterType === fsFilterType.File && !isDir
          )
            targetDetails.push({
              name: file,
              isDirectory: stats.isDirectory(),
              size: stats.size,
              mtime: stats.mtime,
            });
        }
        completedFiles++;
        if(completedFiles === files.length) {
          // console.log(targetDetails);
          event.sender.send('files-list', targetDetails);
        }
      })
    })
  })
});

ipcMain.on("dialog:openFile", async(event) => {
  let parentWindow = BrowserWindow.getFocusedWindow(); // 获取当前聚焦的窗口
  if (parentWindow) parentWindow.setEnabled(false);
  const options = {
    title: "Select a Folder",
    properties: ["openDirectory"],
    parent: parentWindow, // 获取当前聚焦的窗口作为父窗口
    modal: true, // 设置为模态，对 macOS 是必要的
  };
  const { canceled, filePaths } = await dialog.showOpenDialog(options);

  const focusListener = () => {
    dialog.focus();
    // todo 闪烁的效果？
  }
  parentWindow.on("focus", focusListener);
  // 恢复父窗口
  if (parentWindow) {
    parentWindow.setEnabled(true);
    parentWindow.removeListener("focus", focusListener); // 移除监听器以避免内存泄漏
    parentWindow.focus();
  }
  
  if (canceled) {
    console.log(1);
    event.sender.send("dialog:fileSelected", -1);
  } else {
    console.log(2, filePaths);
    event.sender.send("dialog:fileSelected", filePaths[0]);
  }
});
