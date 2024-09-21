<template>
  <div class="outside-contanier flex-center-both">
    <div class="card-container flex-center-vertical">
      <Header v-if="Header" entry/>
      <div class="drag-able flex-center-vertical">
        <h1 class="no-warp pdb-15 unselectable">FTP-Hav：Client</h1>
        <v-img class="img unselectable no-pointer" aspect-ratio="1/1" :width="200" cover src="/Avatars.png"/>
      </div>
      <v-text-field class="block"
        label="User Name" variant="outlined"
        density="compact"
        v-model="loginForm.username">
      </v-text-field>
     <v-text-field class="block mt--5"
        label="Password" variant="outlined" type="password"
        density="compact"
        v-model="loginForm.password">
      </v-text-field>
      <v-btn class="block btn mt-5" :class="{'btn-finished': loginFinish}"
        :disabled="!loginFinish"
        :loading="loginLoading"
        @click="login">
        <p class="btn-text">登录</p>
      </v-btn>
      <!-- link -->
      <a class="mt-10 unselectable" @click="gotoRegister">注册账号</a>
    </div>
  </div>
</template>

<script setup lang="ts">

import { ref, reactive, computed, defineAsyncComponent, onMounted } from 'vue';
import type { IpcRenderer } from 'electron';
import axios from "@/axios";

// tip 借助动态导入的方式，同时支持窗口 && WEB；
const isElectron = navigator.userAgent.includes('Electron');
const Header = isElectron ? defineAsyncComponent(() => import('@/components/overall/header.vue')) : null;

const ipcRenderer = ref<IpcRenderer | null>(null);

import createEmial from "@/stores/emailStore";
const email = createEmial();
import { useRouter } from 'vue-router';
const router = useRouter();

// DATA
const loginForm = reactive({
  username: "",
  password: ""
})
const loginLoading = ref(false);

const loginFinish = computed(() => {
  return loginForm.username != "" && loginForm.password != "";
})

// FUNC
onMounted(async() => {
  if(isElectron) {
    const {ipcRenderer: ipc} = await import('electron');
    ipcRenderer.value = ipc;
  }
})

async function login() {
  // console.info("Login", loginForm);
  axios.post( '/tcp/login', loginForm)
    .then(res => {
      if(res.status != 200) {
        email.sendEmail(false, res.data.err_msg);
      } else {
        email.sendEmail(true, res.data.msg);
        // console.info(res)
        setTimeout(() => {
          if(isElectron)
            ipcRenderer.value?.send('entry_mainWindow');
          else
            router.push("/main");
        }, 500);
      }
    })
    .catch(err => {
      console.error(err);
      email.sendEmail(false, err);
    });
}

  // router
function gotoRegister() {
  // router.push('/register');
  // console.log("跳转执行 /register");

  email.sendEmail(false, "此功能待开发。"); // test
}

</script>

<style scoped>

.outside-contanier {
  /* position: absolute;
  top: 0px;
  bottom: 0px;
  left: 0px;
  right: 0px;
  margin: auto; */

  margin: 20px;
  margin-top: 80px;

  /* position: absolute;
  top: 50%;
  transform: translateY(-50%); */
}

.card-container {
  width: 50%;
  min-width: 250px;
  h1 {
    color: #f0f0f0;
    font-weight: 800;
  }
  .btn {
    background-color: #BFFCF9;
  }
  .btn-finished {
    box-shadow: 0px 0px 7px 1px #ffffff80;  
    /* 白色阴影 */
  }
  .btn-text {
    font-weight: 600;
    font-size: 15px;
  }
}

/* container */

.pdb-15 {
    padding-bottom: 20%;
}

/* Tag && animation */
@keyframes opacity-cycle {
  0%, 100% {
    opacity: .5; 
  }
  50% {
    opacity: .1; 
  }
}

.img {
  position: absolute;
  top: 0px;
  animation: opacity-cycle 4s infinite; /* 4秒完成一个周期，无限循环 */
}

</style>
