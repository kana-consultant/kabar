// Hapus import Toast langsung
// import { Toast } from "@kana-consultant/ui-kit";
import type { ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing
import { getAPIKeys } from "@/services/apiKey";
import { generateArticle, generateImage } from "@/services/generate";
import { createDraft, publishDraft } from "@/services/draft";
import type { CreateDraftRequest } from "@/services/draft";

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export async function quickGenerate(
    topic: string,
    selectedProducts: string[],
    article: string,
    imageUrl: string,
    setPublishing: (val: boolean) => void,
    setCurrentDraftId: (val: string | null) => void,
    resetForm: () => void,
    toast: ToastContextType //   Tambahkan parameter toast
) {
    if (!topic) {
        toast.error("Generate artikel terlebih dahulu"); //   Ganti toast.error dengan toast.error
        return;
    }

    if (selectedProducts.length === 0) {
        toast.error("Pilih minimal 1 produk"); //   Ganti toast.error dengan toast.error
        return;
    }

    setPublishing(true);

    try {
        const data = await getAPIKeys();
        const availableModels = data as any[];

        if (availableModels.length < 2) {
            toast.error("Model belum cukup"); //   Ganti toast.error dengan toast.error
            return;
        }

        const articleModel = availableModels.find(model => model.service === 'text');
        const imageModel = availableModels.find(model => model.service === 'image');

        if (!articleModel) {
            toast.error("Model untuk generate artikel tidak ditemukan"); //   Ganti toast.error dengan toast.error
            return;
        }

        if (!imageModel) {
            toast.error("Model untuk generate gambar tidak ditemukan"); //   Ganti toast.error dengan toast.error
            return;
        }

        const articleResult = await generateArticle({
            topic: topic,
            modelId: articleModel.id,
            tone: "professional",
            length: "medium",
            language: "id",
        });

        const generatedArticle = articleResult?.content || article;
        const imagePrompt = `Ilustrasi tentang ${topic}`;

        await delay(2000);

        const requestImage = {
            prompt: imagePrompt,
            modelId: imageModel.id
        };

        const imageResult = await generateImage(requestImage);
        const generatedImageUrl = imageResult?.imageUrl || imageUrl;

        await delay(1000);

        const draftData = {
            title: topic,
            topic: topic,
            article: generatedArticle,
            image_url: generatedImageUrl || undefined,
            image_prompt: imagePrompt,
            target_products: selectedProducts,
        };

        const draft = await createDraft(draftData as CreateDraftRequest);

        await delay(1500);

        await publishDraft(draft.id, null);

        toast.success("Publish berhasil"); //   Ganti toast.success dengan toast.success
        resetForm();
        return draft.id;

    } catch (error: any) {
        console.error(error);
        toast.error("Gagal melakukan quick generate", { //   Ganti toast.error dengan toast.error
            description: error?.message || "Terjadi kesalahan",
        });
        return undefined;
    } finally {
        setPublishing(false);
    }
}