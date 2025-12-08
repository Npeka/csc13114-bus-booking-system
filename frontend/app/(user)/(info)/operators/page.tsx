import { OperatorCard } from "./_components/operator-card";
import { Building2 } from "lucide-react";

const operators = [
  {
    name: "Phương Trang FutaBus",
    rating: 4.8,
    totalTrips: 856,
    routes: ["HCM - Đà Lạt", "HCM - Nha Trang", "HCM - Vũng Tàu"],
    verified: true,
  },
  {
    name: "Thành Bưởi",
    rating: 4.7,
    totalTrips: 642,
    routes: ["HCM - Đà Lạt", "HCM - Phan Thiết", "HCM - Mũi Né"],
    verified: true,
  },
  {
    name: "Hoàng Long",
    rating: 4.6,
    totalTrips: 534,
    routes: ["Hà Nội - Hải Phòng", "Hà Nội - Ninh Bình", "Hà Nội - Sa Pa"],
    verified: true,
  },
  {
    name: "Kumho Samco",
    rating: 4.5,
    totalTrips: 428,
    routes: ["HCM - Đà Lạt", "HCM - Buôn Ma Thuột", "HCM - Pleiku"],
    verified: true,
  },
  {
    name: "Mai Linh Express",
    rating: 4.4,
    totalTrips: 392,
    routes: ["HCM - Vũng Tàu", "HCM - Phan Thiết", "Đà Nẵng - Hội An"],
    verified: true,
  },
  {
    name: "Hà Linh",
    rating: 4.3,
    totalTrips: 315,
    routes: ["Hà Nội - Đà Nẵng", "Hà Nội - Vinh", "Hà Nội - Huế"],
    verified: false,
  },
];

export const metadata = {
  title: "Nhà xe | BusTicket.vn",
  description: "Danh sách các nhà xe uy tín tại Việt Nam",
};

export default function OperatorsPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <Building2 className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Nhà xe đối tác</h1>
        <p className="text-lg text-muted-foreground">
          Hợp tác với các nhà xe uy tín, chất lượng cao
        </p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {operators.map((operator, index) => (
          <OperatorCard key={index} {...operator} />
        ))}
      </div>
    </div>
  );
}
