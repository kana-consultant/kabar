import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, Plus } from "lucide-react";

export  function ProductHeader({
    searchQuery,
    setSearchQuery
}: any) {

    const navigate = useNavigate();

    return (
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">

            <div>
                <h2 className="text-2xl font-bold">Produk</h2>
                <p className="text-sm text-muted-foreground">
                    Kelola produk dan koneksi API
                </p>
            </div>

            <div className="flex gap-2">

                <div className="relative">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2" />

                    <Input
                        placeholder="Cari produk..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-9 w-64"
                    />
                </div>

                <Button onClick={() => navigate({ to: "/products/add" })}>
                    <Plus className="mr-2 h-4 w-4" />
                    Tambah
                </Button>

            </div>

        </div>
    );
}