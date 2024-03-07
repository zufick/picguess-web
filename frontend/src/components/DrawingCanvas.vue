<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useWebsocketStore } from "@/stores/WebsocketStore";
import type { 
    VirtualCanvas, 
    VirtualCanvasLineData, 
    VirtualCanvasDrawPoint 
} from '@/types/VirtualCanvas'

const wsStore = useWebsocketStore();
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();
const drawColor = ref("#ffffff")
const drawWidth = ref(5)

const virtualCanvas: VirtualCanvas = {
    canvas: undefined,
    lines: {},
    undoLines: {},
    setCanvas(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
    },
    startNewLine(id: string, data: VirtualCanvasLineData) {
        if (this.lines[id] == undefined) {
            this.lines[id] = []
        }

        let totalLayers = Object.values(this.lines).flat(1).length
        data.layerIndex = totalLayers
        this.lines[id].push(data);

        if (this.undoLines[id] != undefined) {
            this.undoLines[id] = []
            console.log(this.undoLines)
        }
    },
    startNewLocalLine(data: VirtualCanvasLineData) {
        this.startNewLine("local", data)
    },
    drawLocalPoint(point: VirtualCanvasDrawPoint) {
        this.drawPoint("local", point)
    },
    drawPoint(id: string, point: VirtualCanvasDrawPoint) {
        let line = this.lines[id];
        if(line[line.length - 1].points == undefined) {
            line[line.length - 1].points = [];
        }
        line[line.length - 1].points.push(point);

        this.redrawCanvas();
    },
    redrawCanvas() {
        let ctx = this.canvas?.getContext("2d");
        ctx!.clearRect(0, 0, this.canvas!.width, this.canvas!.height);

        let lines = Object.values(this.lines).flat(1)
        lines.sort((a,b) => a.layerIndex! - b.layerIndex!)
        lines.forEach(line => {
            ctx!.lineWidth = line.width;
            ctx!.lineCap = "round";
            ctx!.strokeStyle = line.color;
            
            if (line.points != undefined) {
                line.points.forEach(point => {
                    ctx?.lineTo(point.x, point.y);
                });
            }
            ctx?.stroke();
            ctx?.beginPath();
        });
    },
    undoLine(id: string) {
        if(this.lines[id] == undefined) {
            return;
        }
        if(this.undoLines[id] == undefined) {
            this.undoLines[id] = [];
        }

        let line = this.lines[id].pop()!;
        if (line != undefined) {
            this.undoLines[id].push(line);
        }

        this.redrawCanvas();
    },
    redoLine(id: string) {
        if(this.undoLines[id] == undefined) {
            return;
        }
        if(this.lines[id] == undefined) {
            this.lines[id] = [];
        }

        let line = this.undoLines[id].pop()!;
        if (line != undefined) {
            this.lines[id].push(line);
        }

        this.redrawCanvas();
    },
}

type DrawData = {
    x: number;
    y: number;
    color?: string;
    width?: number;
}

let isMouseDown = false;

function mouseMove(event: MouseEvent) {
    if(!isMouseDown) { return; }

    drawEventPoint(event);
}

function mouseDown(event: any) {
    isMouseDown = true;
    
    virtualCanvas.startNewLocalLine({
        color: drawColor.value,
        width: drawWidth.value,
        points: []
    })

    wsStore.sendDrawNewLine({ 
        color: drawColor.value, 
        width: drawWidth.value 
    });

    drawEventPoint(event);
}

function drawEventPoint(event: any) {
    let rect = drawingCanvas.value?.getBoundingClientRect();
    let point = { 
        x: event.clientX - (rect?.left ?? 0),
        y: event.clientY - (rect?.top ?? 0),
    }

    virtualCanvas.drawLocalPoint(point)
    wsStore.sendDrawPoint(point)
}

function mouseUp(event: any) {
    isMouseDown = false;
    let ctx = drawingCanvas.value?.getContext("2d");
}

function undoClicked() {
    virtualCanvas.undoLine("local");
    wsStore.sendDrawUndo();
}

function redoClicked() {
    virtualCanvas.redoLine("local");
    wsStore.sendDrawRedo();
}


wsStore.setVirtualCanvas(virtualCanvas)

onMounted(() => {
    virtualCanvas.setCanvas(drawingCanvas.value!)
})
</script>


<template>
    <div class="flex gap-4">
        <canvas class="border-2 rounded-lg border-white/20" ref="drawingCanvas" height="768" width="1024" @mousemove="mouseMove" @mousedown="mouseDown" @mouseup="mouseUp"></canvas>
        <div class="flex flex-col items-center gap-4">
            <input v-model="drawColor" value="#ffffff" type="color">
            <input v-model="drawWidth" type="range" min="1" max="50" value="5">
            <button @click="undoClicked">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75 3 12m0 0 3.75-3.75M3 12h18" />
                </svg>
            </button>
            <button @click="redoClicked">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M17.25 8.25 21 12m0 0-3.75 3.75M21 12H3" />
                </svg>
            </button>
        </div>
    </div>
</template>

<style scoped>
    input[type="range"] {
        writing-mode: vertical-lr;
        width: 8px;
    }
</style>