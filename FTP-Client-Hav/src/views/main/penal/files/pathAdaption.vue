<template>
    <div class="root-container flex-horizontal gap-5 block mr-1">
        <template v-for="(item, index) in dirList" :key="index">
            <div class="chip-container flex-horizontal btn" @click="emits('addTab', item.path, item.name)">
                <div class="">{{ item.name }}</div>
                <v-icon class="v-btn" icon="mdi-chevron-right"/>
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { splitPath } from "@/utils/path";
    // store
// DATA
const props = defineProps<{ fullPath: string }>();
const emits = defineEmits(['addTab']);

type dirItem = {
    name: string,
    path: string
}

const dirList = ref<dirItem[]>([]);

// FUNC
onMounted(() => {
    init();
    // adjustAlignment();
})

function init() {
    let path = props.fullPath;
    let name = "";
    while(path.includes('/')) {
        [path, name] = splitPath(path);
        dirList.value.unshift({
            name: name,
            path: path
        })
        path = path.slice(0, -1);
        // console.info(path, name);   // test
    }
}

// todo 更好的自适应方法。
// function adjustAlignment() {
//     const container = document.querySelector('.root-container');
//     if(!container)
//         return;
//     const totalWidth = Array.from(container.children).reduce((total, child) => total + child.offsetWidth, 0);
//     if (totalWidth > container.offsetWidth) {
//         container.style.justifyContent = 'flex-end'; // 越界时右对齐
//     } else {
//         container.style.justifyContent = 'flex-start'; // 未越界时左对齐
//     }
// }

// // 在页面加载和窗口大小变化时调用
// window.onload = adjustAlignment;
// window.onresize = adjustAlignment;

</script>

<style scoped>

.root-container {
    overflow: hidden;
    height: 40px;
    justify-content: flex-end;
}

.chip-container {
    padding: 0 5px;
    border: 1px solid transparent;
}

.chip-container:hover {
    height: 30px;
    border: 1px solid #333;
    border-radius: 5px;
    /* box-sizing: border-box; */
}


</style>        