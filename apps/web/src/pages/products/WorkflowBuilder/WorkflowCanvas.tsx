// src/pages/products/WorkflowBuilder/WorkflowCanvas.tsx
import { useCallback, useRef } from "react";
import { GripVertical, X, ArrowRight } from "lucide-react";
import { Button } from "@kana-consultant/ui-kit";
import { useWorkflowNodes } from "@/hooks/useWorkflow";
import type { CanvasNode, CanvasEdge } from "./types";

interface WorkflowCanvasProps {
    workflowId: string;
    adapterConfigs: any[];
    nodes: CanvasNode[];
    edges: CanvasEdge[];
    setNodes: (nodes: CanvasNode[]) => void;
    setEdges: (edges: CanvasEdge[]) => void;
    onNodeSelect: (node: CanvasNode | null) => void;
    selectedNodeId?: string;
}

export function WorkflowCanvas({
    workflowId,
    adapterConfigs,
    nodes,
    edges,
    setNodes,
    setEdges,
    onNodeSelect,
    selectedNodeId,
}: WorkflowCanvasProps) {
    const canvasRef = useRef<HTMLDivElement>(null);
    const { createNode, deleteNode, reorderNodes, loading } = useWorkflowNodes(workflowId);

    const handleDrop = useCallback(
        async (e: React.DragEvent) => {
            e.preventDefault();
            const adapterConfigId = e.dataTransfer.getData("adapterConfigId");
            if (!adapterConfigId || !canvasRef.current) return;

            const rect = canvasRef.current.getBoundingClientRect();
            const position = {
                x: e.clientX - rect.left - 75,
                y: e.clientY - rect.top - 25,
            };

            const config = adapterConfigs.find((c) => c.id === adapterConfigId);
            if (!config) return;

            const stepOrder = nodes.length + 1;

            try {
                const newNode = await createNode({
                    adapterConfigId,
                    stepOrder,
                    inputMapping: {},
                });

                if (newNode) {
                    const canvasNode: CanvasNode = {
                        id: newNode.id,
                        type: "adapterNode",
                        position,
                        data: {
                            label: config.endpointPath || `Node ${stepOrder}`,
                            adapterConfigId,
                            endpointPath: config.endpointPath,
                            httpMethod: config.httpMethod,
                            node: newNode,
                        },
                    };
                    setNodes([...nodes, canvasNode]);
                }
            } catch (error) {
                console.error("Failed to create node:", error);
            }
        },
        [adapterConfigs, nodes, createNode, setNodes]
    );

    const handleDragOver = (e: React.DragEvent) => {
        e.preventDefault();
    };

    const handleRemoveNode = async (nodeId: string) => {
        await deleteNode(nodeId);
        setNodes(nodes.filter((n) => n.id !== nodeId));
        setEdges(edges.filter((e) => e.source !== nodeId && e.target !== nodeId));
    };

    const handleConnectNodes = async (sourceId: string, targetId: string) => {
        const edgeExists = edges.some(
            (e) => e.source === sourceId && e.target === targetId
        );
        if (edgeExists) return;

        const newEdge: CanvasEdge = {
            id: `e-${sourceId}-${targetId}`,
            source: sourceId,
            target: targetId,
            sourceHandle: "bottom",
            targetHandle: "top",
        };

        setEdges([...edges, newEdge]);

        // Update node order
        const updatedNodes = [...nodes];
        const targetIndex = updatedNodes.findIndex((n) => n.id === targetId);
        if (targetIndex !== -1) {
            updatedNodes[targetIndex].data.node.stepOrder =
                updatedNodes.filter((_, i) => i <= targetIndex).length;
            setNodes(updatedNodes);
        }
    };

    return (
        <div className="h-full flex flex-col">
            {/* Toolbar */}
            <div className="flex items-center gap-2 p-2 border-b">
                <span className="text-xs text-muted-foreground">
                    Drag adapter dari panel kiri ke canvas
                </span>
            </div>

            {/* Canvas */}
            <div
                ref={canvasRef}
                className="flex-1 relative bg-muted/30"
                onDrop={handleDrop}
                onDragOver={handleDragOver}
            >
                {nodes.map((node) => (
                    <div
                        key={node.id}
                        className={`absolute bg-card border rounded-lg p-3 shadow-sm cursor-pointer min-w-[150px] ${selectedNodeId === node.id ? "ring-2 ring-primary" : ""
                            }`}
                        style={{ left: node.position.x, top: node.position.y }}
                        onClick={() => onNodeSelect(node)}
                        draggable
                        onDragStart={(e) => {
                            e.dataTransfer.setData("nodeId", node.id);
                        }}
                        onDragEnd={(e) => {
                            if (!canvasRef.current) return;
                            const rect = canvasRef.current.getBoundingClientRect();
                            const newNodes = nodes.map((n) =>
                                n.id === node.id
                                    ? {
                                        ...n,
                                        position: {
                                            x: e.clientX - rect.left - 75,
                                            y: e.clientY - rect.top - 25,
                                        },
                                    }
                                    : n
                            );
                            setNodes(newNodes);
                        }}
                    >
                        <div className="flex items-center justify-between mb-1">
                            <div className="flex items-center gap-1">
                                <GripVertical className="h-3 w-3 text-muted-foreground" />
                                <span className="text-xs font-medium">
                                    {node.data.httpMethod}
                                </span>
                            </div>
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-5 w-5"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleRemoveNode(node.id);
                                }}
                            >
                                <X className="h-3 w-3" />
                            </Button>
                        </div>
                        <p className="text-xs text-muted-foreground truncate max-w-[130px]">
                            {node.data.endpointPath || "No endpoint"}
                        </p>

                        {/* Connection handles */}
                        <div className="flex justify-between mt-2">
                            <div
                                className="h-3 w-3 rounded-full bg-green-500 cursor-crosshair"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    // Start connection logic
                                }}
                            />
                            <div
                                className="h-3 w-3 rounded-full bg-blue-500 cursor-crosshair"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    // Complete connection logic
                                }}
                            />
                        </div>
                    </div>
                ))}

                {nodes.length === 0 && (
                    <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
                        Drop adapter di sini untuk memulai flow
                    </div>
                )}
            </div>

            {/* Adapter Pool */}
            <div className="border-t p-2">
                <p className="text-xs font-medium mb-2">Adapter Pool</p>
                <div className="flex gap-2 overflow-x-auto">
                    {adapterConfigs.map((config) => (
                        <div
                            key={config.id}
                            className="bg-muted rounded px-3 py-1 text-xs cursor-grab whitespace-nowrap"
                            draggable
                            onDragStart={(e) => {
                                e.dataTransfer.setData("adapterConfigId", config.id);
                            }}
                        >
                            {config.httpMethod} {config.endpointPath}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}