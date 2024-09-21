<template>
    <div class="flex-center-horizontal task-container block gap-20">
        <span class="no-warp file-name">
            {{ props.name }}
            <v-tooltip activator="parent" location="end" offset="5">
                {{ props.name }}
            </v-tooltip>
        </span>
        <v-progress-linear v-model="process" :reverse="props.reverse" stream>
            <!-- <template v-slot:default>
                <div class="time-container">
                    {{ timeDifference }}
                </div>
            </template> -->
        </v-progress-linear>
        <span class="no-warp">{{ process + "%" }}</span>
        <div class="flex-center-horizontal gap-5">
            <component :is="sending ? Loading : PauseOne" 
                class="btn"
                :class="{
                    'loading': sending,
                    'finish': finish
                }"
                theme="outline" size="25" fill="#ddd" 
                @click="pauseOrResume"
            />
            <close-one class="btn" theme="outline" size="25" fill="#ddd"/>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from "vue";
import { Loading, PauseOne, CloseOne } from "@icon-park/vue-next";
import { log } from "console";

    // store
// DATA
const props = defineProps({
    path: String,
    name: String,
    size: Number,
    sessionId: String,
    reverse: Boolean, // info 1: local->server; 0: server->local ;; 因为 stream 的效果。
});
const emits = defineEmits(['finish']);

const sending = ref(true);
const sentLength = ref(0);

const startTime = ref(null);
const endTime = ref(null);

// FUNC

const finish = computed(() => {
    return sentLength.value >= props.size;
})

const process = computed(() => {
    if(props.size == 0)
        return (finish.value * 100);
    return Math.round(sentLength.value / props.size * 100);
})

// const buffer = computed(() => { //  :buffer-value="buffer"
//     if(props.size == 0)
//         return (finish.value * 100);
//     return Math.round(1024 / props.size * 100);
// })

// const timeDifference = computed(() => {
//     if (!startTime.value || !endTime.value) {
//         return '00:00:00';
//     }
//     // 计算差异（毫秒）
//     const diff = endTime.value - startTime.value;

//     // 将毫秒转换为小时、分钟和秒
//     let seconds = Math.floor((diff / 1000) % 60);
//     let minutes = Math.floor((diff / (1000 * 60)) % 60);
//     let hours = Math.floor((diff / (1000 * 60 * 60)) % 24);

//     // 前导零格式化
//     const pad = (num) => num.toString().padStart(2, '0');

//     // 返回格式化的时间差
//     return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
// });

let ws;
onMounted(() => {
    startTime.value = new Date();   // 获取创建时间；
    connectWebSocket();
})

// 每一个子任务交给每一个组件独立完成。
function connectWebSocket() { // question 异步？
    ws = new WebSocket(`ws://127.0.0.1:8080/v1/ws/${props.reverse ? 'download' : 'upload'}?sessionId=${props.sessionId}`);
    ws.onmessage = function(event) {
        const data = JSON.parse(event.data)
        if(data.state == 2) {
            sending.value = false;
            endTime.value = new Date();
            emits('finish', props.path, props.name, props.size, props.reverse);
        } else {
            sentLength.value = data.sent_byte;
        }
    };
    ws.onerror = function(event) {
        console.error("WebSocket error:", error)
    }
    ws.onclose = function() {
        console.log("WebSocket connection closed");
    };
}

function pauseOrResume() {
    if(ws.readyState === WebSocket.OPEN) {
        const command = sending.value ? "pause" : "resume";
        ws.send(JSON.stringify({command: command}))
        sending.value = !sending.value;
    }
}


onUnmounted(() => {
    ws.close();  // 组件卸载时关闭 WebSocket 连接
});

</script>

<style scoped>

.task-container {
    span {
        font-size: 20px;
        font-weight: bold;
    }

    .file-name {
        min-width: 100px;
        max-width: 150px;
        white-space: nowrap; /* 保持文本不换行 */
        overflow: hidden; /* 隐藏超出部分 */
        text-overflow: ellipsis; /* 使用省略号表示超出的文本 */
    }

    .loading {
        animation: rotate 2s linear infinite;
    }

    .finish {
        display: none;
    }

    .time-container {
        color: #fff;
        font-size: 20px;
        font-weight: bold;
    }
}

@keyframes rotate {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

</style>        