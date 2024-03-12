<script setup lang="ts">
import { defineProps, ref, watch, type Ref } from "vue";
import type { 
    VirtualCanvas, 
    VirtualCanvasLine, 
    VirtualCanvasDrawPoint 
} from '@/types/VirtualCanvas'

const props = defineProps<{
    line: VirtualCanvasLine
}>()
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();

// watch works directly on a ref

const canvasLayer = {
    line: {} as VirtualCanvasLine,
    pendingDrawPoints: [] as VirtualCanvasDrawPoint[],
    lastDrawPoint: {} as VirtualCanvasDrawPoint,
    drawPendingPoints() {
        let ctx = drawingCanvas.value?.getContext("2d");

        ctx!.lineWidth = this.line.width;
        ctx!.lineCap = "round";
        ctx!.strokeStyle = this.line.color;

        ctx?.beginPath();


        for (let i = 0; i < this.pendingDrawPoints.length; i++) {
            let point = this.pendingDrawPoints[i];
            if (this.lastDrawPoint.x && this.lastDrawPoint.y) {
                ctx?.moveTo(this.lastDrawPoint.x, this.lastDrawPoint.y);    
            }
            ctx?.lineTo(point.x, point.y);
            this.lastDrawPoint = point;

        }

        this.pendingDrawPoints = [];

        ctx?.stroke();
    },
    clearCanvas() {
        let ctx = drawingCanvas.value?.getContext("2d");
        ctx?.clearRect(0,0, drawingCanvas.value!.width, drawingCanvas.value!.height);
        this.lastDrawPoint = {} as VirtualCanvasDrawPoint;
    }
}


watch(() => props.line, async (newLine: VirtualCanvasLine, oldLine: VirtualCanvasLine) => {
    if(newLine != oldLine)
        canvasLayer.clearCanvas();

    canvasLayer.line = newLine;
    canvasLayer.pendingDrawPoints.push(...newLine.newPoints!);
    canvasLayer.drawPendingPoints();
},   {deep: true})
</script>


<template>
    <canvas class="layer" ref="drawingCanvas" height="768" width="1024" v-show="!props.line.undo">
    </canvas>
</template>

<style scoped>
    .layer {
        position: absolute;
        top: 0;
        left: 0;
        pointer-events: none;
        z-index: 100;
        color: black;
    }
</style>