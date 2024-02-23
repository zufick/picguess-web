<script setup lang="ts">
import { useWebsocketStore } from "@/stores/WebsocketStore";
import { ref } from "vue";
import VueDrawingCanvas from "vue-drawing-canvas";

const wsStore = useWebsocketStore();
const vueCanvasDrawing = ref<InstanceType<typeof HTMLCanvasElement>>();
let isMouseDown = false;

function mouseMove(event: MouseEvent) {
    if(!isMouseDown) { return; }

    let rect = vueCanvasDrawing.value?.getBoundingClientRect();
    let coords = { 
        x: event.clientX - (rect?.left ?? 0),
        y: event.clientY - (rect?.top ?? 0) 
    }

    let ctx = vueCanvasDrawing.value?.getContext("2d");
    ctx!.lineWidth = 10;
    ctx!.lineCap = "round";
    ctx!.strokeStyle = "white";
    ctx!.fill
    ctx?.lineTo(coords.x, coords.y);
    ctx?.stroke();
    ctx?.closePath();
    wsStore.sendDraw(coords)
}

function mousedown(event: any) {
    isMouseDown = true;
}

function mouseup(event: any) {
    isMouseDown = false;
}

function redrawCanvas(data: any) {
    let ctx = vueCanvasDrawing.value?.getContext("2d");
    ctx!.lineWidth = 10;
    ctx!.lineCap = "round";
    ctx!.strokeStyle = "white";
    ctx?.lineTo(data.x, data.y);
    ctx?.stroke();
    ctx?.closePath();
}

wsStore.setCanvas(redrawCanvas)



</script>


<template>
    <div>
        <canvas ref="vueCanvasDrawing" height="400" width="600" @mousemove="mouseMove" @mousedown="mousedown" @mouseup="mouseup"></canvas>
    </div>
</template>
