<template>
    <div class="root-container">
        <v-list lines="one" density="compact" style="background: transparent;">
            <v-list-item class="btn list-item" v-for="(path, index) in props.paths" :key="index" @click="emits('addTab', path.path, path.name)">
                <span class="unselectable">
                    <span v-if="$props.device">~</span>
                    {{ path.path + path.name }}
                </span>
            </v-list-item>
            <!--add path-->
            <v-list-item v-if="!props.device" class="btn" @click="selectFolderEntry">
                <v-hover>
                    <template v-slot:default="{ isHovering, props }">
                        <div v-bind="props" class="flex-horizontal gap-10">
                            <add-one theme="outline" size="25" :fill="isHovering ? '#fff' : '#555'"/>
                            <span>Add Path</span>
                        </div>
                    </template>
                </v-hover>
            </v-list-item>
        </v-list>
    </div>
</template>

<script setup lang="ts">
import { AddOne } from "@icon-park/vue-next";
import { selectDirectory } from "@/utils/files";
import { splitPath } from "@/utils/path";
    // store
// DATA
type pathItem = {
    path: string;
    name: string;
}

const props = defineProps<{ paths: pathItem[], device: boolean }>(); // tip ts 语法要求
const emits = defineEmits(['addTab']);

// FUNC
async function selectFolderEntry() {
    try {
        const fullPath: string = await selectDirectory();
        const [path, name] = splitPath(fullPath);
        emits('addTab', path, name);
    } catch (error) {
        console.error("Error selected folder", error);
    }
}

</script>

<style scoped>

.root-container {
    background-color: #222;
    border-radius: 10px;
    box-shadow: 0px 0px 7px 1px #ffffff30;
    transform: translateY(5px);

    span {
        color: #ccc;
        font-weight: bold;
    }
}

.list-item:hover {
    background-color: #ffffff30;
}

</style>        