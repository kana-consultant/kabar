// src/components/generate/ModelSelector.tsx
import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Loader2, Cpu, Key } from "lucide-react";
import { getAPIKeys, type APIKey, type APIKeyDetail } from "@/services/apiKeyService";

interface ModelSelectorProps {
  selectedModelId: string;
  onModelChange: (modelId: string) => void;
}

export function ModelSelector({
  selectedModelId,
  onModelChange,
}: ModelSelectorProps) {
  const [models, setModels] = useState<APIKeyDetail[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadModels();
  }, []);

  const loadModels = async () => {
    try {
      setLoading(true);
      setError(null);

      const data = await getAPIKeys();
      setModels(data as any);
      console.log(data);
    } catch (err) {
      console.error(err);
      setError("Failed to load AI models");
    } finally {
      setLoading(false);
    }
  };

  const getProviderColor = (provider: string) => {
    switch (provider?.toLowerCase()) {
      case "openai":
        return "bg-green-100 text-green-700";
      case "anthropic":
        return "bg-purple-100 text-purple-700";
      case "google":
        return "bg-blue-100 text-blue-700";
      case "openrouter":
        return "bg-orange-100 text-orange-700";
      default:
        return "bg-gray-100 text-gray-700";
    }
  };

  const getServiceBadge = (service: string) => {
    return (
      <Badge variant="outline" className="text-xs">
        {service}
      </Badge>
    );
  };

  if (loading) {
    return (
      <Card>
        <CardContent className="flex justify-center py-6">
          <Loader2 className="h-5 w-5 animate-spin text-slate-500" />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="py-4 text-center text-red-500">
          {error}
        </CardContent>
      </Card>
    );
  }

  if (!models) {
    return (
      <Card>
        <CardContent className="py-6 text-center">
          <Key className="mx-auto mb-2 h-8 w-8 text-slate-400" />
          <p className="text-sm text-slate-500">
            Belum ada model yang dikonfigurasi
          </p>
          <p className="mt-1 text-xs text-slate-400">
            Tambahkan API key di halaman settings
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center gap-2">
          <Cpu className="h-4 w-4 text-slate-500" />
          <CardTitle className="text-sm font-medium">
            AI Model
          </CardTitle>
        </div>
      </CardHeader>

      <CardContent>
        <Select
          value={selectedModelId}
          onValueChange={onModelChange}
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Pilih AI Model" />
          </SelectTrigger>

          <SelectContent className="w-[650px] max-w-[95vw]">
            {models.map((model) => (
              <SelectItem
                key={model.id}
                value={model.id}
                className="py-3"
              >
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-medium text-sm">
                    {model.modelDisplayName || model.modelName}
                  </span>

                  {getServiceBadge(model.service)}

                  <Badge
                    className={getProviderColor(
                      model.providerName
                    )}
                  >
                    {model.providerDisplayName}
                  </Badge>

                  {model.isActive && (
                    <Badge
                      variant="outline"
                      className="text-blue-500"
                    >
                      Active
                    </Badge>
                  )}
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </CardContent>
    </Card>
  );
}