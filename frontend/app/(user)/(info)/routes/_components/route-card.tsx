import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { MapPin, Clock, DollarSign } from "lucide-react";
import Link from "next/link";

interface RouteCardProps {
  origin: string;
  destination: string;
  duration: string;
  priceFrom: number;
  operators: number;
  popular?: boolean;
}

export function RouteCard({
  origin,
  destination,
  duration,
  priceFrom,
  operators,
  popular,
}: RouteCardProps) {
  return (
    <Card className="transition-shadow hover:shadow-lg">
      <CardContent className="pt-6">
        <div className="mb-4 flex items-start justify-between">
          <div className="flex-1">
            <div className="mb-2 flex items-center gap-2">
              <MapPin className="h-5 w-5 text-primary" />
              <h3 className="text-lg font-semibold">
                {origin} → {destination}
              </h3>
              {popular && (
                <Badge
                  variant="secondary"
                  className="bg-primary/10 text-primary"
                >
                  Phổ biến
                </Badge>
              )}
            </div>
            <div className="space-y-1 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <Clock className="h-4 w-4" />
                <span>Thời gian: {duration}</span>
              </div>
              <div className="flex items-center gap-2">
                <DollarSign className="h-4 w-4" />
                <span>Từ {priceFrom.toLocaleString()}đ</span>
              </div>
              <p>{operators} nhà xe phục vụ tuyến này</p>
            </div>
          </div>
        </div>
        <Button asChild className="w-full">
          <Link href="/trips">Tìm chuyến xe</Link>
        </Button>
      </CardContent>
    </Card>
  );
}
