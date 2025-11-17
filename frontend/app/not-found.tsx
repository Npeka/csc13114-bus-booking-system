import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Bus } from "lucide-react";

export default function NotFound() {
  return (
    <main className="flex min-h-[70vh] flex-col items-center justify-center px-6 text-center">
      <div className="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-full bg-brand-primary/10 text-brand-primary">
        <Bus className="h-10 w-10" />
      </div>
      <p className="text-sm font-semibold uppercase tracking-widest text-brand-primary">
        404 · Không tìm thấy trang
      </p>
      <h1 className="mt-4 text-3xl font-bold tracking-tight text-foreground md:text-4xl">
        Trang bạn đang tìm không tồn tại
      </h1>
      <p className="mt-3 max-w-xl text-base text-muted-foreground">
        Có thể đường dẫn đã bị thay đổi hoặc trang đã bị xóa. Vui lòng quay lại
        trang chủ hoặc tiếp tục tìm chuyến xe phù hợp với bạn.
      </p>
      <div className="mt-8 flex flex-col gap-3 sm:flex-row">
        <Button asChild>
          <Link href="/">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Về trang chủ
          </Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/trips">Tiếp tục tìm chuyến</Link>
        </Button>
      </div>
    </main>
  );
}

