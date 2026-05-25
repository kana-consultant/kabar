import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@kana-consultant/ui-kit";
import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import type { AdapterConfig } from "@/services/product";

interface ProductApiConfigProps {
    config: AdapterConfig;
    onChange: (updates: Partial<AdapterConfig>) => void;
}

export function ProductApiConfig({ config, onChange }: ProductApiConfigProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Konfigurasi API</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div>
                    <Label>Endpoint Path</Label>
                    <Input
                        value={config.endpointPath}
                        onChange={(e : any) => onChange({ endpointPath: e.target.value })}
                        placeholder="/wp-json/wp/v2/posts"
                    />
                </div>
                <div>
                    <Label>HTTP Method</Label>
                    <Select
                        value={config.httpMethod}
                        onValueChange={(v : any) => onChange({ customHeaders: v as any })}
                    >
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="GET">GET</SelectItem>
                            <SelectItem value="POST">POST</SelectItem>
                            <SelectItem value="PUT">PUT</SelectItem>
                            <SelectItem value="PATCH">PATCH</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <div>
                    <Label>Custom Headers</Label>
                    <textarea
                        value={JSON.stringify(config.customHeaders, null, 2)}
                        onChange={(e : any) => {
                            try {
                                const customHeaders = JSON.parse(e.target.value);
                                onChange({ customHeaders });
                            } catch (err) {
                                // Invalid JSON, ignore
                            }
                        }}
                        className="w-full rounded-md border p-2 font-mono text-sm dark:bg-slate-900"
                        rows={4}
                    />
                    <p className="mt-1 text-xs text-slate-500">JSON format</p>
                </div>
            </CardContent>
        </Card>
    );
}