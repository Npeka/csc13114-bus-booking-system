import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Mail, Linkedin } from "lucide-react";

const team = [
  {
    name: "Nguyễn Duy Hoàng",
    role: "Fullstack Developer",
    bio: "Đam mê xây dựng sản phẩm công nghệ chất lượng cao. Chuyên gia về React, Next.js và kiến trúc hệ thống.",
    email: "duyhoangnguyen198@gmail.com",
    linkedin: "https://www.linkedin.com/in/nduyhoang/",
    initials: "DH",
  },
  {
    name: "Nguyễn Phúc Khang",
    role: "Backend & DevOps Engineer",
    bio: "Chuyên sâu về hệ thống phân tán, cloud infrastructure và tối ưu hóa hiệu năng. Đảm bảo hệ thống vận hành ổn định.",
    email: "npkhang287@gmail.com",
    linkedin: "https://www.linkedin.com/in/npkhang287",
    initials: "PK",
  },
];

export function TeamSection() {
  return (
    <section className="mb-12">
      <h2 className="mb-6 text-center text-3xl font-bold">Đội ngũ sáng lập</h2>
      <div className="grid gap-6 md:grid-cols-2">
        {team.map((member, index) => (
          <Card key={index}>
            <CardContent className="pt-6">
              <div className="flex items-start gap-4">
                <Avatar className="h-16 w-16">
                  <AvatarFallback className="bg-primary/10 text-lg font-semibold text-primary">
                    {member.initials}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <h3 className="text-xl font-semibold">{member.name}</h3>
                  <p className="mb-2 text-sm text-primary">{member.role}</p>
                  <p className="mb-3 text-sm leading-relaxed text-muted-foreground">
                    {member.bio}
                  </p>
                  <div className="flex gap-3">
                    <a
                      href={`mailto:${member.email}`}
                      className="text-muted-foreground transition-colors hover:text-primary"
                    >
                      <Mail className="h-4 w-4" />
                    </a>
                    <a
                      href={member.linkedin}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-muted-foreground transition-colors hover:text-primary"
                    >
                      <Linkedin className="h-4 w-4" />
                    </a>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  );
}
