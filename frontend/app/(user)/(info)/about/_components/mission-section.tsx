import { Card, CardContent } from "@/components/ui/card";
import { Target, Eye, Heart } from "lucide-react";

const missions = [
  {
    icon: Target,
    title: "Sứ mệnh",
    description:
      "Kết nối hành khách với các nhà xe uy tín, tạo ra hệ sinh thái vận tải hành khách hiện đại và tiện lợi nhất Việt Nam.",
  },
  {
    icon: Eye,
    title: "Tầm nhìn",
    description:
      "Trở thành nền tảng đặt vé xe khách số 1 Việt Nam, phục vụ hàng triệu lượt khách mỗi năm với chất lượng dịch vụ vượt trội.",
  },
  {
    icon: Heart,
    title: "Giá trị cốt lõi",
    description:
      "Đặt khách hàng làm trung tâm, cam kết minh bạch, an toàn và không ngừng đổi mới để mang lại trải nghiệm tốt nhất.",
  },
];

export function MissionSection() {
  return (
    <section className="mb-12">
      <h2 className="mb-6 text-center text-3xl font-bold">
        Sứ mệnh & Tầm nhìn
      </h2>
      <div className="grid gap-6 md:grid-cols-3">
        {missions.map((mission, index) => {
          const Icon = mission.icon;
          return (
            <Card key={index}>
              <CardContent className="pt-6">
                <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <Icon className="h-6 w-6 text-primary" />
                </div>
                <h3 className="mb-2 text-xl font-semibold">{mission.title}</h3>
                <p className="text-sm leading-relaxed text-muted-foreground">
                  {mission.description}
                </p>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </section>
  );
}
