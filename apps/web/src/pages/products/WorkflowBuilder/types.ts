// src/pages/products/WorkflowBuilder/types.ts
import type { WorkflowNode } from "@/types/product";

export interface CanvasNode {
    id: string;
    type: 'adapterNode';
    position: { x: number; y: number };
    data: {
        label: string;
        adapterConfigId: string;
        endpointPath: string;
        httpMethod: string;
        node: WorkflowNode;
    };
}

export interface CanvasEdge {
    id: string;
    source: string;
    target: string;
    sourceHandle: string;
    targetHandle: string;
}