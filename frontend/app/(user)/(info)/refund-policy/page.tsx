import { PolicyTOC } from "../_components/policy-toc";
import { PolicySection } from "../_components/policy-section";
import {
  FileText,
  Scale,
  GitPullRequest,
  Clock,
  AlertCircle,
  CheckCircle2,
  XCircle,
} from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

const tocItems = [
  { id: "overview", title: "Tổng quan" },
  { id: "conditions", title: "Điều kiện hoàn tiền" },
  { id: "process", title: "Quy trình hoàn tiền" },
  { id: "timeline", title: "Thời gian xử lý" },
  { id: "exceptions", title: "Trường hợp ngoại lệ" },
];

export const metadata = {
  title: "Chính sách hoàn tiền | BusTicket.vn",
  description: "Chính sách hoàn tiền khi hủy vé tại BusTicket.vn",
};

export default function RefundPolicyPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <Scale className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Chính sách hoàn tiền</h1>
        <p className="text-lg text-muted-foreground">
          Cam kết minh bạch và bảo vệ quyền lợi khách hàng
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr,300px]">
        <div className="space-y-8">
          <PolicySection title="Tổng quan" id="overview" icon={FileText}>
            <p>
              Chính sách hoàn tiền này áp dụng cho tất cả các giao dịch đặt vé
              thông qua BusTicket.vn. Chúng tôi cam kết xử lý yêu cầu hoàn tiền
              một cách nhanh chóng và minh bạch.
            </p>
          </PolicySection>

          <PolicySection
            title="Điều kiện hoàn tiền"
            id="conditions"
            icon={Scale}
          >
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="rounded-lg border bg-green-50/50 p-4 dark:bg-green-900/10">
                <h4 className="mb-2 flex items-center gap-2 font-semibold text-green-700 dark:text-green-400">
                  <CheckCircle2 className="h-5 w-5" />
                  Được hoàn tiền
                </h4>
                <ul className="space-y-1 text-sm text-muted-foreground">
                  <li>• Hủy vé trong thời gian cho phép</li>
                  <li>• Chuyến xe bị hủy do nhà xe</li>
                  <li>• Lỗi kỹ thuật dẫn đến thanh toán trùng</li>
                  <li>• Thông tin chuyến xe không chính xác</li>
                </ul>
              </div>
              <div className="rounded-lg border bg-red-50/50 p-4 dark:bg-red-900/10">
                <h4 className="mb-2 flex items-center gap-2 font-semibold text-red-700 dark:text-red-400">
                  <XCircle className="h-5 w-5" />
                  Không được hoàn tiền
                </h4>
                <ul className="space-y-1 text-sm text-muted-foreground">
                  <li>• Hủy vé dưới 6 giờ trước giờ khởi hành</li>
                  <li>• Không đến điểm đón đúng giờ</li>
                  <li>• Vi phạm điều khoản sử dụng</li>
                </ul>
              </div>
            </div>
          </PolicySection>

          <PolicySection
            title="Quy trình hoàn tiền"
            id="process"
            icon={GitPullRequest}
          >
            <div className="relative ml-2 space-y-6 border-l-2 border-muted pl-6">
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  1
                </span>
                <h4 className="font-semibold">Gửi yêu cầu</h4>
                <p className="text-sm text-muted-foreground">
                  Đăng nhập, chọn vé cần hủy và xác nhận hủy vé.
                </p>
              </div>
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  2
                </span>
                <h4 className="font-semibold">Xác nhận</h4>
                <p className="text-sm text-muted-foreground">
                  Hệ thống tự động tính toán số tiền hoàn và gửi email xác nhận.
                </p>
              </div>
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  3
                </span>
                <h4 className="font-semibold">Hoàn tiền</h4>
                <p className="text-sm text-muted-foreground">
                  Tiền được chuyển về tài khoản của bạn theo thời gian quy định.
                </p>
              </div>
            </div>
          </PolicySection>

          <PolicySection title="Thời gian xử lý" id="timeline" icon={Clock}>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="rounded-lg border p-4 text-center">
                <div className="text-2xl font-bold text-primary">5-7</div>
                <div className="text-xs tracking-wider text-muted-foreground uppercase">
                  Ngày làm việc
                </div>
                <div className="mt-2 text-sm font-medium">Chuyển khoản</div>
              </div>
              <div className="rounded-lg border p-4 text-center">
                <div className="text-2xl font-bold text-primary">3-5</div>
                <div className="text-xs tracking-wider text-muted-foreground uppercase">
                  Ngày làm việc
                </div>
                <div className="mt-2 text-sm font-medium">Ví điện tử</div>
              </div>
              <div className="rounded-lg border p-4 text-center">
                <div className="text-2xl font-bold text-primary">7-14</div>
                <div className="text-xs tracking-wider text-muted-foreground uppercase">
                  Ngày làm việc
                </div>
                <div className="mt-2 text-sm font-medium">Thẻ tín dụng</div>
              </div>
            </div>
            <p className="mt-4 text-center text-sm text-muted-foreground">
              * Trong trường hợp đặc biệt, thời gian có thể kéo dài hơn. Chúng
              tôi sẽ thông báo nếu có thay đổi.
            </p>
          </PolicySection>

          <PolicySection
            title="Trường hợp ngoại lệ"
            id="exceptions"
            icon={AlertCircle}
          >
            <Alert className="mb-4 border-primary/20 bg-primary/5">
              <AlertCircle className="h-4 w-4 text-primary" />
              <AlertTitle className="font-semibold text-primary">
                Hoàn tiền 100%
              </AlertTitle>
              <AlertDescription>
                Áp dụng cho các trường hợp: Chuyến xe bị hủy do nhà xe, sự cố kỹ
                thuật từ phía BusTicket.vn, hoặc thiên tai/dịch bệnh.
              </AlertDescription>
            </Alert>
            <div className="rounded-lg bg-muted p-4">
              <h4 className="mb-2 font-semibold">Liên hệ hỗ trợ</h4>
              <p className="mb-2 text-sm text-muted-foreground">
                Nếu bạn gặp vấn đề với việc hoàn tiền, vui lòng liên hệ bộ phận
                chăm sóc khách hàng:
              </p>
              <ul className="space-y-1 text-sm font-medium">
                <li>Email: support@busticket.vn</li>
                <li>Hotline: 1900 989 901</li>
              </ul>
            </div>
          </PolicySection>
        </div>

        <aside className="hidden lg:block">
          <PolicyTOC items={tocItems} />
        </aside>
      </div>
    </div>
  );
}
