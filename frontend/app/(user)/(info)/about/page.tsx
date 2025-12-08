import { HeroSection } from "./_components/hero-section";
import { MissionSection } from "./_components/mission-section";
import { TeamSection } from "./_components/team-section";
import { StatsSection } from "./_components/stats-section";

export const metadata = {
  title: "Về chúng tôi | BusTicket.vn",
  description:
    "Tìm hiểu về BusTicket.vn - Nền tảng đặt vé xe khách trực tuyến hàng đầu Việt Nam",
};

export default function AboutPage() {
  return (
    <div className="space-y-12">
      <HeroSection />
      <MissionSection />
      <TeamSection />
      <StatsSection />
    </div>
  );
}
