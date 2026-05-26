//   Hapus import Toast langsung, nanti akan menerima toast sebagai parameter
import type {  ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing

export function handleAddKeyword(
    keywordInput: string,
    keywords: string[],
    setKeywords: (val: string[]) => void,
    setKeywordInput: (val: string) => void
) {
    if (keywordInput.trim() && !keywords.includes(keywordInput.trim())) {
        setKeywords([...keywords, keywordInput.trim()]);
        setKeywordInput("");
    }
}

export function handleRemoveKeyword(
    keyword: string,
    keywords: string[],
    setKeywords: (val: string[]) => void
) {
    setKeywords(keywords.filter(k => k !== keyword));
}

export function handleProductToggle(
    product: string,
    selectedProducts: string[],
    postToAll: boolean,
    setSelectedProducts: (val: string[]) => void
) {
    if (postToAll) return;
    setSelectedProducts(
        selectedProducts.includes(product) 
            ? selectedProducts.filter(p => p !== product) 
            : [...selectedProducts, product]
    );
}

//   Tambahkan parameter toast
export function handleSelectAll(
    postToAll: boolean,
    productNames: string[],
    setPostToAll: (val: boolean) => void,
    setSelectedProducts: (val: string[]) => void,
    toast: ToastContextType //   Terima toast sebagai parameter
) {
    if (postToAll) {
        setPostToAll(false);
        setSelectedProducts([]);
        toast.info("Semua produk dibatalkan"); //   Pakai toast dari parameter
    } else {
        setPostToAll(true);
        setSelectedProducts([...productNames]);
        toast.success("Semua produk dipilih"); //   Pakai toast dari parameter
    }
}

//   Tambahkan parameter toast untuk resetForm
export function resetForm(
    setTopic: (val: string) => void,
    setArticle: (val: string) => void,
    setArticleResponse: (val: any) => void,
    setImageUrl: (val: string) => void,
    setSelectedProducts: (val: string[]) => void,
    setPostMode: (val: "instant" | "scheduled" | "draft") => void,
    setCurrentDraftId: (val: string | null) => void,
    setKeywords: (val: string[]) => void,
    setKeywordInput: (val: string) => void,
    setSeoScore: (val: number | null) => void,
    setReadabilityScore: (val: number | null) => void,
    setWordCount: (val: number | null) => void,
    setTone: (val: any) => void,
    setArticleLength: (val: any) => void,
    setLanguage: (val: any) => void,
    toast?: ToastContextType //   Optional toast parameter (opsional)
) {
    setTopic("");
    setArticle("");
    setArticleResponse(null);
    setImageUrl("");
    setSelectedProducts([]);
    setPostMode("instant");
    setCurrentDraftId(null);
    setKeywords([]);
    setKeywordInput("");
    setSeoScore(null);
    setReadabilityScore(null);
    setWordCount(null);
    setTone("professional");
    setArticleLength("medium");
    setLanguage("id");
    
    // Optional: show success toast if toast is provided
    if (toast) {
        toast.success("Form berhasil direset");
    }
}