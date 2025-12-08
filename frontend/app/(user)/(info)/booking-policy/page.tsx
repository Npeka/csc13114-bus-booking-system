import { PolicyTOC } from "../_components/policy-toc";
import { PolicySection } from "../_components/policy-section";
import {
  Info,
  ClipboardList,
  CreditCard,
  RefreshCw,
  Undo2,
  CheckCircle2,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

const tocItems = [
  { id: "overview", title: "Tổng quan" },
  { id: "booking", title: "Quy trình đặt vé" },
  { id: "payment", title: "Thanh toán" },
  { id: "changes", title: "Thay đổi & Hủy vé" },
  { id: "refund", title: "Hoàn tiền" },
];

export const metadata = {
  title: "Chính sách đặt vé | BusTicket.vn",
  description: "Tìm hiểu về chính sách đặt vé tại BusTicket.vn",
};

export default function BookingPolicyPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <ClipboardList className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Chính sách đặt vé</h1>
        <p className="text-lg text-muted-foreground">
          Quy định và hướng dẫn đặt vé tại BusTicket.vn
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr,300px]">
        <div className="space-y-8">
          <PolicySection title="Tổng quan" id="overview" icon={Info}>
            <p>
              Chính sách đặt vé này quy định các điều khoản và điều kiện áp dụng
              khi bạn đặt vé xe thông qua nền tảng BusTicket.vn.
            </p>
            <Alert className="border-primary/20 bg-primary/5">
              <CheckCircle2 className="h-4 w-4 text-primary" />
              <AlertTitle>Đồng ý điều khoản</AlertTitle>
              <AlertDescription>
                Bằng việc sử dụng dịch vụ của chúng tôi, bạn đồng ý tuân thủ các
                điều khoản được nêu trong chính sách này.
              </AlertDescription>
            </Alert>
          </PolicySection>

          <PolicySection
            title="Quy trình đặt vé"
            id="booking"
            icon={ClipboardList}
          >
            <div className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="rounded-lg border p-4">
                  <h4 className="font-semibold text-primary">1. Tìm kiếm</h4>
                  <p className="text-sm text-muted-foreground">
                    Chọn điểm đi, điểm đến và ngày khởi hành phù hợp.
                  </p>
                </div>
                <div className="rounded-lg border p-4">
                  <h4 className="font-semibold text-primary">2. Chọn ghế</h4>
                  <p className="text-sm text-muted-foreground">
                    Lựa chọn vị trí ghế ngồi mong muốn trên xe.
                  </p>
                </div>
                <div className="rounded-lg border p-4">
                  <h4 className="font-semibold text-primary">
                    3. Điền thông tin
                  </h4>
                  <p className="text-sm text-muted-foreground">
                    Cung cấp thông tin hành khách chính xác.
                  </p>
                </div>
                <div className="rounded-lg border p-4">
                  <h4 className="font-semibold text-primary">4. Thanh toán</h4>
                  <p className="text-sm text-muted-foreground">
                    Hoàn tất thanh toán an toàn qua PayOS.
                  </p>
                </div>
              </div>
              <div className="rounded-lg bg-muted p-4">
                <h4 className="mb-2 font-semibold">Thông tin bắt buộc:</h4>
                <ul className="list-disc space-y-1 pl-4">
                  <li>Họ và tên đầy đủ (như trên giấy tờ tùy thân)</li>
                  <li>Số điện thoại liên hệ (để nhận thông báo)</li>
                  <li>Email nhận vé (để nhận vé điện tử)</li>
                </ul>
              </div>
            </div>
          </PolicySection>

          <PolicySection title="Thanh toán" id="payment" icon={CreditCard}>
            <p>BusTicket.vn chấp nhận các phương thức thanh toán sau:</p>
            <ul className="not-prose mb-4 grid gap-2 sm:grid-cols-3">
              <li className="flex items-center gap-2 rounded-md border p-3 text-sm">
                <span className="font-semibold">QR Code</span> (VietQR)
              </li>
              <li className="flex items-center gap-2 rounded-md border p-3 text-sm">
                <span className="font-semibold">Thẻ ATM</span> Nội địa
              </li>
              <li className="flex items-center gap-2 rounded-md border p-3 text-sm">
                <span className="font-semibold">Chuyển khoản</span> 24/7
              </li>
            </ul>
            <p className="text-sm text-muted-foreground">
              * Vé sẽ được xác nhận ngay sau khi thanh toán thành công. Thời
              gian xử lý giao dịch thường từ 1-5 phút.
            </p>
          </PolicySection>

          <PolicySection
            title="Thay đổi & Hủy vé"
            id="changes"
            icon={RefreshCw}
          >
            <div className="mb-4">
              <h4 className="font-semibold">Thay đổi thông tin</h4>
              <p>
                Bạn có thể thay đổi thông tin hành khách (họ tên, số điện thoại)
                miễn phí trước 24 giờ khởi hành.
              </p>
            </div>
            <div>
              <h4 className="mb-2 font-semibold">Chính sách hủy vé</h4>
              <div className="overflow-hidden rounded-lg border">
                <table className="m-0 w-full">
                  <thead className="bg-muted">
                    <tr>
                      <th className="p-3 text-left">Thời gian hủy</th>
                      <th className="p-3 text-left">Mức hoàn tiền</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y">
                    <tr>
                      <td className="p-3">Trước 24 giờ</td>
                      <td className="p-3 font-semibold text-green-600">70%</td>
                    </tr>
                    <tr>
                      <td className="p-3">12 - 24 giờ</td>
                      <td className="p-3 font-semibold text-yellow-600">50%</td>
                    </tr>
                    <tr>
                      <td className="p-3">6 - 12 giờ</td>
                      <td className="p-3 font-semibold text-orange-600">30%</td>
                    </tr>
                    <tr>
                      <td className="p-3">Dưới 6 giờ</td>
                      <td className="p-3 font-semibold text-red-600">0%</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </PolicySection>

          <PolicySection title="Hoàn tiền" id="refund" icon={Undo2}>
            <ul className="space-y-2">
              <li className="flex items-start gap-2">
                <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                <span>
                  Thời gian xử lý hoàn tiền: <strong>5-7 ngày làm việc</strong>{" "}
                  kể từ ngày hủy vé thành công.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                <span>
                  Tiền sẽ được hoàn về tài khoản/phương thức thanh toán ban đầu.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                <span>
                  Phí hoàn tiền có thể phát sinh tùy thuộc vào chính sách của
                  ngân hàng/nhà xe.
                </span>
              </li>
            </ul>
          </PolicySection>
        </div>

        <aside className="hidden lg:block">
          <PolicyTOC items={tocItems} />
        </aside>
      </div>
    </div>
  );
}
