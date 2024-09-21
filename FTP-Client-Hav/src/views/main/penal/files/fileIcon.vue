<template>
    <component 
        :is="fileIcon[props.type]" 
        :theme="props.theme" 
        :size="props.size"
        :fill="fileIconColor[props.type].selected"
    />
</template>

<script setup>
import { SeoFolder, FileTxt, FileQuestion, FileWord } from "@icon-park/vue-next";

import { fileIconColor } from "@/const";
// DATA
const props = defineProps({
    type: Number,
    theme: {
        type: String,
        default: "outline"
    },
    size: {
        type: String,
        default: "25"
    },
});
const emits = defineEmits([]);

const fileIcon = [SeoFolder, FileTxt, FileWord, FileQuestion];
// FUNC

</script>

<style scoped>

</style>        