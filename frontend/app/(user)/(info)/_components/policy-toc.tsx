"use client";

import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { FileText } from "lucide-react";

interface TOCItem {
  id: string;
  title: string;
}

interface PolicyTOCProps {
  items: TOCItem[];
}

export function PolicyTOC({ items }: PolicyTOCProps) {
  return (
    <Card className="sticky top-20">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <FileText className="h-5 w-5" />
          Mục lục
        </CardTitle>
      </CardHeader>
      <CardContent>
        <nav>
          <ul className="space-y-2 text-sm">
            {items.map((item) => (
              <li key={item.id}>
                <Link
                  href={`#${item.id}`}
                  className="text-muted-foreground transition-colors hover:text-primary"
                >
                  {item.title}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      </CardContent>
    </Card>
  );
}
