import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Star, MapPin, Bus } from "lucide-react";
import Link from "next/link";

interface OperatorCardProps {
  name: string;
  rating: number;
  totalTrips: number;
  routes: string[];
  verified?: boolean;
}

export function OperatorCard({
  name,
  rating,
  totalTrips,
  routes,
  verified,
}: OperatorCardProps) {
  return (
    <Card className="transition-shadow hover:shadow-lg">
      <CardContent className="pt-6">
        <div className="mb-4">
          <div className="mb-2 flex items-start justify-between">
            <h3 className="text-lg font-semibold">{name}</h3>
            {verified && (
              <Badge variant="secondary" className="bg-success/10 text-success">
                Đã xác minh
              </Badge>
            )}
          </div>
          <div className="flex items-center gap-4 text-sm">
            <div className="flex items-center gap-1">
              <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
              <span className="font-semibold">{rating}</span>
            </div>
            <div className="flex items-center gap-1 text-muted-foreground">
              <Bus className="h-4 w-4" />
              <span>{totalTrips} chuyến</span>
            </div>
          </div>
        </div>

        <div className="mb-4">
          <div className="mb-2 flex items-center gap-2 text-sm font-medium">
            <MapPin className="h-4 w-4 text-primary" />
            <span>Tuyến đường chính:</span>
          </div>
          <div className="flex flex-wrap gap-2">
            {routes.slice(0, 3).map((route, index) => (
              <Badge key={index} variant="outline">
                {route}
              </Badge>
            ))}
          </div>
        </div>

        <Button asChild variant="outline" className="w-full">
          <Link href="/trips">Xem chuyến xe</Link>
        </Button>
      </CardContent>
    </Card>
  );
}
