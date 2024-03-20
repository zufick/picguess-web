<script setup lang="ts">
import { defineProps, ref, watch, type Ref, onMounted, inject, Transition } from "vue";
import type { 
    VirtualCanvas, 
    VirtualCanvasLine, 
    VirtualCanvasDrawPoint 
} from '@/types/VirtualCanvas'

const props = defineProps<{
    line: VirtualCanvasLine,
    lastUserLine: VirtualCanvasLine,
}>()
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();

// watch works directly on a ref

const canvasLayer = {
    line: {} as VirtualCanvasLine,
    pendingDrawPoints: [] as VirtualCanvasDrawPoint[],
    lastDrawPoint: {} as VirtualCanvasDrawPoint,
    drawingInterval: 0,
    isOpponentDrawing: ref(false),
    opponentDrawingTimeout: 0,
    isBeingDrawn: true,
    drawPendingPoints() {
        if(this.line.userId != "local") {
            this.drawWithInterval();
            return;
        }

        let ctx = drawingCanvas.value?.getContext("2d");

        if (ctx) {
            ctx.lineWidth = this.line.width;
            ctx.lineCap = "round";
            ctx.strokeStyle = this.line.color;
        }

        ctx?.beginPath();

        for (let i = 0; i < this.pendingDrawPoints.length; i++) {
            let point = this.pendingDrawPoints[i];
            if (!this.lastDrawPoint.x && !this.lastDrawPoint.y) {
                this.lastDrawPoint = point;
            }
            ctx?.moveTo(this.lastDrawPoint.x, this.lastDrawPoint.y);    
            ctx?.lineTo(point.x, point.y);
            this.lastDrawPoint = point;
        }
        this.pendingDrawPoints = [];
        ctx?.stroke();
    },
    drawWithInterval() {
        let ctx = drawingCanvas.value?.getContext("2d");

        if (ctx) {
            ctx.lineWidth = this.line.width;
            ctx.lineCap = "round";
            ctx.strokeStyle = this.line.color;
        }

        ctx?.beginPath();

        if (this.drawingInterval) {
            return;
        }
        if(!this.drawingInterval && this.pendingDrawPoints.length > 0) {
            this.drawingInterval = setInterval(() => {

                if (this.pendingDrawPoints.length <= 0) {
                    clearInterval(this.drawingInterval);
                    this.drawingInterval = 0;
                    this.opponentDrawingTimeout = setTimeout(() => {
                        this.isOpponentDrawing.value = false;
                    }, 250)
                    return;
                }

                clearTimeout(this.opponentDrawingTimeout);
                this.isOpponentDrawing.value = true;
                let point = this.pendingDrawPoints.shift();

                if (this.lastDrawPoint.x && this.lastDrawPoint.y) {
                    ctx?.moveTo(this.lastDrawPoint.x, this.lastDrawPoint.y);    
                }
                if(point) {
                    ctx?.lineTo(point.x, point.y);
                    ctx?.stroke();
                    this.lastDrawPoint = point;

                    // Draw also closest points without a delay
                    if (this.pendingDrawPoints.length <= 0) return;

                    while(this.pendingDrawPoints.length > 0 && ((Math.hypot(this.pendingDrawPoints[0].x-point.x, this.pendingDrawPoints[0].y-point.y) <= 10) || !this.isBeingDrawn)) {
                        let anotherPoint = this.pendingDrawPoints.shift()!;

                        ctx?.moveTo(this.lastDrawPoint.x, this.lastDrawPoint.y);
                        ctx?.lineTo(anotherPoint.x, anotherPoint.y);
                        ctx?.stroke();
                        this.lastDrawPoint = anotherPoint;
                    }
                }

            }, 0.1)
        }
    },
    clearCanvas() {
        clearInterval(this.drawingInterval);
        this.drawingInterval = 0;
        this.pendingDrawPoints = [];

        let ctx = drawingCanvas.value?.getContext("2d");
        ctx?.clearRect(0,0, drawingCanvas.value!.width, drawingCanvas.value!.height);
        this.lastDrawPoint = {} as VirtualCanvasDrawPoint;
    },
    mergeCanvases() {
        if (!this.line.mergeCanvases) return;

        for (let i = 0; i < this.line.mergeCanvases.length; i++) {
            let ctx = drawingCanvas.value?.getContext("2d");
            ctx?.drawImage(this.line.mergeCanvases[i], 0, 0)                
        }
        
        virtualCanvas?.deleteReadyToDeleteLines(this.line)
    },
    setNewLine(newLine: VirtualCanvasLine) {
        if(newLine != this.line) {
            canvasLayer.clearCanvas();
        }

        canvasLayer.line = newLine;

        if(newLine.newPoints) {
            canvasLayer.pendingDrawPoints.push(...newLine.newPoints);
            canvasLayer.drawPendingPoints();
        }

        if (newLine.mergeCanvases) {
            this.mergeCanvases();
        }
    }
}

watch(() => props.line, async (newLine: VirtualCanvasLine, oldLine: VirtualCanvasLine) => {
    canvasLayer.setNewLine(newLine);
},   {deep: true})

watch(() => props.lastUserLine, async (newLastUserLine: VirtualCanvasLine, oldLastUserLine: VirtualCanvasLine) => {
    if (canvasLayer.line != newLastUserLine) {
        if(drawingCanvas.value && canvasLayer.pendingDrawPoints.length <= 0) {
            virtualCanvas?.addReadyToMergeCanvas(canvasLayer.line, drawingCanvas.value);
        }
        canvasLayer.isBeingDrawn = false;
    }
})

onMounted(() => {
    if (props.line.userId == "local" && Object.keys(canvasLayer.line).length == 0) {
        canvasLayer.setNewLine(props.line)
    }
})

const virtualCanvas = inject<VirtualCanvas>("virtualCanvas")
</script>


<template>
    <div>
        <canvas class="layer" ref="drawingCanvas" height="768" width="1024" v-show="!props.line.undo">
        </canvas>
        <!--
        <Transition name="brush-icon">
        <div v-if="canvasLayer.isOpponentDrawing.value && props.line.newPoints && canvasLayer.pendingDrawPoints.length > 0 && props.line.newPoints[0]" class="brush-icon" :style="{
            left: props.line.newPoints[0].x - 10 + 'px',
            top: props.line.newPoints[0].y - 40 + 'px',
        }">
            <svg v-if="canvasLayer.pendingDrawPoints.length" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="black" class="w-6 h-6" style="width: 100%; height: 100%;">
                <path fill="white" stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 0 0-5.78 1.128 2.25 2.25 0 0 1-2.4 2.245 4.5 4.5 0 0 0 8.4-2.245c0-.399-.078-.78-.22-1.128Zm0 0a15.998 15.998 0 0 0 3.388-1.62m-5.043-.025a15.994 15.994 0 0 1 1.622-3.395m3.42 3.42a15.995 15.995 0 0 0 4.764-4.648l3.876-5.814a1.151 1.151 0 0 0-1.597-1.597L14.146 6.32a15.996 15.996 0 0 0-4.649 4.763m3.42 3.42a6.776 6.776 0 0 0-3.42-3.42" />
            </svg>
        </div>
        </Transition>
        -->
    </div>
</template>

<style scoped>
    .layer {
        position: absolute;
        top: 0;
        left: 0;
        pointer-events: none;
        z-index: 1;
    }
    .brush-icon {
        position: absolute;
        top: 0;
        left: 0;
        pointer-events: none;
        width: 48px;
        height: 48px;
        z-index: 10;
    }

    .brush-icon-enter-active,
    .brush-icon-leave-active {
    transition: opacity 0.25s ease;
    }

    .brush-icon-enter-from,
    .brush-icon-leave-to {
    opacity: 0;
    }

</style>