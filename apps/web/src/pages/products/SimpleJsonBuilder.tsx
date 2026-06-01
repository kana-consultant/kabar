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

import {
    type Field,
    type SimpleJsonBuilderProps,
} from "@/types/JsonBuilder";

import {
    genId,
    jsonToFields,
    fieldsToJson,
    getNextFieldNumber,
    getNextObjectNumber,
} from "@/utils/SimpleJsonBuilder";

const PLACEHOLDERS = [
    { value: "{id}", label: "Unique ID" },
    { value: "{title}", label: "Article Title" },
    { value: "{slug}", label: "Article Slug" },
    { value: "{tags}", label: "Article Tags (array)" },
    { value: "{keywords}", label: "Article Keywords (array)" },
    { value: "{topic}", label: "Article Topic" },
    { value: "{timestamp}", label: "Current Timestamp" },
    { value: "{content}", label: "Article Content (HTML)" },
    { value: "{content_text}", label: "Article Content (Plain Text)" },
    { value: "{content_with_image}", label: "Article Content + Image after H1 (HTML)" },
    { value: "{excerpt}", label: "Article Excerpt" },
    { value: "{image_url}", label: "Image URL (Plain)" },
    { value: "{image_content_html}", label: "Image HTML Tag" },
    { value: "{scheduled_for}", label: "Scheduled Time" },
    { value: "{meta_title}", label: "Meta Title" },
    { value: "{meta_description}", label: "Meta Description" },
    { value: "{meta_keywords}", label: "Meta Keywords (JSON)" },
    { value: "{og_title}", label: "OG Title" },
    { value: "{og_description}", label: "OG Description" },
    { value: "{og_image}", label: "OG Image URL" },
    { value: "{twitter_title}", label: "Twitter Title" },
    { value: "{twitter_description}", label: "Twitter Description" },
    { value: "{twitter_image}", label: "Twitter Image URL" },
    { value: "{sitemap_priority}", label: "Sitemap Priority" },
    { value: "{sitemap_changefreq}", label: "Sitemap Change Frequency" },
];

function PlaceholderModal({
    onSelect,
    onClose,
}: {
    onSelect: (value: string) => void;
    onClose: () => void;
}) {
    const [search, setSearch] = useState("");

    const filtered = PLACEHOLDERS.filter(
        (p) =>
            p.value.toLowerCase().includes(search.toLowerCase()) ||
            p.label.toLowerCase().includes(search.toLowerCase())
    );

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
            onMouseDown={onClose}
        >
            <div
                className="bg-white dark:bg-slate-900 rounded-xl shadow-xl w-[480px] max-h-[520px] flex flex-col"
                onMouseDown={(e) => e.stopPropagation()}
            >
                <div className="p-4 border-b flex items-center justify-between">
                    <span className="font-semibold text-sm">Pilih Placeholder</span>
                    <button
                        onClick={onClose}
                        className="text-slate-400 hover:text-slate-600 text-lg leading-none"
                    >
                        ✕
                    </button>
                </div>

                <div className="p-3 border-b">
                    <Input
                        autoFocus
                        placeholder="Cari placeholder..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="h-8"
                    />
                </div>

                <div className="overflow-y-auto flex-1 p-2">
                    {filtered.length === 0 ? (
                        <div className="text-center py-6 text-slate-400 text-sm">
                            Tidak ditemukan
                        </div>
                    ) : (
                        <div className="grid grid-cols-2 gap-1">
                            {filtered.map((p) => (
                                <button
                                    key={p.value}
                                    onClick={() => onSelect(p.value)}
                                    className="text-left px-3 py-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700"
                                >
                                    <div className="font-mono text-blue-600 text-xs truncate">
                                        {p.value}
                                    </div>
                                    <div className="text-xs text-slate-500 truncate">
                                        {p.label}
                                    </div>
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export function SimpleJsonBuilder({
    value,
    onChange,
}: SimpleJsonBuilderProps) {

    const [fields, setFields] = useState<Field[]>([]);
    const [showPlaceholderFor, setShowPlaceholderFor] = useState<string | null>(null);
    const [isInternalUpdate, setIsInternalUpdate] = useState(false);

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
        const isPlaceholder = PLACEHOLDERS.some((p) => p.value === field.value);

        return (
            <div key={field.id} style={{ marginLeft: indent }} className="mb-2">
                <div className="flex items-center gap-2 p-2 rounded-lg border bg-white dark:bg-slate-900">

                    <GripVertical className="w-4 h-4 text-slate-400 shrink-0" />

                    {isObject ? (
                        <button onClick={() => toggleExpand(field.id)}>
                            {field.expanded ? (
                                <ChevronDown className="w-4 h-4" />
                            ) : (
                                <ChevronRight className="w-4 h-4" />
                            )}
                        </button>
                    ) : (
                        <div className="w-5" />
                    )}

                    <Input
                        value={field.key}
                        onChange={(e) => updateField(field.id, e.target.value, field.value)}
                        className="w-36 h-8 shrink-0"
                    />

                    {!isObject && (
                        <>
                            <span>:</span>
                            <Input
                                value={field.value}
                                onChange={(e) => updateField(field.id, field.key, e.target.value)}
                                className={`flex-1 h-8 ${isPlaceholder ? "font-mono text-blue-600" : ""}`}
                            />
                            <Button
                                size="sm"
                                variant="outline"
                                onClick={() => setShowPlaceholderFor(field.id)}
                                className="shrink-0"
                            >
                                📋
                            </Button>
                        </>
                    )}

                    {isObject && (
                        <div className="flex gap-1">
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
                    <div className="ml-4 pl-4 mt-1 border-l-2 border-blue-200">
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
                    onSelect={(val) => insertPlaceholder(showPlaceholderFor, val)}
                    onClose={() => setShowPlaceholderFor(null)}
                />
            )}

            <div className="flex gap-2">
                <Button size="sm" onClick={() => addField()}>
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Field
                </Button>
                <Button size="sm" variant="outline" onClick={() => addObject()}>
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Object
                </Button>
            </div>

            {fields.length ? (
                fields.map((field) => renderField(field))
            ) : (
                <div className="text-center py-8 text-slate-500">
                    Klik tambah field
                </div>
            )}

        </div>
    );
}