// src/types/workflow.ts

export interface WorkflowDefinition {
  id: string;
  productId: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface WorkflowNode {
  id: string;
  workflowId: string;
  adapterConfigId: string;
  stepOrder: number;
  inputMapping: Record<string, any>;
  nextNodeId: string | null;
  createdAt: string;
}

export interface AdapterConfig {
  id: string;
  productId: string;
  endpointPath: string;
  httpMethod: string;
  customHeaders: string;
  fieldMapping: string;
  responseMapping: string;
  metaConfig: string;
  sitemapConfig: string;
  timeoutSeconds: number;
  retryCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface WorkflowNodeWithConfig extends WorkflowNode {
  adapterConfig?: AdapterConfig;
}

export interface WorkflowWithNodes extends WorkflowDefinition {
  nodes: WorkflowNodeWithConfig[];
}