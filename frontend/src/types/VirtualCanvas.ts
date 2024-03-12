import type { Ref } from "vue";

export type VirtualCanvas = {
    canvas?: HTMLCanvasElement,
    drawCalls: number,
    sendPointsBuffer: number[],
    sendPointsBufferLimit: number,
    lines: Ref<VirtualCanvasLine[]>,
    lastUserLines: {
        [key: string] : VirtualCanvasLine // userId -> line
    },
    setCanvas: (canvas: HTMLCanvasElement) => void,
    startNewLine: (id: string, data: VirtualCanvasLine) => void,
    startNewLocalLine: (data: VirtualCanvasLine) => void,
    drawLocalPoint: (point: VirtualCanvasDrawPoint) => void,
    drawPoint: (id: string, point: VirtualCanvasDrawPoint) => void,
    drawPoints: (id: string, points: VirtualCanvasDrawPoint[]) => void,
    drawPointArray: (id: string, points: number[]) => void,
    redrawCanvas: () => void,
    redrawCanvasLoop: () => void,
    undoLine: (id: string) => void,
    redoLine: (id: string) => void,
    drawClear: (id: string) => void
};


export type VirtualCanvasDrawPoint = { x: number, y: number };

export type VirtualCanvasLine = {
    userId?: string,
    color: string,
    width: number,
    undo?: boolean,
    newPoints?: VirtualCanvasDrawPoint[]
};