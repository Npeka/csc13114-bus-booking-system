import { TripSearchForm } from "@/components/search/trip-search-form";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { 
  Shield, 
  Clock, 
  CreditCard, 
  HeadphonesIcon, 
  Star, 
  TrendingUp,
  Users,
  Bus
} from "lucide-react";

export default function Home() {
  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-brand-primary/10 via-brand-primary-light/20 to-background py-16 md:py-24">
        <div className="container">
          <div className="mx-auto max-w-3xl text-center mb-12">
            <Badge variant="secondary" className="mb-4">
              🚌 Nền tảng đặt vé #1 Việt Nam
            </Badge>
            <h1 className="text-4xl font-bold tracking-tight text-foreground md:text-5xl lg:text-6xl">
              Đặt vé xe khách
              <br />
              <span className="text-brand-primary">nhanh chóng & tiện lợi</span>
            </h1>
            <p className="mt-6 text-lg text-muted-foreground md:text-xl">
              Hàng trăm tuyến đường khắp Việt Nam. Đặt vé online, thanh toán an toàn, 
              lên xe không lo.
            </p>
          </div>

          {/* Search Form */}
          <div className="flex justify-center">
            <TripSearchForm />
          </div>

          {/* Trust Indicators */}
          <div className="mt-12 grid grid-cols-2 gap-6 text-center md:grid-cols-4">
            <div className="space-y-2">
              <div className="text-3xl font-bold text-brand-primary">500K+</div>
              <div className="text-sm text-muted-foreground">Vé đã đặt</div>
            </div>
            <div className="space-y-2">
              <div className="text-3xl font-bold text-brand-primary">200+</div>
              <div className="text-sm text-muted-foreground">Nhà xe</div>
            </div>
            <div className="space-y-2">
              <div className="text-3xl font-bold text-brand-primary">1000+</div>
              <div className="text-sm text-muted-foreground">Tuyến đường</div>
            </div>
            <div className="space-y-2">
              <div className="text-3xl font-bold text-brand-primary">4.8/5</div>
              <div className="text-sm text-muted-foreground">Đánh giá</div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 md:py-24">
        <div className="container">
          <div className="mx-auto max-w-2xl text-center mb-12">
            <h2 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
              Tại sao chọn BusTicket.vn?
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Chúng tôi cam kết mang đến trải nghiệm đặt vé tốt nhất cho bạn
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
            <Card className="border-2 hover:border-brand-primary transition-colors">
              <CardContent className="pt-6">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-brand-primary/10">
                  <Shield className="h-6 w-6 text-brand-primary" />
                </div>
                <h3 className="mb-2 text-xl font-semibold">An toàn & Bảo mật</h3>
                <p className="text-sm text-muted-foreground">
                  Thanh toán được mã hóa SSL. Thông tin cá nhân được bảo vệ tuyệt đối.
                </p>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-brand-primary transition-colors">
              <CardContent className="pt-6">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-success/10">
                  <Clock className="h-6 w-6 text-success" />
                </div>
                <h3 className="mb-2 text-xl font-semibold">Đặt vé nhanh</h3>
                <p className="text-sm text-muted-foreground">
                  Chỉ 3 bước đơn giản. Nhận vé điện tử ngay lập tức qua email và SMS.
                </p>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-brand-primary transition-colors">
              <CardContent className="pt-6">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-info/10">
                  <CreditCard className="h-6 w-6 text-info" />
                </div>
                <h3 className="mb-2 text-xl font-semibold">Thanh toán linh hoạt</h3>
                <p className="text-sm text-muted-foreground">
                  Hỗ trợ MoMo, ZaloPay, PayOS và các phương thức phổ biến khác.
                </p>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-brand-primary transition-colors">
              <CardContent className="pt-6">
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-warning/10">
                  <HeadphonesIcon className="h-6 w-6 text-warning" />
                </div>
                <h3 className="mb-2 text-xl font-semibold">Hỗ trợ 24/7</h3>
                <p className="text-sm text-muted-foreground">
                  Đội ngũ CSKH luôn sẵn sàng hỗ trợ bạn mọi lúc mọi nơi.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Popular Routes Section */}
      <section className="bg-neutral-50 py-16 md:py-24">
        <div className="container">
          <div className="mx-auto max-w-2xl text-center mb-12">
            <h2 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
              Tuyến đường phổ biến
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Các tuyến xe được khách hàng lựa chọn nhiều nhất
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {popularDestinations.map((route) => (
              <Card key={route.id} className="card-hover cursor-pointer">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold">{route.from}</h3>
                      <div className="flex items-center text-sm text-muted-foreground my-2">
                        <Bus className="h-4 w-4 mr-1" />
                        <span>→</span>
                      </div>
                      <h3 className="text-lg font-semibold">{route.to}</h3>
                    </div>
                    <Badge variant="secondary" className="ml-2">
                      <TrendingUp className="h-3 w-3 mr-1" />
                      Phổ biến
                    </Badge>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">
                      {route.operators} nhà xe
                    </span>
                    <span className="font-semibold text-brand-primary">
                      Từ {route.priceFrom.toLocaleString()}đ
                    </span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Reviews Section */}
      <section className="py-16 md:py-24">
        <div className="container">
          <div className="mx-auto max-w-2xl text-center mb-12">
            <h2 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
              Khách hàng nói gì về chúng tôi
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Hơn 10,000 đánh giá 5 sao từ khách hàng hài lòng
            </p>
          </div>

          <div className="grid gap-6 md:grid-cols-3">
            {reviews.map((review) => (
              <Card key={review.id}>
                <CardContent className="p-6">
                  <div className="mb-4 flex">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className="h-4 w-4 fill-warning text-warning"
                      />
                    ))}
                  </div>
                  <p className="mb-4 text-sm text-muted-foreground">
                    &ldquo;{review.content}&rdquo;
                  </p>
                  <div className="flex items-center">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-brand-primary/10">
                      <Users className="h-5 w-5 text-brand-primary" />
                    </div>
                    <div className="ml-3">
                      <p className="text-sm font-semibold">{review.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {review.route}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-brand-primary py-16 text-white">
        <div className="container text-center">
          <h2 className="text-3xl font-bold md:text-4xl">
            Sẵn sàng cho chuyến đi tiếp theo?
          </h2>
          <p className="mt-4 text-lg text-white/90">
            Tải app ngay để nhận ưu đãi độc quyền và trải nghiệm tốt hơn
          </p>
          <div className="mt-8 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <div className="flex h-14 items-center space-x-2 rounded-lg bg-white px-6 text-foreground">
              <span className="text-2xl">📱</span>
              <div className="text-left">
                <div className="text-xs">Tải về trên</div>
                <div className="text-sm font-semibold">App Store</div>
              </div>
            </div>
            <div className="flex h-14 items-center space-x-2 rounded-lg bg-white px-6 text-foreground">
              <span className="text-2xl">🤖</span>
              <div className="text-left">
                <div className="text-xs">Tải về trên</div>
                <div className="text-sm font-semibold">Google Play</div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}

const popularDestinations = [
  { id: 1, from: "Hà Nội", to: "Đà Nẵng", operators: 25, priceFrom: 350000 },
  { id: 2, from: "TP. Hồ Chí Minh", to: "Đà Lạt", operators: 30, priceFrom: 180000 },
  { id: 3, from: "Hà Nội", to: "Sa Pa", operators: 15, priceFrom: 250000 },
  { id: 4, from: "TP. Hồ Chí Minh", to: "Nha Trang", operators: 28, priceFrom: 220000 },
  { id: 5, from: "Hà Nội", to: "Hạ Long", operators: 20, priceFrom: 150000 },
  { id: 6, from: "TP. Hồ Chí Minh", to: "Phan Thiết", operators: 22, priceFrom: 120000 },
];

const reviews = [
  {
    id: 1,
    name: "Nguyễn Văn A",
    route: "Hà Nội → Đà Nẵng",
    content:
      "Đặt vé rất nhanh và tiện lợi. Nhân viên hỗ trợ nhiệt tình. Sẽ tiếp tục sử dụng dịch vụ.",
  },
  {
    id: 2,
    name: "Trần Thị B",
    route: "TP.HCM → Đà Lạt",
    content:
      "Giao diện đẹp, dễ sử dụng. Thanh toán qua MoMo rất tiện. Vé điện tử được gửi ngay lập tức.",
  },
  {
    id: 3,
    name: "Lê Văn C",
    route: "Hà Nội → Sa Pa",
    content:
      "Giá cả hợp lý, nhiều nhà xe để lựa chọn. Đã giới thiệu cho bạn bè và gia đình.",
  },
];
