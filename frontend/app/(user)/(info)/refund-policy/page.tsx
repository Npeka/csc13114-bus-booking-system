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
            <Alert className="mb-4">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Yêu cầu quan trọng</AlertTitle>
              <AlertDescription>
                Bạn cần cung cấp thông tin tài khoản ngân hàng trong Hồ sơ của
                mình trước khi yêu cầu hoàn tiền.
              </AlertDescription>
            </Alert>
            <div className="relative ml-2 space-y-6 border-l-2 border-muted pl-6">
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  1
                </span>
                <h4 className="font-semibold">Gửi yêu cầu hủy vé</h4>
                <p className="text-sm text-muted-foreground">
                  Đăng nhập, vào Hồ sơ để thêm thông tin tài khoản ngân hàng
                  (nếu chưa có), sau đó chọn vé cần hủy và xác nhận.
                </p>
              </div>
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  2
                </span>
                <h4 className="font-semibold">Quản trị viên xem xét</h4>
                <p className="text-sm text-muted-foreground">
                  Đội ngũ quản trị sẽ kiểm tra yêu cầu, xác nhận điều kiện hoàn
                  tiền và số tiền được hoàn lại.
                </p>
              </div>
              <div className="relative">
                <span className="absolute -left-[31px] flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs text-primary-foreground">
                  3
                </span>
                <h4 className="font-semibold">Chuyển khoản ngân hàng</h4>
                <p className="text-sm text-muted-foreground">
                  Sau khi duyệt, tiền sẽ được chuyển khoản vào tài khoản ngân
                  hàng bạn đã đăng ký trong vòng 7-14 ngày làm việc.
                </p>
              </div>
            </div>
          </PolicySection>

          <PolicySection title="Thời gian xử lý" id="timeline" icon={Clock}>
            <div className="overflow-hidden rounded-lg border">
              <table className="w-full">
                <thead className="bg-muted/50">
                  <tr>
                    <th className="px-4 py-3 text-left text-sm font-semibold">
                      Phương thức hoàn tiền
                    </th>
                    <th className="px-4 py-3 text-left text-sm font-semibold">
                      Thời gian hoàn tiền
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  <tr>
                    <td className="px-4 py-3 text-sm">
                      Chuyển khoản ngân hàng
                    </td>
                    <td className="px-4 py-3 text-sm font-medium">
                      7-14 ngày làm việc
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <Alert className="mt-4">
              <Clock className="h-4 w-4" />
              <AlertDescription>
                Hiện tại, chúng tôi chỉ hỗ trợ hoàn tiền qua chuyển khoản ngân
                hàng. Vui lòng cung cấp thông tin tài khoản chính xác trong Hồ
                sơ để nhận tiền hoàn.
              </AlertDescription>
            </Alert>
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
