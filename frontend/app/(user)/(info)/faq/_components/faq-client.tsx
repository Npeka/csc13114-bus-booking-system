"use client";

import { useState } from "react";
import { FAQSearch } from "./faq-search";
import { FAQAccordion } from "./faq-accordion";
import { HelpCircle, Phone, Mail, MessageSquare } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const faqData = {
  "Đặt vé": [
    {
      question: "Làm thế nào để đặt vé xe?",
      answer:
        "Bạn có thể đặt vé bằng cách: 1) Tìm kiếm chuyến xe phù hợp, 2) Chọn ghế ngồi, 3) Điền thông tin, 4) Thanh toán. Vé điện tử sẽ được gửi qua email ngay sau khi thanh toán thành công.",
    },
    {
      question: "Tôi có thể đặt vé cho người khác không?",
      answer:
        "Có, bạn hoàn toàn có thể đặt vé cho người khác. Chỉ cần điền đúng thông tin của hành khách khi đặt vé.",
    },
    {
      question: "Có thể đặt bao nhiêu ghế trong một lần?",
      answer:
        "Bạn có thể đặt tối đa 5 ghế trong một lần đặt vé. Nếu cần đặt nhiều hơn, vui lòng đặt nhiều lần hoặc liên hệ hotline.",
    },
  ],
  "Thanh toán": [
    {
      question: "Có những phương thức thanh toán nào?",
      answer:
        "Chúng tôi chấp nhận thanh toán qua PayOS, bao gồm chuyển khoản ngân hàng, quét mã QR VietQR, và thẻ ATM nội địa.",
    },
    {
      question: "Thanh toán có an toàn không?",
      answer:
        "Tất cả giao dịch đều được mã hóa SSL và xử lý thông qua cổng thanh toán PayOS đạt chuẩn bảo mật quốc tế.",
    },
    {
      question: "Tôi đã thanh toán nhưng chưa nhận được vé?",
      answer:
        "Vé thường được gửi trong vòng 5 phút. Vui lòng kiểm tra hòm thư spam. Nếu vẫn chưa nhận được, liên hệ hotline 1900 989 901.",
    },
  ],
  "Hủy và đổi vé": [
    {
      question: "Tôi có thể hủy vé không?",
      answer:
        "Có, bạn có thể hủy vé trước giờ khởi hành. Phí hủy tùy thuộc vào thời gian hủy: trước 24h hoàn 70%, 12-24h hoàn 50%, 6-12h hoàn 30%.",
    },
    {
      question: "Làm thế nào để đổi vé?",
      answer:
        "Đổi vé tương đương với hủy vé cũ và đặt vé mới. Bạn cần hủy vé hiện tại (chịu phí hủy nếu có) và đặt vé mới cho chuyến khác.",
    },
    {
      question: "Khi nào tôi nhận được tiền hoàn?",
      answer:
        "Tiền hoàn sẽ được xử lý trong 5-7 ngày làm việc và chuyển về tài khoản/phương thức thanh toán ban đầu.",
    },
  ],
  "Vấn đề khác": [
    {
      question: "Tôi cần mang theo giấy tờ gì khi lên xe?",
      answer:
        "Bạn cần mang theo CMND/CCCD và vé điện tử (có thể hiển thị trên điện thoại hoặc in ra).",
    },
    {
      question: "Tôi có thể mang hành lý bao nhiêu kg?",
      answer:
        "Mỗi hành khách được mang tối đa 20kg hành lý miễn phí. Hành lý vượt quá sẽ được tính phí theo quy định của nhà xe.",
    },
    {
      question: "Làm sao để liên hệ hỗ trợ?",
      answer:
        "Bạn có thể liên hệ qua: Hotline 1900 989 901, Email support@busticket.vn, hoặc chat trực tuyến trên website.",
    },
  ],
};

export function FAQClient() {
  const [searchQuery, setSearchQuery] = useState("");

  const filterFAQs = () => {
    if (!searchQuery.trim()) return faqData;

    const filtered: Partial<typeof faqData> = {};
    (Object.keys(faqData) as Array<keyof typeof faqData>).forEach(
      (category) => {
        const items = faqData[category];
        const matchedItems = items.filter(
          (item) =>
            item.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
            item.answer.toLowerCase().includes(searchQuery.toLowerCase()),
        );
        if (matchedItems.length > 0) {
          filtered[category] = matchedItems;
        }
      },
    );
    return filtered;
  };

  const filteredFAQs = filterFAQs();

  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <HelpCircle className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Câu hỏi thường gặp</h1>
        <p className="text-lg text-muted-foreground">
          Tìm câu trả lời cho các thắc mắc của bạn
        </p>
      </div>

      <FAQSearch onSearch={setSearchQuery} />

      <div className="space-y-8">
        {Object.entries(filteredFAQs).map(
          ([category, items]) =>
            items && (
              <FAQAccordion key={category} category={category} items={items} />
            ),
        )}

        {Object.keys(filteredFAQs).length === 0 && (
          <div className="py-12 text-center text-muted-foreground">
            Không tìm thấy câu hỏi phù hợp. Vui lòng thử từ khóa khác hoặc liên
            hệ hỗ trợ.
          </div>
        )}
      </div>

      <div className="mt-16">
        <h2 className="mb-6 text-center text-2xl font-bold">
          Vẫn cần sự hỗ trợ?
        </h2>
        <div className="grid gap-6 md:grid-cols-3">
          <Card className="text-center transition-shadow hover:shadow-lg">
            <CardContent className="pt-6">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Phone className="h-6 w-6 text-primary" />
              </div>
              <h3 className="mb-2 font-semibold">Hotline</h3>
              <p className="mb-2 text-sm text-muted-foreground">Hỗ trợ 24/7</p>
              <a
                href="tel:1900989901"
                className="font-bold text-primary hover:underline"
              >
                1900 989 901
              </a>
            </CardContent>
          </Card>
          <Card className="text-center transition-shadow hover:shadow-lg">
            <CardContent className="pt-6">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Mail className="h-6 w-6 text-primary" />
              </div>
              <h3 className="mb-2 font-semibold">Email</h3>
              <p className="mb-2 text-sm text-muted-foreground">
                Phản hồi trong 24h
              </p>
              <a
                href="mailto:support@busticket.vn"
                className="font-bold text-primary hover:underline"
              >
                support@busticket.vn
              </a>
            </CardContent>
          </Card>
          <Card className="text-center transition-shadow hover:shadow-lg">
            <CardContent className="pt-6">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <MessageSquare className="h-6 w-6 text-primary" />
              </div>
              <h3 className="mb-2 font-semibold">Live Chat</h3>
              <p className="mb-2 text-sm text-muted-foreground">
                Chat trực tiếp với nhân viên
              </p>
              <span className="cursor-pointer font-bold text-primary hover:underline">
                Bắt đầu chat
              </span>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
