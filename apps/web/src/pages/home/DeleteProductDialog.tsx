import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogClose,
  Button,
} from "@kana-consultant/ui-kit";
import type { Product } from "@/services/product";

interface DeleteProductDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    product: Product | null;
    onDelete: (id: string, name: string) => void;
}

export function DeleteProductDialog({ open, onOpenChange, product, onDelete }: DeleteProductDialogProps) {
    if (!product) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Hapus Produk?</DialogTitle>
                    <DialogDescription>
                        Apakah Anda yakin ingin menghapus produk "{product.name}"?
                        <br />
                        Semua konfigurasi mapping akan hilang. Tindakan ini tidak dapat dibatalkan.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <DialogClose asChild>
                        <Button variant="outline">Batal</Button>
                    </DialogClose>
                    <Button 
                        onClick={() => onDelete(product.id, product.name)} 
                        variant="destructive"
                    >
                        Hapus
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}