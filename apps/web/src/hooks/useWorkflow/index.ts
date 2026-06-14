// src/hooks/useWorkflow.ts

import { useState, useCallback } from "react";
import {
  getWorkflowsByProductId,
  updateWorkflow,
  deleteWorkflow,
  getWorkflowNodes,
  createWorkflowNode,
  updateWorkflowNode,
  deleteWorkflowNode,
  reorderWorkflowNodes,
  updateWorkflowNodes as updateWorkflowNodesService,
} from "@/services/workflow/workflowService";
import type { WorkflowDefinition, WorkflowNode } from "@/types/product";
import { useToast } from "../use-toast";

export function useWorkflow(productId: string) {
  const [workflows, setWorkflows] = useState<WorkflowDefinition[]>([]);
  const [currentWorkflow, setCurrentWorkflow] = useState<WorkflowDefinition | null>(null);
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const fetchWorkflows = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getWorkflowsByProductId(productId);
      setWorkflows(data);
    } catch (error: any) {
      toast.error(error.message || "Failed to fetch workflows");
    } finally {
      setLoading(false);
    }
  }, [productId, toast]);

  // Hapus fetchWorkflowWithNodes karena tidak menggunakan getWorkflowWithNodes

  const create = useCallback(async (name: string) => {
    setLoading(true);
    try {
      // Create temporary workflow object
      const tempId = `temp_wf_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      const tempWorkflow: WorkflowDefinition = {
        id: tempId,
        product_id: productId,
        name: name,
        nodes: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      setWorkflows((prev) => [...prev, tempWorkflow]);
      toast.success("Workflow created locally");

      return tempWorkflow;
    } catch (error: any) {
      toast.error(error.message || "Failed to create workflow locally");
      return null;
    } finally {
      setLoading(false);
    }
  }, [productId, toast]);

  const update = useCallback(async (id: string, updates: Partial<WorkflowDefinition>) => {
    setLoading(true);
    try {
      await updateWorkflow(id, updates);
      // Update local state
      setWorkflows((prev) => 
        prev.map((w) => w.id === id ? { ...w, ...updates, updated_at: new Date().toISOString() } : w)
      );
      // Also update currentWorkflow if it's the one being updated
      setCurrentWorkflow((prev) => 
        prev?.id === id ? { ...prev, ...updates, updated_at: new Date().toISOString() } : prev
      );
      toast.success("Workflow updated");
      return true;
    } catch (error: any) {
      toast.error(error.message || "Failed to update workflow");
      return false;
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const remove = useCallback(async (id: string) => {
    setLoading(true);
    try {
      await deleteWorkflow(id);
      setWorkflows((prev) => prev.filter((w) => w.id !== id));
      // Clear currentWorkflow if it's the one being deleted
      setCurrentWorkflow((prev) => prev?.id === id ? null : prev);
      toast.success("Workflow deleted");
      return true;
    } catch (error: any) {
      toast.error(error.message || "Failed to delete workflow");
      return false;
    } finally {
      setLoading(false);
    }
  }, [toast]);

  // Optional: Set current workflow manually (without fetching nodes)
  const setCurrentWorkflowById = useCallback((workflow: WorkflowDefinition | null) => {
    setCurrentWorkflow(workflow);
  }, []);

  return {
    workflows,
    currentWorkflow,
    loading,
    fetchWorkflows,
    createWorkflow: create,
    updateWorkflow: update,
    deleteWorkflow: remove,
    setCurrentWorkflow: setCurrentWorkflowById,
  };
}

export function useWorkflowNodes(workflowId: string) {
  const [nodes, setNodes] = useState<WorkflowNode[]>([]);
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  const fetchNodes = useCallback(async () => {
    if (!workflowId) return;
    setLoading(true);
    try {
      const data = await getWorkflowNodes(workflowId);
      setNodes(data);
    } catch (error: any) {
      toast.error(error.message || "Failed to fetch nodes");
    } finally {
      setLoading(false);
    }
  }, [workflowId, toast]);

  const createNode = useCallback(async (data: Partial<WorkflowNode>) => {
    setLoading(true);
    try {
      const node = await createWorkflowNode(workflowId, data);
      setNodes((prev) => [...prev, node]);
      toast.success("Node added");
      return node;
    } catch (error: any) {
      toast.error(error.message || "Failed to add node");
      return null;
    } finally {
      setLoading(false);
    }
  }, [workflowId, toast]);

  const updateNode = useCallback(async (nodeId: string, updates: Partial<WorkflowNode>) => {
    setLoading(true);
    try {
      await updateWorkflowNode(nodeId, updates);
      // Update local state
      setNodes((prev) => 
        prev.map((n) => n.id === nodeId ? { ...n, ...updates } : n)
      );
      toast.success("Node updated");
      return true;
    } catch (error: any) {
      toast.error(error.message || "Failed to update node");
      return false;
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const removeNode = useCallback(async (nodeId: string) => {
    setLoading(true);
    try {
      await deleteWorkflowNode(nodeId);
      setNodes((prev) => prev.filter((n) => n.id !== nodeId));
      toast.success("Node deleted");
      return true;
    } catch (error: any) {
      toast.error(error.message || "Failed to delete node");
      return false;
    } finally {
      setLoading(false);
    }
  }, [toast]);

  const reorderNodes = useCallback(async (nodeIds: string[]) => {
    setLoading(true);
    try {
      await reorderWorkflowNodes(workflowId, nodeIds);
      // Reorder local state
      setNodes((prev) => {
        const orderedNodes: WorkflowNode[] = [];
        nodeIds.forEach(id => {
          const node = prev.find(n => n.id === id);
          if (node) orderedNodes.push(node);
        });
        return orderedNodes;
      });
      toast.success("Nodes reordered");
      return true;
    } catch (error: any) {
      toast.error(error.message || "Failed to reorder nodes");
      return false;
    } finally {
      setLoading(false);
    }
  }, [workflowId, toast]);

  // Batch update multiple nodes
  const updateWorkflowNodes = useCallback(async (
    updates: Array<{ id: string; updates: Partial<WorkflowNode> }>
  ) => {
    setLoading(true);
    try {
      const updatedNodes = await updateWorkflowNodesService(workflowId, updates);
      
      // Update local state
      setNodes((prevNodes) => {
        const updatedMap = new Map(updatedNodes.map(node => [node.id, node]));
        return prevNodes.map(node => updatedMap.get(node.id) || node);
      });
      
      toast.success(`${updates.length} node(s) updated successfully`);
      return updatedNodes;
    } catch (error: any) {
      toast.error(error.message || "Failed to update nodes");
      return null;
    } finally {
      setLoading(false);
    }
  }, [workflowId, toast]);

  // Update same property for multiple nodes
  const updateNodesProperty = useCallback(async (
    nodeIds: string[],
    property: keyof WorkflowNode,
    value: any
  ) => {
    setLoading(true);
    try {
      const updates = nodeIds.map(id => ({
        id,
        updates: { [property]: value } as Partial<WorkflowNode>
      }));
      
      const updatedNodes = await updateWorkflowNodesService(workflowId, updates);
      
      setNodes((prevNodes) => {
        const updatedMap = new Map(updatedNodes.map(node => [node.id, node]));
        return prevNodes.map(node => updatedMap.get(node.id) || node);
      });
      
      toast.success(`Updated ${nodeIds.length} node(s)`);
      return updatedNodes;
    } catch (error: any) {
      toast.error(error.message || "Failed to update nodes");
      return null;
    } finally {
      setLoading(false);
    }
  }, [workflowId, toast]);

  return {
    nodes,
    loading,
    fetchNodes,
    createNode,
    updateNode,
    deleteNode: removeNode,
    reorderNodes,
    updateWorkflowNodes,    // Batch update
    updateNodesProperty,     // Update property for multiple nodes
  };
}