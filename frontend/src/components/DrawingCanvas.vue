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
const drawColor = ref("#000000")
const drawWidth = ref(5)
const drawClearDialog = ref(false)

const virtualCanvas: VirtualCanvas = {
    canvas: undefined,
    lines: {},
    undoLines: {},
    drawCalls: 0,
    sendPointsBuffer: [],
    sendPointsBufferLimit: 5,
    setCanvas(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
    },
    startNewLine(id: string, data: VirtualCanvasLineData) {
        this.sendPointsBuffer = []

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
        this.sendPointsBuffer.push(point.x, point.y)
        this.drawPoint("local", point)

        if (this.sendPointsBuffer.length >= 2 * this.sendPointsBufferLimit) {
            wsStore.sendDrawPoints(this.sendPointsBuffer)
            this.sendPointsBuffer = []
        }
    },
    drawPoint(id: string, point: VirtualCanvasDrawPoint) {
        let line = this.lines[id];
        if(line[line.length - 1].points == undefined) {
            line[line.length - 1].points = [];
        }
        line[line.length - 1].points.push(point);

        this.redrawCanvas();
    },
    drawPoints(id: string, points: number[]) {
        let line = this.lines[id];
        if(line[line.length - 1].points == undefined) {
            line[line.length - 1].points = [];
        }

        for (let i = 0; i < points.length; i+=2) {
            let point = {
                x: points[i],
                y: points[i+1]
            }
            line[line.length - 1].points.push(point);
        }


        this.redrawCanvas();
    },
    redrawCanvas() {
        window.requestAnimationFrame(() => this.redrawCanvasLoop())
    },
    redrawCanvasLoop() {
        this.drawCalls = 0;

        let ctx = this.canvas?.getContext("2d");
        ctx!.clearRect(0, 0, this.canvas!.width, this.canvas!.height);


        let lines = Object.values(this.lines).flat(1)
        lines.sort((a,b) => a.layerIndex! - b.layerIndex!)

        for (let i = 0; i < lines.length; i++) {
            let line = lines[i]

            ctx!.lineWidth = line.width;
            ctx!.lineCap = "round";
            ctx!.strokeStyle = line.color;
            
            if (line.points != undefined) {
                line.points.forEach(point => {
                    //if (lastPoint.x != 0) {
                    //    ctx?.moveTo(lastPoint.x, lastPoint.y);    
                    //}
                    ctx?.lineTo(point.x, point.y);
                    this.drawCalls++;
                    //lastPoint = point;
                });

            }
            ctx?.stroke();
            ctx?.beginPath();
        }
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
    drawClear(id: string) {
        if (id == "") {
            this.lines = {};
            this.undoLines = {};
        } else {
            this.lines[id] = [];
        }

        this.redrawCanvas();
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
        x: Math.floor(event.clientX - rect!.left ?? 0),
        y: Math.floor(event.clientY - rect!.top ?? 0),
    }

    virtualCanvas.drawLocalPoint(point)

}

function mouseUp(event: any) {
    isMouseDown = false;

    if(virtualCanvas.sendPointsBuffer.length > 0) {
        wsStore.sendDrawPoints(virtualCanvas.sendPointsBuffer)
        virtualCanvas.sendPointsBuffer = []
    }
    
}

function undoClicked() {
    virtualCanvas.undoLine("local");
    wsStore.sendDrawUndo();
}

function redoClicked() {
    virtualCanvas.redoLine("local");
    wsStore.sendDrawRedo();
}

function drawClearClicked() {
    virtualCanvas.drawClear("local");
    wsStore.sendDrawClear();
    drawClearDialog.value = false
}

wsStore.setVirtualCanvas(virtualCanvas)

onMounted(() => {
    virtualCanvas.setCanvas(drawingCanvas.value!)
})
</script>


<template>
    <div class="flex gap-4">
        <div class="canvas-container border-2 rounded-lg border-white/20 overflow-hidden">
            <div class="canvas-bg"></div>
            <canvas ref="drawingCanvas" height="768" width="1024" @mousemove="mouseMove" @mousedown="mouseDown" @mouseup="mouseUp" @mouseout="mouseUp">
            </canvas>
        </div>
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
            <button @click="drawClearDialog = true">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                </svg>
            </button>
            <div v-if="drawClearDialog" class="flex flex-col gap-4 appear-anim" style="margin-top: -60px;">
                <button @click="drawClearClicked()" class="danger-btn">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                    </svg>
                </button>

                <button @click="drawClearDialog = false">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
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
</style>