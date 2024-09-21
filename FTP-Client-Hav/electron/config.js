import path from "path";
const basePath = path.resolve(process.cwd(), "./");

const DEBUG = false;

export const entryConfig = {
  windowConfig: {
    width: 320,
    height: 450,
    minWidth: 320,
    minHeight: 450,

    resizable: false | DEBUG, // 禁止改变窗口大小
    frame: false | DEBUG,
    icon: path.join(basePath, "/public/icon.png"),

    webPreferences: {
      // preload: path.join(basePath, "/electron/preload.js"),
      nodeIntegration: true,
      contextIsolation: false,
      // enableRemoteModule: true, // 如果需要使用 remote 模块
    },
  },
  resizeable: false,
  startRoute: ""
};

export const mainConfig = {
  windowConfig: {
    width: 1200,
    height: 700,
    minWidth: 1150,
    minHeight: 500,

    frame: false | DEBUG,
    icon: path.join(basePath, "/public/icon.png"),

    webPreferences: {
      // preload: path.join(basePath, "/electron/preload.js"),
      nodeIntegration: true,
      contextIsolation: false,
      // enableRemoteModule: true, // 如果需要使用 remote 模块
    },
  },
  resizeable: true,
  startRoute: "/main",
};
