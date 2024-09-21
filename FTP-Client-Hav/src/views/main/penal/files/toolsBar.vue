<template>
    <div class="flex-horizontal gap-5 root-container">
        <v-btn variant="text" size="25">
            <find theme="outline" size="25" :fill="selected(0)" @mouseover="setHover(0)" @mouseleave="clearHover"></find>
            <v-menu v-model="menuOpen" activator="parent" location="center" :close-on-content-click="false"> <!--info open-on-hover 考虑不用-->
                <div class="search-container flex-horizontal top gap-5">
                    <search class="btn" theme="outline" size="25" fill="#ddd" @click="emits('searchUpdate', searchInput)"/>
                    <v-text-field v-model="searchInput" label="Label" variant="underlined" 
                        clearable @click:clear="emits('searchUpdate', '')"
                        @blur="searchUpdate" @keyup.enter="searchUpdate"
                    />
                </div>
            </v-menu>
        </v-btn>
       
        <component :is="$props.floderState ? HamburgerButton : WaterfallsV" 
            theme="outline" size="25" :fill="selected(1)"
            @mouseover="setHover(1)" @mouseleave="clearHover"
        />
        <component :is="props.device ? DownloadTwo : UploadTwo"
            theme="outline" size="25" :fill="selected(2)"
            @mouseover="setHover(2)" @mouseleave="clearHover"
            @click="emits('getSelectedFiles')"
        />

    </div>
</template>

<script setup>
import { ref, computed } from "vue";

import { Find, HamburgerButton, WaterfallsV, UploadTwo, DownloadTwo, Search } from "@icon-park/vue-next";
import { iconHoverColor } from "@/const";

    // store
// DATA
const props = defineProps({
    searchOutside: String, // 好像还是没用到
    device: Boolean,
    floderState: Boolean,
});
const emits = defineEmits(['searchUpdate', 'getSelectedFiles']);
const hoverIndex = ref(-1);
const searchInput = ref("");
const menuOpen = ref(false);

// FUNC
const selected = computed(() => {
    return (index) => {
        return index == hoverIndex.value ? iconHoverColor.selected : iconHoverColor.unselected;
    }
});

function setHover(index) {
    hoverIndex.value = index;
}

function clearHover() {
    hoverIndex.value = -1;
}

function searchUpdate() {
    emits("searchUpdate", searchInput.value);
    menuOpen.value = false;
    // searchInput.value = "";
}

</script>

<style scoped>

.root-container > * {
    cursor: pointer;
}

.search-container {
    background-color: #222;
    border-radius: 10px;
    box-shadow: 0px 0px 7px 1px #ffffff30;
    transform: translateY(5px);

    padding: 0 10px;
    height: 50px;
    min-width: 200px;

    span {
        color: #ccc;
        font-weight: bold;
    }
}


</style>        