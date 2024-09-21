<!--info 主要处理任务列表之间的关系-->
<template>
    <div class="flex-center-vertical root-container"
        :style="{
            'height': props.height + 'px'
        }">
        <template v-for="(item, index) in taskList" :key="'task-'+index">
            <task-item :path="item.path" :name="item.fileName" :size="item.size" :session-id="item.sessionId" :reverse="item.type"
                @finish="(path: string, fileName: string, size: number, type: boolean) => emits('finish', path, fileName, size, type)"/>
        </template>
        <div v-if="taskList.length == 0" class="message-container flex-center-vertical">
            <span>No Task, Waiting...</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

import taskItem from "./terminal/taskItem.vue";
import type { WSprepareType } from "@/models/props";
    // store
// DATA
const props = defineProps<{
    height: number;

    wsPrepare: WSprepareType | null;
    startWs: boolean;
}>();
const emits = defineEmits(['finish']);

const taskList = ref<WSprepareType[]>([]);

// FUNC
// question 并发时会发生什么
watch(() => props.startWs, () => {
    if(props.wsPrepare == null) {
        console.error("Error when startWS = true, but wsPrepare = null ?");
        return;
    }
    taskList.value.push(props.wsPrepare);
});


</script>

<style scoped>

.root-container {
    padding: 5px 10px;
    overflow-y: auto;
}

.message-container {
    height: 100%;
    span {
        font-size: 20px;
        font-weight: 900;
    }
}

</style>        