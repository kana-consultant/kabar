import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { AdapterConfig } from "@/types/product";

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
                        onChange={(e) => onChange({ endpointPath: e.target.value })}
                        placeholder="/wp-json/wp/v2/posts"
                    />
                </div>
                <div>
                    <Label>HTTP Method</Label>
                    <Select
                        value={config.httpMethod}
                        onValueChange={(v) => onChange({ customHeaders: v as any })}
                    >
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
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
                        onChange={(e) => {
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