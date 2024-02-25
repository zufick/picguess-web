export type VirtualCanvas = {
    canvas?: HTMLCanvasElement,
    lines: {
        [key: string] : { // Key is user index in current room
            color: string,
            width: number,
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
};

export type VirtualCanvasLineData = { color: string, width: number, points: [] };

export type VirtualCanvasDrawPoint = { x: number, y: number };