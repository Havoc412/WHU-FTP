<template>
    <div v-show="show" 
        class="snackbar-container flex-horizontal" 
        :style="{ backgroundColor: setBackgroundColor,
                color: color ? '#333' : '#fff' }"
        :class="{'snackbar-container-leave': leave_flag}">
        <p>{{ message }}</p>
        <div class="spacer"/>
        <v-btn class="close" :style="{
                color: color ? '#333' : '#fff'
            }"
            size="x-small" icon="mdi-close" variant="text" 
            @click="close"
        />
    </div>
</template>

<script>
export default {
    name: 'm-snackbar',
    emits: ['update:modelValue'],
    props: {
        modelValue: {
            type: Boolean,
            required: true
        },
        message: {
            type: String,
            required: true
        },
        color: {
            type: [Number, Boolean], // 向上兼容
            // required: true
        },
        timeout: {
            type: Number,
        },
    },
    data() {
        return {
            show: false,
            leave_flag: false,
        }
    },
    computed: {
        setBackgroundColor() {
            let color;
            if(typeof this.color == 'boolean')
                color = this.color ? 1: 0;
            else
                color = this.color;
            switch(color) {
                case 0: return '#FB4F95';
                case 1: return '#8BE8A1';
                case 2: return '#F9F871';
            }
        }
    },
    methods: {
        closeSelf(TIMEOUT = 1000) {
            setTimeout(() => {
                this.leave_flag = true;
                setTimeout(() => {
                    this.show = false;
                    this.leave_flag = false;
                    this.$emit('update:modelValue', false);
                }, 500);
            }, TIMEOUT);
        },
        close() {
            this.show = false;
        }
    },
    watch: {
        modelValue(newValue, oldValue) {
            if(newValue == true) {
                this.show = true;
                this.closeSelf(this.timeout);
            }
        }
    },
    // mounted() {
    //     this.show = true;
    //     this.closeSelf(this.timeout);
    // },
}

</script>

<style>

.snackbar-container {
    position: fixed;
    top: 35px;
    left: 50%;
    /* bottom: 20px; */
    transform: translateX(-50%);

    width: 80%;
    max-width: 400px;
    border-radius: 5px;
    padding: 7px;

    /* 动画部分 */
    animation: enter 0.5s;
    transition: opacity 0.5s;

    z-index: 9999;

    box-shadow: 0px 0px 7px 1px #ffffff80;  
}

.snackbar-container p {
    font-size: 14px;
    font-weight: 600;
}

/* 定义渐入渐出动画效果 */
@keyframes enter {
    0% {
        opacity: 0;
    }
    20% {
        opacity: 0.5;   
    }
    100% {
        opacity: 1;
    }
}

.snackbar-container-leave {
    opacity: 0;
}

</style>