export interface AdapterConfig {
    endpoint: string;
    method: 'POST' | 'PUT' | 'PATCH';
    headers: Record<string, string>;
    // Bisa berupa array (standar) ATAU raw JSON string
    fieldMapping: FieldMapping[] | string;  // ← string untuk raw JSON
}

export interface FieldMapping {
    id: string;
    sourceField: string;
    targetField: string;
    isRequired: boolean;
    defaultValue?: string;
}