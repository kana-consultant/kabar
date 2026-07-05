import { useEffect } from "react";
// Hapus import Toast langsung
// import { Toast } from "@kana-consultant/ui-kit";
import type { ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing
import { getDraftById } from "@/services/draft";
import { getHistoryById } from "@/services/history";

export function useLoadDraft(
    editId: string | undefined,
    TopicId: string | undefined,
    setTopic: (val: string) => void,
    setArticle: (val: string) => void,
    setImageUrl: (val: string) => void,
    setSelectedProducts: (val: string[]) => void,
    setCurrentDraftId: (val: string | null) => void,
    setSlug: (val: string) => void,
    setKeywords: (val: string[]) => void,
    setExcerpt : (val : string) => void,
    toast: ToastContextType //   Tambahkan parameter toast
) {
    useEffect(() => {
        const LoadDraft = async () => {
            if (editId) {
                try {
                    const draft = await getDraftById(editId);
                    if (draft) {
                        console.log("===============")
                        console.log(draft.data)
                        setTopic(draft.data.topic);
                        setArticle(draft.data.article);
                        setImageUrl(draft.data.image_url || "");
                        setSelectedProducts(draft.data.target_products || []);
                        setCurrentDraftId(draft.data.id as string);
                        setSlug(draft.data.slug);
                        setKeywords(draft.data.keywords as string[]);
                        setExcerpt(draft.data.excerpt)

                        toast.info("Memuat draft", { //   Ganti Toast.info dengan toast.info
                            description: `"${draft.data.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan"); //   Ganti toast.error dengan toast.error
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft"); //   Ganti toast.error dengan toast.error
                }
            } else if (TopicId) {
                try {
                    const draft = await getHistoryById(TopicId);
                    if (draft) {
                        console.log(draft)
                        setTopic(draft.topic);
                        setArticle(draft.content);
                        setImageUrl(draft.imageUrl || "");
                        setSelectedProducts(draft.targetProducts || []);
                        setKeywords(draft.keywords as string[]);
                        setExcerpt(draft.excerpt)
                        
                        toast.info("Memuat draft", { //   Ganti Toast.info dengan toast.info
                            description: `"${draft.title}" siap diedit`,
                        });
                    } else {
                        toast.error("Draft tidak ditemukan"); //   Ganti toast.error dengan toast.error
                    }
                } catch (error) {
                    console.error("Failed to load draft:", error);
                    toast.error("Gagal memuat draft"); //   Ganti toast.error dengan toast.error
                }
            }
        };
        LoadDraft();
    }, [editId, TopicId]); //   Tambahkan toast ke dependencies
}