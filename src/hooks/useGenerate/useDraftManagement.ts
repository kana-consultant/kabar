import { useEffect } from "react";
import { useLocation } from "@tanstack/react-router";
import { toast } from "sonner";
import { createDraft, updateDraft, getDraftById } from "@/services/draft";
import { getHistoryById } from "@/services/history";

interface UseDraftManagementProps {
    setTopic: (topic: string) => void;
    setArticle: (article: string) => void;
    setImageUrl: (url: string) => void;
    setSelectedProducts: (products: string[]) => void;
    setCurrentDraftId: (id: string | null) => void;
}

export function useDraftManagement(props: UseDraftManagementProps) {
    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const editId = searchParams.get("edit") || undefined;
    const TopicId = searchParams.get("topic") || undefined;

    const loadDraft = async () => {
        if (editId) {
            try {
                const draft = await getDraftById(editId);
                if (draft) {
                    props.setTopic(draft.topic);
                    props.setArticle(draft.article);
                    props.setImageUrl(draft.imageUrl || "");
                    props.setSelectedProducts(draft.targetProducts || []);
                    props.setCurrentDraftId(draft.id);
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
                    props.setTopic(draft.topic);
                    props.setArticle(draft.content);
                    props.setImageUrl(draft.imageUrl || "");
                    props.setSelectedProducts(draft.targetProducts || []);
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

    useEffect(() => {
        loadDraft();
    }, [editId, TopicId]);

    const saveAsDraft = async (
        topic: string,
        article: string,
        imageUrl: string,
        selectedProducts: string[],
        currentDraftId: string | null,
        setCurrentDraftId: (id: string) => void
    ) => {
        if (!article) {
            toast.error("Generate artikel terlebih dahulu");
            return false;
        }

        const draftData = {
            Title: topic,
            Topic: topic,
            Article: article,
            ImageUrl: imageUrl || undefined,
            ImagePrompt: topic,
            TargetProducts: selectedProducts,
            hasImage: !!imageUrl,
        };

        try {
            if (currentDraftId) {
                await updateDraft(currentDraftId, { ...draftData, status: 'draft' });
                toast.success("Draft diperbarui!");
            } else {
                const response = await createDraft(draftData);
                setCurrentDraftId(response.id);
                toast.success("Draft tersimpan!");
            }
            return true;
        } catch (error) {
            toast.error("Gagal menyimpan draft");
            return false;
        }
    };

    return {
        saveAsDraft,
        editId,
        TopicId
    };
}