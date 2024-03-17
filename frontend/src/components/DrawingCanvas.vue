<script setup lang="ts">
import { nextTick, onMounted, provide, ref } from "vue";
import { useWebsocketStore } from "@/stores/WebsocketStore";
import type { 
    VirtualCanvas, 
    VirtualCanvasLine, 
    VirtualCanvasDrawPoint, 
    VirtualCanvasBrush
} from '@/types/VirtualCanvas'
import DrawingCanvasLayer from "./DrawingCanvasLayer.vue";
import DrawingBrushPreview from "./DrawingBrushPreview.vue";
import DrawingCanvasClear from './DrawingCanvasClear.vue'

const wsStore = useWebsocketStore();
const drawingCanvas = ref<InstanceType<typeof HTMLCanvasElement>>();
const drawClearDialog = ref(false)

const virtualCanvas: VirtualCanvas = {
    canvas: undefined,
    brush: ref(<VirtualCanvasBrush>{
        color: "#000000",
        width: 5,
        point: {
            x: 0,
            y: 0,
        },
        visible: false,
    }),
    lines: ref([]),
    linesToDelete: [],
    lastUserLines: {},
    sendPointsBuffer: [],
    sendPointsBufferLimit: 5,
    setCanvas(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
    },
    startNewLine(id: string, data: VirtualCanvasLine) {
        this.sendPointsBuffer = []

        data.newPoints = []
        data.userId = id;
        data.id = Date.now() + "";
        this.lines.value.push(data);

        if (this.lastUserLines[id] != undefined && this.lastUserLines[id].undo) {
            this.lines.value.splice(this.lines.value.indexOf(this.lastUserLines[id]), 1);
        }

        this.lastUserLines[id] = data;
        this.mergeReadyCanvases();
    },
    startNewLocalLine(data: VirtualCanvasLine) {
        this.startNewLine("local", data)
    },
    drawLocalPoint(point: VirtualCanvasDrawPoint) {
        this.sendPointsBuffer.push(point.x, point.y)

        try {
            this.drawPoint("local", point)


        } catch (e) {
            console.error(e)
        }

        if (this.sendPointsBuffer.length >= 2 * this.sendPointsBufferLimit) {
                wsStore.sendDrawPoints(this.sendPointsBuffer)
                this.sendPointsBuffer = []
            }
    },
    drawPoint(id: string, point: VirtualCanvasDrawPoint) {
        this.drawPoints(id, [point])
    },
    drawPointArray(id: string, points: number[]) {
        let pointObjects = [];

        for (let i = 0; i < points.length; i += 2) {
            pointObjects.push({
                x: points[i],
                y: points[i+1],
            })
        }
        this.drawPoints(id, pointObjects);
    },
    drawPoints(id: string, points: VirtualCanvasDrawPoint[]) {
        let line = this.lastUserLines[id];
        if(this.lines.value.indexOf(line) == -1) {
            console.error("line -1", this.lines.value, line)
            return;
        }
        this.lines.value[this.lines.value.indexOf(line)].newPoints = points;
    },
    undoLine(id: string) {
        if(this.lastUserLines[id] == undefined) {
            return;
        }

        this.lines.value[this.lines.value.indexOf(this.lastUserLines[id])].undo = true;
    },
    redoLine(id: string) {
        if(this.lastUserLines[id] == undefined) {
            return;
        }

        this.lines.value[this.lines.value.indexOf(this.lastUserLines[id])].undo = false;
    },
    drawClear(id: string) {
        if (id != "") {
            wsStore.roomState.clients.forEach((u) => {
                if (u.id == id) {
                    alert(`User ${u.username} cleared canvas`)
                }
            });
        }

        this.lines.value = []
    },
    addReadyToMergeCanvas(line: VirtualCanvasLine, canvas: HTMLCanvasElement) {
        if(line.readyToMergeCanvas) { return; }
        this.lines.value[this.lines.value.indexOf(line)].readyToMergeCanvas = canvas;
        this.mergeReadyCanvases();
    },
    mergeReadyCanvases() {
        for(let i = this.lines.value.length - 1; i >= 0; i--) {
            let line = this.lines.value[i];
            if(!line || !line.readyToMergeCanvas)
                continue;
            
            let mergeBottomLayers = [] as HTMLCanvasElement[];
            let j = 0;
            while (this.lines.value[i-j] && this.lines.value[i-j].readyToMergeCanvas) {
                mergeBottomLayers.unshift(this.lines.value[i-j].readyToMergeCanvas!)
                this.linesToDelete.push(this.lines.value[i-j])
                j++;
            }
            mergeBottomLayers.shift();
            this.lines.value[i-j+1].mergeCanvases = mergeBottomLayers;
            this.lines.value[i-j+1].newPoints = [];
            this.linesToDelete.splice(this.linesToDelete.indexOf(this.lines.value[i-j+1]), 1)
            i = i-j;
        }
    },
    deleteReadyToDeleteLines(mergedOnto: VirtualCanvasLine) {
        let deletedLines = []
        for (let i = 0; i < this.linesToDelete.length; i++) {
            let lineToDelete = this.linesToDelete[i];
            if (mergedOnto.mergeCanvases && lineToDelete.readyToMergeCanvas && mergedOnto.mergeCanvases?.includes(lineToDelete.readyToMergeCanvas)) {
                this.lines.value[this.lines.value.indexOf(mergedOnto)].mergeCanvases?.splice(mergedOnto.mergeCanvases?.indexOf(lineToDelete.readyToMergeCanvas), 1)
                this.lines.value.splice(this.lines.value.indexOf(lineToDelete), 1);
                deletedLines.push(i);
            }
        }

        deletedLines.forEach((l) => {
            this.linesToDelete.splice(l, 1);
        })
    }
}

type DrawData = {
    x: number;
    y: number;
    color?: string;
    width?: number;
}

let isMouseDown = false;

function mouseMove(event: MouseEvent) {
    let rect = drawingCanvas.value?.getBoundingClientRect();
    virtualCanvas.brush.value.point = {
        x: Math.floor(event.clientX - rect!.left ?? 0),
        y: Math.floor(event.clientY - rect!.top ?? 0),
    };

    if(!isMouseDown) { return; }

    drawEventPoint(event);
}

function mouseDown(event: MouseEvent) {
    isMouseDown = true;
    
    virtualCanvas.startNewLocalLine({
        color: virtualCanvas.brush.value.color,
        width: virtualCanvas.brush.value.width,
    })

    wsStore.sendDrawNewLine({ 
        color: virtualCanvas.brush.value.color, 
        width: virtualCanvas.brush.value.width 
    });

    drawEventPoint(event);
}

function drawEventPoint(event: MouseEvent) {
    let rect = drawingCanvas.value?.getBoundingClientRect();
    let point = virtualCanvas.brush.value.point;

    virtualCanvas.drawLocalPoint(point)

}

function mouseUp(event: MouseEvent) {
    isMouseDown = false;

    if(virtualCanvas.sendPointsBuffer.length > 0) {
        wsStore.sendDrawPoints(virtualCanvas.sendPointsBuffer)
        virtualCanvas.sendPointsBuffer = []
    }
    
}

function mouseEnter(event: MouseEvent) {
    virtualCanvas.brush.value.visible = true;
}

function mouseLeave(event: MouseEvent) {
    virtualCanvas.brush.value.visible = false;
}

function undoClicked() {
    virtualCanvas.undoLine("local");
    wsStore.sendDrawUndo();
}

function redoClicked() {
    virtualCanvas.redoLine("local");
    wsStore.sendDrawRedo();
}

function eraserClicked() {
    virtualCanvas.brush.value.color = "#ffffff"
}

wsStore.setVirtualCanvas(virtualCanvas)
onMounted(() => {
    virtualCanvas.setCanvas(drawingCanvas.value!)
})


provide('virtualCanvas', virtualCanvas)
</script>


<template>
    <div class="flex gap-4">
        <div class="canvas-container border-2 rounded-lg border-white/20 overflow-hidden">
            <div class="canvas-bg"></div>
            <canvas ref="drawingCanvas" height="768" width="1024" 
                @mousemove="mouseMove" 
                @mousedown="mouseDown" 
                @mouseup="mouseUp" 
                @mouseout="mouseUp"
                @mouseenter="mouseEnter"
                @mouseleave="mouseLeave">
            </canvas>
            <DrawingCanvasLayer v-for="line in virtualCanvas.lines.value" :line="line" :key="line.id + '_' + line.color" :lastUserLine="virtualCanvas.lastUserLines[line.userId!]"/>

            <DrawingBrushPreview :brush="virtualCanvas.brush.value"/>
        </div>
        <div class="flex flex-col items-center gap-4">
            <input v-model="virtualCanvas.brush.value.color" value="#ffffff" type="color">
            <input v-model="virtualCanvas.brush.value.width" type="range" min="1" max="100" value="5">
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
            <DrawingCanvasClear/>
            <button @click="eraserClicked()">E</button>
        </div>
    </div>
</template>

<style scoped>
    input[type="range"] {
        writing-mode: vertical-lr;
        width: 8px;
    }

    .canvas-container {
        width: 1024px;
        height: 768px;
        position: relative;
        user-select: none; /* Standard syntax */
    }

    .canvas-bg {
        position: absolute;
        z-index: -1;
        background-color: white;
        width: 100%;
        height: 100%;
    }

    canvas {
        cursor: none;
    }
</style>