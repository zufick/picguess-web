export type VirtualCanvas = {
    canvas?: HTMLCanvasElement,
    drawCalls: number,
    sendPointsBuffer: number[],
    sendPointsBufferLimit: number,
    lines: {
        [key: string] : { // Key is user index in current room
            color: string,
            width: number,
            layerIndex?: number,
            points: VirtualCanvasDrawPoint[]
        }[]
    },
    undoLines: {
        [key: string] : { // Key is user index in current room
            color: string,
            width: number,
            layerIndex?: number,
            points: VirtualCanvasDrawPoint[]
        }[]
    },
    setCanvas: (canvas: HTMLCanvasElement) => void,
    startNewLine: (id: string, data: VirtualCanvasLineData) => void,
    startNewLocalLine: (data: VirtualCanvasLineData) => void,
    drawLocalPoint: (point: VirtualCanvasDrawPoint) => void,
    drawPoint: (id: string, point: VirtualCanvasDrawPoint) => void,
    drawPoints: (id: string, points: number[]) => void,
    redrawCanvas: () => void,
    redrawCanvasLoop: () => void,
    undoLine: (id: string) => void,
    redoLine: (id: string) => void,
    drawClear: (id: string) => void
};

export type VirtualCanvasLineData = { color: string, width: number, points: [], layerIndex?: number };

export type VirtualCanvasDrawPoint = { x: number, y: number };
