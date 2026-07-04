import { useState, useEffect } from "react";
import {type Paragraph, LABEL_MAP, isListType, wrapParagraph, stripQuillDefaultWrapper, nextParagraphId } from "./types";
const BLOCK_ELEMENTS = ['P', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6', 'BLOCKQUOTE', 'UL', 'OL', 'PRE'];

const getTypeByTag = (tag: string): Paragraph['type'] => {
    const typeMap: Record<string, Paragraph['type']> = {
        h1: 'h1', h2: 'h2', h3: 'h3', h4: 'h4', h5: 'h5', h6: 'h6',
        blockquote: 'blockquote', ul: 'ul', ol: 'ol', pre: 'pre',
    };
    return typeMap[tag.toLowerCase()] || 'p';
};

export const parseArticleToParagraphs = (html: string): Paragraph[] => {
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = html;

    const paragraphs: Paragraph[] = [];

    const processNode = (node: Node) => {
        if (node.nodeType === Node.ELEMENT_NODE) {
            const el = node as HTMLElement;
            const tag = el.tagName;

            if (BLOCK_ELEMENTS.includes(tag) && el.textContent?.trim()) {
                const type = getTypeByTag(tag);
                const innerHTML = isListType(type) ? el.outerHTML : el.innerHTML;

                paragraphs.push({
                    id: nextParagraphId(),
                    innerHTML,
                    type,
                    label: LABEL_MAP[type],
                });
            } else if (el.children.length > 0) {
                Array.from(el.children).forEach((child) => processNode(child));
            }
        }
    };

    Array.from(tempDiv.children).forEach((child) => processNode(child));

    if (paragraphs.length === 0) {
        const text = tempDiv.textContent || '';
        paragraphs.push({
            id: nextParagraphId(),
            innerHTML: text.trim() || '<br>',
            type: 'p',
            label: LABEL_MAP.p,
        });
    }

    return paragraphs;
};

interface UseArticleEditorParams {
    article: string;
    onArticleUpdate?: (newArticle: string) => void;
}

export function useArticleEditor({ article, onArticleUpdate }: UseArticleEditorParams) {
    const [isEditing, setIsEditing] = useState(false);
    const [paragraphs, setParagraphs] = useState<Paragraph[]>([]);
    const [isSaving, setIsSaving] = useState(false);

    useEffect(() => {
        if (article) {
            setParagraphs(parseArticleToParagraphs(article));
        } else {
            setParagraphs([]);
        }
    }, [article]);

    const toggleEdit = () => {
        if (isEditing) {
            setParagraphs(parseArticleToParagraphs(article));
        }
        setIsEditing(!isEditing);
    };

    const handleParagraphChange = (paragraphId: string, rawValue: string) => {
        setParagraphs(prev =>
            prev.map(p => {
                if (p.id !== paragraphId) return p;
                const innerHTML = isListType(p.type) ? rawValue : stripQuillDefaultWrapper(rawValue);
                return { ...p, innerHTML };
            })
        );
    };

    const changeParagraphType = (paragraphId: string, newType: Paragraph['type']) => {
        setParagraphs(prev =>
            prev.map(p => {
                if (p.id !== paragraphId) return p;

                let innerHTML = p.innerHTML;
                const wasList = isListType(p.type);
                const willBeList = isListType(newType);

                if (wasList && !willBeList) {
                    const temp = document.createElement('div');
                    temp.innerHTML = p.innerHTML;
                    innerHTML = temp.textContent?.trim() || '<br>';
                } else if (!wasList && willBeList) {
                    const temp = document.createElement('div');
                    temp.innerHTML = p.innerHTML;
                    const text = temp.textContent?.trim();
                    innerHTML = `<${newType}><li>${text || '<br>'}</li></${newType}>`;
                } else if (wasList && willBeList && p.type !== newType) {
                    const temp = document.createElement('div');
                    temp.innerHTML = p.innerHTML;
                    const liItems = temp.children[0]?.innerHTML || '';
                    innerHTML = `<${newType}>${liItems}</${newType}>`;
                }

                return { ...p, type: newType, label: LABEL_MAP[newType], innerHTML };
            })
        );
    };

    const addParagraph = (afterIndex: number, type: Paragraph['type'] = 'p') => {
        const newParagraph: Paragraph = {
            id: nextParagraphId(),
            innerHTML: isListType(type) ? `<${type}><li><br></li></${type}>` : '<br>',
            type,
            label: LABEL_MAP[type],
        };

        setParagraphs(prev => {
            const newParagraphs = [...prev];
            newParagraphs.splice(afterIndex + 1, 0, newParagraph);
            return newParagraphs;
        });
    };

    const removeParagraph = (paragraphId: string, onRemoved?: (id: string) => void) => {
        setParagraphs(prev => {
            if (prev.length <= 1) return prev;
            return prev.filter(p => p.id !== paragraphId);
        });
        onRemoved?.(paragraphId);
    };

    const handleSave = async () => {
        if (!onArticleUpdate) {
            console.warn("onArticleUpdate tidak tersedia");
            return;
        }

        setIsSaving(true);
        try {
            const combinedHTML = paragraphs.map(p => wrapParagraph(p)).join('\n');

            const cleanHTML = combinedHTML
                .replace(/<p>\s*<br\s*\/?>\s*<\/p>/g, '')
                .replace(/<p>\s*<\/p>/g, '')
                .replace(/<(h[1-6]|blockquote|pre|ul|ol)>\s*<br\s*\/?>\s*<\/(h[1-6]|blockquote|pre|ul|ol)>/g, '');

            await onArticleUpdate(cleanHTML);
            setIsEditing(false);
        } catch (error) {
            console.error("Gagal menyimpan:", error);
        } finally {
            setIsSaving(false);
        }
    };

    return {
        isEditing,
        paragraphs,
        isSaving,
        toggleEdit,
        handleParagraphChange,
        changeParagraphType,
        addParagraph,
        removeParagraph,
        handleSave,
    };
}
