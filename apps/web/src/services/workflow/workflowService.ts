// src/services/workflowService.ts

import { apiClient } from '../api';
import type {
  WorkflowDefinition,
  WorkflowNode,
  WorkflowWithNodes,
} from '../../types/workflow';

// ==================== Workflow Definitions ====================

export async function getWorkflowsByProductId(productId: string): Promise<WorkflowDefinition[]> {
  return apiClient.get(`/products/${productId}/workflows`);
}

export async function getWorkflowById(id: string): Promise<WorkflowDefinition> {
  return apiClient.get(`/workflows/${id}`);
}

export async function createWorkflow(data: { productId: string; name: string }): Promise<WorkflowDefinition> {
  return apiClient.post('/workflows', data);
}

export async function updateWorkflow(id: string, updates: Partial<WorkflowDefinition>): Promise<void> {
  await apiClient.put(`/workflows/${id}`, updates);
}

export async function deleteWorkflow(id: string): Promise<void> {
  await apiClient.delete(`/workflows/${id}`);
}

// ==================== Workflow Nodes ====================

export async function getWorkflowNodes(workflowId: string): Promise<WorkflowNode[]> {
  return apiClient.get(`/workflows/${workflowId}/nodes`);
}

export async function getWorkflowNodeById(id: string): Promise<WorkflowNode> {
  return apiClient.get(`/nodes/${id}`);
}

export async function createWorkflowNode(workflowId: string, data: Partial<WorkflowNode>): Promise<WorkflowNode> {
  return apiClient.post(`/workflows/${workflowId}/nodes`, data);
}

export async function updateWorkflowNode(id: string, updates: Partial<WorkflowNode>): Promise<void> {
  await apiClient.put(`/nodes/${id}`, updates);
}

export async function deleteWorkflowNode(id: string): Promise<void> {
  await apiClient.delete(`/nodes/${id}`);
}

export async function reorderWorkflowNodes(workflowId: string, nodeIds: string[]): Promise<void> {
  await apiClient.put(`/workflows/${workflowId}/nodes/reorder`, nodeIds);
}

// ==================== Batch Operations ====================

// Batch create workflow nodes
export interface BatchCreateNode {
  tempId: string;
  adapterConfigId: string;
  stepOrder: number;
  inputMapping: Record<string, any>;
  nextNodeId?: string | null;
}

export interface BatchCreateResponse {
  tempId: string;
  id: string;
  workflowId: string;
  adapterConfigId: string;
  stepOrder: number;
  inputMapping: Record<string, any>;
  nextNodeId?: string | null;
  createdAt: string;
  updatedAt: string;
}

export async function createWorkflowNodesBatch(
  workflowId: string,
  nodes: BatchCreateNode[]
): Promise<BatchCreateResponse[]> {
  const response = await apiClient.post<BatchCreateResponse[]>(`/workflows/${workflowId}/nodes/batch`, { nodes });
  return response;
}

// Batch update workflow nodes (for multiple nodes at once)
export interface BatchUpdateNode {
  id: string;
  updates: Partial<WorkflowNode>;
}

export async function updateWorkflowNodesBatch(
  workflowId: string,
  updates: BatchUpdateNode[]
): Promise<void> {
  await apiClient.put(`/workflows/${workflowId}/nodes/batch`, { updates });
}

// Batch delete workflow nodes
export async function deleteWorkflowNodesBatch(
  workflowId: string,
  nodeIds: string[]
): Promise<void> {
  await apiClient.delete(`/workflows/${workflowId}/nodes/batch`);
}

// Save all nodes (create + update + delete) in one request
export interface SaveNodesPayload {
  toCreate: BatchCreateNode[];
  toUpdate: BatchUpdateNode[];
  toDelete: string[];
}

export interface SaveNodesResponse {
  created: BatchCreateResponse[];
  updated: WorkflowNode[];
  deleted: string[];
}

export async function saveAllWorkflowNodes(
  workflowId: string,
  payload: SaveNodesPayload
): Promise<SaveNodesResponse> {
  const response = await apiClient.post<SaveNodesResponse>(`/workflows/${workflowId}/nodes`, payload);
  return response;
}

// ==================== Full Workflow with Nodes ====================

export async function getWorkflowWithNodes(workflowId: string): Promise<WorkflowWithNodes> {
  const [workflow, nodes] = await Promise.all([
    getWorkflowById(workflowId),
    getWorkflowNodes(workflowId),
  ]);
  return { ...workflow, nodes };
}