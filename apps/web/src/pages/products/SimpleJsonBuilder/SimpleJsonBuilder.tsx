import { useState, useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, Trash2, GripVertical, ChevronRight, ChevronDown } from "lucide-react";
import { type Field, type SimpleJsonBuilderProps } from "@/types/JsonBuilder";
import { genId, jsonToFields, fieldsToJson, getNextFieldNumber, getNextObjectNumber } from "@/utils/SimpleJsonBuilder";

export function SimpleJsonBuilder({ value, onChange }: SimpleJsonBuilderProps) {
    const [fields, setFields] = useState<Field[]>([]);
    const isInitialLoad = useRef(true);
    const isInternalUpdate = useRef(false);

    // Fungsi untuk parse value (bisa string JSON atau object)
    const parseValue = (val: any): any => {
        if (typeof val === 'string') {
            try {
                return JSON.parse(val);
            } catch (e) {
                console.error("Failed to parse JSON string:", e);
                return {};
            }
        }
        return val;
    };

    // Load data dari value HANYA saat pertama kali atau saat value dari luar berubah (bukan dari internal edit)
    useEffect(() => {
        // Jangan reset jika ini adalah update dari internal (user edit)
        if (isInternalUpdate.current) {
            isInternalUpdate.current = false;
            return;
        }

        const parsedValue = parseValue(value);
        if (parsedValue && typeof parsedValue === "object") {
            const hasData = Object.keys(parsedValue).length > 0;
            if (hasData || isInitialLoad.current) {
                const newFields = jsonToFields(parsedValue);
                setFields(newFields);
                isInitialLoad.current = false;
            }
        } else {
            if (isInitialLoad.current) {
                setFields([]);
                isInitialLoad.current = false;
            }
        }
    }, [value]);

    const addField = (parentId?: string) => {
        const nextNumber = getNextFieldNumber(fields);
        const newFieldName = `field${nextNumber}`;
        
        const newField: Field = {
            id: genId(),
            key: newFieldName,
            value: "",
            type: "field",
            children: [],
            expanded: true,
        };

        isInternalUpdate.current = true;

        if (parentId) {
            const addToParent = (items: Field[]): Field[] => {
                return items.map(item => {
                    if (item.id === parentId) {
                        return { ...item, children: [...item.children, newField] };
                    }
                    if (item.children.length > 0) {
                        return { ...item, children: addToParent(item.children) };
                    }
                    return item;
                });
            };
            const newFields = addToParent(fields);
            setFields(newFields);
            onChange(JSON.stringify(fieldsToJson(newFields), null, 2));
        } else {
            const newFields = [...fields, newField];
            setFields(newFields);
            onChange(JSON.stringify(fieldsToJson(newFields), null, 2));
        }
    };

    const addObject = (parentId?: string) => {
        const nextNumber = getNextObjectNumber(fields);
        const newObjectName = `object${nextNumber}`;
        
        const newObject: Field = {
            id: genId(),
            key: newObjectName,
            value: "",
            type: "object",
            children: [],
            expanded: true,
        };

        isInternalUpdate.current = true;

        if (parentId) {
            const addToParent = (items: Field[]): Field[] => {
                return items.map(item => {
                    if (item.id === parentId) {
                        return { ...item, children: [...item.children, newObject] };
                    }
                    if (item.children.length > 0) {
                        return { ...item, children: addToParent(item.children) };
                    }
                    return item;
                });
            };
            const newFields = addToParent(fields);
            setFields(newFields);
            onChange(JSON.stringify(fieldsToJson(newFields), null, 2));
        } else {
            const newFields = [...fields, newObject];
            setFields(newFields);
            onChange(JSON.stringify(fieldsToJson(newFields), null, 2));
        }
    };

    const updateField = (id: string, key: string, value: string) => {
        const updateRecursive = (items: Field[]): Field[] => {
            return items.map(item => {
                if (item.id === id) {
                    return { ...item, key, value };
                }
                if (item.children.length > 0) {
                    return { ...item, children: updateRecursive(item.children) };
                }
                return item;
            });
        };
        
        isInternalUpdate.current = true;
        const newFields = updateRecursive(fields);
        setFields(newFields);
        onChange(JSON.stringify(fieldsToJson(newFields), null, 2));
    };

    const deleteField = (id: string) => {
        const deleteRecursive = (items: Field[]): Field[] => {
            return items.filter(item => {
                if (item.id === id) return false;
                if (item.children.length > 0) {
                    item.children = deleteRecursive(item.children);
                }
                return true;
            });
        };
        
        const newFields = deleteRecursive(fields);
        
        // Re-generate nama field setelah delete (reorder)
        const reorderFields = (items: Field[], startNumber: number = 1): { items: Field[], nextNumber: number } => {
            let currentNumber = startNumber;
            const newItems: Field[] = [];
            
            for (const item of items) {
                if (item.type === "field") {
                    const newItem = { ...item, key: `field${currentNumber}` };
                    if (item.children.length > 0) {
                        const result = reorderFields(item.children, 1);
                        newItem.children = result.items;
                    }
                    newItems.push(newItem);
                    currentNumber++;
                } else if (item.type === "object") {
                    const newItem = { ...item, key: `object${currentNumber}` };
                    if (item.children.length > 0) {
                        const result = reorderFields(item.children, 1);
                        newItem.children = result.items;
                    }
                    newItems.push(newItem);
                    currentNumber++;
                } else {
                    newItems.push(item);
                    currentNumber++;
                }
            }
            
            return { items: newItems, nextNumber: currentNumber };
        };
        
        isInternalUpdate.current = true;
        const { items: reorderedFields } = reorderFields(newFields);
        setFields(reorderedFields);
        onChange(JSON.stringify(fieldsToJson(reorderedFields), null, 2));
    };

    const toggleExpand = (id: string) => {
        const toggleRecursive = (items: Field[]): Field[] => {
            return items.map(item => {
                if (item.id === id) {
                    return { ...item, expanded: !item.expanded };
                }
                if (item.children.length > 0) {
                    return { ...item, children: toggleRecursive(item.children) };
                }
                return item;
            });
        };
        setFields(toggleRecursive(fields));
    };

    const renderField = (field: Field, level: number = 0) => {
        const indent = level * 20;
        const isObject = field.type === "object";

        return (
            <div key={field.id} style={{ marginLeft: indent }} className="mb-2">
                <div className="flex items-center gap-2 p-2 border rounded-lg bg-white dark:bg-slate-900">
                    <div className="cursor-move text-slate-400">
                        <GripVertical className="h-4 w-4" />
                    </div>

                    {isObject && (
                        <button onClick={() => toggleExpand(field.id)} className="w-5">
                            {field.expanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                        </button>
                    )}
                    {!isObject && <div className="w-5" />}

                    <Input
                        placeholder="nama field"
                        value={field.key}
                        onChange={(e) => updateField(field.id, e.target.value, field.value)}
                        className="w-48 h-8 text-sm"
                    />

                    {!isObject && (
                        <>
                            <span className="text-slate-400">:</span>
                            <Input
                                placeholder="nilai"
                                value={field.value}
                                onChange={(e) => updateField(field.id, field.key, e.target.value)}
                                className="flex-1 h-8 text-sm"
                            />
                        </>
                    )}

                    {isObject && (
                        <div className="flex gap-1">
                            <Button size="sm" variant="outline" onClick={() => addField(field.id)}>
                                <Plus className="h-3 w-3 mr-1" />
                                Field
                            </Button>
                            <Button size="sm" variant="outline" onClick={() => addObject(field.id)}>
                                <Plus className="h-3 w-3 mr-1" />
                                Object
                            </Button>
                        </div>
                    )}

                    <Button size="sm" variant="ghost" onClick={() => deleteField(field.id)}>
                        <Trash2 className="h-4 w-4 text-red-500" />
                    </Button>
                </div>

                {isObject && field.expanded && field.children.length > 0 && (
                    <div className="ml-4 pl-4 border-l-2 border-blue-200 mt-1">
                        {field.children.map(child => renderField(child, level + 1))}
                    </div>
                )}
            </div>
        );
    };

    return (
        <div className="space-y-4">
            <div className="flex gap-2 mb-4">
                <Button size="sm" onClick={() => addField()}>
                    <Plus className="h-3 w-3 mr-1" />
                    + Tambah Field (field{getNextFieldNumber(fields)})
                </Button>
                <Button size="sm" variant="outline" onClick={() => addObject()}>
                    <Plus className="h-3 w-3 mr-1" />
                    + Tambah Object (object{getNextObjectNumber(fields)})
                </Button>
            </div>

            <div className="space-y-2">
                {fields.map(field => renderField(field))}
            </div>

            {fields.length === 0 && (
                <div className="text-center text-slate-500 py-8">
                    Klik tombol di atas untuk mulai membuat JSON
                </div>
            )}
        </div>
    );
}