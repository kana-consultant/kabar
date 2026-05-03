import { toast } from "sonner";
import { getAPIKeys } from "@/services/apiKey";
import { generateArticle, generateImage } from "@/services/generate";
import { createDraft, publishDraft } from "@/services/draft";

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export async function quickGenerate(
    topic: string,
    selectedProducts: string[],
    article: string,
    imageUrl: string,
    setPublishing: (val: boolean) => void,
    setCurrentDraftId: (val: string | null) => void,
    resetForm: () => void
) {
    if (!topic) {
        toast.error("Generate artikel terlebih dahulu");
        return;
    }

    if (selectedProducts.length === 0) {
        toast.error("Pilih minimal 1 produk");
        return;
    }

    setPublishing(true);

    try {
        const data = await getAPIKeys();
        const availableModels = data as any[];

        if (availableModels.length < 2) {
            toast.error("Model belum cukup");
            return;
        }

        const articleModel = availableModels.find(model => model.service === 'text');
        const imageModel = availableModels.find(model => model.service === 'image');

        if (!articleModel) {
            toast.error("Model untuk generate artikel tidak ditemukan");
            return;
        }

        if (!imageModel) {
            toast.error("Model untuk generate gambar tidak ditemukan");
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
            Title: topic,
            Topic: topic,
            Article: generatedArticle,
            ImageURL: generatedImageUrl || undefined,
            ImagePrompt: imagePrompt,
            TargetProducts: selectedProducts,
        };

        const draft = await createDraft(draftData);

        await delay(1500);

        await publishDraft(draft.id);

        toast.success("Publish berhasil");
        resetForm();
        return draft.id;

    } catch (error: any) {
        console.error(error);
        toast.error("Gagal melakukan quick generate", {
            description: error?.message || "Terjadi kesalahan",
        });
        return undefined;
    } finally {
        setPublishing(false);
    }
}