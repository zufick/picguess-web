<script setup lang="ts">
import type { 
    VirtualCanvasBrush
} from '@/types/VirtualCanvas'
import { ref, watch } from 'vue';

const props = defineProps<{
    brush: VirtualCanvasBrush,
}>()
const invertedColor = ref("")

watch(props.brush, (newValue, oldValue) => {
    invertedColor.value = invertHex(newValue.color)
})

function invertHex(hex: string) {
    hex = hex.substring(1)
    let inverted = (Number(`0x1${hex}`) ^ 0xFFFFFF).toString(16).substr(1).toUpperCase();
      return  "#" + inverted
}



</script>


<template>
    <div class="brush-preview" :style="{ 
        'border-color': props.brush.color, 
        width: props.brush.width + 'px',
        height: props.brush.width + 'px',
        left: (props.brush.point.x - props.brush.width / 2) + 'px',
        top: (props.brush.point.y - props.brush.width / 2) + 'px',
        visibility: props.brush.visible ? 'visible' : 'hidden'}">
            <div class="brush-diff" :style="{
            'border-color': invertedColor, 
            }"></div>

    </div>
</template>

<style scoped>
    .brush-preview {
        width: 25px;
        height: 25px;
        position: absolute;
        top: 0;
        left: 0;
        pointer-events: none;
        z-index: 10;
        background: transparent;
        border: 1px solid black;
        border-radius: 100%;
    }
    
    .brush-diff {
        width: 100%;
        height: 100%;
        border-radius: 100%;
        border: 1px solid black;
        border-color: rgb(255, 255, 255);
        pointer-events: none;

    }
</style>