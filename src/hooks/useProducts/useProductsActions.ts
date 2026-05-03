import { toast } from "sonner";
import { deleteProduct, testConnection, syncProduct } from "@/services/product";

export function useProductsActions(
    loadProducts: () => Promise<void>,
    fetchProducts: () => Promise<void>,
    setShowDeleteDialog: (val: boolean) => void,
    setSelectedProduct: (val: any) => void,
    setTestingId: (val: string | null) => void,
    setSyncingId: (val: string | null) => void
) {
    const handleDelete = async (id: string, name: string) => {
        try {
            await deleteProduct(id);
            toast.success('Produk dihapus', { description: `"${name}" telah dihapus` });
            setShowDeleteDialog(false);
            setSelectedProduct(null);
            await loadProducts();
            await fetchProducts(); // Refresh both
        } catch (error) {
            toast.error('Gagal menghapus produk');
        }
    };

    const handleTestConnection = async (id: string) => {
        setTestingId(id);
        try {
            const result = await testConnection(id);
            if (result.success) {
                toast.success('Koneksi berhasil');
                await loadProducts();
                await fetchProducts();
            } else {
                toast.error(result.message || 'Koneksi gagal');
            }
        } catch (error) {
            toast.error('Gagal menguji koneksi');
        } finally {
            setTestingId(null);
        }
    };

    const handleSync = async (id: string) => {
        setSyncingId(id);
        try {
            const result = await syncProduct(id);
            if (result.success) {
                toast.success(result.message || 'Sinkronisasi berhasil');
                await loadProducts();
                await fetchProducts();
            } else {
                toast.error(result.message || 'Sinkronisasi gagal');
            }
        } catch (error) {
            toast.error('Gagal sinkronisasi');
        } finally {
            setSyncingId(null);
        }
    };

    return {
        handleDelete,
        handleTestConnection,
        handleSync,
    };
}