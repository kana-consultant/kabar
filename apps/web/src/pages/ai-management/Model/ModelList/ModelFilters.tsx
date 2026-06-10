import { Input, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import { Search } from "lucide-react";
import { useState } from "react";

export function ModelFilters() {
    const [search, setSearch] = useState("");
    const [provider, setProvider] = useState("all");
    const [status, setStatus] = useState("all");

    return (
        <div className="flex gap-4">
            <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                    placeholder="Search models..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    className="pl-9"
                />
            </div>
            <Select value={provider} onValueChange={setProvider}>
                <SelectTrigger className="w-48">
                    <SelectValue placeholder="All Providers" />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="all">All Providers</SelectItem>
                    <SelectItem value="openai">OpenAI</SelectItem>
                    <SelectItem value="anthropic">Anthropic</SelectItem>
                    <SelectItem value="google">Google</SelectItem>
                </SelectContent>
            </Select>
            <Select value={status} onValueChange={setStatus}>
                <SelectTrigger className="w-36">
                    <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="active">Active</SelectItem>
                    <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
            </Select>
        </div>
    );
}