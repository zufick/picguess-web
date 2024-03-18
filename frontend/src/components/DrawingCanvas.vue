<script setup lang="ts">
import { nextTick, onMounted, provide, ref, watch } from "vue";
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
        width: 25,
        point: {
            x: 0,
            y: 0,
        },
        visible: false,
    }),
    paintBrush: ref("#000000"),
    eraserBrush: "#ffffff",
    lines: ref([]),
    linesToDelete: [],
    lastUserLines: {},
    sendPointsBuffer: [],
    sendPointsBufferLimit: 5,
    setCanvas(canvas: HTMLCanvasElement) {
        this.canvas = canvas;

        watch(this.paintBrush, (newVal, oldVal) => {
            if (this.brush.value.color != this.eraserBrush) {
                this.brush.value.color = newVal;
            }
        })
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
    },
    selectPaintMode() {
        this.brush.value.color = this.paintBrush.value;
    },
    selectEraseMode() {
        this.brush.value.color = this.eraserBrush;
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
            <input v-model="virtualCanvas.paintBrush.value" value="#ffffff" type="color">
            <button @click="virtualCanvas.selectPaintMode()" :class="{'selected-btn': virtualCanvas.brush.value.color == virtualCanvas.paintBrush.value}">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 0 0-5.78 1.128 2.25 2.25 0 0 1-2.4 2.245 4.5 4.5 0 0 0 8.4-2.245c0-.399-.078-.78-.22-1.128Zm0 0a15.998 15.998 0 0 0 3.388-1.62m-5.043-.025a15.994 15.994 0 0 1 1.622-3.395m3.42 3.42a15.995 15.995 0 0 0 4.764-4.648l3.876-5.814a1.151 1.151 0 0 0-1.597-1.597L14.146 6.32a15.996 15.996 0 0 0-4.649 4.763m3.42 3.42a6.776 6.776 0 0 0-3.42-3.42" />
                </svg>
            </button>
            <button :class="{'selected-btn': virtualCanvas.brush.value.color == virtualCanvas.eraserBrush}" @click="virtualCanvas.selectEraseMode()">
                <svg width="20px" height="20px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5.50506 11.4096L6.03539 11.9399L5.50506 11.4096ZM3 14.9522H2.25H3ZM12.5904 18.4949L12.0601 17.9646L12.5904 18.4949ZM9.04776 21V21.75V21ZM11.4096 5.50506L10.8792 4.97473L11.4096 5.50506ZM13.241 17.8444C13.5339 18.1373 14.0088 18.1373 14.3017 17.8444C14.5946 17.5515 14.5946 17.0766 14.3017 16.7837L13.241 17.8444ZM7.21629 9.69832C6.9234 9.40543 6.44852 9.40543 6.15563 9.69832C5.86274 9.99122 5.86274 10.4661 6.15563 10.759L7.21629 9.69832ZM16.073 16.073C16.3659 15.7801 16.3659 15.3053 16.073 15.0124C15.7801 14.7195 15.3053 14.7195 15.0124 15.0124L16.073 16.073ZM18.4676 11.5559C18.1759 11.8499 18.1777 12.3248 18.4718 12.6165C18.7658 12.9083 19.2407 12.9064 19.5324 12.6124L18.4676 11.5559ZM6.03539 11.9399L11.9399 6.03539L10.8792 4.97473L4.97473 10.8792L6.03539 11.9399ZM6.03539 17.9646C5.18538 17.1146 4.60235 16.5293 4.22253 16.0315C3.85592 15.551 3.75 15.2411 3.75 14.9522H2.25C2.25 15.701 2.56159 16.3274 3.03 16.9414C3.48521 17.538 4.1547 18.2052 4.97473 19.0253L6.03539 17.9646ZM4.97473 10.8792C4.1547 11.6993 3.48521 12.3665 3.03 12.9631C2.56159 13.577 2.25 14.2035 2.25 14.9522H3.75C3.75 14.6633 3.85592 14.3535 4.22253 13.873C4.60235 13.3752 5.18538 12.7899 6.03539 11.9399L4.97473 10.8792ZM12.0601 17.9646C11.2101 18.8146 10.6248 19.3977 10.127 19.7775C9.64651 20.1441 9.33665 20.25 9.04776 20.25V21.75C9.79649 21.75 10.423 21.4384 11.0369 20.97C11.6335 20.5148 12.3008 19.8453 13.1208 19.0253L12.0601 17.9646ZM4.97473 19.0253C5.79476 19.8453 6.46201 20.5148 7.05863 20.97C7.67256 21.4384 8.29902 21.75 9.04776 21.75V20.25C8.75886 20.25 8.449 20.1441 7.9685 19.7775C7.47069 19.3977 6.88541 18.8146 6.03539 17.9646L4.97473 19.0253ZM17.9646 6.03539C18.8146 6.88541 19.3977 7.47069 19.7775 7.9685C20.1441 8.449 20.25 8.75886 20.25 9.04776H21.75C21.75 8.29902 21.4384 7.67256 20.97 7.05863C20.5148 6.46201 19.8453 5.79476 19.0253 4.97473L17.9646 6.03539ZM19.0253 4.97473C18.2052 4.1547 17.538 3.48521 16.9414 3.03C16.3274 2.56159 15.701 2.25 14.9522 2.25V3.75C15.2411 3.75 15.551 3.85592 16.0315 4.22253C16.5293 4.60235 17.1146 5.18538 17.9646 6.03539L19.0253 4.97473ZM11.9399 6.03539C12.7899 5.18538 13.3752 4.60235 13.873 4.22253C14.3535 3.85592 14.6633 3.75 14.9522 3.75V2.25C14.2035 2.25 13.577 2.56159 12.9631 3.03C12.3665 3.48521 11.6993 4.1547 10.8792 4.97473L11.9399 6.03539ZM14.3017 16.7837L7.21629 9.69832L6.15563 10.759L13.241 17.8444L14.3017 16.7837ZM15.0124 15.0124L12.0601 17.9646L13.1208 19.0253L16.073 16.073L15.0124 15.0124ZM19.5324 12.6124C20.1932 11.9464 20.7384 11.3759 21.114 10.8404C21.5023 10.2869 21.75 9.71511 21.75 9.04776H20.25C20.25 9.30755 20.1644 9.58207 19.886 9.979C19.5949 10.394 19.1401 10.8781 18.4676 11.5559L19.5324 12.6124Z" fill="#FFFFFF"/>
                    <path d="M9 21H21" stroke="#FFFFFF" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
            </button>
            <input type="radio" style="width: 15px; height: 15px;" v-model="virtualCanvas.brush.value.width" value="15">
            <input type="radio" style="width: 25px; height: 25px;" v-model="virtualCanvas.brush.value.width" value="25">
            <input type="radio" style="width: 50px; height: 50px;"v-model="virtualCanvas.brush.value.width" value="50">
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