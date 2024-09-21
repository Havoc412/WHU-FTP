// preload.js
import { fileURLToPath } from "url";
import { dirname } from "path";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// contextBridge.exposeInMainWorld("electron", {
//   send: (channel, data) => ipcRenderer.send(channel, data),
//   on: (channel, func) => {
//     ipcRenderer.on(channel, (event, ...args) => func(...args));
//   },
// });
