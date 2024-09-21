<template>
    <div class="header flex-end-horizontal top drag-unable">
        <v-btn class="common" size="x-small" icon="mdi-minus" variant="text" @click="minimize"/>
        <v-btn v-if="!entry" class="common" size="x-small" icon="mdi-checkbox-blank-outline" variant="text" @click="toggleMaximize"/>
        <v-btn class="special" size="x-small" icon="mdi-close" variant="text" @click="close"/>
    </div>
</template>

<script>
// import { ref } from 'vue';
import axios from '@/axios';
import emailStore from '@/stores/emailStore';
import { ipcRenderer } from 'electron';

const email = emailStore();

export default {
  data() {
    return {
      fullscreen: false
    }
  },
  props: {
    entry: Boolean
  },
  methods: {
    toggleMaximize() {
      ipcRenderer.send('toggle-maximize')
    },
    minimize() {
      ipcRenderer.send('minimize-window')
    },
    close() {
      if(!this.entry) {
        axios.get('tcp/exit')
          .then(res => {
            if(res.status != 200)
              email.sendEmail(res.status != 200, res.data.err_msg);
            else
              email.sendEmail(res.status != 200, res.data.msg)
          })
          .catch(err => {
            console.error(err);
          });
        setTimeout(() => {
          ipcRenderer.send('close-window')
        },)
      }
      
    }
  }
}
</script>

<style scoped>

.header {
    position: absolute;
    top: 2px;
    right: 2px;
}

.common:hover {
    background-color: #93939340 !important;
    border-radius: 5px 10px;
}

.special:hover {
    background-color: #ff0000c9 !important;
    border-radius: 5px 10px;
}

</style>        