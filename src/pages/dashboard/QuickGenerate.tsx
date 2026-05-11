import { useState, useEffect } from "react";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, Sparkles } from "lucide-react";
import { useGenerate } from "@/hooks/useGenerate";
import { useNavigate } from "@tanstack/react-router";

export function QuickGenerate() {
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const {
        products,
        topic,
        setSelectedProducts,
        selectedProducts,
        setTopic,
        quickGenerate
    } = useGenerate();

    // Auto select first product
    useEffect(() => {
        if (products.length > 0 && selectedProducts.length === 0) {
            setSelectedProducts([products[0].id]);
        }
    }, [products, selectedProducts.length, setSelectedProducts]);

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const draftId = await quickGenerate();
            if (draftId) {
                navigate({ to: "/history" });
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card className="w-full max-w-2xl mx-auto">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Sparkles className="h-5 w-5" />
                    Quick Generate
                </CardTitle>
                <CardDescription>
                    Masukkan topik, sistem akan generate artikel + gambar otomatis
                </CardDescription>
            </CardHeader>

            <CardContent className="space-y-6">
                <div className="space-y-2">
                    <Label htmlFor="topic" className="text-sm font-medium">
                        Topik / Keyword
                    </Label>
                    <Input
                        id="topic"
                        value={topic}
                        placeholder="Contoh: Cara Memilih Sepatu Gunung untuk Pemula"
                        onChange={(e) => setTopic(e.target.value)}
                        className="w-full"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="product" className="text-sm font-medium">
                        Target Produk
                    </Label>
                    <Select
                        value={selectedProducts[0] || ""}
                        onValueChange={(value) => setSelectedProducts([value])}
                    >
                        <SelectTrigger id="product" className="w-full">
                            <SelectValue placeholder="Pilih produk target" />
                        </SelectTrigger>
                        <SelectContent>
                            {products.map((item) => (
                                <SelectItem key={item.id} value={item.id}>
                                    {item.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>
            </CardContent>

            <CardFooter>
                <Button
                    onClick={handleGenerate}
                    disabled={loading || !topic.trim()}
                    className="w-full"
                    size="lg"
                >
                    {loading ? (
                        <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            Generating...
                        </>
                    ) : (
                        <>
                            <Sparkles className="mr-2 h-4 w-4" />
                            Generate Konten
                        </>
                    )}
                </Button>
            </CardFooter>
        </Card>
    );
}