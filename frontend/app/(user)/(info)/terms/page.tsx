import { PolicyTOC } from "../_components/policy-toc";
import { PolicySection } from "../_components/policy-section";
import {
  CheckCircle2,
  Briefcase,
  UserCheck,
  Copyright,
  ShieldAlert,
  History,
} from "lucide-react";

const tocItems = [
  { id: "acceptance", title: "Chấp nhận điều khoản" },
  { id: "services", title: "Dịch vụ" },
  { id: "user-obligations", title: "Nghĩa vụ người dùng" },
  { id: "intellectual-property", title: "Quyền sở hữu trí tuệ" },
  { id: "liability", title: "Giới hạn trách nhiệm" },
  { id: "changes", title: "Thay đổi điều khoản" },
];

export const metadata = {
  title: "Điều khoản dịch vụ | BusTicket.vn",
  description: "Điều khoản sử dụng dịch vụ tại BusTicket.vn",
};

export default function TermsPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <CheckCircle2 className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Điều khoản dịch vụ</h1>
        <p className="text-lg text-muted-foreground">
          Quy định sử dụng dịch vụ tại BusTicket.vn
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr,300px]">
        <div className="space-y-8">
          <PolicySection
            title="Chấp nhận điều khoản"
            id="acceptance"
            icon={CheckCircle2}
          >
            <p>
              Bằng việc truy cập và sử dụng BusTicket.vn, bạn đồng ý tuân thủ và
              bị ràng buộc bởi các điều khoản và điều kiện sau đây.
            </p>
            <div className="rounded-lg bg-muted p-4 text-sm text-muted-foreground">
              Nếu bạn không đồng ý với bất kỳ phần nào của các điều khoản này,
              vui lòng không sử dụng dịch vụ của chúng tôi.
            </div>
          </PolicySection>

          <PolicySection title="Dịch vụ" id="services" icon={Briefcase}>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <h4 className="mb-2 font-semibold text-primary">
                  Chúng tôi cung cấp
                </h4>
                <ul className="space-y-1 text-sm">
                  <li>• Nền tảng tìm kiếm và đặt vé xe</li>
                  <li>• Thanh toán trực tuyến an toàn</li>
                  <li>• Thông tin minh bạch về nhà xe</li>
                  <li>• Hỗ trợ khách hàng 24/7</li>
                </ul>
              </div>
              <div>
                <h4 className="mb-2 font-semibold text-muted-foreground">
                  Chúng tôi KHÔNG
                </h4>
                <ul className="space-y-1 text-sm text-muted-foreground">
                  <li>• Vận hành xe khách trực tiếp</li>
                  <li>• Chịu trách nhiệm chất lượng xe (do nhà xe đảm nhận)</li>
                </ul>
              </div>
            </div>
          </PolicySection>

          <PolicySection
            title="Nghĩa vụ người dùng"
            id="user-obligations"
            icon={UserCheck}
          >
            <p>Khi sử dụng dịch vụ, bạn cam kết:</p>
            <ul className="space-y-2">
              <li className="flex items-start gap-2">
                <CheckCircle2 className="mt-1 h-4 w-4 shrink-0 text-green-500" />
                <span>
                  Cung cấp thông tin chính xác và đầy đủ khi đăng ký/đặt vé.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 className="mt-1 h-4 w-4 shrink-0 text-green-500" />
                <span>Bảo mật thông tin tài khoản và mật khẩu của mình.</span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 className="mt-1 h-4 w-4 shrink-0 text-green-500" />
                <span>
                  Không sử dụng dịch vụ cho mục đích bất hợp pháp hoặc phá hoại.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <CheckCircle2 className="mt-1 h-4 w-4 shrink-0 text-green-500" />
                <span>
                  Tuân thủ quy định của nhà xe khi sử dụng dịch vụ vận chuyển.
                </span>
              </li>
            </ul>
          </PolicySection>

          <PolicySection
            title="Quyền sở hữu trí tuệ"
            id="intellectual-property"
            icon={Copyright}
          >
            <p>
              Tất cả nội dung trên BusTicket.vn bao gồm văn bản, hình ảnh, logo,
              và mã nguồn đều thuộc quyền sở hữu của chúng tôi hoặc đối tác.
            </p>
            <p className="mt-4 font-medium">Nghiêm cấm các hành vi:</p>
            <ul className="list-disc space-y-1 pl-4">
              <li>Sao chép, phân phối nội dung mà không có sự cho phép</li>
              <li>
                Sử dụng logo, thương hiệu của chúng tôi cho mục đích thương mại
              </li>
              <li>Tạo sản phẩm phái sinh từ dịch vụ của chúng tôi</li>
            </ul>
          </PolicySection>

          <PolicySection
            title="Giới hạn trách nhiệm"
            id="liability"
            icon={ShieldAlert}
          >
            <p>
              BusTicket.vn không chịu trách nhiệm cho các thiệt hại gián tiếp,
              sự cố kỹ thuật không mong muốn, hoặc hành vi của nhà xe đối tác.
            </p>
            <div className="mt-4 rounded-lg border border-yellow-200 bg-yellow-50 p-4 text-yellow-800 dark:border-yellow-900/50 dark:bg-yellow-900/20 dark:text-yellow-200">
              <p className="text-sm font-medium">
                Trách nhiệm tối đa của chúng tôi giới hạn ở số tiền bạn đã thanh
                toán cho vé trong giao dịch liên quan.
              </p>
            </div>
          </PolicySection>

          <PolicySection
            title="Thay đổi điều khoản"
            id="changes"
            icon={History}
          >
            <p>
              Chúng tôi có quyền cập nhật, thay đổi điều khoản này bất kỳ lúc
              nào. Mọi thay đổi sẽ có hiệu lực ngay khi được đăng tải trên
              website.
            </p>
            <p>
              Việc bạn tiếp tục sử dụng dịch vụ sau khi có thay đổi đồng nghĩa
              với việc bạn chấp nhận các điều khoản mới.
            </p>
          </PolicySection>
        </div>

        <aside className="hidden lg:block">
          <PolicyTOC items={tocItems} />
        </aside>
      </div>
    </div>
  );
}
