// src/components/generate/ModelSelector.tsx
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Loader2, Cpu, Key, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useModels } from "@/hooks/useGenerate/useModel";
import { useEffect } from "react";

interface ModelSelectorProps {
  selectedModelId: string;
  onModelChange: (modelId: string) => void;
  filterByService?: "text" | "image" | "all";
}

export function ModelSelector({
  selectedModelId,
  onModelChange,
  filterByService = "all",
}: ModelSelectorProps) {
  const { 
    models,
    loadingModels,
    isLoading,
    isError,
    error,
    refetch,
    getProviderColor,
    getServiceBadgeColor,
    getTextAndImageModels,
  } = useModels();

  useEffect(()=>{
    getTextAndImageModels()
  },[])

  // Filter models based on service
  const getFilteredModels = () => {
    return models;
  };

  const filteredModels = getFilteredModels();

  const getServiceBadge = (service: string) => {
    const colors = getServiceBadgeColor(service);
    const labels: Record<string, string> = {
      text: "Text Generation",
      image: "Image Generation",
      embedding: "Embedding",
    };
    
    return (
      <Badge variant="outline" className={`text-xs ${colors}`}>
        {labels[service] || service}
      </Badge>
    );
  };

  // Use loadingModels or isLoading (both available)
  const isProcessing = loadingModels || isLoading;

  if (isProcessing) {
    return (
      <Card>
        <CardContent className="flex justify-center py-6">
          <Loader2 className="h-5 w-5 animate-spin text-slate-500" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="py-6 text-center">
          <p className="text-sm text-red-500 mb-2">
            {error?.message || "Failed to load AI models"}
          </p>
          <Button
            variant="outline"
            size="sm"
            onClick={() => refetch()}
            className="mt-2"
          >
            <RefreshCw className="mr-2 h-3 w-3" />
            Retry
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (!filteredModels || filteredModels.length === 0) {
    return (
      <Card>
        <CardContent className="py-6 text-center">
          <Key className="mx-auto mb-2 h-8 w-8 text-slate-400" />
          <p className="text-sm text-slate-500">
            {filterByService === "text" 
              ? "Belum ada model teks yang dikonfigurasi"
              : filterByService === "image"
                ? "Belum ada model gambar yang dikonfigurasi"
                : "Belum ada model yang dikonfigurasi"}
          </p>
          <p className="mt-1 text-xs text-slate-400">
            Tambahkan API key di halaman settings
          </p>
        </CardContent>
      </Card>
    );
  }

  // Find selected model name for display
  const selectedModel = filteredModels.find(m => m.id === selectedModelId);

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Cpu className="h-4 w-4 text-slate-500" />
            <CardTitle className="text-sm font-medium">
              AI Model {filterByService !== "all" && `(${filterByService === "text" ? "Text" : "Image"})`}
            </CardTitle>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => refetch()}
            className="h-6 px-2"
            title="Refresh models"
          >
            <RefreshCw className="h-3 w-3" />
          </Button>
        </div>
      </CardHeader>

      <CardContent>
        <Select
          value={selectedModelId}
          onValueChange={(value)=>{
            onModelChange(value)
          }}
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Pilih AI Model" />
          </SelectTrigger>

          <SelectContent className="w-[650px] max-w-[95vw]">
            {filteredModels.map((model) => (
              <SelectItem
                key={model.id}
                value={model.id}
                className="py-3"
              >
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-medium text-sm">
                    {model.modelDisplayName}
                  </span>

                  {getServiceBadge(model.service)}

                  <Badge
                    className={getProviderColor(model.providerDisplayName)}
                  >
                    {model.providerDisplayName}
                  </Badge>

                  {model.isActive && (
                    <Badge
                      variant="outline"
                      className="bg-green-50 text-green-600 border-green-200"
                    >
                      Active
                    </Badge>
                  )}
                </div>
                
                {/* Show system prompt preview if exists */}
                {model.systemPrompt && (
                  <div className="mt-1 text-xs text-slate-400 truncate">
                    System: {model.systemPrompt.substring(0, 50)}...
                  </div>
                )}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        
        {/* Show selected model info */}
        {selectedModelId && selectedModel && (
          <div className="mt-3 text-xs text-slate-400 flex items-center gap-2">
            <span>✅</span>
            <span>{selectedModel.modelDisplayName} ready to use</span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}