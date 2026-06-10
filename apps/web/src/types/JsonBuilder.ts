// types/JsonBuilder.ts

export interface Field {
  id: string;
  key: string;
  value: any;
  type: "field" | "object";
  children: Field[];
  expanded?: boolean;
}

// 1. Generate ID unik untuk key React & Tracker Drag-and-Drop
export function genId(): string {
  return Math.random().toString(36).substring(2, 9);
}

export function jsonToFields(obj: any): Field[] {
  // Default kosong - return array kosong
  if (!obj || Object.keys(obj).length === 0) {
    return [];
  }

  if (typeof obj !== "object") return [];

  return Object.entries(obj).map(([key, value]) => {
    if (Array.isArray(value)) {
      return {
        id: genId(),
        key,
        value: "",
        type: "object",
        expanded: true,
        children: value.map((item, index) => {
          if (typeof item === "object" && item !== null) {
            return {
              id: genId(),
              key: `[${index}]`,
              value: "",
              type: "object",
              expanded: true,
              children: jsonToFields(item),
            };
          }

          return {
            id: genId(),
            key: `[${index}]`,
            value: item,
            type: "field",
            children: [],
          };
        }),
      };
    }

    if (typeof value === "object" && value !== null) {
      return {
        id: genId(),
        key,
        value: "",
        type: "object",
        expanded: true,
        children: jsonToFields(value),
      };
    }

    return {
      id: genId(),
      key,
      value,
      type: "field",
      children: [],
    };
  });
}

// 3. Mengubah kembali dari Array Berstruktur (Field[]) menjadi JSON Object murni sebelum di-save
export function fieldsToJson(fields: Field[]): any {
  if (!fields || fields.length === 0) {
    return {};
  }

  // Cek apakah ini bentuk array (cirinya child pertama memiliki key bernilai '[0]')
  const isArray = fields.length > 0 && /^\[\d+\]$/.test(fields[0].key);

  if (isArray) {
    return fields.map((field) => {
      if (field.type === "object") {
        return fieldsToJson(field.children);
      }
      return field.value;
    });
  }

  // Jika berupa object biasa {}
  const obj: any = {};
  fields.forEach((field) => {
    if (!field.key) return; // Skip jika key kosong saat diketik

    if (field.type === "object") {
      obj[field.key] = fieldsToJson(field.children);
    } else {
      obj[field.key] = field.value;
    }
  });

  return obj;
}

// 4. Mencari penomoran otomatis untuk field baru agar tidak bentrok (e.g., field1, field2)
export function getNextFieldNumber(fields: Field[]): number {
  let max = 0;
  fields.forEach((f) => {
    const match = f.key.match(/^field(\d+)$/);
    if (match) {
      const num = parseInt(match[1], 10);
      if (num > max) max = num;
    }
  });
  return max + 1;
}

// 5. Mencari penomoran otomatis untuk object baru (e.g., object1, object2)
export function getNextObjectNumber(fields: Field[]): number {
  let max = 0;
  fields.forEach((f) => {
    const match = f.key.match(/^object(\d+)$/);
    if (match) {
      const num = parseInt(match[1], 10);
      if (num > max) max = num;
    }
  });
  return max + 1;
}