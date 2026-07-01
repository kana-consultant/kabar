// src/pages/products/WorkflowBuilder/WorkflowBuilder.tsx
import { useState, useCallback, useRef, useEffect } from "react";
import { ChevronLeft, ChevronRight as ChevronRightIcon, Pencil, Trash2, Maximize2, Minimize2 } from "lucide-react";
import { Button } from "@kana-consultant/ui-kit";
import ReactFlow, {
  Controls,
  Background,
  addEdge,
  useNodesState,
  useEdgesState,
  BackgroundVariant,
} from "reactflow";
import type { Connection, Node, Edge, NodeChange } from "reactflow";
import "reactflow/dist/style.css";
import { WorkflowList } from "./WorkflowList";
import { NodePanel } from "./NodePanel";
import { WorkflowNode as WorkflowNodeComponent } from "./WorkflowNode";
import type { Product, WorkflowDefinition, WorkflowNode as WorkflowNodeType, AdapterConfig, AdapterConfigNode } from "@/types/product";

interface WorkflowBuilderProps {
  productId: string;
  product: Product;
  selectedWorkflowId?: string;
  onWorkflowSelect: (workflowId: string) => void;
  onWorkflowDelete: (workflowId: string) => void;
  onWorkflowCreate?: (name: string) => void;
  onNodeAdd: (node: WorkflowNodeType) => void;
  onNodeUpdate: (nodeId: string, updates: Partial<WorkflowNodeType>) => void;
  onNodeDelete: (nodeId: string) => void;
  onChange?: (workflowId: string) => void;
  isFullscreen?: boolean;
  onToggleFullscreen?: () => void;
}

interface FlowNodeData {
  label: string;
  workflowNode: WorkflowNodeType;
  adapterConfig: AdapterConfigNode;
}

type FlowNode = Node<FlowNodeData>;

const nodeTypes = {
  workflowNode: WorkflowNodeComponent,
};

// Key untuk menyimpan posisi node di localStorage
const getPositionStorageKey = (workflowId: string) => `workflow-node-positions-${workflowId}`;

export function WorkflowBuilder({
  productId,
  product,
  selectedWorkflowId: externalSelectedWorkflowId,
  onWorkflowSelect,
  onWorkflowDelete,
  onWorkflowCreate,
  onNodeAdd,
  onNodeUpdate,
  onNodeDelete,
  onChange,
  isFullscreen = false,
  onToggleFullscreen
}: WorkflowBuilderProps) {
  const [selectedWorkflowId, setSelectedWorkflowId] = useState<string | null>(externalSelectedWorkflowId || null);
  const [selectedNode, setSelectedNode] = useState<FlowNode | null>(null);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [nodes, setNodes, onNodesChange] = useNodesState<FlowNodeData>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const reactFlowWrapper = useRef<HTMLDivElement>(null);

  // Get current workflow from product
  const workflows = product?.workflows || [];
  const currentWorkflow = workflows.find(w => w.id === selectedWorkflowId);

  // Simpan posisi node ke localStorage
  const saveNodePositions = useCallback((workflowId: string, currentNodes: FlowNode[]) => {
    const positions: Record<string, { x: number; y: number }> = {};
    currentNodes.forEach(node => {
      positions[node.id] = { x: node.position.x, y: node.position.y };
    });
    localStorage.setItem(getPositionStorageKey(workflowId), JSON.stringify(positions));
  }, []);

  // Ambil posisi node dari localStorage
  const getNodePositions = useCallback((workflowId: string): Record<string, { x: number; y: number }> => {
    try {
      const saved = localStorage.getItem(getPositionStorageKey(workflowId));
      return saved ? JSON.parse(saved) : {};
    } catch {
      return {};
    }
  }, []);

  // Override onNodesChange untuk menyimpan posisi setelah drag
  const handleNodesChange = useCallback((changes: NodeChange[]) => {
    onNodesChange(changes);

    const hasPositionChange = changes.some(
      change => change.type === 'position' && change.dragging === false
    );

    if (hasPositionChange && selectedWorkflowId) {
      setTimeout(() => {
        setNodes(currentNodes => {
          saveNodePositions(selectedWorkflowId, currentNodes);
          return currentNodes;
        });
      }, 0);
    }
  }, [onNodesChange, selectedWorkflowId, saveNodePositions]);



  // Auto-select first workflow when workflows change

 
  useEffect(() => {
    console.log(workflows, "workflows======================");
    if (workflows.length > 0) {
      const firstWorkflowId = workflows[0].id;
      console.log(firstWorkflowId, "firstWorkflowId");
      handleSelectWorkflow(firstWorkflowId);
      // setSelectedWorkflowId(firstWorkflowId);
      // onWorkflowSelect(firstWorkflowId);
      // console.log(selectedWorkflowId,firstWorkflowId, "selectedWorkflowId");
      // onWorkflowSelect(firstWorkflowId);
      // onChange?.(firstWorkflowId);
    } else {
      console.log("No workflows available");
      if (!workflows || workflows.length === 0) {
        const defaultWorkflowName = "Default Workflow";
        onWorkflowCreate?.(defaultWorkflowName);
      }
    }
  }, [workflows]);

  // Sync selected workflow dengan external
  useEffect(() => {
    if (externalSelectedWorkflowId && externalSelectedWorkflowId !== selectedWorkflowId) {
      setSelectedWorkflowId(externalSelectedWorkflowId);
      onWorkflowSelect(externalSelectedWorkflowId);
      onChange?.(externalSelectedWorkflowId);
    }
  }, [externalSelectedWorkflowId, selectedWorkflowId, onWorkflowSelect, onChange]);

  // Load workflow nodes dari product
  useEffect(() => {
    if (!selectedWorkflowId || !currentWorkflow) {
      setNodes([]);
      setEdges([]);
      return;
    }

    const workflowNodes = currentWorkflow.nodes || [];
    const savedPositions = getNodePositions(selectedWorkflowId);

    // Sort nodes by step_order
    const sortedNodes = [...workflowNodes].sort((a, b) => (a.step_order || 0) - (b.step_order || 0));

    // Transform WorkflowNode ke FlowNode
    const flowNodes: FlowNode[] = sortedNodes.map((node, index) => {
      const adapterConfig = node.adapter_config!;

      const defaultPosition = {
        x: index * 350 + 100,
        y: 100,
      };

      const position = savedPositions[node.id as string] || defaultPosition;

      return {
        id: node.id as string,
        type: "workflowNode",
        position: position,
        data: {
          label: `Step ${node.step_order}: ${adapterConfig.endpoint_path}`,
          workflowNode: node,
          adapterConfig: adapterConfig,
        },
        draggable: true,
        selectable: true,
      };
    });

    setNodes(flowNodes);

    // Create edges from node relationships
    // ONLY use next_node_ids and previous_node_ids (MULTIPLE)
    const flowEdges: Edge[] = [];

    sortedNodes.forEach((node) => {
      // Create edges from next_node_ids (MULTIPLE)
      if (node.next_node_ids && node.next_node_ids.length > 0) {
        node.next_node_ids.forEach(nextNodeId => {
          const edgeId = `e-${node.id}-${nextNodeId}`;
          if (!flowEdges.some(e => e.id === edgeId)) {
            flowEdges.push({
              id: edgeId,
              source: node.id as string,
              target: nextNodeId,
              animated: true,
            });
          }
        });
      }

      // Create edges from previous_node_ids (MULTIPLE)
      if (node.previous_node_ids && node.previous_node_ids.length > 0) {
        node.previous_node_ids.forEach(prevNodeId => {
          const edgeId = `e-${prevNodeId}-${node.id}`;
          if (!flowEdges.some(e => e.id === edgeId)) {
            flowEdges.push({
              id: edgeId,
              source: prevNodeId,
              target: node.id as string,
              animated: true,
            });
          }
        });
      }
    });

    setEdges(flowEdges);
  }, [selectedWorkflowId, currentWorkflow, setNodes, setEdges, getNodePositions]);

  const handleSelectWorkflow = (workflowId: string) => {
    console.log(workflowId, "handleSelectWorkflow");
    setSelectedWorkflowId(workflowId);
    onWorkflowSelect(workflowId);
    onChange?.(workflowId);
  };

  console.log(selectedWorkflowId, "selectedWorkflowId");

  // Update node data
  const handleUpdateNode = useCallback(async (nodeId: string, updates: Partial<WorkflowNodeType> & { adapter_config?: Partial<AdapterConfig> }) => {
    const nodeToUpdate = nodes.find(n => n.id === nodeId);
    if (!nodeToUpdate) {
      console.warn(`Node ${nodeId} not found`);
      return;
    }

    const { adapter_config, ...nodeUpdates } = updates;

    setNodes((nds) =>
      nds.map((n) => {
        if (n.id === nodeId) {
          let updatedNode = { ...n.data.workflowNode, ...nodeUpdates };
          let updatedAdapterConfig = n.data.adapterConfig;
          let endpointPath = n.data.adapterConfig?.endpoint_path || '/unknown';

          if (adapter_config) {
            updatedAdapterConfig = {
              ...n.data.adapterConfig,
              ...adapter_config,
            };
            endpointPath = adapter_config.endpoint_path || endpointPath;
          }

          const newLabel = `Step ${updatedNode.step_order}: ${endpointPath}`;

          return {
            ...n,
            data: {
              ...n.data,
              workflowNode: updatedNode,
              adapterConfig: updatedAdapterConfig,
              label: newLabel,
            },
          };
        }
        return n;
      })
    );

    if (adapter_config) {
      onNodeUpdate(nodeId, {
        ...nodeUpdates,
        adapter_config: adapter_config
      });
    } else {
      onNodeUpdate(nodeId, nodeUpdates);
    }
  }, [nodes, onNodeUpdate]);

  // Handle connection between nodes
  const onConnect = useCallback(
    async (connection: Connection) => {
      if (!connection.source || !connection.target || !selectedWorkflowId) return;

      // Cek apakah edge sudah ada
      const edgeExists = edges.some(
        e => e.source === connection.source && e.target === connection.target
      );

      if (edgeExists) {
        setError("Connection already exists between these nodes");
        return;
      }

      const newEdge = {
        id: `e-${connection.source}-${connection.target}`,
        source: connection.source,
        target: connection.target,
        animated: true,
      };

      setEdges((eds) => addEdge(newEdge, eds));

      // Update source node's next_node_ids (MULTIPLE)
      const sourceNode = nodes.find(n => n.id === connection.source);
      if (sourceNode) {
        const currentNextIds = sourceNode.data.workflowNode.next_node_ids || [];

        if (!currentNextIds.includes(connection.target)) {
          await handleUpdateNode(connection.source, {
            next_node_ids: [...currentNextIds, connection.target]
          });
        }
      }

      // Update target node's previous_node_ids (MULTIPLE)
      const targetNode = nodes.find(n => n.id === connection.target);
      if (targetNode) {
        const currentPrevIds = targetNode.data.workflowNode.previous_node_ids || [];
        if (!currentPrevIds.includes(connection.source)) {
          await handleUpdateNode(connection.target, {
            previous_node_ids: [...currentPrevIds, connection.source]
          });
        }
      }
    },
    [handleUpdateNode, nodes, selectedWorkflowId, setEdges, edges]
  );

  const handleNodeClick = useCallback((_event: React.MouseEvent, node: FlowNode) => {
    setSelectedNode(node);
  }, []);

  const handleDeleteNode = useCallback(
    async (nodeId: string) => {
      if (!selectedWorkflowId) return;

      const connectedEdges = edges.filter(e => e.source === nodeId || e.target === nodeId);

      // Remove from UI immediately
      setNodes((nds) => nds.filter((n) => n.id !== nodeId));
      setEdges((eds) => eds.filter((e) => e.source !== nodeId && e.target !== nodeId));

      if (selectedNode?.id === nodeId) {
        setSelectedNode(null);
      }

      // Update connected nodes
      for (const edge of connectedEdges) {
        // If this node was the source, remove from target's previous_node_ids
        if (edge.source === nodeId && edge.target) {
          const targetNode = nodes.find(n => n.id === edge.target);
          if (targetNode?.data.workflowNode.previous_node_ids) {
            const updatedPrevIds = targetNode.data.workflowNode.previous_node_ids.filter(id => id !== nodeId);
            await handleUpdateNode(edge.target, { previous_node_ids: updatedPrevIds });
          }
        }

        // If this node was the target, remove from source's next_node_ids
        if (edge.target === nodeId && edge.source) {
          const sourceNode = nodes.find(n => n.id === edge.source);
          if (sourceNode?.data.workflowNode.next_node_ids) {
            const updatedNextIds = sourceNode.data.workflowNode.next_node_ids.filter(id => id !== nodeId);
            await handleUpdateNode(edge.source, { next_node_ids: updatedNextIds });
          }
        }
      }

      // Cleanup posisi yang disimpan di localStorage
      if (selectedWorkflowId) {
        const savedPositions = getNodePositions(selectedWorkflowId);
        delete savedPositions[nodeId];
        localStorage.setItem(getPositionStorageKey(selectedWorkflowId), JSON.stringify(savedPositions));
      }

      onNodeDelete(nodeId);
    },
    [selectedWorkflowId, setNodes, setEdges, nodes, edges, selectedNode, handleUpdateNode, onNodeDelete, getNodePositions]
  );

  // Handle tambah node
  const handleAddNode = async () => {
    if (!selectedWorkflowId) {
      setError("Pilih workflow terlebih dahulu");
      return;
    }

    if (!currentWorkflow) {
      setError("Workflow tidak ditemukan");
      return;
    }

    try {
      setIsSaving(true);
      setError(null);

      const stepOrder = (currentWorkflow.nodes?.length || 0) + 1;
      const nodeId = `temp-node-${Date.now()}-${Math.random()}`;
      const adapterConfigId = `temp-adapter-${Date.now()}-${Math.random()}`;

      const newAdapterConfig: AdapterConfigNode = {
        id: adapterConfigId,
        product_id: productId,
        endpoint_path: "/api/default",
        http_method: "GET",
      };

      const newNode: WorkflowNodeType = {
        id: nodeId,
        workflow_id: selectedWorkflowId,
        adapter_config_id: adapterConfigId,
        step_order: stepOrder,
        input_mapping: {},
        next_node_ids: [],
        previous_node_ids: [],
        created_at: new Date().toISOString(),
        adapter_config: newAdapterConfig,
      };

      onNodeAdd(newNode);

      const horizontalSpacing = 350;
      const startX = 100;
      const defaultY = 100;

      const flowNode: FlowNode = {
        id: nodeId,
        type: "workflowNode",
        position: {
          x: startX + (nodes.length * horizontalSpacing),
          y: defaultY,
        },
        data: {
          label: `Step ${stepOrder}: ${newAdapterConfig.endpoint_path}`,
          workflowNode: newNode,
          adapterConfig: newAdapterConfig,
        },
      };

      setNodes(nds => [...nds, flowNode]);

    } catch (err) {
      console.error(err);
      setError("Gagal membuat node");
    } finally {
      setIsSaving(false);
    }
  };

  // Handle delete edge
  const handleEdgeDelete = useCallback(
    async (edgeId: string) => {
      const edgeToDelete = edges.find(e => e.id === edgeId);
      if (!edgeToDelete) return;

      // Update source node - hapus dari next_node_ids
      const sourceNode = nodes.find(n => n.id === edgeToDelete.source);
      if (sourceNode?.data.workflowNode.next_node_ids) {
        const updatedNextIds = sourceNode.data.workflowNode.next_node_ids.filter(id => id !== edgeToDelete.target);
        await handleUpdateNode(edgeToDelete.source, { next_node_ids: updatedNextIds });
      }

      // Update target node - hapus dari previous_node_ids
      const targetNode = nodes.find(n => n.id === edgeToDelete.target);
      if (targetNode?.data.workflowNode.previous_node_ids) {
        const updatedPrevIds = targetNode.data.workflowNode.previous_node_ids.filter(id => id !== edgeToDelete.source);
        await handleUpdateNode(edgeToDelete.target, { previous_node_ids: updatedPrevIds });
      }

      setEdges(eds => eds.filter(e => e.id !== edgeId));
    },
    [edges, nodes, handleUpdateNode, setEdges]
  );

  return (
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-background' : ''}`}>
      <div className={`mt-2 flex gap-4 relative ${isFullscreen ? 'h-screen' : 'h-[600px]'}`}>
        {/* Error Alert */}
        {error && (
          <div className="absolute top-2 left-1/2 transform -translate-x-1/2 z-50 bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded shadow-lg">
            <div className="flex items-center gap-2">
              <span className="text-sm">{error}</span>
              <button
                onClick={() => setError(null)}
                className="ml-2 text-red-700 hover:text-red-900"
              >
                ×
              </button>
            </div>
          </div>
        )}

        {/* Loading Overlay */}
        {isSaving && (
          <div className="absolute inset-0 bg-background/50 z-50 flex items-center justify-center">
            <div className="bg-background rounded-lg shadow-lg p-4 flex items-center gap-3">
              <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-primary"></div>
              <span className="text-sm">Menyimpan...</span>
            </div>
          </div>
        )}

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
                workflows={product?.workflows as WorkflowDefinition[]}
                selectedId={selectedWorkflowId}
                onSelect={handleSelectWorkflow}
                onDelete={onWorkflowDelete}
                loading={false}
              />
            </div>

            <Button
              onClick={handleAddNode}
              className="w-full"
              size="sm"
              disabled={!selectedWorkflowId}
            >
              + Tambah Node
            </Button>
            {!selectedWorkflowId && (
              <p className="text-xs text-muted-foreground text-center -mt-2">
                Pilih workflow terlebih dahulu
              </p>
            )}

            {selectedWorkflowId && currentWorkflow && currentWorkflow.nodes && currentWorkflow.nodes.length > 0 && (
              <div className="pt-3 border-t">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold text-sm">WORKFLOW NODES</h3>
                  <span className="text-xs text-muted-foreground">
                    {currentWorkflow.nodes.length} nodes
                  </span>
                </div>
                <div className="space-y-1 max-h-64 overflow-y-auto">
                  {nodes.map((node) => {
                    const workflowNode = node.data.workflowNode;
                    const endpointPath = node.data.adapterConfig?.endpoint_path || '/unknown';

                    return (
                      <div
                        key={node.id}
                        className={`group rounded px-2 py-1.5 text-xs cursor-pointer hover:bg-primary/10 transition-colors flex items-center justify-between ${selectedNode?.id === node.id ? 'bg-primary/20 border border-primary/50' : 'bg-muted'
                          }`}
                        onClick={() => handleNodeClick({} as React.MouseEvent, node)}
                      >
                        <span className="truncate flex items-center gap-1 flex-1">
                          <span className="font-mono font-bold text-primary min-w-[20px]">
                            {workflowNode.step_order}
                          </span>
                          <span className="truncate">{endpointPath}</span>
                        </span>
                        <div className="hidden group-hover:flex gap-1">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              setSelectedNode(node);
                            }}
                            className="text-slate-500 hover:text-slate-700"
                          >
                            <Pencil className="h-3 w-3" />
                          </button>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDeleteNode(node.id);
                            }}
                            className="text-red-500 hover:text-red-700"
                          >
                            <Trash2 className="h-3 w-3" />
                          </button>
                        </div>
                      </div>
                    );
                  })}
                </div>
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
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={handleNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeClick={handleNodeClick}
            onEdgeClick={(_event, edge) => {
              if (window.confirm('Hapus koneksi ini?')) {
                handleEdgeDelete(edge.id);
              }
            }}
            nodeTypes={nodeTypes}
            fitView
            deleteKeyCode="Delete"
            onNodesDelete={(deleted) => {
              deleted.forEach((n) => handleDeleteNode(n.id));
            }}
            onEdgesDelete={(deleted) => {
              deleted.forEach((e) => handleEdgeDelete(e.id));
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
              product={product}
              onClose={() => setSelectedNode(null)}
              onDelete={handleDeleteNode}
              onUpdate={handleUpdateNode}
            />
          </div>
        )}
      </div>
    </div>
  );
}