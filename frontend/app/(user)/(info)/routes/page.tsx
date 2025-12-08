import { RouteCard } from "./_components/route-card";
import { Bus } from "lucide-react";

const popularRoutes = [
  {
    origin: "Hồ Chí Minh",
    destination: "Đà Lạt",
    duration: "6-7 giờ",
    priceFrom: 180000,
    operators: 25,
    popular: true,
  },
  {
    origin: "Hà Nội",
    destination: "Hải Phòng",
    duration: "2-3 giờ",
    priceFrom: 120000,
    operators: 18,
    popular: true,
  },
  {
    origin: "Hồ Chí Minh",
    destination: "Vũng Tàu",
    duration: "2-3 giờ",
    priceFrom: 100000,
    operators: 30,
    popular: true,
  },
  {
    origin: "Hà Nội",
    destination: "Sa Pa",
    duration: "5-6 giờ",
    priceFrom: 200000,
    operators: 15,
  },
  {
    origin: "Hồ Chí Minh",
    destination: "Nha Trang",
    duration: "8-9 giờ",
    priceFrom: 220000,
    operators: 20,
  },
  {
    origin: "Đà Nẵng",
    destination: "Hội An",
    duration: "45 phút",
    priceFrom: 50000,
    operators: 12,
  },
  {
    origin: "Hà Nội",
    destination: "Ninh Bình",
    duration: "2-3 giờ",
    priceFrom: 100000,
    operators: 10,
  },
  {
    origin: "Hồ Chí Minh",
    destination: "Phan Thiết",
    duration: "4-5 giờ",
    priceFrom: 150000,
    operators: 22,
  },
  {
    origin: "Hà Nội",
    destination: "Đà Nẵng",
    duration: "14-16 giờ",
    priceFrom: 350000,
    operators: 15,
  },
];

export const metadata = {
  title: "Tuyến đường phổ biến | BusTicket.vn",
  description: "Các tuyến xe phổ biến nhất tại Việt Nam",
};

export default function RoutesPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <Bus className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Tuyến đường phổ biến</h1>
        <p className="text-lg text-muted-foreground">
          Khám phá các tuyến xe được yêu thích nhất trên toàn quốc
        </p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {popularRoutes.map((route, index) => (
          <RouteCard key={index} {...route} />
        ))}
      </div>
    </div>
  );
}
