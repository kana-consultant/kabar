// src/pages/products/WorkflowBuilder/WorkflowBuilder.tsx
import { useState, useCallback, useRef, useEffect } from "react";
import { Plus, ChevronDown, ChevronRight, ChevronLeft, ChevronRight as ChevronRightIcon, Pencil, Trash2, Maximize2, Minimize2 } from "lucide-react";
import { Button, Input, Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import ReactFlow, {
  Controls,
  Background,
  addEdge,
  useNodesState,
  useEdgesState,
  BackgroundVariant,
} from "reactflow";
import type { Connection, Node, Edge } from "reactflow";
import "reactflow/dist/style.css";
import { useWorkflow } from "@/hooks/useWorkflow";
import { WorkflowList } from "./WorkflowList";
import { NodePanel } from "./NodePanel";
import { AdapterNode } from "./AdapterNode";
import { getWorkflowNodes } from "@/services/workflow/workflowService";
import type { AdapterConfig, Product } from "@/types/product";

interface WorkflowBuilderProps {
  productId: string;
  adapterConfigs: AdapterConfig[];
  product: Product;
  onUpdateAdapterConfigs?: (adapterConfigs: AdapterConfig[]) => void;
  onChange?: (workflowId: string) => void;
  initialWorkflowId?: string;
  isFullscreen?: boolean;
  onToggleFullscreen?: () => void;
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
  previousNodeIds?: string[];
  nextNodeId?: string | null;
  createdAt: string;
  updatedAt: string;
}

export function WorkflowBuilder({
  productId,
  adapterConfigs,
  product,
  onUpdateAdapterConfigs,
  onChange,
  initialWorkflowId,
  isFullscreen = false,
  onToggleFullscreen
}: WorkflowBuilderProps) {
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
  const [isInitialLoading, setIsInitialLoading] = useState(false);
  const [showAdapterModal, setShowAdapterModal] = useState(false);
  const [editingAdapter, setEditingAdapter] = useState<AdapterConfig | null>(null);
  const [newAdapter, setNewAdapter] = useState<Partial<AdapterConfig>>({
    endpointPath: "",
    httpMethod: "POST",
    customHeaders: "{}",
    fieldMapping: "{}",
    responseMapping: {},
    metaConfig: "{}",
    sitemapConfig: "{}",
    timeoutSeconds: 30,
    retryCount: 3,
  });

  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const reactFlowWrapper = useRef<HTMLDivElement>(null);

  const [tempNodes, setTempNodes] = useState<Record<string, TempWorkflowNode[]>>({});
  const [deletedNodes, setDeletedNodes] = useState<Record<string, string[]>>({});
  const [updatedNodes, setUpdatedNodes] = useState<Record<string, Record<string, any>>>({});

  // Helper function to generate readable temporary node ID
  const generateTempId = (endpointPath: string = "node") => {
    const cleanName = endpointPath
      .replace(/^\//, '') // Remove leading slash
      .replace(/[^a-zA-Z0-9]/g, '_') // Replace special chars with underscore
      .substring(0, 20); // Limit length
    const timestamp = Date.now();
    const random = Math.random().toString(36).substr(2, 4);
    return `draft_${cleanName}_${timestamp}_${random}`;
  };

  // Handle save adapter
  const handleSaveAdapter = () => {
    if (!newAdapter.endpointPath) return;

    if (editingAdapter) {
      const updated = adapterConfigs.map(a =>
        a.id === editingAdapter.id
          ? { ...a, ...newAdapter, id: a.id } as AdapterConfig
          : a
      );
      onUpdateAdapterConfigs?.(updated);
    } else {
      const adapter: AdapterConfig = {
        id: `adapter_${Date.now()}`,
        productId,
        endpointPath: newAdapter.endpointPath || "",
        httpMethod: newAdapter.httpMethod as any || "POST",
        customHeaders: newAdapter.customHeaders || "{}",
        fieldMapping: newAdapter.fieldMapping || "{}",
        responseMapping: newAdapter.responseMapping || {},
        metaConfig: newAdapter.metaConfig || "{}",
        sitemapConfig: newAdapter.sitemapConfig || "{}",
        timeoutSeconds: newAdapter.timeoutSeconds || 30,
        retryCount: newAdapter.retryCount || 3,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      onUpdateAdapterConfigs?.([...adapterConfigs, adapter]);
    }

    setShowAdapterModal(false);
    setEditingAdapter(null);
    setNewAdapter({
      endpointPath: "",
      httpMethod: "POST",
      customHeaders: "{}",
      fieldMapping: "{}",
      responseMapping: {},
      metaConfig: "{}",
      sitemapConfig: "{}",
      timeoutSeconds: 30,
      retryCount: 3,
    });
  };

  const handleEditAdapter = (adapter: AdapterConfig) => {
    setEditingAdapter(adapter);
    setNewAdapter(adapter);
    setShowAdapterModal(true);
  };

  const handleDeleteAdapter = (id: string) => {
    onUpdateAdapterConfigs?.(adapterConfigs.filter(a => a.id !== id));
  };

  // Load initial workflow jika ada initialWorkflowId
  useEffect(() => {
    if (initialWorkflowId && workflows.length > 0) {
      const workflowExists = workflows.some(w => w.id === initialWorkflowId);
      if (workflowExists) {
        setSelectedWorkflowId(initialWorkflowId);
        const workflow = workflows.find(w => w.id === initialWorkflowId);
        if (workflow) {
          setWorkflowName(workflow.name);
        }
        loadWorkflowNodes(initialWorkflowId);
      }
    }
  }, [initialWorkflowId, workflows]);

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

      const edgeMap = new Map<string, Edge>();

      backendNodes.forEach((node: any) => {
        // ===== RELASI DARI nextNodeId =====
        if (node.nextNodeId) {
          const edgeId = `e-${node.id}-${node.nextNodeId}`;

          edgeMap.set(edgeId, {
            id: edgeId,
            source: node.id,
            target: node.nextNodeId,
            sourceHandle: null,
            targetHandle: null,
            animated: true,
          });
        }

        // ===== RELASI DARI previousNodeIds =====
        if (
          Array.isArray(node.previousNodeIds) &&
          node.previousNodeIds.length > 0
        ) {
          node.previousNodeIds.forEach((prevId: string) => {
            const edgeId = `e-${prevId}-${node.id}`;

            edgeMap.set(edgeId, {
              id: edgeId,
              source: prevId,
              target: node.id,
              sourceHandle: null,
              targetHandle: null,
              animated: true,
            });
          });
        }
      });

      const flowEdges = Array.from(edgeMap.values());

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
    const workflow = workflows.find(w => w.id === workflowId);
    if (workflow) {
      setWorkflowName(workflow.name);
    }
    await loadWorkflowNodes(workflowId);

    if (onChange) {
      onChange(workflowId);
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
        onChange(wf.id);
      }
    }
  };

  const handleUpdateNode = useCallback(
    async (nodeId: string, updates: Record<string, any>) => {
      console.log("📝 Updating node:", nodeId, updates); // ← DEBUG

      if (nodeId.startsWith('draft_') && selectedWorkflowId) {
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
      } else if (selectedWorkflowId) {
        // Update real node
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
      }

      // 🔥 UPDATE NODE DI STATE LANGSUNG
      setNodes((nds) =>
        nds.map((n) =>
          n.id === nodeId
            ? {
              ...n,
              data: {
                ...n.data,
                workflowNode: {
                  ...n.data.workflowNode,
                  ...updates
                },
              },
            }
            : n
        )
      );
    },
    [setNodes, selectedWorkflowId]
  );

  const updateNodeRelationship = (
    sourceId: string,
    targetId: string
  ) => {
    setNodes((nds) =>
      nds.map((node) => {
        const workflowNode = node.data.workflowNode || {};

        if (node.id === sourceId) {
          return {
            ...node,
            data: {
              ...node.data,
              workflowNode: {
                ...workflowNode,
                nextNodeId: targetId,
              },
            },
          };
        }

        if (node.id === targetId) {
          const prevIds = workflowNode.previousNodeIds || [];

          return {
            ...node,
            data: {
              ...node.data,
              workflowNode: {
                ...workflowNode,
                previousNodeIds: prevIds.includes(sourceId)
                  ? prevIds
                  : [...prevIds, sourceId],
              },
            },
          };
        }

        return node;
      })
    );
  };


  const updateRealNodeRelationship = (
    workflowId: string,
    sourceId: string,
    targetId: string
  ) => {
    setUpdatedNodes((prev) => ({
      ...prev,
      [workflowId]: {
        ...prev[workflowId],

        [sourceId]: {
          ...prev[workflowId]?.[sourceId],
          nextNodeId: targetId,
        },

        [targetId]: {
          ...prev[workflowId]?.[targetId],
          previousNodeIds: [
            ...(prev[workflowId]?.[targetId]
              ?.previousNodeIds || []),
            sourceId,
          ].filter(
            (v, i, arr) => arr.indexOf(v) === i
          ),
        },
      },
    }));
  };


  const updateTempNodeRelationship = (
    workflowId: string,
    sourceId: string,
    targetId: string
  ) => {
    setTempNodes((prev) => {
      const workflowTempNodes = [
        ...(prev[workflowId] || [])
      ];

      const updated = workflowTempNodes.map((node) => {
        if (node.id === sourceId) {
          return {
            ...node,
            nextNodeId: targetId,
          };
        }

        if (node.id === targetId) {
          return {
            ...node,
            previousNodeIds: [
              ...(node.previousNodeIds || []),
              sourceId,
            ].filter(
              (v, i, arr) => arr.indexOf(v) === i
            ),
          };
        }

        return node;
      });

      return {
        ...prev,
        [workflowId]: updated,
      };
    });
  };

  const onConnect = useCallback(
    (connection: Connection) => {
      if (!connection.source || !connection.target) return;

      setEdges((eds) =>
        addEdge(
          {
            ...connection,
            id: `e-${connection.source}-${connection.target}`,
            animated: true,
          },
          eds
        )
      );

      // ReactFlow state
      updateNodeRelationship(
        connection.source,
        connection.target
      );

      if (!selectedWorkflowId) return;

      const sourceIsTemp =
        connection.source.startsWith("draft_");

      const targetIsTemp =
        connection.target.startsWith("draft_");

      if (sourceIsTemp || targetIsTemp) {
        updateTempNodeRelationship(
          selectedWorkflowId,
          connection.source,
          connection.target
        );
      }

      updateRealNodeRelationship(
        selectedWorkflowId,
        connection.source,
        connection.target
      );
    },
    [selectedWorkflowId]
  );

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  const onDrop = useCallback(
    async (event: React.DragEvent) => {
      event.preventDefault();
      const adapterConfigId = event.dataTransfer.getData("adapterConfigId");
      if (!adapterConfigId || !reactFlowWrapper.current) return;

      const config = adapterConfigs.find((c) => c.id === adapterConfigId);
      if (!config) return;

      let currentWorkflowId = selectedWorkflowId;

      if (!currentWorkflowId) {
        const defaultWorkflowName = `Workflow ${new Date().toLocaleString()}`;
        const newWorkflow = await createWorkflow(defaultWorkflowName);
        if (newWorkflow) {
          currentWorkflowId = newWorkflow.id;
          setSelectedWorkflowId(newWorkflow.id);
          setWorkflowName(defaultWorkflowName);
          setTempNodes(prev => ({ ...prev, [newWorkflow.id]: [] }));

          if (onChange) {
            onChange(newWorkflow.id);
          }
        } else {
          return;
        }
      }

      const rect = reactFlowWrapper.current.getBoundingClientRect();
      const position = {
        x: event.clientX - rect.left - 75,
        y: event.clientY - rect.top - 25,
      };

      const tempId = generateTempId(config.endpointPath);
      const stepOrder = nodes.filter(n => !n.id.startsWith('draft_')).length +
        (tempNodes[currentWorkflowId]?.length || 0) + 1;

      const tempWorkflowNode: TempWorkflowNode = {
        id: tempId,
        workflowId: currentWorkflowId,
        adapterConfigId,
        stepOrder,
        inputMapping: {},
        previousNodeIds: [], // 
        nextNodeId: null,    // 
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      setTempNodes(prev => ({
        ...prev,
        [currentWorkflowId]: [...(prev[currentWorkflowId] || []), tempWorkflowNode]
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
    [adapterConfigs, nodes.length, selectedWorkflowId, tempNodes, setNodes, createWorkflow, onChange]
  );

  const handleNodeClick = useCallback((_event: React.MouseEvent, node: Node) => {
    setSelectedNode(node);
  }, []);

  const handleDeleteNode = useCallback(
    async (nodeId: string) => {
      if (nodeId.startsWith('draft_') && selectedWorkflowId) {
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


  const hasChanges = useCallback(() => {
    if (!selectedWorkflowId) return false;
    return (tempNodes[selectedWorkflowId]?.length > 0 ||
      deletedNodes[selectedWorkflowId]?.length > 0 ||
      Object.keys(updatedNodes[selectedWorkflowId] || {}).length > 0);
  }, [selectedWorkflowId, tempNodes, deletedNodes, updatedNodes]);

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
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-background' : ''}`}>
      <div className={`flex gap-4 relative ${isFullscreen ? 'h-screen' : 'h-[600px]'}`}>
        {/* Fullscreen Toggle Button */}
        {onToggleFullscreen && (
          <button
            onClick={onToggleFullscreen}
            className="absolute top-2 right-2 z-20 bg-background border rounded-md p-1.5 hover:bg-muted transition-colors"
          >
            {isFullscreen ? (
              <Minimize2 className="h-4 w-4" />
            ) : (
              <Maximize2 className="h-4 w-4" />
            )}
          </button>
        )}

        {/* Left Sidebar */}
        {isSidebarOpen && (
          <div className="w-64 border-r pr-4 flex flex-col gap-4 transition-all duration-300 p-2">
            <div>
              <div className="flex items-center justify-between mb-2">
                <h3 className="font-semibold text-sm">WORKFLOWS</h3>
              </div>

              <WorkflowList
                workflows={workflows}
                selectedId={selectedWorkflowId}
                onSelect={handleSelectWorkflow}
                onDelete={deleteWorkflow}
                loading={workflowLoading}
              />
            </div>

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
                        className="group bg-muted rounded px-2 py-1.5 text-xs cursor-grab hover:bg-primary/10 transition-colors flex items-center justify-between"
                        draggable
                        onDragStart={(e) => {
                          e.dataTransfer.setData("adapterConfigId", config.id);
                          e.dataTransfer.effectAllowed = "move";
                        }}
                      >
                        <span className="truncate">
                          <span className="font-mono font-bold text-primary">{config.httpMethod}</span>{" "}
                          <span className="truncate">{config.endpointPath}</span>
                        </span>
                        <div className="hidden group-hover:flex gap-1">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleEditAdapter(config);
                            }}
                            className="text-slate-500 hover:text-slate-700"
                          >
                            <Pencil className="h-3 w-3" />
                          </button>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDeleteAdapter(config.id);
                            }}
                            className="text-red-500 hover:text-red-700"
                          >
                            <Trash2 className="h-3 w-3" />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>

                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    className="mt-2 w-full"
                    onClick={() => {
                      setEditingAdapter(null);
                      setNewAdapter({
                        endpointPath: "",
                        httpMethod: "POST",
                        customHeaders: "{}",
                        fieldMapping: "{}",
                        responseMapping: {},
                        metaConfig: "{}",
                        sitemapConfig: "{}",
                        timeoutSeconds: 30,
                        retryCount: 3,
                      });
                      setShowAdapterModal(true);
                    }}
                  >
                    <Plus className="h-3 w-3 mr-1" />
                    Tambah Adapter
                  </Button>
                </>
              )}
            </div>
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

            <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
          </ReactFlow>
        </div>

        {/* Right Panel */}
        {selectedNode && (
          <div className="w-150 border-l pl-4 overflow-y-auto">
            <NodePanel
              node={selectedNode}
              allNodes={nodes}
              adapterConfigs={adapterConfigs}
              product={product}
              onClose={() => setSelectedNode(null)}
              onDelete={handleDeleteNode}
              onUpdate={handleUpdateNode}
            />
          </div>
        )}

        {/* Modal Add/Edit Adapter */}
        <Dialog open={showAdapterModal} onOpenChange={setShowAdapterModal}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingAdapter ? "Edit Adapter" : "Tambah Adapter Baru"}</DialogTitle>
              <DialogDescription>
                Definisikan endpoint API yang akan digunakan di workflow
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>HTTP Method</Label>
                <Select
                  value={newAdapter.httpMethod as string}
                  onValueChange={(v) => setNewAdapter({ ...newAdapter, httpMethod: v as any })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="GET">GET</SelectItem>
                    <SelectItem value="POST">POST</SelectItem>
                    <SelectItem value="PUT">PUT</SelectItem>
                    <SelectItem value="PATCH">PATCH</SelectItem>
                    <SelectItem value="DELETE">DELETE</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label>Endpoint Path</Label>
                <Input
                  placeholder="/posts"
                  value={newAdapter.endpointPath}
                  onChange={(e) => setNewAdapter({ ...newAdapter, endpointPath: e.target.value })}
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Contoh: /posts, /media, /categories
                </p>
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="outline" onClick={() => setShowAdapterModal(false)}>
                  Batal
                </Button>
                <Button onClick={handleSaveAdapter}>
                  Simpan
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}