import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface RawJsonMappingProps {
    value: any; // bisa array atau string
    onChange: (value: any) => void;
}

export function RawJsonMapping({ value, onChange }: RawJsonMappingProps) {
    // Jika value adalah array (dari form), konversi ke JSON string
    const getInitialJson = () => {
        if (typeof value === 'string') {
            return value;
        }
        return JSON.stringify(value, null, 2);
    };

    const [jsonText, setJsonText] = useState(getInitialJson);
    const [error, setError] = useState("");

    useEffect(() => {
        setJsonText(getInitialJson());
    }, [value]);

    const handleSave = () => {
        // Simpan PERSIS seperti yang user tulis, tanpa parsing!
        onChange(jsonText);
        setError("");
    };

    const handleFormat = () => {
        try {
            const parsed = JSON.parse(jsonText);
            setJsonText(JSON.stringify(parsed, null, 2));
            setError("");
        } catch (err) {
            setError("JSON tidak valid, tidak bisa format");
        }
    };

   

    return (
        <div className="space-y-4">
            {error && (
                <div className="bg-red-50 text-red-600 p-2 rounded text-sm">
                    ⚠️ {error}
                </div>
            )}

            <div>
                <label className="block text-sm font-medium mb-2">
                    Raw JSON Mapping (bebas apapun)
                </label>
                <textarea
                    value={jsonText}
                    onChange={(e) => setJsonText(e.target.value)}
                    className="w-full h-96 rounded border p-3 font-mono text-sm dark:bg-slate-900"
                    placeholder="Input JSON bebas sesuai kebutuhan API Anda"
                />
            </div>

            <div className="flex gap-2">
                <Button onClick={handleSave}>
                    💾 Simpan JSON Ini
                </Button>
                <Button variant="outline" onClick={handleFormat}>
                    Format Ulang (jika valid)
                </Button>
            </div>

            <div className="bg-slate-100 p-3 rounded text-xs dark:bg-slate-800">
                <p className="font-medium mb-1">📝 JSON bebas, contoh:</p>
                <pre>{`{
  "post": {
    "title": "{title}",
    "content": "{body}"
  },
  "meta": {
    "image": "{imageUrl}"
  }
}`}</pre>
                <p className="mt-2 text-slate-500">
                    ✅ Simpan JSON APAPUN. Tidak akan diubah atau divalidasi berlebihan.
                </p>
            </div>
        </div>
    );
}