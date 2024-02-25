export type VirtualCanvas = {
    canvas?: HTMLCanvasElement,
    lines: {
        [key: string] : { // Key is user index in current room
            color: string,
            width: number,
            layerIndex?: number,
            points: { 
                x: number, 
                y: number 
            }[]
        }[]
    },
    undoLines: {
        [key: string] : { // Key is user index in current room
            color: string,
            width: number,
            layerIndex?: number,
            points: { 
                x: number, 
                y: number 
            }[]
        }[]
    },
    setCanvas: (canvas: HTMLCanvasElement) => void,
    startNewLine: (id: string, data: VirtualCanvasLineData) => void,
    startNewLocalLine: (data: VirtualCanvasLineData) => void,
    drawLocalPoint: (point: VirtualCanvasDrawPoint) => void,
    drawPoint: (id: string, point: VirtualCanvasDrawPoint) => void,
    redrawCanvas: () => void,
    undoLine: (id: string) => void,
    redoLine: (id: string) => void,
};

export type VirtualCanvasLineData = { color: string, width: number, points: [], layerIndex?: number };

export type VirtualCanvasDrawPoint = { x: number, y: number };
