<template>
    <div class="flex-vertical"
        :style="{
            'width': props.width + 'px',
            'height': props.height + 'px'
        }">
        <div class="flex-horizontal tab-container">
            <v-menu v-model="showRecommandFlag" open-on-hover :close-on-content-click="false">
                <template v-slot:activator="{ props }">
                    <div class="flex-center-horizontal mr-2" v-bind="props">
                        <component :is="$props.device ? DataServer : NewComputer" 
                            theme="outline" size="25" fill="#ddd"
                        />
                    </div>
                </template>
                <PathRecommand :paths="pathRecommand" :device="props.device" @add-tab="addTabListValue"/>
            </v-menu>
            <!-- update 原方案 <div class="flex-center-horizontal mr-2 btn"
                @click="showRecommandFlag = !showRecommandFlag">
                <component :is="$props.device ? DataServer : NewComputer" 
                    theme="outline" size="25" fill="#ddd"
                />
            </div>
            <PathRecommand v-show="showRecommandFlag" 
                class="recommand-container top" 
                :paths="pathRecommand" @add-tab="addTabListValue"
            /> -->
            <v-tabs v-model="tab"
                show-arrows
                center-active>
                <template v-for="(item, index) in tabListValue" :key="index">
                    <v-hover>
                        <template v-slot:default="{ isHovering, props }">
                            <v-tab v-bind="props" :value="item.path+item.name">
                                {{ item.name == "" ? "root" : item.name }}
                                <close-small v-show="isHovering" class="tab-close" theme="outline" size="20" fill="#aaa" 
                                    @click.stop="closeTab(index, item.path+item.name)"
                                />
                            </v-tab>
                        </template>
                    </v-hover>
                </template>
            </v-tabs>
            <v-spacer/>
            <!--info 工具栏-->
            <ToolsBar :device="$props.device" :search-outside="search" 
                @search-update="(searchInput: string) => { search = searchInput }"
                @get-selected-files="getSelectedFiles"
            />
        </div>
        <v-tabs-window v-model="tab">
            <template v-for="(itemTab, index) in tabListContent" :key="index">
                <v-tabs-window-item :value="itemTab.value">
                    <!--show full path-->
                    <div class="flex-divide-horizontal no-warp">
                        <!-- <div>{{ itemTab.value }}</div> -->
                         <PathAdaption :fullPath="itemTab.value" @add-tab="addTabListValue"/>
                    </div>
                    <!--tables-->
                    <v-data-table-virtual 
                        v-model:search="search"                        
                        v-model:sort-by="sortBy"

                        v-model="selected"
                        return-object

                        :headers="tableHeader"
                        :item-value="item => `${item.name}-${item.size}`"
                        density="compact"
                        show-select
                        hover

                        class="data-container"

                        :height="$props.height - 90"
                        :items="itemTab.list">
                        <!--默认的排序算法-->
                        <!-- <template v-slot:header.id="{ column }">
                            {{ column.title.toUpperCase() }}
                        </template> -->
                        <template v-slot:headers="{ someSelected, allSelected, selectAll, toggleSort, columns, getSortIcon }">
                            <div style="display: none;">{{ selecteAllInBox = selectAll }}</div><!--tip 截获！-->
                             <tr>
                                <td>
                                    <v-icon class="checkbox" size="24" @click="selectAll(!allSelected)">
                                        {{ allSelected ? 'mdi-checkbox-marked' : 
                                            (someSelected ? 'mdi-minus-box' : 'mdi-checkbox-blank-outline')}}
                                    </v-icon>
                                </td>
                                <td class="no-warp"
                                    v-for="(item, index) in tableHeader"
                                    @click="toggleSort(columns[index+1])">
                                    {{ item.title }}
                                    <v-icon v-show="sortBy[0]?.key==item.key">
                                        {{ getSortIcon(columns[index+1]) }}
                                    </v-icon>
                                </td>
                             </tr>
                        </template>
                        
                        <template v-slot:item.name="{ item }">
                            <div class="flex-horizontal gap-10">
                                <FileIcon :type="item.type" size="20"/>
                                <span>{{ item.name }}</span>
                            </div>
                        </template>

                        <template v-slot:item.type="{ item }">
                            <span class="no-warp">{{ fileTypeDescribe[item.type] }}</span>
                        </template>

                        <template v-slot:item.size="{ item }">
                            <span class="no-warp" v-if="!item.isDirectory">{{ formatBytes(item.size) }}</span>
                        </template>

                        <template v-slot:item.actions="{ item }">
                            <increase v-if="item.isDirectory" class="btn" theme="outline" size="20" fill="#ddd" @click="addTabListValue(itemTab.value+'/', item.name, $event)"/>
                            <Editor v-else class="btn" theme="outline" size="20" fill="#ddd"/>
                        </template>

                    </v-data-table-virtual>
                </v-tabs-window-item>
            </template>
        </v-tabs-window>     
        <!--test <pre>{{ sortBy }}</pre>    -->
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import axios from "@/axios";
import createEmial from "@/stores/emailStore";
const email = createEmial();

import { fileType, fileTypeDescribe, localPath, ftpServerPath } from "@/const";
import { fsFilterType } from "@/const/electron.js";
import { FileDetail } from "@/models/electron";
import { type listRes } from "@/models/api";
import { type fileItemType } from "@/models/props";
import { fetchFiles, formatBytes } from "@/utils/files";
    // com
import { NewComputer, DataServer, CloseSmall, Increase, Editor } from "@icon-park/vue-next";
import FileIcon from "./files/fileIcon.vue";
import ToolsBar from "./files/toolsBar.vue";
import PathRecommand from "./files/pathRecommand.vue";
import PathAdaption from "./files/pathAdaption.vue";
// import HoverBtn from "@/components/overall/hoverBtn.vue";

// DATA
type dirItem = {
    path: string;
    name: string;
}

type fileList = {
    value: string;
    list: FileDetail[];
}

const props = defineProps<{
    tab: string,
    width: number,
    height: number,
    device: boolean, // update 现阶段只有，0：local；1：server
}>();
const emits = defineEmits(['update:tab', 'exchangeFiles']);

const pathRecommand = ref<dirItem[]>([]);
const search = ref('')
const sortBy = ref([{key: 'type', order: 'asc'}]);
const tab = ref<string | null>(null);
const tabListValue = ref<dirItem[]>([]);
const tabListContent = ref<fileList[]>([])
const selected = ref<fileItemType[]>([]); // info 存储多选框选中的文件名。
const tableHeader = [
    { title: "文件名", align: "start", key: "name" },
    { title: "修改时间", align: "start", key: "time" },
    { title: "类型", align: "start", key: "type" },
    { title: "大小", align: "start", key: "size" },
    { title: "", key: "actions", sortable: false },
]

const showRecommandFlag = ref(false);

// FUNC
onMounted(() => {
    // pathRecommand init
    pathRecommand.value = props.device ? ftpServerPath : localPath;
    if(!props.device) {
        tabListValue.value.push({ path: "C:/Users/Havoc/Desktop/HavProject/WHU-homework/FTP/", name: "FTP-client-store" });
        tab.value = "C:/Users/Havoc/Desktop/HavProject/WHU-homework/FTP/FTP-client-store";
    }
    else {
        tabListValue.value.push({ path: "", name: "" });
        tab.value = "";
    }

    tabListValue.value.forEach((item: dirItem) => {
        fetchFilesEntrance(item.path + item.name);
    })
})

watch(() => tab.value, () => {
    emits('update:tab', tab.value);
})

// todo 因为异步，需要修改删除部分。
async function fetchFilesEntrance(path: string) {
    console.info("fetch files", path);
    if(!props.device) {
        try {
            const filesList: FileDetail[] = await fetchFiles(fsFilterType.ALL, path);
            tabListContent.value.push({
                value: path,
                list: filesList
            })
        } catch (error) {
            // todo 错误处理
            console.error("Error fetching files:", error);
            email.sendEmail(false, "Error fetching files: " + path);
        }  
    } else {
        try {
            axios.get(`/ftp/list/${path}`)
                .then(res => {
                    if(res.status != 200)
                        email.sendEmail(false, res.data.err_msg);
                    else {
                        let filesList: FileDetail[] = [];
                        const listres: listRes = res.data;
                        if(listres.dir_list)
                            listres.dir_list.forEach(dir => {
                                filesList.push(
                                    new FileDetail(dir.Name, dir.Size, true, dir.Modified)
                                )
                            });
                        if(listres.file_list)
                            listres.file_list.forEach(file => {
                                filesList.push(
                                    new FileDetail(file.Name, file.Size, false, file.Modified)
                                )
                            });
                        tabListContent.value.push({
                            value: path,
                            list: filesList
                        })
                    }
                })
                .catch(err => {
                    console.error(err);
                    email.sendEmail(false, err);
                });
        } catch (error) {
            // todo 错误处理
            console.error("Error fetching files:", error);
            email.sendEmail(false, "Error fetching files: " + path);
        }
    }
}

function addTabListValue(path: string, name: string, event: Object) {
    console.info("file penal", path, name);
    const temp: dirItem = {
        path: path,
        name: name
    }
    const value = path + name;
    if(tabListValue.value.some(item => item.path === temp.path && item.name === temp.name)) {
        tab.value = value;
        showRecommandFlag.value = false;
        return;
    }
    tabListValue.value.push(temp);
    fetchFilesEntrance(value);

    if(!event?.ctrlKey) // update 目前只有 行末按钮有效。
        tab.value = value;

    if(showRecommandFlag.value)
        showRecommandFlag.value = false;
}

function closeTab(index: number, value: string) {
    tabListValue.value.splice(index, 1);
    tabListContent.value = tabListContent.value.filter(content => content.value != value);

    console.info(value, tab.value);
    if(value != tab.value)
        return
    else
        tab.value = tabListContent.value[index-1]?.value;
}

function getSelectedFiles() {
    if(selected.value.length == 0) {
        email.sendEmail(2, "请先选中转发文件");
        return;
    }
    emits('exchangeFiles', selected.value);

    // clear checkbox
    triggerSelectAll(false); // !!! 成功 !!!
}

// info 直接交互添加项目
const selecteAllInBox = ref<((value: boolean) => void) | null>(null);
const triggerSelectAll = (value : boolean) => {
    if(selecteAllInBox.value) {
        selecteAllInBox.value(value);
    } else {
        console.error("The function is not yet available");
    }
}

// info 暴露给外部，用于检查文件名冲突与否
function checkFileNames(selectedFiles: fileItemType[]): [boolean, null | string[] | string] {
    let existedFileNames: string[] = [];

    const targetList = tabListContent.value.find(item => item.value == tab.value);
    if(targetList) {
        targetList.list.forEach((fileDetailItem: FileDetail) => {
            if(selectedFiles.some(selectedFile => selectedFile.name === fileDetailItem.name)) {
                existedFileNames.push(fileDetailItem.name);
            }
        })
        if(existedFileNames.length == 0)
            return [true, null];
        else
            return [false, existedFileNames];

    } else {
        console.error("No fileList found with the value !?")
        return [false, "fail"];
    }
}

function addItem(tabValue: string, fileName: string, size: number) {
    const targetList = tabListContent.value.find(item => item.value == tabValue);
    if(targetList) {
        targetList.list.push(
            new FileDetail(fileName, size)
        )
    } else {
        console.error("No fileList found with the value !?")
        // return [false, "fail"];
    }
} 

defineExpose({checkFileNames, addItem});
</script>

<style scoped>

.tab-container {
    margin: 0 10px;
}

.data-container {
    background-color: transparent;
    color: #ccc !important;

    td {
        cursor: pointer;
        user-select: none;
    }
    td:hover {
        color: #fff;
    }
}

.recommand-container {
    position: absolute;
    top: 80px;
}

.tab-close {
    position: absolute;
    top: 0;
    right:0;
}
</style>        