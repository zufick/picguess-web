<script setup lang="ts">
import { nextTick, onMounted, provide, ref, watch } from "vue";
import { useWebsocketStore } from "@/stores/WebsocketStore";
import { 
    type VirtualCanvas, 
    type VirtualCanvasLine, 
    type VirtualCanvasDrawPoint, 
    type VirtualCanvasBrush,
    VirtualCanvasBrushMode
} from '@/types/VirtualCanvas'
import DrawingCanvasLayer from "./DrawingCanvasLayer.vue";
import DrawingBrushPreview from "./DrawingBrushPreview.vue";
import DrawingCanvasClear from './DrawingCanvasClear.vue'
import DrawingCanvasBrushSize from "./DrawingCanvasBrushSize.vue";
import ColorPicker from "./ColorPicker.vue";

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
    brushMode: ref(VirtualCanvasBrushMode.Paint),
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
        this.brushMode.value = VirtualCanvasBrushMode.Paint;
    },
    selectEraseMode() {
        this.brush.value.color = this.eraserBrush;
        this.brushMode.value = VirtualCanvasBrushMode.Erase;
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
            <ColorPicker v-model="virtualCanvas.paintBrush.value"/>
            <button @click="virtualCanvas.selectPaintMode()" :class="{'selected-btn': virtualCanvas.brushMode.value == VirtualCanvasBrushMode.Paint}">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17.66 17.66" class="w-5 h-5">
                    <g>
                        <path fill="white" class="cls-1" d="m16.91,17.66c.41,0,.75-.34.75-.75s-.34-.75-.75-.75H7.53s.03-.03.05-.04c0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0l9.42-9.67c.75-.75.87-1.93.27-2.81-.88-1.28-1.97-2.37-3.25-3.25-.87-.6-2.05-.48-2.8.27-.97.97-3.09,3.02-5.13,5-1.91,1.85-3.76,3.64-4.55,4.43,0,0,0,0,0,0,0,0,0,0,0,0-.46.46-.78,1.04-.91,1.69L.05,14.55c-.18.86.09,1.73.7,2.35.49.49,1.15.76,1.82.76h14.33ZM13.18,1.62c1.13.77,2.09,1.73,2.86,2.86.19.28.15.66-.1.91l-.81.83-3.68-3.68c.33-.32.61-.6.83-.82.24-.24.62-.28.9-.09Zm-6.05,5.11c1.13-1.09,2.27-2.2,3.24-3.15l3.71,3.71-6.06,6.22c-.05-.1-.1-.2-.17-.29-.92-1.34-2.07-2.49-3.41-3.41-.1-.07-.2-.12-.3-.17.86-.84,1.91-1.86,2.99-2.9ZM1.52,14.86l.58-2.78c.07-.36.25-.68.5-.93.27-.27.69-.31.99-.1,1.19.81,2.21,1.83,3.02,3.02.21.31.17.72-.1.99,0,0,0,0,0,0-.26.25-.58.43-.93.5l-2.78.58c-.36.08-.72-.03-.98-.29-.26-.26-.37-.62-.29-.99Z"/>
                    </g>
                </svg>
            </button>
            <button :class="{'selected-btn': virtualCanvas.brushMode.value == VirtualCanvasBrushMode.Erase}" @click="virtualCanvas.selectEraseMode()">
                <svg width="20px" height="20px" viewBox="0 0 17.96 16.06" fill="white" xmlns="http://www.w3.org/2000/svg">
                    <path class="cls-1" d="m2.58,14.56H.75c-.41,0-.75.34-.75.75s.34.75.75.75h1.83c.41,0,.75-.34.75-.75s-.34-.75-.75-.75Z"/>
      <path class="cls-1" d="m17.21,14.56h-6.32l6.23-6.23c.94-.94,1.08-2.4.34-3.48-1.18-1.72-2.65-3.19-4.37-4.37-1.08-.74-2.55-.6-3.48.34l-5.4,5.4s0,0,0,0,0,0,0,0l-1.52,1.52c-1.07,1.07-1.07,2.81,0,3.88l3.64,3.64c.52.52,1.21.8,1.94.8,0,0,.01,0,.02,0,0,0,0,0,0,0h8.93c.41,0,.75-.34.75-.75s-.34-.75-.75-.75ZM10.67,1.87c.25-.25.57-.37.89-.37.24,0,.48.07.68.21,1.57,1.07,2.91,2.41,3.98,3.98h0c.33.49.27,1.15-.16,1.58l-4.88,4.88-5.4-5.4L10.67,1.87ZM3.75,10.56c-.24-.24-.36-.55-.36-.88s.13-.65.36-.88l.99-.99,5.4,5.4-.99.99c-.47.47-1.29.47-1.76,0l-3.64-3.64Z"/>
                </svg>
            </button>
            <DrawingCanvasBrushSize v-model="virtualCanvas.brush.value.width"/>
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
        user-select: none;
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