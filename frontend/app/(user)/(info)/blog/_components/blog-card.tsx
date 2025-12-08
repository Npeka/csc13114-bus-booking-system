import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar, User, ArrowRight } from "lucide-react";
import Link from "next/link";

interface BlogCardProps {
  title: string;
  excerpt: string;
  category: string;
  date: string;
  author: string;
  slug: string;
}

export function BlogCard({
  title,
  excerpt,
  category,
  date,
  author,
  slug,
}: BlogCardProps) {
  return (
    <Card className="flex flex-col transition-shadow hover:shadow-lg">
      <CardContent className="flex-1 pt-6">
        <Badge variant="secondary" className="mb-3">
          {category}
        </Badge>
        <h3 className="mb-2 line-clamp-2 text-xl font-semibold">{title}</h3>
        <p className="mb-4 line-clamp-3 text-sm text-muted-foreground">
          {excerpt}
        </p>
        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <div className="flex items-center gap-1">
            <Calendar className="h-3 w-3" />
            <span>{date}</span>
          </div>
          <div className="flex items-center gap-1">
            <User className="h-3 w-3" />
            <span>{author}</span>
          </div>
        </div>
      </CardContent>
      <CardFooter>
        <Button asChild variant="ghost" className="w-full">
          <Link href={`/blog/${slug}`}>
            Đọc thêm
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
