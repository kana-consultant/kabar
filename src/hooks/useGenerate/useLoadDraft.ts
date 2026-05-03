import { useEffect } from "react";
import { toast } from "sonner";
import { getDraftById } from "@/services/draft";
import { getHistoryById } from "@/services/history";

export function useLoadDraft(
    editId: string | undefined,
    TopicId: string | undefined,
    setTopic: (val: string) => void,
    setArticle: (val: string) => void,
    setImageUrl: (val: string) => void,
    setSelectedProducts: (val: string[]) => void,
    setCurrentDraftId: (val: string | null) => void
) {
    useEffect(() => {
        const LoadDraft = async () => {
            if (editId) {
                try {
                    const draft = await getDraftById(editId);
                    if (draft) {
                        setTopic(draft.topic);
                        setArticle(draft.article);
                        setImageUrl(draft.imageUrl || "");
                        setSelectedProducts(draft.targetProducts || []);
                        setCurrentDraftId(draft.id);
                        toast.info("Memuat draft", {
                            description: `"${draft.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan");
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft");
                }
            } else if (TopicId) {
                try {
                    const draft = await getHistoryById(TopicId);
                    if (draft) {
                        setTopic(draft.topic);
                        setArticle(draft.content);
                        setImageUrl(draft.imageUrl || "");
                        setSelectedProducts(draft.targetProducts || []);
                        toast.info("Memuat draft", {
                            description: `"${draft.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan");
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft");
                }
            }
        };
        LoadDraft();
    }, [editId, TopicId]);
}