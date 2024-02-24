<script setup lang="ts">
import { useWebsocketStore } from "@/stores/WebsocketStore";
import { ref } from "vue";

const wsStore = useWebsocketStore();
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();
const drawColor = ref("#ffffff")
const drawWidth = ref(5)
let isMouseDown = false;

type DrawData = {
    x: number;
    y: number;
    color?: string;
    width?: number;
}

function mouseMove(event: MouseEvent) {
    if(!isMouseDown) { return; }

    let rect = drawingCanvas.value?.getBoundingClientRect();
    let drawData = { 
        x: event.clientX - (rect?.left ?? 0),
        y: event.clientY - (rect?.top ?? 0),
        color: drawColor.value,
        width: drawWidth.value
    }

    redrawCanvas(drawData, true)
    wsStore.sendDraw(drawData)
}

function mousedown(event: any) {
    isMouseDown = true;
}

function mouseup(event: any) {
    isMouseDown = false;
    let ctx = drawingCanvas.value?.getContext("2d");
    ctx?.beginPath();
    wsStore.sendDrawStop();
}

function redrawCanvas(data: DrawData, isLocal = false) {
    let ctx = drawingCanvas.value?.getContext("2d");
    ctx!.lineWidth = data.width!;
    ctx!.lineCap = "round";
    ctx!.strokeStyle = data.color!;
    ctx?.lineTo(data.x, data.y);
    ctx?.stroke();
}

function drawStop() {
    let ctx = drawingCanvas.value?.getContext("2d");
    ctx?.beginPath();
    console.log("draw stop")
}

wsStore.setCanvasFunctions(redrawCanvas, drawStop)
</script>


<template>
    <div>
        <div>
            <input v-model="drawColor" value="#ffffff" type="color">
            <input v-model="drawWidth" type="range" min="1" max="50" value="5">
        </div>
        <canvas ref="drawingCanvas" height="400" width="600" @mousemove="mouseMove" @mousedown="mousedown" @mouseup="mouseup"></canvas>
    </div>
</template>
