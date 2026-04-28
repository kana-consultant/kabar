import type { FieldMapping } from "@/types/product";

interface InternalData {
    title: string;
    body: string;
    imageUrl: string;
    status: string;
    category?: string;
    tags?: string[];
    slug?: string;
    excerpt?: string;
    author?: string;
    date?: string;
    [key: string]: any;
}

export function transformPayload(
    internalData: InternalData,
    fieldMapping: FieldMapping[]
): Record<string, any> {
    const result: Record<string, any> = {};

    fieldMapping.forEach((field) => {
        // Ambil nilai dari source
        let value = internalData[field.sourceField];

        // Jika tidak ada nilai, pakai default
        if (value === undefined || value === null) {
            if (field.defaultValue) {
                value = field.defaultValue;
            } else if (field.isRequired) {
                console.warn(`Required field ${field.sourceField} is missing`);
                return;
            } else {
                return;
            }
        }

        // Handle nested object (parent.child)
        if (field.targetField.includes('.')) {
            const parts = field.targetField.split('.');
            let current = result;

            for (let i = 0; i < parts.length - 1; i++) {
                if (!current[parts[i]]) {
                    current[parts[i]] = {};
                }
                current = current[parts[i]];
            }
            current[parts[parts.length - 1]] = value;
        } else {
            // Field biasa (top level)
            result[field.targetField] = value;
        }
    });

    return result;
}

// Contoh penggunaan:
// const internalData = {
//     title: "Judul Artikel",
//     body: "<p>Isi artikel</p>",
//     imageUrl: "https://example.com/image.jpg",
//     status: "publish"
// };
//
// const mapping = [
//     { sourceField: "title", targetField: "post.title", isRequired: true },
//     { sourceField: "body", targetField: "post.content", isRequired: true },
//     { sourceField: "imageUrl", targetField: "meta.thumbnail", isRequired: false },
// ];
//
// const payload = transformPayload(internalData, mapping);
// // Hasil:
// // {
// //     post: { title: "Judul Artikel", content: "<p>Isi artikel</p>" },
// //     meta: { thumbnail: "https://example.com/image.jpg" }
// // }