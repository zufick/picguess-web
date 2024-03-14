import type { Ref } from "vue";

export type VirtualCanvas = {
    canvas?: HTMLCanvasElement,
    sendPointsBuffer: number[],
    sendPointsBufferLimit: number,
    lines: Ref<VirtualCanvasLine[]>,
    linesToDelete: VirtualCanvasLine[],
    brush: Ref<VirtualCanvasBrush>,
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
    undoLine: (id: string) => void,
    redoLine: (id: string) => void,
    drawClear: (id: string) => void,
    addReadyToMergeCanvas: (line: VirtualCanvasLine, canvas: HTMLCanvasElement) => void,
    mergeReadyCanvases: () => void,
    deleteReadyToDeleteLines: (mergedOnto: VirtualCanvasLine) => void,
};


export type VirtualCanvasDrawPoint = { x: number, y: number };

export type VirtualCanvasLine = {
    id?: string,
    userId?: string,
    color: string,
    width: number,
    undo?: boolean,
    newPoints?: VirtualCanvasDrawPoint[],
    readyToMergeCanvas?: HTMLCanvasElement,
    mergeCanvases?: HTMLCanvasElement[],
};

export type VirtualCanvasBrush = {
    color: string,
    width: number,
    point: VirtualCanvasDrawPoint,
    visible: boolean
}