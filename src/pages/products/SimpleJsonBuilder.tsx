import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
    { value: "{title}", label: "Judul Artikel" },
    { value: "{topic}", label: "Topik Artikel" },
    { value: "{content}", label: "Isi Artikel (HTML)" },
    { value: "{excerpt}", label: "Ringkasan" },
    { value: "{image_url}", label: "URL Gambar" },
    { value: "{scheduled_for}", label: "Waktu Terjadwal" },
];

export function SimpleJsonBuilder({
    value,
    onChange,
}: SimpleJsonBuilderProps) {
    const [fields, setFields] = useState<Field[]>([]);
    const [showPlaceholderSelect, setShowPlaceholderSelect] =
        useState<string | null>(null);

    /*
     Guard:
     true = perubahan dari builder sendiri
     false = perubahan dari parent/external
    */
    const [isInternalUpdate, setIsInternalUpdate] =
        useState(false);

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

    /*
      Sync external value ke builder
      Tapi skip kalau update dari typing sendiri
    */
    useEffect(() => {
        if (isInternalUpdate) {
            setIsInternalUpdate(false);
            return;
        }

        const parsed = parseValue(value);

        if (
            parsed &&
            typeof parsed === "object"
        ) {
            setFields(jsonToFields(parsed));
        }
    }, [value]);

    /*
     Debounce update ke parent
    */
    useEffect(() => {
        const t = setTimeout(() => {
            setIsInternalUpdate(true);

            onChange(
                fieldsToJson(fields)
            );
        }, 250);

        return () => clearTimeout(t);
    }, [fields]);

    const updateTree = (
        items: Field[],
        id: string,
        updater: (item: Field) => Field
    ): Field[] => {
        return items.map((item) => {
            if (item.id === id) {
                return updater(item);
            }

            if (item.children?.length) {
                return {
                    ...item,
                    children: updateTree(
                        item.children,
                        id,
                        updater
                    ),
                };
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
            setFields((prev) => [
                ...prev,
                newField,
            ]);
            return;
        }

        setFields((prev) =>
            updateTree(
                prev,
                parentId,
                (item) => ({
                    ...item,
                    children: [
                        ...item.children,
                        newField,
                    ],
                })
            )
        );
    };

    const addObject = (
        parentId?: string
    ) => {
        const newObject: Field = {
            id: genId(),
            key: `object${getNextObjectNumber(
                fields
            )}`,
            value: "",
            type: "object",
            children: [],
            expanded: true,
        };

        if (!parentId) {
            setFields((prev) => [
                ...prev,
                newObject,
            ]);
            return;
        }

        setFields((prev) =>
            updateTree(
                prev,
                parentId,
                (item) => ({
                    ...item,
                    children: [
                        ...item.children,
                        newObject,
                    ],
                })
            )
        );
    };

    const updateField = (
        id: string,
        key: string,
        value: string
    ) => {
        setFields((prev) =>
            updateTree(
                prev,
                id,
                (item) => ({
                    ...item,
                    key,
                    value,
                })
            )
        );
    };

    const deleteRecursive = (
        items: Field[],
        id: string
    ): Field[] =>
        items
            .filter(
                (item) =>
                    item.id !== id
            )
            .map((item) => ({
                ...item,
                children: deleteRecursive(
                    item.children,
                    id
                ),
            }));

    const deleteField = (
        id: string
    ) => {
        setFields((prev) =>
            deleteRecursive(
                prev,
                id
            )
        );
    };

    const toggleExpand = (
        id: string
    ) => {
        setFields((prev) =>
            updateTree(
                prev,
                id,
                (item) => ({
                    ...item,
                    expanded:
                        !item.expanded,
                })
            )
        );
    };

    const insertPlaceholder = (
        fieldId: string,
        placeholder: string
    ) => {
        setFields((prev) =>
            updateTree(
                prev,
                fieldId,
                (item) => ({
                    ...item,
                    value: placeholder,
                })
            )
        );

        setShowPlaceholderSelect(
            null
        );
    };

    const renderField = (
        field: Field,
        level = 0
    ) => {
        const indent =
            level * 20;

        const isObject =
            field.type === "object";

        const isPlaceholder =
            PLACEHOLDERS.some(
                (p) =>
                    p.value ===
                    field.value
            );

        return (
            <div
                key={field.id}
                style={{
                    marginLeft: indent,
                }}
                className="mb-2"
            >
                <div className="flex items-center gap-2 p-2 rounded-lg border bg-white dark:bg-slate-900">

                    <GripVertical className="w-4 h-4 text-slate-400" />

                    {isObject ? (
                        <button
                            onClick={() =>
                                toggleExpand(
                                    field.id
                                )
                            }
                        >
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
                        onChange={(e) =>
                            updateField(
                                field.id,
                                e.target.value,
                                field.value
                            )
                        }
                        className="w-48 h-8"
                    />

                    {!isObject && (
                        <>
                            <span>:</span>

                            <div className="flex-1 relative">

                                <div className="flex gap-1">
                                    <Input
                                        value={
                                            field.value
                                        }
                                        onChange={(e) =>
                                            updateField(
                                                field.id,
                                                field.key,
                                                e.target
                                                    .value
                                            )
                                        }
                                        className={`h-8 ${isPlaceholder
                                                ? "font-mono text-blue-600"
                                                : ""
                                            }`}
                                    />

                                    <Button
                                        size="sm"
                                        variant="outline"
                                        onClick={() =>
                                            setShowPlaceholderSelect(
                                                showPlaceholderSelect ===
                                                    field.id
                                                    ? null
                                                    : field.id
                                            )
                                        }
                                    >
                                        📋
                                    </Button>
                                </div>

                                {showPlaceholderSelect ===
                                    field.id && (
                                        <div className="absolute top-full mt-1 left-0 w-72 rounded-lg border bg-white shadow-lg z-50 dark:bg-slate-800">

                                            <div className="p-2 border-b text-xs font-medium">
                                                Pilih Placeholder
                                            </div>

                                            {PLACEHOLDERS.map(
                                                (p) => (
                                                    <button
                                                        key={
                                                            p.value
                                                        }
                                                        onClick={() =>
                                                            insertPlaceholder(
                                                                field.id,
                                                                p.value
                                                            )
                                                        }
                                                        className="w-full text-left px-3 py-2 hover:bg-slate-100 dark:hover:bg-slate-700"
                                                    >
                                                        <div className="font-mono text-blue-600">
                                                            {p.value}
                                                        </div>

                                                        <div className="text-xs text-slate-500">
                                                            {p.label}
                                                        </div>

                                                    </button>
                                                )
                                            )}

                                        </div>
                                    )}

                            </div>
                        </>
                    )}

                    {isObject && (
                        <div className="flex gap-1">
                            <Button
                                size="sm"
                                variant="outline"
                                onClick={() =>
                                    addField(
                                        field.id
                                    )
                                }
                            >
                                <Plus className="w-3 h-3 mr-1" />
                                Field
                            </Button>

                            <Button
                                size="sm"
                                variant="outline"
                                onClick={() =>
                                    addObject(
                                        field.id
                                    )
                                }
                            >
                                <Plus className="w-3 h-3 mr-1" />
                                Object
                            </Button>
                        </div>
                    )}

                    <Button
                        size="sm"
                        variant="ghost"
                        onClick={() =>
                            deleteField(
                                field.id
                            )
                        }
                    >
                        <Trash2 className="w-4 h-4 text-red-500" />
                    </Button>
                </div>

                {isObject &&
                    field.expanded &&
                    field.children.length >
                    0 && (
                        <div className="ml-4 pl-4 mt-1 border-l-2 border-blue-200">
                            {field.children.map(
                                (child) =>
                                    renderField(
                                        child,
                                        level + 1
                                    )
                            )}
                        </div>
                    )}
            </div>
        );
    };

    return (
        <div className="space-y-4">

            <div className="flex gap-2">
                <Button
                    size="sm"
                    onClick={() =>
                        addField()
                    }
                >
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Field
                </Button>

                <Button
                    size="sm"
                    variant="outline"
                    onClick={() =>
                        addObject()
                    }
                >
                    <Plus className="w-3 h-3 mr-1" />
                    Tambah Object
                </Button>
            </div>

            <div className="p-3 rounded-lg bg-blue-50 dark:bg-blue-950 text-xs">
                <p className="font-medium mb-2">
                    Placeholder tersedia
                </p>

                <div className="grid grid-cols-2 gap-2">
                    {PLACEHOLDERS.map(
                        (p) => (
                            <div
                                key={p.value}
                                className="font-mono text-blue-600"
                            >
                                {p.value}
                            </div>
                        )
                    )}
                </div>
            </div>

            {fields.length ? (
                fields.map((field) =>
                    renderField(field)
                )
            ) : (
                <div className="text-center py-8 text-slate-500">
                    Klik tambah field
                </div>
            )}

        </div>
    );
}