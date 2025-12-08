import { BlogCard } from "./_components/blog-card";
import { BookOpen } from "lucide-react";

const blogPosts = [
  {
    title: "10 điểm đến du lịch không thể bỏ qua trong năm 2025",
    excerpt:
      "Khám phá những địa điểm du lịch tuyệt vời nhất Việt Nam dành cho những chuyến đi cuối tuần hoặc kỳ nghỉ dài.",
    category: "Địa điểm",
    date: "15/12/2024",
    author: "Nguyễn Văn A",
    slug: "10-diem-den-du-lich-2025",
  },
  {
    title: "Hướng dẫn đặt vé xe khách online cho người mới",
    excerpt:
      "Cẩm nang chi tiết từ A đến Z về cách đặt vé xe online một cách nhanh chóng và an toàn.",
    category: "Hướng dẫn",
    date: "10/12/2024",
    author: "Trần Thị B",
    slug: "huong-dan-dat-ve-xe-online",
  },
  {
    title: "So sánh giá vé xe các tuyến phổ biến",
    excerpt:
      "Tổng hợp và so sánh giá vé xe khách các tuyến Sài Gòn - Đà Lạt, Hà Nội - Hải Phòng và nhiều tuyến khác.",
    category: "So sánh",
    date: "05/12/2024",
    author: "Nguyễn Văn A",
    slug: "so-sanh-gia-ve-xe",
  },
  {
    title: "Kinh nghiệm đi xe khách đường dài",
    excerpt:
      "Những mẹo hay để có chuyến đi xe khách đường dài thoải mái và an toàn nhất.",
    category: "Kinh nghiệm",
    date: "01/12/2024",
    author: "Trần Thị B",
    slug: "kinh-nghiem-di-xe-duong-dai",
  },
  {
    title: "Khuyến mãi xe khách tháng 12",
    excerpt:
      "Cập nhật các chương trình khuyến mãi, giảm giá vé xe hấp dẫn trong tháng 12.",
    category: "Khuyến mãi",
    date: "28/11/2024",
    author: "Nguyễn Văn A",
    slug: "khuyen-mai-thang-12",
  },
  {
    title: "Review top 5 nhà xe chất lượng cao",
    excerpt:
      "Đánh giá chi tiết về dịch vụ, chất lượng của 5 nhà xe được yêu thích nhất hiện nay.",
    category: "Review",
    date: "25/11/2024",
    author: "Trần Thị B",
    slug: "review-top-5-nha-xe",
  },
];

export const metadata = {
  title: "Blog | BusTicket.vn",
  description: "Tin tức, hướng dẫn và kinh nghiệm đi xe khách",
};

export default function BlogPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <BookOpen className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Blog</h1>
        <p className="text-lg text-muted-foreground">
          Tin tức, hướng dẫn và kinh nghiệm du lịch
        </p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {blogPosts.map((post, index) => (
          <BlogCard key={index} {...post} />
        ))}
      </div>
    </div>
  );
}
