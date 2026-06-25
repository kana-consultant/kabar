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
import { Button } from "@kana-consultant/ui-kit";
import { Wifi } from "lucide-react";
import type { Product } from "@/services/product";

interface ProductBasicInfoProps {
    product: Partial<Product>;
    onUpdate: (updates: Partial<Product>) => void;
    onTestConnection: () => void;
    isTesting: boolean;
}

export function ProductBasicInfo({ product, onUpdate, onTestConnection, isTesting }: ProductBasicInfoProps) {
    return (
        <Card className="shadow-sm">
            <CardHeader className="pb-3">
                <CardTitle className="text-lg">Informasi Dasar</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="name" className="text-sm font-medium">
                        Nama Produk <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="name"
                        value={product.name || ''}  // ← FIX: default empty string
                        onChange={(e) => onUpdate({ name: e.target.value })}
                        placeholder="Contoh: Toko Saya"
                        className="h-10"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="platform" className="text-sm font-medium">
                        Platform
                    </Label>
                    <Select
                        value={product.platform || ''}  // ← FIX: default empty string
                        onValueChange={(v) => onUpdate({ platform: v as any })}
                    >
                        <SelectTrigger className="h-10">
                            <SelectValue placeholder="Pilih platform" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="wordpress">📝 WordPress</SelectItem>
                            <SelectItem value="shopify">🛍️ Shopify</SelectItem>
                            <SelectItem value="custom">🔧 Custom API</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="apiEndpoint" className="text-sm font-medium">
                        Base URL <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="apiEndpoint"
                        value={product.api_endpoint || ''}  // ← FIX: default empty string
                        onChange={(e) => onUpdate({ api_endpoint: e.target.value })}
                        placeholder="https://domain.com/wp-json/wp/v2"
                        className="h-10 font-mono text-sm"
                    />
                    <p className="text-xs text-slate-400">
                        Base URL tanpa endpoint path (contoh: /posts, /media akan ditambahkan di adapter)
                    </p>
                </div>

               

                {/* Tombol Test Koneksi */}
                <div className="pt-2">
                    <Button
                        type="button"
                        onClick={onTestConnection}
                        disabled={isTesting || !product.api_endpoint}  // ← FIX: disable if no endpoint
                        className="w-full"
                        variant="outline"
                    >
                        <Wifi className="mr-2 h-4 w-4" />
                        {isTesting ? "Menguji Koneksi..." : "Test Koneksi API"}
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}