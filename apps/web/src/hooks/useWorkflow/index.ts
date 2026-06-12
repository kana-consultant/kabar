// src/hooks/useWorkflow.ts

import { useState, useCallback } from "react";
import {
  getWorkflowsByProductId,
  getWorkflowWithNodes,
  createWorkflow,
  updateWorkflow,
  deleteWorkflow,
  getWorkflowNodes,
  createWorkflowNode,
  updateWorkflowNode,
  deleteWorkflowNode,
  reorderWorkflowNodes,
} from "@/services/workflow/workflowService";
import type {
  WorkflowDefinition,
  WorkflowNode,
  WorkflowWithNodes,
} from "@/types/workflow";
import { useToast } from "../use-toast";

export function useWorkflow(productId: string) {
  const [workflows, setWorkflows] = useState<WorkflowDefinition[]>([]);
  const [currentWorkflow, setCurrentWorkflow] = useState<WorkflowWithNodes | null>(null);
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
  }, [productId]);

  const fetchWorkflowWithNodes = useCallback(async (workflowId: string) => {
    setLoading(true);
    try {
      const data = await getWorkflowWithNodes(workflowId);
      setCurrentWorkflow(data);
    } catch (error: any) {
      toast.error(error.message || "Failed to fetch workflow");
    } finally {
      setLoading(false);
    }
  }, []);

  // const create = useCallback(async (name: string) => {
  //   setLoading(true);
  //   try {
  //     const data = await createWorkflow({ productId, name });
  //     setWorkflows((prev) => [...prev, data]);
  //     toast.success("Workflow created");
  //     return data;
  //   } catch (error: any) {
  //     toast.error(error.message || "Failed to create workflow");
  //   } finally {
  //     setLoading(false);
  //   }
  // }, [productId]);


  const create = useCallback(async (name: string) => {
    setLoading(true);
    try {
      // Create temporary workflow object
      const tempId = `temp_wf_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      const tempWorkflow = {
        id: tempId,
        productId: productId,
        name: name,
        isActive: false,
        isTemp: true, // Mark as temporary
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      // Add to local state without API call
      setWorkflows((prev) => [...prev, tempWorkflow]);
      toast.success("Workflow created locally");

      // Optional: Initialize temp nodes storage for this workflow
      // You might want to add this to a separate temp storage system
      // For example: setTempWorkflows(prev => ({ ...prev, [tempId]: tempWorkflow }));

      return tempWorkflow;
    } catch (error: any) {
      toast.error(error.message || "Failed to create workflow locally");
    } finally {
      setLoading(false);
    }
  }, [productId]);
  const update = useCallback(async (id: string, updates: Partial<WorkflowDefinition>) => {
    setLoading(true);
    try {
      await updateWorkflow(id, updates);
      toast.success("Workflow updated");
      fetchWorkflows();
    } catch (error: any) {
      toast.error(error.message || "Failed to update workflow");
    } finally {
      setLoading(false);
    }
  }, [fetchWorkflows]);

  const remove = useCallback(async (id: string) => {
    setLoading(true);
    try {
      await deleteWorkflow(id);
      setWorkflows((prev) => prev.filter((w) => w.id !== id));
      toast.success("Workflow deleted");
    } catch (error: any) {
      toast.error(error.message || "Failed to delete workflow");
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    workflows,
    currentWorkflow,
    loading,
    fetchWorkflows,
    fetchWorkflowWithNodes,
    createWorkflow: create,
    updateWorkflow: update,
    deleteWorkflow: remove,
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
  }, [workflowId]);

  const createNode = useCallback(async (data: Partial<WorkflowNode>) => {
    setLoading(true);
    try {
      const node = await createWorkflowNode(workflowId, data);
      setNodes((prev) => [...prev, node]);
      toast.success("Node added");
      return node;
    } catch (error: any) {
      toast.error(error.message || "Failed to add node");
    } finally {
      setLoading(false);
    }
  }, [workflowId]);

  const updateNode = useCallback(async (nodeId: string, updates: Partial<WorkflowNode>) => {
    setLoading(true);
    try {
      await updateWorkflowNode(nodeId, updates);
      toast.success("Node updated");
      fetchNodes();
    } catch (error: any) {
      toast.error(error.message || "Failed to update node");
    } finally {
      setLoading(false);
    }
  }, [fetchNodes]);

  const removeNode = useCallback(async (nodeId: string) => {
    setLoading(true);
    try {
      await deleteWorkflowNode(nodeId);
      setNodes((prev) => prev.filter((n) => n.id !== nodeId));
      toast.success("Node deleted");
    } catch (error: any) {
      toast.error(error.message || "Failed to delete node");
    } finally {
      setLoading(false);
    }
  }, []);

  const reorderNodes = useCallback(async (nodeIds: string[]) => {
    setLoading(true);
    try {
      await reorderWorkflowNodes(workflowId, nodeIds);
      fetchNodes();
    } catch (error: any) {
      toast.error(error.message || "Failed to reorder nodes");
    } finally {
      setLoading(false);
    }
  }, [workflowId, fetchNodes]);

  return {
    nodes,
    loading,
    fetchNodes,
    createNode,
    updateNode,
    deleteNode: removeNode,
    reorderNodes,
  };
}