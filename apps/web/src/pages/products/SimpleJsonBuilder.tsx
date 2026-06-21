// SimpleJsonBuilder.tsx
import { useState, useEffect } from "react";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import {
    Plus,
    Trash2,
    GripVertical,
    ChevronRight,
    ChevronDown,
} from "lucide-react";

// ==================== TYPES ====================

interface Field {
    id: string;
    key: string;
    value: string;
    type: "field" | "object";
    children: Field[];
    expanded: boolean;
}

interface PlaceholderItem {
    value: string;
    label: string;
    group?: string;
}

interface SimpleJsonBuilderProps {
    value: any;
    onChange: (value: any) => void;
    placeholders?: PlaceholderItem[]; // MODE 1: Props placeholders (Payload/Node response)
}

// ==================== MODE 2: DEFAULT PLACEHOLDERS (BUILT-IN) ====================

const DEFAULT_PLACEHOLDERS: PlaceholderItem[] = [
    { value: "{id}", label: "Unique ID", group: "Basic" },
    { value: "{title}", label: "Article Title", group: "Article" },
    { value: "{slug}", label: "Article Slug", group: "Article" },
    { value: "{tags}", label: "Article Tags (array)", group: "Article" },
    { value: "{keywords}", label: "Article Keywords (array)", group: "Article" },
    { value: "{topic}", label: "Article Topic", group: "Article" },
    { value: "{timestamp}", label: "Current Timestamp", group: "Utility" },
    { value: "{content}", label: "Article Content (HTML)", group: "Content" },
    { value: "{content_text}", label: "Article Content (Plain Text)", group: "Content" },
    { value: "{content_with_image}", label: "Article Content + Image after H1 (HTML)", group: "Content" },
    { value: "{excerpt}", label: "Article Excerpt", group: "Content" },
    { value: "{image_url}", label: "Image URL (Plain)", group: "Image" },
    { value: "{image_content_html}", label: "Image HTML Tag", group: "Image" },
    { value: "{scheduled_for}", label: "Scheduled Time", group: "Schedule" },
    { value: "{meta_title}", label: "Meta Title", group: "SEO" },
    { value: "{meta_description}", label: "Meta Description", group: "SEO" },
    { value: "{meta_keywords}", label: "Meta Keywords (JSON)", group: "SEO" },
    { value: "{og_title}", label: "OG Title", group: "Social Media" },
    { value: "{og_description}", label: "OG Description", group: "Social Media" },
    { value: "{og_image}", label: "OG Image URL", group: "Social Media" },
    { value: "{twitter_title}", label: "Twitter Title", group: "Social Media" },
    { value: "{twitter_description}", label: "Twitter Description", group: "Social Media" },
    { value: "{twitter_image}", label: "Twitter Image URL", group: "Social Media" },
    { value: "{sitemap_priority}", label: "Sitemap Priority", group: "SEO" },
    { value: "{sitemap_changefreq}", label: "Sitemap Change Frequency", group: "SEO" },
];

// ==================== HELPER FUNCTIONS ====================

const genId = () => `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

const getNextFieldNumber = (fields: Field[]): number => {
    let max = 0;
    const traverse = (items: Field[]) => {
        items.forEach(item => {
            const match = item.key.match(/^field(\d+)$/);
            if (match) max = Math.max(max, parseInt(match[1], 10));
            if (item.children.length) traverse(item.children);
        });
    };
    traverse(fields);
    return max + 1;
};

const getNextObjectNumber = (fields: Field[]): number => {
    let max = 0;
    const traverse = (items: Field[]) => {
        items.forEach(item => {
            const match = item.key.match(/^object(\d+)$/);
            if (match) max = Math.max(max, parseInt(match[1], 10));
            if (item.children.length) traverse(item.children);
        });
    };
    traverse(fields);
    return max + 1;
};

const jsonToFields = (obj: any, parentKey?: string): Field[] => {
    if (!obj || typeof obj !== "object") return [];

    return Object.entries(obj).map(([key, value]) => {
        const isObject = value && typeof value === "object" && !Array.isArray(value);
        return {
            id: genId(),
            key,
            value: isObject ? "" : String(value),
            type: isObject ? "object" : "field",
            children: isObject ? jsonToFields(value, key) : [],
            expanded: true,
        };
    });
};

const fieldsToJson = (fields: Field[]): any => {
    const result: any = {};
    fields.forEach(field => {
        if (!field.key.trim()) return;

        if (field.type === "object") {
            result[field.key] = fieldsToJson(field.children);
        } else {
            const val = field.value.trim();
            if ((val.startsWith("{") && val.endsWith("}")) ||
                (val.startsWith("[") && val.endsWith("]"))) {
                try {
                    result[field.key] = JSON.parse(val);
                } catch {
                    result[field.key] = val;
                }
            } else {
                result[field.key] = val;
            }
        }
    });
    return result;
};

// ==================== PLACEHOLDER MODAL (DYNAMIC) ====================

function PlaceholderModal({
    placeholders,
    onSelect,
    onClose,
}: {
    placeholders: PlaceholderItem[] | undefined;
    onSelect: (value: string) => void;
    onClose: () => void;
}) {
    const [search, setSearch] = useState("");

    const filtered = placeholders?.filter(
        (p) =>
            p.value.toLowerCase().includes(search.toLowerCase()) ||
            p.label.toLowerCase().includes(search.toLowerCase())
    );

    // Group by group
    const grouped = filtered?.reduce((acc, p) => {
        const group = p.group || "Lainnya";
        if (!acc[group]) acc[group] = [];
        acc[group].push(p);
        return acc;
    }, {} as Record<string, PlaceholderItem[]>);

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center "
            onMouseDown={onClose}
        >
            <div
                className="bg-white dark:bg-slate-900 rounded-xl shadow-xl w-[500px] max-h-[550px] flex flex-col"
                onMouseDown={(e) => e.stopPropagation()}
            >
                <div className="p-4 border-b flex items-center justify-between dark:border-slate-700">
                    <div>
                        <span className="font-semibold text-sm">Pilih Placeholder</span>
                        <p className="text-xs text-slate-500 mt-0.5">
                            Pilih nilai dari placeholder yang tersedia
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="text-slate-400 hover:text-slate-600 text-xl leading-none"
                    >
                        ✕
                    </button>
                </div>

                <div className="p-3 border-b dark:border-slate-700">
                    <Input
                        autoFocus
                        placeholder="Cari placeholder..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="h-8"
                    />
                </div>

                <div className="overflow-y-auto flex-1 p-2">
                    {Object.keys(grouped as Record<string, PlaceholderItem[]>)?.length === 0 ? (
                        <div className="text-center py-6 text-slate-400 text-sm">
                            Tidak ditemukan
                        </div>
                    ) : (
                        Object.entries(grouped as Record<string, PlaceholderItem[]>).map(([group, items]) => (
                            <div key={group} className="mb-4">
                                <div className="text-xs font-semibold text-slate-500 px-2 py-1 uppercase tracking-wide">
                                    {group}
                                </div>
                                <div className="grid grid-cols-2 gap-1">
                                    {items.map((p) => (
                                        <button
                                            key={p.value}
                                            onClick={() => onSelect(p.value)}
                                            className="text-left px-3 py-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                                        >
                                            <div className="font-mono text-blue-600 dark:text-blue-400 text-xs truncate">
                                                {p.value}
                                            </div>
                                            <div className="text-xs text-slate-500 truncate">
                                                {p.label}
                                            </div>
                                        </button>
                                    ))}
                                </div>
                            </div>
                        ))
                    )}
                </div>

                <div className="p-3 border-t dark:border-slate-700 flex justify-end">
                    <Button variant="outline" size="sm" onClick={onClose}>
                        Tutup
                    </Button>
                </div>
            </div>
        </div>
    );
}

// ==================== MAIN COMPONENT ====================

export function SimpleJsonBuilder({
    value,
    onChange,
    placeholders, // OPTIONAL: kalau ada, pake mode 1 (Payload/Node response)
}: SimpleJsonBuilderProps) {

    const [fields, setFields] = useState<Field[]>([]);
    const [showPlaceholderFor, setShowPlaceholderFor] = useState<string | null>(null);
    const [isInternalUpdate, setIsInternalUpdate] = useState(false);

    // ==================== MODE LOGIC ====================
    // MODE 1: Jika props placeholders ada (tidak undefined), pakai dari props
    // MODE 2: Jika props placeholders tidak ada, pakai DEFAULT_PLACEHOLDERS (built-in)
    const activePlaceholders = [...DEFAULT_PLACEHOLDERS, ...(placeholders ?? [])];
    // Atau pakai nullish coalescing: const activePlaceholders = placeholders ?? DEFAULT_PLACEHOLDERS;

    const parseValue = (val: any) => {
        if (!val) return {};
        if (typeof val === "string") {
            try {
                return JSON.parse(val);
            } catch {
                return {};
            }
        }
        return val;
    };

    useEffect(() => {
        if (isInternalUpdate) {
            setIsInternalUpdate(false);
            return;
        }

        const parsed = parseValue(value);
        if (parsed && typeof parsed === "object") {
            const incoming = JSON.stringify(parsed);
            const current = JSON.stringify(fieldsToJson(fields));
            if (incoming !== current) {
                setFields(jsonToFields(parsed));
            }
        }
    }, [value]);

    useEffect(() => {
        const t = setTimeout(() => {
            setIsInternalUpdate(true);
            onChange(fieldsToJson(fields));
        }, 250);
        return () => clearTimeout(t);
    }, [fields]);

    const updateTree = (
        items: Field[],
        id: string,
        updater: (item: Field) => Field
    ): Field[] => {
        return items.map((item) => {
            if (item.id === id) return updater(item);
            if (item.children?.length) {
                return { ...item, children: updateTree(item.children, id, updater) };
            }
            return item;
        });
    };

    const addField = (parentId?: string) => {
        const newField: Field = {
            id: genId(),
            key: `field${getNextFieldNumber(fields)}`,
            value: "",
            type: "field",
            children: [],
            expanded: true,
        };
        if (!parentId) {
            setFields((prev) => [...prev, newField]);
            return;
        }
        setFields((prev) =>
            updateTree(prev, parentId, (item) => ({
                ...item,
                children: [...item.children, newField],
            }))
        );
    };

    const addObject = (parentId?: string) => {
        const newObject: Field = {
            id: genId(),
            key: `object${getNextObjectNumber(fields)}`,
            value: "",
            type: "object",
            children: [],
            expanded: true,
        };
        if (!parentId) {
            setFields((prev) => [...prev, newObject]);
            return;
        }
        setFields((prev) =>
            updateTree(prev, parentId, (item) => ({
                ...item,
                children: [...item.children, newObject],
            }))
        );
    };

    const updateField = (id: string, key: string, value: string) => {
        setFields((prev) =>
            updateTree(prev, id, (item) => ({ ...item, key, value }))
        );
    };

    const deleteRecursive = (items: Field[], id: string): Field[] =>
        items
            .filter((item) => item.id !== id)
            .map((item) => ({
                ...item,
                children: deleteRecursive(item.children, id),
            }));

    const deleteField = (id: string) => {
        setFields((prev) => deleteRecursive(prev, id));
    };

    const toggleExpand = (id: string) => {
        setFields((prev) =>
            updateTree(prev, id, (item) => ({ ...item, expanded: !item.expanded }))
        );
    };

    const insertPlaceholder = (fieldId: string, placeholder: string) => {
        setFields((prev) =>
            updateTree(prev, fieldId, (item) => ({ ...item, value: placeholder }))
        );
        setShowPlaceholderFor(null);
    };

    const renderField = (field: Field, level = 0) => {
        const indent = level * 20;
        const isObject = field.type === "object";
        const isPlaceholder = activePlaceholders?.some((p) => p.value === field.value);

        return (
            <div key={field.id} style={{ marginLeft: indent }} className="mb-2">
                <div className="flex items-center gap-2 p-2 rounded-lg border bg-white dark:bg-slate-900 dark:border-slate-700">

                    <GripVertical className="w-4 h-4 text-slate-400 shrink-0 cursor-move" />

                    {isObject ? (
                        <button onClick={() => toggleExpand(field.id)} className="shrink-0">
                            {field.expanded ? (
                                <ChevronDown className="w-4 h-4" />
                            ) : (
                                <ChevronRight className="w-4 h-4" />
                            )}
                        </button>
                    ) : (
                        <div className="w-5 shrink-0" />
                    )}

                    <Input
                        value={field.key}
                        onChange={(e) => updateField(field.id, e.target.value, field.value)}
                        className="w-36 h-8 shrink-0 font-mono text-sm"
                        placeholder="field_name"
                    />

                    {!isObject && (
                        <>
                            <span className="text-slate-400">:</span>
                            <Input
                                value={field.value}
                                onChange={(e) => updateField(field.id, field.key, e.target.value)}
                                className={`flex-1 h-8 font-mono text-sm ${isPlaceholder ? "text-blue-600 dark:text-blue-400" : ""}`}
                                placeholder="value atau placeholder..."
                            />
                            <Button
                                size="sm"
                                variant="outline"
                                onClick={() => setShowPlaceholderFor(field.id)}
                                className="shrink-0"
                                title="Pilih dari placeholder"
                            >
                                📋
                            </Button>
                        </>
                    )}

                    {isObject && (
                        <div className="flex gap-1 ml-auto">
                            <Button size="sm" variant="outline" onClick={() => addField(field.id)}>
                                <Plus className="w-3 h-3 mr-1" />
                                Field
                            </Button>
                            <Button size="sm" variant="outline" onClick={() => addObject(field.id)}>
                                <Plus className="w-3 h-3 mr-1" />
                                Object
                            </Button>
                        </div>
                    )}

                    <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => deleteField(field.id)}
                        className="shrink-0"
                    >
                        <Trash2 className="w-4 h-4 text-red-500" />
                    </Button>

                </div>

                {isObject && field.expanded && field.children.length > 0 && (
                    <div className="ml-6 pl-4 mt-1 border-l-2 border-blue-200 dark:border-blue-800">
                        {field.children.map((child) => renderField(child, level + 1))}
                    </div>
                )}
            </div>
        );
    };

    return (
        <div className="space-y-4">

            {showPlaceholderFor && (
                <PlaceholderModal
                    placeholders={activePlaceholders}
                    onSelect={(val) => insertPlaceholder(showPlaceholderFor, val)}
                    onClose={() => setShowPlaceholderFor(null)}
                />
            )}

            <div className="flex gap-2">
                <Button size="sm" onClick={() => addField()} variant="secondary">
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Field
                </Button>
                <Button size="sm" variant="outline" onClick={() => addObject()}>
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Object
                </Button>
            </div>

            {fields.length ? (
                <div className="space-y-2">
                    {fields.map((field) => renderField(field))}
                </div>
            ) : (
                <div className="text-center py-8 text-slate-500 dark:text-slate-400 border-2 border-dashed rounded-lg">
                    <p>Belum ada field</p>
                    <p className="text-sm mt-1">Klik "Tambah Field" untuk memulai</p>
                </div>
            )}

            {/* Preview JSON */}
            {fields.length > 0 && (
                <div className="mt-4 p-3 bg-slate-50 dark:bg-slate-800 rounded-lg">
                    <div className="text-xs font-semibold text-slate-500 mb-2">Preview Output:</div>
                    <pre className="text-xs font-mono overflow-auto max-h-40 p-2 bg-white dark:bg-slate-900 rounded border dark:border-slate-700">
                        {JSON.stringify(fieldsToJson(fields), null, 2)}
                    </pre>
                </div>
            )}
        </div>
    );
}