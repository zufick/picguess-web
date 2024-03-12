<script setup lang="ts">
import { defineProps, ref, watch, type Ref, onMounted } from "vue";
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
    drawingInterval: 0,
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
                    return;
                }

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

                    console.log(Math.hypot(this.pendingDrawPoints[0].x-point.x, this.pendingDrawPoints[0].y-point.y))

                    while(this.pendingDrawPoints.length > 0 && Math.hypot(this.pendingDrawPoints[0].x-point.x, this.pendingDrawPoints[0].y-point.y) <= 20) {
                        let anotherPoint = this.pendingDrawPoints.shift()!;

                        ctx?.lineTo(anotherPoint.x, anotherPoint.y);
                        ctx?.stroke();
                        this.lastDrawPoint = anotherPoint;

                        console.log("too close")
                    }
                }

            }, 0.1)
        }
    },
    clearCanvas() {
        let ctx = drawingCanvas.value?.getContext("2d");
        ctx?.clearRect(0,0, drawingCanvas.value!.width, drawingCanvas.value!.height);
        this.lastDrawPoint = {} as VirtualCanvasDrawPoint;
    },
    setNewLine(newLine: VirtualCanvasLine) {
        if(newLine != this.line) {
            canvasLayer.clearCanvas();
        }

        canvasLayer.line = newLine;
        console.log("new line")

        if(newLine.newPoints) {
            canvasLayer.pendingDrawPoints.push(...newLine.newPoints);
            canvasLayer.drawPendingPoints();
        }
    }
}

watch(() => props.line, async (newLine: VirtualCanvasLine, oldLine: VirtualCanvasLine) => {
    canvasLayer.setNewLine(newLine);
},   {deep: true})

onMounted(() => {
    if (props.line.userId == "local" && Object.keys(canvasLayer.line).length == 0) {
        canvasLayer.setNewLine(props.line)
    }
})
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