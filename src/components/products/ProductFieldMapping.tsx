import { useState, useEffect } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { SimpleJsonBuilder } from "./SimpleJsonBuilder";
import { Eye } from "lucide-react";

interface ProductFieldMappingProps {
    fieldMapping: any; // Bisa string atau object
    onChange: (value: string) => void;
}

export function ProductFieldMapping({ fieldMapping, onChange }: ProductFieldMappingProps) {
    // Konversi input ke object untuk preview
    const getObjectValue = (input: any): any => {
        if (typeof input === 'string') {
            try {
                return JSON.parse(input);
            } catch (e) {
                return {};
            }
        }
        return input || {};
    };

    const [rawJson, setRawJson] = useState(() => {
        const obj = getObjectValue(fieldMapping);
        return JSON.stringify(obj, null, 2);
    });

    const [activeTab, setActiveTab] = useState<"form" | "raw">(() => {
        const obj = getObjectValue(fieldMapping);
        if (typeof obj === "object" && obj !== null) {
            return "form";
        }
        return "raw";
    });

    const [previewObject, setPreviewObject] = useState<any>(() => getObjectValue(fieldMapping));

    useEffect(() => {
        const obj = getObjectValue(fieldMapping);
        setRawJson(JSON.stringify(obj, null, 2));
        setPreviewObject(obj);
    }, [fieldMapping]);

    const handleBuilderChange = (newValue: any) => {
        // newValue adalah object dari SimpleJsonBuilder
        const jsonString = JSON.stringify(newValue, null, 2);
        console.log("🔄 Builder onChange (object):", newValue);
        console.log("🔄 Builder onChange (string):", jsonString);
        onChange(jsonString);
        setPreviewObject(newValue);
    };

    const handleRawJsonChange = (value: string) => {
        setRawJson(value);
        try {
            const parsed = JSON.parse(value);
            setPreviewObject(parsed);
        } catch (e) {
            setPreviewObject({});
        }
    };

    const saveRawJson = () => {
        console.log("💾 Saving raw JSON:", rawJson);
        onChange(rawJson);
        try {
            const parsed = JSON.parse(rawJson);
            setPreviewObject(parsed);
            if (typeof parsed === "object") {
                setActiveTab("form");
            }
        } catch (e) { }
    };

    // Value untuk SimpleJsonBuilder (harus object)
    const builderValue = getObjectValue(fieldMapping);

    return (
        <div className="border rounded-lg p-4 space-y-4">
            <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as "form" | "raw")} className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="form">Drag & Drop</TabsTrigger>
                    <TabsTrigger value="raw">Raw JSON</TabsTrigger>
                </TabsList>

                <TabsContent value="form" className="pt-4">
                    <SimpleJsonBuilder
                        value={builderValue}
                        onChange={handleBuilderChange}
                    />
                </TabsContent>

                <TabsContent value="raw" className="pt-4 space-y-3">
                    <textarea
                        value={rawJson}
                        onChange={(e) => handleRawJsonChange(e.target.value)}
                        className="w-full h-96 rounded border p-4 font-mono text-sm dark:bg-slate-900"
                        placeholder="Input JSON bebas..."
                    />
                    <div className="flex justify-end">
                        <Button onClick={saveRawJson}>💾 Simpan JSON</Button>
                    </div>
                </TabsContent>
            </Tabs>

            {/* PREVIEW - dalam format OBJECT */}
            <div className="border-t pt-4">
                <div className="flex items-center gap-2 mb-2">
                    <Eye className="h-4 w-4 text-slate-500" />
                    <h4 className="text-sm font-medium">Preview (Object)</h4>
                </div>
                <pre className="w-full h-64 overflow-auto rounded-lg bg-slate-100 p-4 text-xs font-mono dark:bg-slate-900">
                    {JSON.stringify(previewObject, null, 2)}
                </pre>
            </div>
        </div>
    );
}