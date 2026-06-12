// src/pages/products/WorkflowBuilder/WorkflowBuilder.tsx
import { useState, useCallback, useRef, useEffect } from "react";
import { Plus, ChevronDown, ChevronRight, ChevronLeft, ChevronRight as ChevronRightIcon } from "lucide-react";
import { Button, Input } from "@kana-consultant/ui-kit";
import ReactFlow, {
  Controls,
  Background,
  MiniMap,
  addEdge,
  useNodesState,
  useEdgesState,
  BackgroundVariant,
} from "reactflow";
import type { Connection, Node, Edge } from "reactflow";
import "reactflow/dist/style.css";
import { useWorkflow, useWorkflowNodes } from "@/hooks/useWorkflow";
import { WorkflowList } from "./WorkflowList";
import { NodePanel } from "./NodePanel";
import { AdapterNode } from "./AdapterNode";
import {
  getWorkflowNodes,
  saveAllWorkflowNodes,
  type BatchCreateNode,
  type BatchUpdateNode,
  type SaveNodesResponse
} from "@/services/workflow/workflowService";
import type { AdapterConfig } from "@/types/workflow";

interface WorkflowBuilderProps {
  productId: string;
  adapterConfigs: AdapterConfig[];
  onChange?: (data: any) => void;
  initialWorkflowId?: string | null; // Untuk edit mode
}

const nodeTypes = {
  adapterNode: AdapterNode,
};

interface TempWorkflowNode {
  id: string;
  workflowId: string;
  adapterConfigId: string;
  stepOrder: number;
  inputMapping: Record<string, any>;
  nextNodeId?: string | null;
  createdAt: string;
  updatedAt: string;
}

export function WorkflowBuilder({ productId, adapterConfigs, onChange, initialWorkflowId }: WorkflowBuilderProps) {
  const {
    workflows,
    loading: workflowLoading,
    createWorkflow,
    deleteWorkflow,
  } = useWorkflow(productId);

  const [selectedWorkflowId, setSelectedWorkflowId] = useState<string | null>(initialWorkflowId || null);
  const [workflowName, setWorkflowName] = useState("");
  const [showNewInput, setShowNewInput] = useState(false);
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);
  const [isDragSectionExpanded, setIsDragSectionExpanded] = useState(true);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isInitialLoading, setIsInitialLoading] = useState(false);

  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const reactFlowWrapper = useRef<HTMLDivElement>(null);

  const [tempNodes, setTempNodes] = useState<Record<string, TempWorkflowNode[]>>({});
  const [deletedNodes, setDeletedNodes] = useState<Record<string, string[]>>({});
  const [updatedNodes, setUpdatedNodes] = useState<Record<string, Record<string, any>>>({});

  // Load initial workflow jika ada initialWorkflowId
  useEffect(() => {
    if (initialWorkflowId && workflows.length > 0) {
      const workflowExists = workflows.some(w => w.id === initialWorkflowId);
      if (workflowExists) {
        setSelectedWorkflowId(initialWorkflowId);
        loadWorkflowNodes(initialWorkflowId);
      }
    }
  }, [initialWorkflowId, workflows]);

  // Trigger onChange when data changes
  useEffect(() => {
    if (onChange && selectedWorkflowId) {
      const workflowData = {
        workflowId: selectedWorkflowId,
        nodes: nodes.map(node => ({
          id: node.id,
          type: node.type,
          position: node.position,
          data: node.data,
        })),
        edges: edges,
        hasChanges: hasChanges(),
      };
      onChange(workflowData);
    }
  }, [nodes, edges, selectedWorkflowId]);

  const loadWorkflowNodes = useCallback(async (workflowId: string) => {
    setIsInitialLoading(true);
    try {
      const backendNodes = await getWorkflowNodes(workflowId);
      if (!backendNodes) return;

      const flowNodes: Node[] = backendNodes.map((wn: any, index: number) => ({
        id: wn.id,
        type: "adapterNode",
        position: { x: 250, y: index * 150 + 50 },
        data: {
          label: `Node ${wn.stepOrder}`,
          workflowNode: wn,
          adapterConfig: adapterConfigs.find((c) => c.id === wn.adapterConfigId),
          isTemp: false,
        },
      }));
      setNodes(flowNodes);

      const flowEdges: Edge[] = backendNodes
        .filter((wn: any) => wn.nextNodeId)
        .map((wn: any) => ({
          id: `e-${wn.id}-${wn.nextNodeId}`,
          source: wn.id,
          target: wn.nextNodeId as string,
          sourceHandle: null,
          targetHandle: null,
          animated: true,
        }));
      setEdges(flowEdges);

      setTempNodes(prev => ({ ...prev, [workflowId]: [] }));
      setDeletedNodes(prev => ({ ...prev, [workflowId]: [] }));
      setUpdatedNodes(prev => ({ ...prev, [workflowId]: {} }));
    } catch (error) {
      console.error("Failed to load workflow nodes:", error);
    } finally {
      setIsInitialLoading(false);
    }
  }, [adapterConfigs, setNodes, setEdges]);

  const handleSelectWorkflow = async (workflowId: string) => {
    setSelectedWorkflowId(workflowId);
    await loadWorkflowNodes(workflowId);

    // Trigger onChange untuk workflow yang dipilih
    if (onChange) {
      onChange({ workflowId, isNew: false });
    }
  };

  const handleCreateWorkflow = async () => {
    if (!workflowName.trim()) return;
    const wf = await createWorkflow(workflowName.trim());
    if (wf) {
      setWorkflowName("");
      setShowNewInput(false);
      setSelectedWorkflowId(wf.id);
      setTempNodes(prev => ({ ...prev, [wf.id]: [] }));

      if (onChange) {
        onChange({ workflowId: wf.id, isNew: true, name: workflowName });
      }
    }
  };

  const generateTempId = () => {
    return `temp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  };

  const onConnect = useCallback(
    async (connection: Connection) => {
      if (!connection.source || !connection.target) return;

      setEdges((eds) => addEdge({
        ...connection,
        id: `e-${connection.source}-${connection.target}`,
        animated: true
      }, eds));

      if (selectedWorkflowId) {
        const isTempSource = connection.source.startsWith('temp_');
        const isTempTarget = connection.target.startsWith('temp_');

        if (isTempSource || isTempTarget) {
          setTempNodes(prev => {
            const workflowTempNodes = [...(prev[selectedWorkflowId] || [])];
            const nodeIndex = workflowTempNodes.findIndex(n => n.id === connection.source);
            if (nodeIndex !== -1) {
              workflowTempNodes[nodeIndex] = {
                ...workflowTempNodes[nodeIndex],
                nextNodeId: connection.target
              };
            }
            return { ...prev, [selectedWorkflowId]: workflowTempNodes };
          });
        } else {
          setUpdatedNodes(prev => ({
            ...prev,
            [selectedWorkflowId]: {
              ...prev[selectedWorkflowId],
              [connection.source as any]: {
                ...prev[selectedWorkflowId]?.[connection.source as any],
                nextNodeId: connection.target
              }
            }
          }));
        }
      }
    },
    [setEdges, selectedWorkflowId]
  );

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop = useCallback(
    async (event: React.DragEvent) => {
      event.preventDefault();
      const adapterConfigId = event.dataTransfer.getData("adapterConfigId");
      if (!adapterConfigId || !reactFlowWrapper.current || !selectedWorkflowId) return;

      const config = adapterConfigs.find((c) => c.id === adapterConfigId);
      if (!config) return;

      const rect = reactFlowWrapper.current.getBoundingClientRect();
      const position = {
        x: event.clientX - rect.left - 75,
        y: event.clientY - rect.top - 25,
      };

      const tempId = generateTempId();
      const stepOrder = nodes.filter(n => !n.id.startsWith('temp_')).length +
        (tempNodes[selectedWorkflowId]?.length || 0) + 1;

      const tempWorkflowNode: TempWorkflowNode = {
        id: tempId,
        workflowId: selectedWorkflowId,
        adapterConfigId,
        stepOrder,
        inputMapping: {},
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      setTempNodes(prev => ({
        ...prev,
        [selectedWorkflowId]: [...(prev[selectedWorkflowId] || []), tempWorkflowNode]
      }));

      const flowNode: Node = {
        id: tempId,
        type: "adapterNode",
        position,
        data: {
          label: config.endpointPath || `Node ${stepOrder}`,
          workflowNode: tempWorkflowNode,
          adapterConfig: config,
          isTemp: true,
        },
      };
      setNodes((nds) => [...nds, flowNode]);
    },
    [adapterConfigs, nodes.length, selectedWorkflowId, tempNodes, setNodes]
  );

  const handleNodeClick = useCallback((_event: React.MouseEvent, node: Node) => {
    setSelectedNode(node);
  }, []);

  const handleDeleteNode = useCallback(
    async (nodeId: string) => {
      if (nodeId.startsWith('temp_') && selectedWorkflowId) {
        setTempNodes(prev => ({
          ...prev,
          [selectedWorkflowId]: (prev[selectedWorkflowId] || []).filter(n => n.id !== nodeId)
        }));
      } else if (selectedWorkflowId) {
        setDeletedNodes(prev => ({
          ...prev,
          [selectedWorkflowId]: [...(prev[selectedWorkflowId] || []), nodeId]
        }));

        setUpdatedNodes(prev => {
          const currentUpdates = prev[selectedWorkflowId] || {};
          const { [nodeId]: _, ...rest } = currentUpdates;
          return {
            ...prev,
            [selectedWorkflowId]: rest
          };
        });
      }

      setNodes((nds) => nds.filter((n) => n.id !== nodeId));
      setEdges((eds) => eds.filter((e) => e.source !== nodeId && e.target !== nodeId));
      setSelectedNode(null);
    },
    [setNodes, setEdges, selectedWorkflowId]
  );

  const handleUpdateNode = useCallback(
    async (nodeId: string, updates: Record<string, any>) => {
      if (nodeId.startsWith('temp_') && selectedWorkflowId) {
        setTempNodes(prev => {
          const workflowTempNodes = [...(prev[selectedWorkflowId] || [])];
          const nodeIndex = workflowTempNodes.findIndex(n => n.id === nodeId);
          if (nodeIndex !== -1) {
            workflowTempNodes[nodeIndex] = {
              ...workflowTempNodes[nodeIndex],
              ...updates
            };
          }
          return { ...prev, [selectedWorkflowId]: workflowTempNodes };
        });

        setNodes((nds) =>
          nds.map((n) =>
            n.id === nodeId
              ? {
                ...n,
                data: {
                  ...n.data,
                  workflowNode: { ...n.data.workflowNode, ...updates },
                },
              }
              : n
          )
        );
      } else if (selectedWorkflowId) {
        setUpdatedNodes(prev => ({
          ...prev,
          [selectedWorkflowId]: {
            ...prev[selectedWorkflowId],
            [nodeId]: {
              ...prev[selectedWorkflowId]?.[nodeId],
              ...updates
            }
          }
        }));

        setNodes((nds) =>
          nds.map((n) =>
            n.id === nodeId
              ? {
                ...n,
                data: {
                  ...n.data,
                  workflowNode: { ...n.data.workflowNode, ...updates },
                },
              }
              : n
          )
        );
      }
    },
    [setNodes, selectedWorkflowId]
  );

  const saveAllChanges = useCallback(async () => {
    if (!selectedWorkflowId) return;

    const unsavedNodes = tempNodes[selectedWorkflowId] || [];
    const nodesToDelete = deletedNodes[selectedWorkflowId] || [];
    const nodesToUpdate = updatedNodes[selectedWorkflowId] || {};

    if (unsavedNodes.length === 0 && nodesToDelete.length === 0 && Object.keys(nodesToUpdate).length === 0) {
      return;
    }

    setIsSaving(true);

    try {
      const toCreate: BatchCreateNode[] = unsavedNodes.map(tempNode => ({
        tempId: tempNode.id,
        adapterConfigId: tempNode.adapterConfigId,
        stepOrder: tempNode.stepOrder,
        inputMapping: tempNode.inputMapping,
        nextNodeId: tempNode.nextNodeId || null,
      }));

      const toUpdate: BatchUpdateNode[] = Object.entries(nodesToUpdate).map(([id, updates]) => ({
        id,
        updates
      }));

      const result: SaveNodesResponse = await saveAllWorkflowNodes(selectedWorkflowId, {
        toCreate,
        toUpdate,
        toDelete: nodesToDelete,
      });

      if (result.created && result.created.length > 0) {
        const tempToRealMap = new Map<string, any>();
        result.created.forEach(saved => {
          tempToRealMap.set(saved.tempId, saved);
        });

        setNodes((nds) =>
          nds.map((n) => {
            const savedNode = tempToRealMap.get(n.id);
            if (savedNode) {
              return {
                ...n,
                id: savedNode.id,
                data: {
                  ...n.data,
                  workflowNode: savedNode,
                  isTemp: false,
                },
              };
            }
            return n;
          })
        );

        setEdges((eds) =>
          eds.map((e) => {
            const sourceNode = tempToRealMap.get(e.source);
            const targetNode = tempToRealMap.get(e.target);

            if (sourceNode || targetNode) {
              const newSource = sourceNode ? sourceNode.id : e.source;
              const newTarget = targetNode ? targetNode.id : e.target;
              return {
                ...e,
                source: newSource,
                target: newTarget,
                id: `e-${newSource}-${newTarget}`,
              };
            }
            return e;
          })
        );
      }

      if (result.updated && result.updated.length > 0) {
        setNodes((nds) =>
          nds.map((n) => {
            const updatedNode = result.updated.find((u: any) => u.id === n.id);
            if (updatedNode) {
              return {
                ...n,
                data: {
                  ...n.data,
                  workflowNode: updatedNode,
                },
              };
            }
            return n;
          })
        );
      }

      if (result.deleted && result.deleted.length > 0) {
        setNodes((nds) => nds.filter(n => !result.deleted.includes(n.id)));
        setEdges((eds) => eds.filter(e =>
          !result.deleted.includes(e.source) && !result.deleted.includes(e.target)
        ));
      }

      setTempNodes(prev => ({ ...prev, [selectedWorkflowId]: [] }));
      setDeletedNodes(prev => ({ ...prev, [selectedWorkflowId]: [] }));
      setUpdatedNodes(prev => ({ ...prev, [selectedWorkflowId]: {} }));

    } catch (error) {
      console.error("Failed to save changes:", error);
    } finally {
      setIsSaving(false);
    }
  }, [selectedWorkflowId, tempNodes, deletedNodes, updatedNodes, setNodes, setEdges]);

  const hasChanges = useCallback(() => {
    if (!selectedWorkflowId) return false;
    return (tempNodes[selectedWorkflowId]?.length > 0 ||
      deletedNodes[selectedWorkflowId]?.length > 0 ||
      Object.keys(updatedNodes[selectedWorkflowId] || {}).length > 0);
  }, [selectedWorkflowId, tempNodes, deletedNodes, updatedNodes]);

  const getChangesCount = () => {
    if (!selectedWorkflowId) return 0;
    return (tempNodes[selectedWorkflowId]?.length || 0) +
      (deletedNodes[selectedWorkflowId]?.length || 0) +
      Object.keys(updatedNodes[selectedWorkflowId] || {}).length;
  };

  if (isInitialLoading) {
    return (
      <div className="flex h-[600px] items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-sm text-muted-foreground">Loading workflow...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-[600px] gap-4 relative">
      {/* Left Sidebar */}
      {isSidebarOpen && (
        <div className="w-64 border-r pr-4 flex flex-col gap-4 transition-all duration-300">
          <div>
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold text-sm">WORKFLOWS</h3>
              <Button variant="ghost" size="sm" onClick={() => setShowNewInput(true)}>
                <Plus className="h-4 w-4" />
              </Button>
            </div>

            {showNewInput && (
              <div className="flex gap-2 mb-2">
                <Input
                  placeholder="Nama workflow"
                  value={workflowName}
                  onChange={(e) => setWorkflowName(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleCreateWorkflow()}
                  autoFocus
                  className="text-sm"
                />
                <Button type="button" size="sm" onClick={handleCreateWorkflow}>OK</Button>
              </div>
            )}

            <WorkflowList
              workflows={workflows}
              selectedId={selectedWorkflowId}
              onSelect={handleSelectWorkflow}
              onDelete={deleteWorkflow}
              loading={workflowLoading}
            />
          </div>

          {selectedWorkflowId && (
            <div className="mt-4">
              <button
                onClick={() => setIsDragSectionExpanded(!isDragSectionExpanded)}
                className="flex items-center justify-between w-full text-xs font-medium text-muted-foreground hover:text-foreground transition-colors mb-2"
              >
                <div className="flex items-center gap-1">
                  {isDragSectionExpanded ? (
                    <ChevronDown className="h-3 w-3" />
                  ) : (
                    <ChevronRight className="h-3 w-3" />
                  )}
                  <span>DRAG TO CANVAS</span>
                </div>
                {!isDragSectionExpanded && (
                  <span className="text-[10px] bg-muted px-1.5 py-0.5 rounded">
                    {adapterConfigs.length} items
                  </span>
                )}
              </button>

              {isDragSectionExpanded && (
                <>
                  <div className="space-y-1 max-h-64 overflow-y-auto">
                    {adapterConfigs.map((config) => (
                      <div
                        key={config.id}
                        className="bg-muted rounded px-2 py-1.5 text-xs cursor-grab hover:bg-primary/10 transition-colors"
                        draggable
                        onDragStart={(e) => {
                          e.dataTransfer.setData("adapterConfigId", config.id);
                          e.dataTransfer.effectAllowed = "move";
                        }}
                      >
                        <span className="font-mono font-bold text-primary">{config.httpMethod}</span>{" "}
                        <span className="truncate">{config.endpointPath}</span>
                      </div>
                    ))}
                  </div>

                  {hasChanges() && (
                    <Button
                      size="sm"
                      className="mt-4 w-full"
                      onClick={saveAllChanges}
                      disabled={isSaving}
                    >
                      {isSaving ? "Saving..." : `Save All Changes (${getChangesCount()})`}
                    </Button>
                  )}
                </>
              )}
            </div>
          )}
        </div>
      )}

      {/* Toggle Sidebar Button */}
      <button
        onClick={() => setIsSidebarOpen(!isSidebarOpen)}
        className="absolute left-0 top-1/2 -translate-y-1/2 z-10 bg-background border rounded-r-md p-1 hover:bg-muted transition-colors"
        style={{ left: isSidebarOpen ? '256px' : '0px' }}
      >
        {isSidebarOpen ? (
          <ChevronLeft className="h-4 w-4" />
        ) : (
          <ChevronRightIcon className="h-4 w-4" />
        )}
      </button>

      {/* Canvas */}
      <div className="flex-1 border rounded-lg" ref={reactFlowWrapper}>
        {selectedWorkflowId ? (
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onDrop={onDrop}
            onDragOver={onDragOver}
            onNodeClick={handleNodeClick}
            nodeTypes={nodeTypes}
            fitView
            deleteKeyCode="Delete"
            onNodesDelete={(deletedNodes) => {
              deletedNodes.forEach((n) => handleDeleteNode(n.id));
            }}
          >
            <Controls />
            <MiniMap />
            <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
          </ReactFlow>
        ) : (
          <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
            Pilih atau buat workflow untuk mulai
          </div>
        )}
      </div>

      {/* Right Panel */}
      {selectedNode && (
        <div className="w-150 border-l pl-4 overflow-y-auto">
          <NodePanel
            node={selectedNode}
            allNodes={nodes}
            adapterConfigs={adapterConfigs}
            onClose={() => setSelectedNode(null)}
            onDelete={handleDeleteNode}
            onUpdate={handleUpdateNode}
          />
        </div>
      )}
    </div>
  );
}