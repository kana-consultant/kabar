import { type Field } from "@/types/JsonBuilder";

export function genId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 6)}`;
}

export function jsonToFields(data: any, parentKey: string = ""): Field[] {
    if (typeof data !== "object" || data === null) return [];
    if (Array.isArray(data)) {
        const obj: any = {};
        data.forEach((item, idx) => {
            obj[idx] = item;
        });
        return jsonToFields(obj, parentKey);
    }

    return Object.entries(data).map(([key, val], idx) => {
        const isObject = typeof val === "object" && val !== null && !Array.isArray(val);
        const isArray = Array.isArray(val);
        return {
            id: `${parentKey}-${idx}-${Date.now()}-${Math.random()}`,
            key: key,
            value: isObject || isArray ? "" : String(val),
            type: (isObject || isArray) ? "object" : "field",
            children: isObject ? jsonToFields(val, `${parentKey}-${key}`) : [],
            expanded: true,
        };
    });
}

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

export function getNextFieldNumber(fields: Field[]): number {
    let maxNumber = 0;
    const regex = /^field(\d+)$/;

    const checkFields = (items: Field[]) => {
        for (const field of items) {
            const match = field.key.match(regex);
            if (match) {
                const num = parseInt(match[1], 10);
                if (num > maxNumber) maxNumber = num;
            }
            if (field.children.length > 0) {
                checkFields(field.children);
            }
        }
    };

    checkFields(fields);
    return maxNumber + 1;
}

export function getNextObjectNumber(fields: Field[]): number {
    let maxNumber = 0;
    const regex = /^object(\d+)$/;

    const checkFields = (items: Field[]) => {
        for (const field of items) {
            const match = field.key.match(regex);
            if (match) {
                const num = parseInt(match[1], 10);
                if (num > maxNumber) maxNumber = num;
            }
            if (field.children.length > 0) {
                checkFields(field.children);
            }
        }
    };

    checkFields(fields);
    return maxNumber + 1;
}