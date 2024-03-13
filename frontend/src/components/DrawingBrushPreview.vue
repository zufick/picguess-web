<script setup lang="ts">
import type { 
    VirtualCanvasBrush
} from '@/types/VirtualCanvas'

const props = defineProps<{
    brush: VirtualCanvasBrush,
}>()
</script>


<template>
    <div class="brush-preview" :style="{ 
        'border-color': props.brush.color, 
        width: props.brush.width + 'px',
        height: props.brush.width + 'px',
        left: (props.brush.point.x - props.brush.width / 2) + 'px',
        top: (props.brush.point.y - props.brush.width / 2) + 'px',
        visibility: props.brush.visible ? 'visible' : 'hidden'}">
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
        border: 2px solid black;
        border-radius: 100%;
    }
</style>