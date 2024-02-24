<script setup lang="ts">
import { useWebsocketStore } from "@/stores/WebsocketStore";
import { ref } from "vue";

const wsStore = useWebsocketStore();
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();
const drawColor = ref("#ffffff")
const drawWidth = ref(5)

const virtualCanvas: VirtualCanvas = {
    lines: {
        "aaa": [
            {
                color: "red",
                width: 1,
                points: [
                    { x: 1, y: 1}
                ]
            }
        ],
    },
    startNewLine(id: string, data: { color: string, width: string, lines: [] }) {
        this.lines[id].push(data);
    },
}

type VirtualCanvas = {
    lines: {
        [key: string] : {
            color: string,
            width: number,
            points: { 
                x: number, 
                y: number 
            }[]
        }[]
    },
    startNewLine: Function,
};

type DrawData = {
    x: number;
    y: number;
    color?: string;
    width?: number;
}

let isMouseDown = false;

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

function mouseDown(event: any) {
    isMouseDown = true;

    virtualCanvas["local"].push({
        color: drawColor.value,
        width: drawWidth.value,
        points: []
    })

    wsStore.sendDrawNewLine();
}

function mouseUp(event: any) {
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
        <canvas ref="drawingCanvas" height="400" width="600" @mousemove="mouseMove" @mousedown="mouseDown" @mouseup="mouseUp"></canvas>
    </div>
</template>
