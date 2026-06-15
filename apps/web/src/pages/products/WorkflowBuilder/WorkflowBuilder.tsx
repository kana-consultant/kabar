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
import type { Connection, Node, Edge } from "reactflow";
import "reactflow/dist/style.css";
import { WorkflowList } from "./WorkflowList";
import { NodePanel } from "./NodePanel";
import { WorkflowNode as WorkflowNodeComponent } from "./WorkflowNode";
import type { Product, WorkflowDefinition, WorkflowNode as WorkflowNodeType, AdapterConfig } from "@/types/product";

interface WorkflowBuilderProps {
  productId: string;
  product: Product;
  selectedWorkflowId?: string;
  onWorkflowSelect: (workflowId: string) => void;
  onWorkflowDelete: (workflowId: string) => void;
  onWorkflowCreate: (name: string) => void;
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
  adapterConfig: AdapterConfig;
}

type FlowNode = Node<FlowNodeData>;

const nodeTypes = {
  workflowNode: WorkflowNodeComponent,
};

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
  const currentWorkflow = product.workflows?.find(w => w.id === selectedWorkflowId);

  // Sync selected workflow dengan external
  useEffect(() => {
    if (externalSelectedWorkflowId) {
      setSelectedWorkflowId(externalSelectedWorkflowId);
    }
  }, [externalSelectedWorkflowId]);

  // Load workflow nodes dari product
  useEffect(() => {
    if (!selectedWorkflowId || !currentWorkflow) {
      setNodes([]);
      setEdges([]);
      return;
    }

    const workflowNodes = currentWorkflow.nodes || [];

    // Transform WorkflowNode ke FlowNode
    const flowNodes: FlowNode[] = workflowNodes.map((node, index) => {
      const adapterConfig = node.adapter_config!;

      return {
        id: node.id as string,
        type: "workflowNode",
        position: { x: 250, y: index * 150 + 50 },
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
    const flowEdges: Edge[] = [];
    workflowNodes.forEach((node) => {
      if (node.next_node_id) {
        flowEdges.push({
          id: `e-${node.id}-${node.next_node_id}`,
          source: node.id as string,
          target: node.next_node_id,
          animated: true,
        });
      }
    });

    setEdges(flowEdges);
  }, [selectedWorkflowId, currentWorkflow, setNodes, setEdges]);

  const handleSelectWorkflow = (workflowId: string) => {
    setSelectedWorkflowId(workflowId);
    onWorkflowSelect(workflowId);
    onChange?.(workflowId);
  };

  const handleCreateWorkflow = () => {
    const name = `Workflow ${new Date().toLocaleDateString()}`;
    onWorkflowCreate(name);
  };

  // Update node data
  const handleUpdateNode = useCallback(async (nodeId: string, updates: Partial<WorkflowNodeType> & { adapter_config?: Partial<AdapterConfig> }) => {
    console.log(`🔄 handleUpdateNode called:`, { nodeId, updates });

    const nodeToUpdate = nodes.find(n => n.id === nodeId);
    if (!nodeToUpdate) {
      console.warn(`Node ${nodeId} not found`);
      return;
    }

    // Pisahkan adapter_config dari updates
    const { adapter_config, ...nodeUpdates } = updates;

    console.log("📦 Updates breakdown:", {
      nodeUpdates,
      adapter_config
    });

    // Update local state immediately
    setNodes((nds) =>
      nds.map((n) => {
        if (n.id === nodeId) {
          let updatedNode = { ...n.data.workflowNode, ...nodeUpdates };
          let updatedAdapterConfig = n.data.adapterConfig;
          let endpointPath = n.data.adapterConfig?.endpoint_path || '/unknown';

          // Update adapter config if provided
          if (adapter_config) {
            updatedAdapterConfig = {
              ...n.data.adapterConfig,
              ...adapter_config,
              updated_at: new Date().toISOString()
            };
            endpointPath = adapter_config.endpoint_path || endpointPath;
            console.log("✅ Adapter config updated:", updatedAdapterConfig);
          }

          const newLabel = `Step ${updatedNode.step_order}: ${endpointPath}`;

          console.log("✅ Node updated:", {
            step_order: updatedNode.step_order,
            endpoint: endpointPath,
            label: newLabel,
            input_mapping: updatedNode.input_mapping
          });

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

    // Call parent handler - kirim kedua updates
    if (adapter_config) {
      // Jika ada adapter_config, update melalui onNodeUpdate dengan kedua data
      onNodeUpdate(nodeId, {
        ...nodeUpdates,
        adapter_config: adapter_config
      });
    } else {
      // Jika hanya node updates
      onNodeUpdate(nodeId, nodeUpdates);
    }

    console.log(`✅ Update complete for node ${nodeId}`);
  }, [nodes, onNodeUpdate]);

  // Handle connection between nodes
  const onConnect = useCallback(
    async (connection: Connection) => {
      if (!connection.source || !connection.target || !selectedWorkflowId) return;

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

      // Update source node's next_node_id
      await handleUpdateNode(connection.source, { next_node_id: connection.target });

      // Update target node's previous_node_ids
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

      // Update connected nodes - remove references from previous_node_ids
      for (const edge of connectedEdges) {
        if (edge.target && edge.target !== nodeId) {
          const targetNode = nodes.find(n => n.id === edge.target);
          if (targetNode && targetNode.data.workflowNode.previous_node_ids) {
            const updatedPrevIds = targetNode.data.workflowNode.previous_node_ids.filter(id => id !== nodeId);
            await handleUpdateNode(edge.target, { previous_node_ids: updatedPrevIds });
          }
        }
      }

      // Call parent handler
      onNodeDelete(nodeId);
    },
    [selectedWorkflowId, setNodes, setEdges, nodes, edges, selectedNode, handleUpdateNode, onNodeDelete]
  );

  // Handle tambah node - langsung dari product
  const handleAddNode = async () => {
    if (!selectedWorkflowId) {
      setError("Pilih atau buat workflow terlebih dahulu");
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

      // Default adapter config
      const newAdapterConfig: AdapterConfig = {
        id: adapterConfigId,
        product_id: productId,
        endpoint_path: "/api/default",
        http_method: "GET",
        timeout_seconds: 30,
        retry_count: 3,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };

      // Default workflow node
      const newNode: WorkflowNodeType = {
        id: nodeId,
        workflow_id: selectedWorkflowId,
        adapter_config_id: adapterConfigId,
        step_order: stepOrder,
        input_mapping: {},
        next_node_id: null,
        previous_node_ids: [],
        created_at: new Date().toISOString(),
        adapter_config: newAdapterConfig,
      };

      // Call parent handler to add node to product
      onNodeAdd(newNode);

      // Tambahkan node ke React Flow state
      const flowNode: FlowNode = {
        id: nodeId,
        type: "workflowNode",
        position: {
          x: 250,
          y: (nodes.length || 0) * 120 + 50,
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

  // Get workflows from product
  const workflows = product.workflows || [];

  return (
    <div className={`${isFullscreen ? 'fixed inset-0 z-50 bg-background' : ''}`}>
      <div className={`flex gap-4 relative ${isFullscreen ? 'h-screen' : 'h-[600px]'}`}>
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
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleCreateWorkflow}
                  className="h-6 px-2 text-xs"
                >
                  + Baru
                </Button>
              </div>
              <WorkflowList
                workflows={product.workflows as WorkflowDefinition[]}
                selectedId={selectedWorkflowId}
                onSelect={handleSelectWorkflow}
                onDelete={onWorkflowDelete}
                loading={false}
              />
            </div>

            {/* Tombol Tambah Node */}
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
                Pilih atau buat workflow terlebih dahulu
              </p>
            )}

            {/* Workflow Nodes List */}
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
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeClick={handleNodeClick}
            nodeTypes={nodeTypes}
            fitView
            deleteKeyCode="Delete"
            onNodesDelete={(deleted) => {
              deleted.forEach((n) => handleDeleteNode(n.id));
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