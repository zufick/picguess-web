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
    },
    startNewLocalLine(data: VirtualCanvasLineData) {
        if (this.lines["local"] == undefined) {
            this.lines["local"] = []
        }

        let totalLayers = Object.values(this.lines).flat(1).length
        data.layerIndex = totalLayers
        this.lines["local"].push(data);
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

        /*for (const [userLineKey, userLines] of Object.entries(this.lines)) {
            for (const [lineKey, line] of Object.entries(userLines)) {
                ctx!.lineWidth = line.width;
                ctx!.lineCap = "round";
                ctx!.strokeStyle = line.color;
                
                line.points.forEach(point => {
                    ctx?.lineTo(point.x, point.y);
                });
                ctx?.stroke();
                ctx?.beginPath();
            }
        }*/

        //ctx!.lineWidth = data.width!;
        //ctx!.lineCap = "round";
        //ctx!.strokeStyle = data.color!;
        //ctx?.lineTo(data.x, data.y);
        //ctx?.stroke();
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

    let rect = drawingCanvas.value?.getBoundingClientRect();
    let point = { 
        x: event.clientX - (rect?.left ?? 0),
        y: event.clientY - (rect?.top ?? 0),
    }

    //redrawCanvas(drawData, true)
    virtualCanvas.drawLocalPoint(point)
    wsStore.sendDrawPoint(point)
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
}

function mouseUp(event: any) {
    isMouseDown = false;
    let ctx = drawingCanvas.value?.getContext("2d");
}

function redrawCanvas(data: DrawData, isLocal = false) {
    let ctx = drawingCanvas.value?.getContext("2d");
    ctx!.lineWidth = data.width!;
    ctx!.lineCap = "round";
    ctx!.strokeStyle = data.color!;
    ctx?.lineTo(data.x, data.y);
    ctx?.stroke();
}

wsStore.setVirtualCanvas(virtualCanvas)

onMounted(() => {
    virtualCanvas.setCanvas(drawingCanvas.value!)
})
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
