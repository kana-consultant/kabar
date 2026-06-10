// components/KeyValueEditor.tsx
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Plus, Trash2 } from "lucide-react";

interface KeyValueEditorProps {
  value: Record<string, string>;
  onChange: (value: Record<string, string>) => void;
}

export function KeyValueEditor({ value, onChange }: KeyValueEditorProps) {
  const entries = Object.entries(value);

  const handleAdd = () => {
    onChange({ ...value, "": "" });
  };

  const handleRemove = (key: string) => {
    const newValue = { ...value };
    delete newValue[key];
    onChange(newValue);
  };

  const handleChange = (oldKey: string, newKey: string, newVal: string) => {
    const newValue: Record<string, string> = {};
    Object.entries(value).forEach(([k, v]) => {
      if (k === oldKey) {
        if (newKey) newValue[newKey] = newVal;
      } else {
        newValue[k] = v;
      }
    });
    onChange(newValue);
  };

  return (
    <div className="space-y-2">
      {entries.map(([key, val], index) => (
        <div key={index} className="flex space-x-2">
          <Input
            value={key}
            onChange={(e) => handleChange(key, e.target.value, val)}
            placeholder="Header Name"
            className="flex-1"
          />
          <Input
            value={val}
            onChange={(e) => handleChange(key, key, e.target.value)}
            placeholder="Header Value"
            className="flex-1"
          />
          <Button
            variant="ghost"
            size="icon"
            onClick={() => handleRemove(key)}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      ))}
      <Button variant="outline" size="sm" onClick={handleAdd}>
        <Plus className="h-4 w-4 mr-2" />
        Add Header
      </Button>
    </div>
  );
}