<!--todo 无效状态-->
<template>
    <div>
        <v-hover>
            <template v-slot:default="{ isHovering, props }">
                <component v-bind="props"
                    class="icon"
                    :is="Increase" 
                    :theme="props.theme" :size="props.size" 
                    :fill="isHovering ? props.selectedColoe : props.unselectedColor"
                    >
                </component>
            </template>
        </v-hover>
    </div>
</template>

<script setup>
import { ref } from "vue";
import { iconHoverColor } from "@/const";
import { NewComputer, DataServer, CloseSmall, Increase } from "@icon-park/vue-next";

    // store
// DATA
const props = defineProps({
    cpt: {
        type: Object,
        require: true
    },
    theme: {
        type: String,
        default: "outline"
    },
    size: {
        type: String,
        default: "25"
    },
    selectedColor: {
        type: String,
        default: iconHoverColor.selected
    },
    unselectedColor: {
        type: String,
        default: iconHoverColor.unselected
    }
});
const emits = defineEmits([]);

// FUNC

</script>

<style scoped>
.icon {
    cursor: pointer;
}
</style>        