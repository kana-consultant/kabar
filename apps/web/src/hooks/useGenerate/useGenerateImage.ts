import { toast } from "sonner";
import { generateImage } from "@/services/generate";

export async function generateImageManually(
    articleResponse: any,
    topic: string,
    selectedModelId: string,
    setLoadingImage: (val: boolean) => void,
    setImageUrl: (val: string) => void
) {
    if (!articleResponse?.imagePrompt && !topic) {
        toast.error("Generate artikel terlebih dahulu");
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
        toast.success("Gambar berhasil di-generate!");
    } catch (error) {
        toast.error("Gagal mengenerate gambar");
    } finally {
        setLoadingImage(false);
    }
}