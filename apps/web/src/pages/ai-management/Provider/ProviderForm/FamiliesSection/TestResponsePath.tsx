// pages/admin/ai-management/components/ProviderForm/FamiliesSection/TestResponsePath.tsx

import { useState } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Textarea } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Beaker, Check, X, ArrowRight } from "lucide-react";

interface TestResponsePathProps {
    responsePath: string | null;
    onValidate?: (isValid: boolean) => void;
}

export function TestResponsePath({ responsePath, onValidate }: TestResponsePathProps) {
    const [open, setOpen] = useState(false);
    const [sampleResponse, setSampleResponse] = useState("");
    const [extractedValue, setExtractedValue] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);

    const extractValue = (path: string, json: any): any => {
        // Parse JSON path like "choices[0].message.content"
        const parts = path.split(/\.|\[|\]/).filter(p => p !== "");
        
        let result = json;
        for (const part of parts) {
            if (result === undefined || result === null) {
                return undefined;
            }
            
            if (part.match(/^\d+$/)) {
                // Array index
                result = result[parseInt(part)];
            } else {
                // Object property
                result = result[part];
            }
        }
        
        return result;
    };

    const handleTest = () => {
        if (!responsePath) {
            setError("Response path is empty");
            return;
        }

        if (!sampleResponse.trim()) {
            setError("Please paste a sample response JSON");
            return;
        }

        try {
            const json = JSON.parse(sampleResponse);
            const result = extractValue(responsePath, json);
            
            if (result === undefined) {
                setError("Path not found in JSON");
                setExtractedValue(null);
                onValidate?.(false);
            } else {
                setExtractedValue(typeof result === "object" ? JSON.stringify(result, null, 2) : String(result));
                setError(null);
                onValidate?.(true);
            }
        } catch (err) {
            setError("Invalid JSON format");
            setExtractedValue(null);
            onValidate?.(false);
        }
    };

    const getSampleTemplate = () => {
        if (responsePath?.includes("content")) {
            return `{
  "choices": [
    {
      "message": {
        "content": "Hello! This is a test response."
      }
    }
  ]
}`;
        }
        if (responsePath?.includes("url") || responsePath?.includes("image")) {
            return `{
  "data": [
    {
      "url": "https://example.com/image.png"
    }
  ]
}`;
        }
        return `{
  "result": "Sample value"
}`;
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    disabled={!responsePath}
                >
                    <Beaker className="h-3 w-3 mr-1" />
                    Test Path
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
                <DialogHeader>
                    <DialogTitle>Test Response Path</DialogTitle>
                    <DialogDescription>
                        Paste a sample response JSON to test if the path extracts correctly
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    <div className="space-y-2">
                        <Label>Response Path</Label>
                        <code className="block p-2 bg-muted rounded text-sm font-mono">
                            {responsePath || "(empty)"}
                        </code>
                    </div>

                    <div className="space-y-2">
                        <div className="flex justify-between items-center">
                            <Label>Sample Response JSON</Label>
                            <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                onClick={() => setSampleResponse(getSampleTemplate())}
                            >
                                Load Template
                            </Button>
                        </div>
                        <Textarea
                            value={sampleResponse}
                            onChange={(e) => setSampleResponse(e.target.value)}
                            placeholder={`{
  "choices": [
    {
      "message": {
        "content": "Hello World!"
      }
    }
  ]
}`}
                            rows={8}
                            className="font-mono text-sm"
                        />
                    </div>

                    <Button onClick={handleTest} className="w-full">
                        Extract Value
                        <ArrowRight className="h-4 w-4 ml-2" />
                    </Button>

                    {error && (
                        <div className="p-3 bg-red-50 dark:bg-red-950 border border-red-200 rounded-lg">
                            <div className="flex items-center gap-2 text-red-600">
                                <X className="h-4 w-4" />
                                <span className="text-sm">{error}</span>
                            </div>
                        </div>
                    )}

                    {extractedValue && (
                        <div className="p-3 bg-green-50 dark:bg-green-950 border border-green-200 rounded-lg">
                            <div className="flex items-center gap-2 text-green-600 mb-2">
                                <Check className="h-4 w-4" />
                                <span className="text-sm font-medium">Extracted Value:</span>
                            </div>
                            <pre className="text-sm font-mono bg-white dark:bg-black p-2 rounded overflow-x-auto">
                                {extractedValue}
                            </pre>
                        </div>
                    )}

                    <div className="text-xs text-muted-foreground border-t pt-3">
                        <p className="font-medium mb-1">How to write JSON path:</p>
                        <ul className="list-disc list-inside space-y-0.5">
                            <li><code className="text-xs">field</code> - access object property</li>
                            <li><code className="text-xs">array[0]</code> - access array index</li>
                            <li><code className="text-xs">parent.child</code> - nested properties</li>
                            <li><code className="text-xs">choices[0].message.content</code> - full path example</li>
                        </ul>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}