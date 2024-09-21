<template>
    <div class="container flex-center-vertical mt-2 gap-15">
        <template v-for="(item, index) in components" :key="index">
            <component :is="item" @mouseover="setHover(index)" @mouseleave="clearHover"
                theme="outline" size="25" :fill="selected(index)"/>
        </template>
    </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { AllApplication, Refresh, SettingTwo, Bookmark, Tips, Tag } from "@icon-park/vue-next";

import { iconHoverColor } from "@/const";
// DATA
const props = defineProps({

    });
const emits = defineEmits([]);

const selectedColor = iconHoverColor.selected;
const unselectedColor = iconHoverColor.unselected;

const selectedIndex = ref(0);
const hoverIndex = ref(-1);

const components = [
    AllApplication, Refresh, Tag, Tips, SettingTwo
];

// FUNC
const selected = computed(() => {
    return (index) => {
        return index == selectedIndex.value || index == hoverIndex.value ? selectedColor : unselectedColor;
    }
});

function selectIcon(index) {
    selectedIndex.value = index;
}

function setHover(index) {
    hoverIndex.value = index;
}

function clearHover() {
    hoverIndex.value = -1;
}

</script>

<style scoped>

.container {
    width: 60px; /* 好像没什么卵用 */
}

.container > * {
    height: 30px;
    width: 40px;
}

.container > * :hover {
    cursor: pointer;
}

</style>        