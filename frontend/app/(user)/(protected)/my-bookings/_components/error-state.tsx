import { Card, CardContent } from "@/components/ui/card";

interface ErrorStateProps {
  error: Error | null;
}

export function ErrorState({ error }: ErrorStateProps) {
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
          <CardContent className="py-12 text-center">
            <p className="text-error">
              Đã xảy ra lỗi khi tải dữ liệu. Vui lòng thử lại sau.
            </p>
            <p className="mt-2 text-sm text-muted-foreground">
              {error instanceof Error ? error.message : "Lỗi không xác định"}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
