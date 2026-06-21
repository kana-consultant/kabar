import type { Product, AdapterConfig, WorkflowDefinition, WorkflowNode } from "@/types/product";

export function useProductFormActions(
    setProduct: (updater: any) => void
) {
    const updateProductInfo = (updates: Partial<Product>) => {
        setProduct((prev: Partial<Product>) => ({ ...prev, ...updates }));
    };

    const updateAdapterConfig = (updates: Partial<AdapterConfig>) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapter_config: { ...prev.adapter_config, ...updates } as AdapterConfig,
        }));
    };

    const updateFieldMapping = (value: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapter_config: {
                ...prev.adapter_config,
                fieldMapping: value,
            } as AdapterConfig,
        }));
    };

    const updateMetaConfig = (metaConfig: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapter_config: {
                ...prev.adapter_config,
                meta_config: metaConfig,
            } as AdapterConfig,
        }));
    };

    const updateSitemapConfig = (sitemapConfig: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapter_config: {
                ...prev.adapter_config,
                sitemap_config: sitemapConfig,
            } as AdapterConfig,
        }));
    };

    // ========== WORKFLOW ACTIONS ==========

    const updateWorkflowId = (workflowId: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflow_id: workflowId,
        }));
    };

    // ========== NEW WORKFLOW MANAGEMENT ACTIONS ==========

    // Add new workflow
    const addWorkflow = (workflow: WorkflowDefinition) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: [...(prev.workflows || []), workflow],
        }));
    };

    // Update existing workflow
    const updateWorkflow = (workflowId: string, updates: Partial<WorkflowDefinition>) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: prev.workflows?.map(workflow =>
                workflow.id === workflowId
                    ? { ...workflow, ...updates, updated_at: new Date().toISOString() }
                    : workflow
            ) || [],
        }));
    };

    // Delete workflow
    const deleteWorkflow = (workflowId: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: prev.workflows?.filter(workflow => workflow.id !== workflowId) || [],
            // If deleted workflow was active, clear active_workflow_id
            active_workflow_id: prev.workflow_id === workflowId ? undefined : prev.workflow_id,
        }));
    };

    // Set active workflow
    const setActiveWorkflow = (workflowId: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            active_workflow_id: workflowId,
        }));
    };

    // ========== NODE MANAGEMENT ACTIONS ==========

    // Add node to workflow
    const addNodeToWorkflow = (node: WorkflowNode) => {
        console.log("processss");
        console.log(node);

        setProduct((prev: any) => {
            const workflows = prev.workflows || [];

            if (workflows.length === 0) {
                return {
                    ...prev,
                    workflows: [
                        {
                            id: node.id,
                            name: "Default Workflow",
                            product_id: prev.id || "",
                            created_at: new Date().toISOString(),
                            updated_at: new Date().toISOString(),
                            nodes: [node],
                        }
                    ]
                };
            }

            return {
                ...prev,
                workflows: workflows.map((workflow: any) => ({
                    ...workflow,
                    nodes: [...(workflow.nodes || []), node],
                }))
            };
        });
    };

    // Update node in workflow
    // Update node in workflow
    // Update node in workflow
    const updateNodeInWorkflow = (
        workflowId: string,
        nodeId: string,
        updates: Partial<WorkflowNode> & { adapter_config?: Partial<AdapterConfig> }
    ) => {
        console.log(`🔄 updateNodeInWorkflow:`, { workflowId, nodeId, updates });

        setProduct((prev: Partial<Product>) => {
            const updatedWorkflows = prev.workflows?.map(workflow =>
                workflow.id === workflowId
                    ? {
                        ...workflow,
                        nodes: workflow.nodes?.map(node => {
                            if (node.id === nodeId) {
                                // Ambil adapter_config dari updates
                                const { adapter_config, ...nodeUpdates } = updates;

                                // Update node
                                let updatedNode = {
                                    ...node,
                                    ...nodeUpdates,
                                    updated_at: new Date().toISOString()
                                };

                                // Update adapter_config jika ada
                                if (adapter_config) {
                                    console.log("📡 Updating adapter_config:", adapter_config);

                                    // Pastikan custom_headers tetap sebagai string
                                    let customHeadersStr = adapter_config.custom_headers;
                                    if (customHeadersStr && typeof customHeadersStr === 'object') {
                                        customHeadersStr = JSON.stringify(customHeadersStr);
                                    }

                                    updatedNode = {
                                        ...updatedNode,
                                        adapter_config: {
                                            ...node.adapter_config,
                                            ...adapter_config,
                                            custom_headers: customHeadersStr || node.adapter_config?.custom_headers,
                                            updated_at: new Date().toISOString()
                                        } as AdapterConfig
                                    };
                                }

                                console.log(`✅ Node ${nodeId} updated:`, {
                                    nodeUpdates,
                                    hasadapter_config: !!adapter_config,
                                    step_order: updatedNode.step_order,
                                });

                                return updatedNode;
                            }
                            return node;
                        }) || [],
                        updated_at: new Date().toISOString(),
                    }
                    : workflow
            ) || [];

            return {
                ...prev,
                workflows: updatedWorkflows,
            };
        });
    };

    // Delete node from workflow
    const deleteNodeFromWorkflow = (workflowId: string, nodeId: string) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: prev.workflows?.map(workflow =>
                workflow.id === workflowId
                    ? {
                        ...workflow,
                        nodes: workflow.nodes?.filter(node => node.id !== nodeId) || [],
                        updated_at: new Date().toISOString(),
                    }
                    : workflow
            ) || [],
        }));
    };

    // Update node connections (next_node_id, previous_node_ids)
    const updateNodeConnections = (
        workflowId: string,
        nodeId: string,
        nextNodeId: string[] | null,
        previousNodeIds?: string[]
    ) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: prev.workflows?.map(workflow =>
                workflow.id === workflowId
                    ? {
                        ...workflow,
                        nodes: workflow.nodes?.map(node => {
                            if (node.id === nodeId) {
                                const updates: Partial<WorkflowNode> = { next_node_ids : nextNodeId as string[]};
                                if (previousNodeIds !== undefined) {
                                    updates.previous_node_ids = previousNodeIds;
                                }
                                return { ...node, ...updates, updated_at: new Date().toISOString() };
                            }
                            return node;
                        }) || [],
                        updated_at: new Date().toISOString(),
                    }
                    : workflow
            ) || [],
        }));
    };

    // Reorder nodes in workflow (update step_order)
    const reorderWorkflowNodes = (workflowId: string, nodes: WorkflowNode[]) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: prev.workflows?.map(workflow =>
                workflow.id === workflowId
                    ? {
                        ...workflow,
                        nodes: nodes.map((node, index) => ({
                            ...node,
                            step_order: index + 1,
                            updated_at: new Date().toISOString(),
                        })),
                        updated_at: new Date().toISOString(),
                    }
                    : workflow
            ) || [],
        }));
    };

    // ========== ADAPTER CONFIG MANAGEMENT (untuk multiple adapters) ==========



    // ========== BULK OPERATIONS ==========

    // Replace entire workflows array
    const setWorkflows = (workflows: WorkflowDefinition[]) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows,
        }));
    };

    // Replace entire adapter_configs array
    const setAdapterConfigs = (adapterConfigs: AdapterConfig[]) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapter_configs: adapterConfigs,
        }));
    };

    // Clear all workflow data
    const clearWorkflows = () => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            workflows: [],
            active_workflow_id: undefined,
        }));
    };

    return {
        // Product info
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        updateMetaConfig,
        updateSitemapConfig,

        // Workflow actions
        updateWorkflowId,
        addWorkflow,
        updateWorkflow,
        deleteWorkflow,
        setActiveWorkflow,
        setWorkflows,
        clearWorkflows,

        // Node actions
        addNodeToWorkflow,
        updateNodeInWorkflow,
        deleteNodeFromWorkflow,
        updateNodeConnections,
        reorderWorkflowNodes,


        setAdapterConfigs,
    };
}