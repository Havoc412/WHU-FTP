<template>
  <div id="drag-region" class="base">
    <router-view class="drag-unable"/>
    <!-- 全局消息条 -->
    <Snackbar class="top" v-model="emailForm.snackbar_switch" 
      :message="emailForm.snackbar_message" :color="emailForm.kind_flag"
      :timeout="emailForm.time_flag ? 15000 : 2000 ">
    </Snackbar>
  </div>
</template>


<script setup>
// com
import Snackbar from "./components/overall/snackbar.vue";

import cache from './stores/routeStore';
const catchStore = cache();
import emial from './stores/emailStore';
const emailForm = emial();

</script>

<style>

/* other */
.top {
  z-index: 9999;
}

#drag-region {
  /* z-index: 1000; */
  border-top: 1px solid transparent;  /* 解决 边界塌陷 */
  -webkit-app-region: drag;
  /* width: 100vw; */
  /* height: 100vh; */

  /* background-color: #ffffff3f; */
}

.drag-able {
  -webkit-app-region: drag;
}

.drag-unable {
  -webkit-app-region: no-drag;

  /* background-color: #ffffff3f; */
}

/* info 直接全局修改样式表 */
::-webkit-scrollbar {
  width: 3px;
  /*滚动条宽度*/
  height: 5px;
  /*滚动条高度*/
  background-color: #333;
}

/*定义滑块 内阴影+圆角*/
::-webkit-scrollbar-thumb {
  -webkit-box-shadow: inset 0 0 6px rgba(202, 202, 202, 0.3);
  /*滚动条的背景颜色*/
  background-color: #ccc;
}

::-webkit-scrollbar-thumb:hover {
  cursor: grab;
}

.mark {
  background-color: #ffffff56;
}

</style>
