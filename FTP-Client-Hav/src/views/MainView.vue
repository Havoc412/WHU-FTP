<template>
    <div class="root-contianer" @mousemove="" ref="splitPanelRef">
        <Header/>
        <!--main window-->
          <!--tools container-->
        <ToolsBar/>
        <div class="handle-y-tools-bar"></div>
          <!--work container-->
        <div class="work-container flex-vertical">
            <!--INFO file box-->
            <div class="flex-horizontal">
              <!--local file-->
              <FilesPenal
                ref="localRef"
                v-model:tab="tabLocal"
                :width="panelList[0].width" 
                :height="panelList[0].height"
                :device="false"
                @exchange-files="(selectedList) => exchangeFiles(selectedList , false)"
              />
              <div class="handle-y"></div>
              <!--FTP file-->
              <FilesPenal
                ref="serverRef"
                v-model:tab="tabServer"
                :width="panelList[1].width"
                :height="panelList[1].height"
                :device="true"
                @exchange-files="(selectedList) => exchangeFiles(selectedList , true)"
              />
            </div>
            <div class="handle-x"></div>
            <!--task & state bar-->
            <TerminalPenal
              :height="panelList[2].height"
              :ws-prepare="wsPrepare"
              :start-ws="startWS"
              @finish="addFileItem"
            />
        </div>
    </div>
</template>

<script lang="ts" setup name="SplitPanel">
import { onMounted, reactive, ref } from "vue";
import { useResizeObserver } from "@vueuse/core";

import axios from "@/axios";

import Header from "./main/header.vue"; // update 之后再合并
import ToolsBar from "./main/toolsBar.vue";
import FilesPenal from "./main/penal/filesPenal.vue";
import TerminalPenal from "./main/penal/terminalPenal.vue";
import { type uploadApiType, type downloadApiType } from "@/models/api";
import { type fileItemType, type WSprepareType } from "@/models/props";
import emailStore from "@/stores/emailStore";
import { Tv } from "@icon-park/vue-next";

const email = emailStore();

// DATA
    // type
type mouseDataType = {  // todo 拖动设计用。
  id: number;
  width: number;
  height: number;
  startX: number;
  clientX: number;
  startY: number;
  clientY: number;
};

type panelType = {
  id: number;
  width: number;
  height: number;
  type: string;
};
    // consts
const MIN_WIDTH = 300;
const MIN_HEIGHT = 200;

const splitPanelRef = ref<HTMLElement | null>(null); // root element
const isDown = ref(false); // todo 拖动设计用。
  // 整个容器的宽高
const container = reactive({
  width: 0,
  height: 0
});

  // core data
const tabLocal = ref("");
const tabServer = ref("");

const localRef = ref(null); // 直接通过 ref 获取控件，然后用于调用内部函数。
const serverRef = ref(null);

// 窗口设定 // info 现阶段就先基本的三个窗口。
const panelList = ref<panelType[]>([
  {id: 1, width: 200, height: 200, type: "local files"},
  {id: 2, width: 200, height: 200, type: "ftp files"},
  {id: 3, width: 200, height: 200, type: "terminal"},
]);

// WS SETTING
const wsPrepare = ref<WSprepareType | null>(null);
const startWS = ref<boolean>(false); // info 通过反转的方式来触发。

// FUNC
onMounted(() => {
  init();
})

// info 监听主窗口的大小。
useResizeObserver(splitPanelRef, (entries) => {
  const entry = entries[0];
  const { width, height } = entry.contentRect;
  container.width = width;
  container.height = height;
  init();
});

function init() { // update 只适用于三个窗口，需要兼容拖动后状态。
  container.width = splitPanelRef.value?.clientWidth || 0;
  container.height = splitPanelRef.value?.clientHeight || 0;
  // 1.
  panelList.value[2].height = MIN_HEIGHT;

  panelList.value[0].height = container.height - panelList.value[2].height;
  panelList.value[1].height = container.height - panelList.value[2].height;

  let width = (container.width - 40 - 4) / 2;
  panelList.value[0].width = Math.max(MIN_WIDTH, width);
  panelList.value[1].width = container.width - 65 - panelList.value[0].width;
}

function exchangeFiles(selectedList: fileItemType[], type: boolean) {
  console.debug(tabLocal.value, tabServer.value, selectedList);
  // 1. check the list exist
  let success, res;
  if(!type && serverRef.value) { // local -> server
    [success, res] = serverRef.value.checkFileNames(selectedList);
  } else if(localRef.value) {
    [success, res] = localRef.value.checkFileNames(selectedList);
  } else {
    email.sendEmail(false, "异常错误，无法鉴定文件存在与否");
    return;
  }

  if(!success) {
    if(typeof res === "string")
      email.sendEmail(false, res);
    else {
      console.info("Existed: ", res);
      email.sendEmail(false, "存在文件名冲突。"); // todo 允许重开小窗来批量改名 && 允许修改本地文件名。
    }
    return;
  }

  // 2. API - HTTP
  if(!type) {  // Local -> server
    selectedList.forEach((file: fileItemType) => { // todo 现阶段还是一个个处理会比较好。
      let uploadForm: uploadApiType = {
        targetpath: tabServer.value.slice(1),
        localfilepath: tabLocal.value + "/" + file.name
      }
      console.debug(uploadForm); // test
      axios.post("/ftp/upload", uploadForm)
        .then(res => {
          if(res.status != 200) {
            email.sendEmail(false, res.data.err_msg); // 如果有问题，那就是容量、分配等等。
          } else {
            email.sendEmail(true, res.data.msg);

            // info WebSocket
            wsPrepare.value = {
              type: false,
              sessionId: res.data.sessionId,

              path: tabServer.value,
              fileName: file.name,
              size: file.size,
            }
            startWS.value = !startWS.value;  // test close                     
          }
        })
        .catch(err => {
          console.error(err);
          email.sendEmail(false, err);
        });
    })
  } else { // server -> local
    selectedList.forEach((file: fileItemType) => { // todo 现阶段还是一个个处理会比较好。
      let downloadForm: downloadApiType = {
        targetpath: tabServer.value + "/" + file.name,
        savepath: tabLocal.value + "/" + file.name
      }
      axios.post("/ftp/download", downloadForm)
        .then(res => {
          if(res.status != 200) {
            email.sendEmail(false, res.data.err_msg);
          } else {
            email.sendEmail(true, res.data.msg);

             // info WebSocket
            wsPrepare.value = {
              type: true,
              sessionId: res.data.sessionId,

              path: tabLocal.value,
              fileName: file.name,
              size: file.size,
            }
            startWS.value = !startWS.value;
          }
        })
        .catch(err => {
          console.error(err);
          email.sendEmail(false, err);
        });
    })
  }
}

function addFileItem(tabValue: string, fileName: string, size: number, type: boolean) {
  if(type && localRef.value) {
    localRef.value.addItem(tabValue, fileName, size);
  } else if(serverRef.value) {
    serverRef.value.addItem(tabValue, fileName, size);
  } else {
    console.info("No ref available.")
  }
} 

</script>

<style scoped>
.root-contianer {
  width: 100%;
  margin-top: 40px;
  height: calc(100vh - 45px); /* 数值待确定 */
  display: flex;
  /* .panel-item {
    position: relative;
    border-right: 1px solid #ccc;
  } */
}

.handle-y-tools-bar {
    width: 4px;
    height: 100%;
    border-right: 1px solid #333;
    /* cursor: e-resize; */
    user-select: none;
}

.work-container {
  width: 100%;

  .handle-y {
    /* position: absolute; */
    width: 2px;
    height: 100%;
    border-right: 1px solid #333;
    cursor: e-resize;
    user-select: none;
  }

  .handle-x {
    width: 100%;
    height: 4px;
    border-top: 1px solid #333;
    cursor: s-resize;
    user-select: none;
  }
}

.file-container {
  padding: 5px;
}

.terminal-container {
  width: 100%;
  padding: 5px;
}

</style>

