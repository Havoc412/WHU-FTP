<template>
    <Icon/>
    <Title/>
    <Header v-if="Header"/>
</template>

<script setup>
import { ref, defineAsyncComponent } from "vue";

import Icon from "./header/icon.vue";
import Title from "./header/title.vue";
// tip 借助动态导入的方式，同时支持窗口 && WEB；
const isElectron = navigator.userAgent.includes('Electron');
const Header = isElectron ? defineAsyncComponent(() => import('@/components/overall/header.vue')) : null;

// DATA
const props = defineProps({

    });
const emits = defineEmits([]);

// FUNC

</script>

<style scoped>

</style>        