// hooks/useModelManagement/useTemplateEditor.ts

import { useState, useCallback, useEffect } from "react";

export type TemplateMode = "visual" | "code";

interface UseTemplateEditorProps {
    initialTemplate: Record<string, any> | null;
    onTemplateChange: (template: string) => void;
}

export function useTemplateEditor({ initialTemplate, onTemplateChange }: UseTemplateEditorProps) {
    const [mode, setMode] = useState<TemplateMode>("visual");
    const [templateString, setTemplateString] = useState<string>("");
    const [templateObject, setTemplateObject] = useState<Record<string, any>>({});
    const [isValid, setIsValid] = useState(true);

    // Parse template string to object while preserving {placeholder}
    const parseTemplateToObject = useCallback((templateStr: string) => {
        // If it's a double stringified JSON, parse it first
        let processed = templateStr;
        try {
            const firstParse = JSON.parse(templateStr);
            if (typeof firstParse === 'string') {
                processed = firstParse;
            } else if (typeof firstParse === 'object') {
                return firstParse;
            }
        } catch {
            // Not double stringified, continue
        }

        // Replace {placeholder} with valid JSON values temporarily, then restore      
        const placeholderMap: Record<string, string> = {};
        let counter = 0;
        
        // Extract all placeholders and replace with temp keys
        let jsonString = processed.replace(/"\{(\w+)\}"/g, (match, key) => {
            const tempKey = `__TEMP_${counter++}__`;
            placeholderMap[tempKey] = `{${key}}`;
            return `"${tempKey}"`;
        }).replace(/\{(\w+)\}/g, (match, key) => {
            const tempKey = `__TEMP_${counter++}__`;
            placeholderMap[tempKey] = `{${key}}`;
            return tempKey;
        });

        try {
            // Parse as JSON
            let parsed = JSON.parse(jsonString);
            
            // Restore placeholders
            const jsonStr = JSON.stringify(parsed);
            const restored = jsonStr.replace(/__TEMP_\d+__/g, (match) => {
                return placeholderMap[match] || match;
            });
            
            return JSON.parse(restored);
        } catch (e) {
            console.error('Parse error:', e);
            return {};
        }
    }, []);

    // Convert object back to string template
    const objectToString = useCallback((obj: Record<string, any>) => {
        return JSON.stringify(obj, null, 2);
    }, []);

    // Initialize on mount or when initialTemplate changes
    useEffect(() => {
        if (initialTemplate && Object.keys(initialTemplate).length > 0) {
            try {
                const obj = initialTemplate;
                setTemplateObject(obj);
                setTemplateString(objectToString(obj));
                setIsValid(true);
            } catch (e: any) {
                console.error('Parse error:', e.message);
                setTemplateString("");
                setTemplateObject({});
                setIsValid(false);
            }
        } else {
            setTemplateString("{}");
            setTemplateObject({});
            setIsValid(true);
        }
    }, [initialTemplate, objectToString]);

    const updateFromObject = useCallback((obj: Record<string, any>) => {
        const jsonString = objectToString(obj);
        setTemplateObject(obj);
        setTemplateString(jsonString);
        setIsValid(true);
        onTemplateChange(jsonString);
    }, [onTemplateChange, objectToString]);

    const updateFromString = useCallback((value: string) => {
        setTemplateString(value);
        if (!value || value.trim() === "") {
            setIsValid(false);
            return;
        }
        try {
            const parsed = JSON.parse(value);
            setTemplateObject(parsed);
            setIsValid(true);
            onTemplateChange(value);
        } catch {
            // Try to parse with placeholders
            const parsed = parseTemplateToObject(value);
            if (Object.keys(parsed).length > 0) {
                setTemplateObject(parsed);
                setIsValid(true);
                onTemplateChange(objectToString(parsed));
            } else {
                setIsValid(false);
            }
        }
    }, [onTemplateChange, parseTemplateToObject, objectToString]);

    const formatJSON = useCallback(() => {
        if (!templateString || templateString.trim() === "") {
            return;
        }
        try {
            const parsed = JSON.parse(templateString);
            const formatted = JSON.stringify(parsed, null, 2);
            setTemplateString(formatted);
            setTemplateObject(parsed);
            setIsValid(true);
            onTemplateChange(formatted);
        } catch {
            const parsed = parseTemplateToObject(templateString);
            if (Object.keys(parsed).length > 0) {
                const formatted = objectToString(parsed);
                setTemplateString(formatted);
                setTemplateObject(parsed);
                setIsValid(true);
                onTemplateChange(formatted);
            } else {
                setIsValid(false);
            }
        }
    }, [templateString, onTemplateChange, parseTemplateToObject, objectToString]);

    const switchToVisualMode = useCallback(() => {
        if (isValid && templateString && templateString.trim() !== "") {
            try {
                const parsed = JSON.parse(templateString);
                setTemplateObject(parsed);
                setMode("visual");
            } catch {
                const parsed = parseTemplateToObject(templateString);
                if (Object.keys(parsed).length > 0) {
                    setTemplateObject(parsed);
                    setMode("visual");
                }
            }
        }
    }, [isValid, templateString, parseTemplateToObject]);

    const switchToCodeMode = useCallback(() => {
        setMode("code");
    }, []);

    return {
        mode,
        templateString,
        templateObject,
        isValid,
        updateFromObject,
        updateFromString,
        formatJSON,
        switchToVisualMode,
        switchToCodeMode,
        setMode,
    };
}