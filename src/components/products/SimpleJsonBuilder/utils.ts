import type { Field } from "@/types/JsonBuilder";

// Generate unique ID
export function genId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

// Convert JSON ke Fields (terima object, bukan string)
export function jsonToFields(data: any, parentKey: string = ""): Field[] {
    if (typeof data !== "object" || data === null) return [];
    
    return Object.entries(data).map(([key, val], idx) => {
        const isObject = typeof val === "object" && val !== null && !Array.isArray(val);
        return {
            id: `${parentKey}-${idx}-${Date.now()}-${Math.random()}`,
            key: key,
            value: isObject ? "" : String(val),
            type: isObject ? "object" : "field",
            children: isObject ? jsonToFields(val, `${parentKey}-${key}`) : [],
            expanded: true,
        };
    });
}

// Convert Fields ke JSON (object, bukan string)
export function fieldsToJson(fields: Field[]): any {
    const result: any = {};
    for (const field of fields) {
        if (field.type === "object") {
            result[field.key] = fieldsToJson(field.children);
        } else {
            result[field.key] = field.value;
        }
    }
    return result;
}

// Get next field number untuk naming
export function getNextFieldNumber(fields: Field[]): number {
    let maxNumber = 0;
    const extractNumber = (items: Field[]) => {
        for (const item of items) {
            if (item.type === "field" && item.key.startsWith("field")) {
                const num = parseInt(item.key.replace("field", ""));
                if (!isNaN(num) && num > maxNumber) maxNumber = num;
            }
            if (item.children.length > 0) {
                extractNumber(item.children);
            }
        }
    };
    extractNumber(fields);
    return maxNumber + 1;
}

// Get next object number untuk naming
export function getNextObjectNumber(fields: Field[]): number {
    let maxNumber = 0;
    const extractNumber = (items: Field[]) => {
        for (const item of items) {
            if (item.type === "object" && item.key.startsWith("object")) {
                const num = parseInt(item.key.replace("object", ""));
                if (!isNaN(num) && num > maxNumber) maxNumber = num;
            }
            if (item.children.length > 0) {
                extractNumber(item.children);
            }
        }
    };
    extractNumber(fields);
    return maxNumber + 1;
}