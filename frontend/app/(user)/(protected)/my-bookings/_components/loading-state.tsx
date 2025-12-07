import { Card, CardContent } from "@/components/ui/card";
import { Loader2 } from "lucide-react";

export function LoadingState() {
  return (
    <div className="min-h-screen bg-secondary/30">
      <div className="container py-4">
        <div className="mb-4">
          <h1 className="text-2xl font-bold">Vé đã đặt</h1>
          <p className="text-muted-foreground">
            Quản lý và theo dõi các chuyến đi của bạn
          </p>
        </div>
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <div className="flex flex-col items-center gap-3">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-muted-foreground">Đang tải vé của bạn...</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
