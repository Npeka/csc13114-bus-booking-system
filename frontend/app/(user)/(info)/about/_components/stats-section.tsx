import { Card, CardContent } from "@/components/ui/card";
import { Users, Bus, MapPin, Award } from "lucide-react";

const stats = [
  {
    icon: Users,
    value: "500K+",
    label: "Người dùng",
  },
  {
    icon: Bus,
    value: "100+",
    label: "Nhà xe đối tác",
  },
  {
    icon: MapPin,
    value: "200+",
    label: "Tuyến đường",
  },
  {
    icon: Award,
    value: "4.8/5",
    label: "Đánh giá",
  },
];

export function StatsSection() {
  return (
    <section>
      <h2 className="mb-6 text-center text-3xl font-bold">
        Thành tích nổi bật
      </h2>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <Card key={index}>
              <CardContent className="pt-6 text-center">
                <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <Icon className="h-6 w-6 text-primary" />
                </div>
                <div className="mb-1 text-3xl font-bold text-primary">
                  {stat.value}
                </div>
                <div className="text-sm text-muted-foreground">
                  {stat.label}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
