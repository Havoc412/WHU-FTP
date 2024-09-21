<template>
  <div class="animated-line">
    <svg width="200" height="200">
      <path :d="pathD" fill="none" stroke="black" stroke-width="2"/>
    </svg>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const pathD = ref('M10 10 C 20 20, 40 20, 50 10');

// 模拟动态路径变化
onMounted(() => {
  setInterval(() => {
    const start = Math.random() * 100;
    const c1 = Math.random() * 100;
    const c2 = Math.random() * 100;
    const end = Math.random() * 100;
    pathD.value = `M10 ${start} C 20 ${c1}, 40 ${c2}, 50 ${end}`;
  }, 1000);
});
</script>

<style>
@keyframes dance {
  0% { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: 1000; }
}

svg path {
  stroke-dasharray: 1000;
  animation: dance 10s linear infinite;
}
</style>