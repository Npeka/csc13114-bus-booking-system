"use client";

import { useState } from "react";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";

interface FAQSearchProps {
  onSearch: (query: string) => void;
}

export function FAQSearch({ onSearch }: FAQSearchProps) {
  const [query, setQuery] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    onSearch(value);
  };

  return (
    <div className="relative mb-8">
      <Search className="absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
      <Input
        type="text"
        placeholder="Tìm kiếm câu hỏi..."
        value={query}
        onChange={handleChange}
        className="pl-10"
      />
    </div>
  );
}
