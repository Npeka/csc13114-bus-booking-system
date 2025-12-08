import { Bus } from "lucide-react";

export function HeroSection() {
  return (
    <div className="mb-12 text-center">
      <div className="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-primary/10">
        <Bus className="h-10 w-10 text-primary" />
      </div>
      <h1 className="mb-4 text-4xl font-bold">Về BusTicket.vn</h1>
      <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
        Nền tảng đặt vé xe khách trực tuyến hàng đầu Việt Nam, mang đến trải
        nghiệm đặt vé nhanh chóng, an toàn và tiện lợi cho hàng triệu hành
        khách.
      </p>
    </div>
  );
}
