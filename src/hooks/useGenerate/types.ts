import type { Product } from "@/types/product";
import type { AIModel } from "@/services/model";
import type { GenerateArticleResponse } from "@/services/generate";
import type { PublishResponse } from "@/services/draft";

export type PostMode = "instant" | "scheduled" | "draft";
export type Tone = "professional" | "casual" | "friendly" | "formal";
export type ArticleLength = "short" | "medium" | "long";
export type Language = "id" | "en";

export interface GenerateState {
    // Content
    topic: string;
    article: string;
    articleResponse: GenerateArticleResponse | null;
    imageUrl: string;
    
    // Loading states
    loadingArticle: boolean;
    loadingImage: boolean;
    publishing: boolean;
    isPosting: boolean;
    
    // Selection
    selectedProducts: string[];
    currentDraftId: string | null;
    
    // Post configuration
    postMode: PostMode;
    scheduleTime: string;
    scheduleDate: string;
    dailySchedule: boolean;
    dailyTime: string;
    autoGenerateImage: boolean;
    postToAll: boolean;
    
    // Model
    selectedModelId: string;
    
    // Article options
    tone: Tone;
    articleLength: ArticleLength;
    keywords: string[];
    keywordInput: string;
    language: Language;
    
    // SEO Metrics
    seoScore: number | null;
    readabilityScore: number | null;
    wordCount: number | null;
}

export interface DataState {
    models: AIModel[];
    products: Product[];
    productNames: string[];
    loadingModels: boolean;
    productsLoading: boolean;
    productsError: string | null;
}

export interface DialogState {
    showResultDialog: boolean;
    publishResults: PublishResponse | null;
}

export interface GenerateActions {
    generateArticle: () => Promise<void>;
    generateImage: () => Promise<void>;
    saveAsDraft: () => Promise<boolean>;
    saveAsSchedule: () => Promise<boolean>;
    postInstant: () => Promise<boolean>;
    handlePost: () => Promise<void>;
    quickGenerate: () => Promise<string | undefined>;
    resetForm: () => void;
}

export interface UIHandlers {
    handleAddKeyword: () => void;
    handleRemoveKeyword: (keyword: string) => void;
    handleProductToggle: (product: string) => void;
    handleSelectAll: () => void;
    closeResultDialog: () => void;
}