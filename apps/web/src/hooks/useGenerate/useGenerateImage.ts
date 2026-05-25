// Hapus import Toast langsung
// import { Toast } from "@kana-consultant/ui-kit";
import type { ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing
import { generateImage } from "@/services/generate";

export async function generateImageManually(
    articleResponse: any,
    topic: string,
    selectedModelId: string,
    setLoadingImage: (val: boolean) => void,
    setImageUrl: (val: string) => void,
    toast: ToastContextType //   Tambahkan parameter toast
) {
    if (!articleResponse?.imagePrompt && !topic) {
        toast.error("Generate artikel terlebih dahulu"); //   Ganti toast.error dengan toast.error
        return;
    }

    setLoadingImage(true);
    const prompt = articleResponse?.imagePrompt || `Ilustrasi tentang ${topic}`;

    try {
        const response = await generateImage({
            prompt: prompt,
            modelId: selectedModelId
        });
        setImageUrl(response.imageUrl);
        toast.success("Gambar berhasil di-generate!"); //   Ganti toast.success dengan toast.success
    } catch (error) {
        console.error("Error generating image:", error);
        toast.error("Gagal mengenerate gambar"); //   Ganti toast.error dengan toast.error
    } finally {
        setLoadingImage(false);
    }
}