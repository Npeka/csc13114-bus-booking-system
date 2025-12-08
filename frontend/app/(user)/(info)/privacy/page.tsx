import { PolicyTOC } from "../_components/policy-toc";
import { PolicySection } from "../_components/policy-section";
import {
  Shield,
  Database,
  Settings,
  Share2,
  Lock,
  UserCog,
} from "lucide-react";

const tocItems = [
  { id: "intro", title: "Giới thiệu" },
  { id: "collection", title: "Thu thập thông tin" },
  { id: "usage", title: "Sử dụng thông tin" },
  { id: "sharing", title: "Chia sẻ thông tin" },
  { id: "security", title: "Bảo mật" },
  { id: "rights", title: "Quyền của bạn" },
];

export const metadata = {
  title: "Chính sách bảo mật | BusTicket.vn",
  description: "Chính sách bảo mật thông tin cá nhân tại BusTicket.vn",
};

export default function PrivacyPage() {
  return (
    <div>
      <div className="mb-8 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
          <Shield className="h-8 w-8 text-primary" />
        </div>
        <h1 className="mb-2 text-4xl font-bold">Chính sách bảo mật</h1>
        <p className="text-lg text-muted-foreground">
          Cam kết bảo vệ thông tin cá nhân của bạn
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr,300px]">
        <div className="space-y-8">
          <PolicySection title="Giới thiệu" id="intro" icon={Shield}>
            <p>
              BusTicket.vn cam kết bảo vệ quyền riêng tư và thông tin cá nhân
              của bạn. Chính sách này giải thích cách chúng tôi thu thập, sử
              dụng và bảo vệ thông tin của bạn khi bạn sử dụng dịch vụ của chúng
              tôi.
            </p>
          </PolicySection>

          <PolicySection
            title="Thu thập thông tin"
            id="collection"
            icon={Database}
          >
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="rounded-lg border p-4">
                <h4 className="mb-2 font-semibold">Thông tin cá nhân</h4>
                <p className="text-sm text-muted-foreground">
                  Họ tên, email, số điện thoại, địa chỉ (khi bạn cung cấp).
                </p>
              </div>
              <div className="rounded-lg border p-4">
                <h4 className="mb-2 font-semibold">Thông tin thanh toán</h4>
                <p className="text-sm text-muted-foreground">
                  Lịch sử giao dịch (chúng tôi không lưu trữ số thẻ tín dụng).
                </p>
              </div>
              <div className="rounded-lg border p-4">
                <h4 className="mb-2 font-semibold">Thông tin kỹ thuật</h4>
                <p className="text-sm text-muted-foreground">
                  IP address, loại trình duyệt, thiết bị sử dụng.
                </p>
              </div>
              <div className="rounded-lg border p-4">
                <h4 className="mb-2 font-semibold">Cookies</h4>
                <p className="text-sm text-muted-foreground">
                  Dữ liệu duyệt web để cải thiện trải nghiệm người dùng.
                </p>
              </div>
            </div>
          </PolicySection>

          <PolicySection title="Sử dụng thông tin" id="usage" icon={Settings}>
            <p>Chúng tôi sử dụng thông tin của bạn để:</p>
            <ul className="not-prose grid gap-2 sm:grid-cols-2">
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Xử lý đặt vé và thanh toán
              </li>
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Gửi xác nhận và vé điện tử
              </li>
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Hỗ trợ khách hàng
              </li>
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Cải thiện dịch vụ
              </li>
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Phòng chống gian lận
              </li>
              <li className="flex items-center gap-2 text-sm">
                <span className="h-1.5 w-1.5 rounded-full bg-primary" />
                Tuân thủ pháp luật
              </li>
            </ul>
          </PolicySection>

          <PolicySection title="Chia sẻ thông tin" id="sharing" icon={Share2}>
            <p>Chúng tôi có thể chia sẻ thông tin với:</p>
            <ul className="space-y-2">
              <li>
                <strong>Nhà xe đối tác:</strong> Để sắp xếp chỗ ngồi và liên hệ
                đón trả.
              </li>
              <li>
                <strong>Đối tác thanh toán:</strong> PayOS và ngân hàng để xử lý
                giao dịch.
              </li>
              <li>
                <strong>Cơ quan pháp luật:</strong> Khi có yêu cầu hợp pháp từ
                cơ quan nhà nước.
              </li>
            </ul>
            <div className="mt-4 rounded-lg bg-green-50 p-4 text-green-800 dark:bg-green-900/20 dark:text-green-300">
              <p className="text-center text-sm font-medium">
                Chúng tôi CAM KẾT KHÔNG bán thông tin cá nhân của bạn cho bên
                thứ ba vì mục đích thương mại.
              </p>
            </div>
          </PolicySection>

          <PolicySection title="Bảo mật thông tin" id="security" icon={Lock}>
            <p>Các biện pháp bảo mật chúng tôi áp dụng:</p>
            <div className="flex flex-wrap gap-2">
              <span className="rounded-full bg-secondary px-3 py-1 text-sm">
                Mã hóa SSL/TLS
              </span>
              <span className="rounded-full bg-secondary px-3 py-1 text-sm">
                Server bảo mật
              </span>
              <span className="rounded-full bg-secondary px-3 py-1 text-sm">
                Kiểm soát truy cập
              </span>
              <span className="rounded-full bg-secondary px-3 py-1 text-sm">
                Giám sát 24/7
              </span>
            </div>
          </PolicySection>

          <PolicySection title="Quyền của bạn" id="rights" icon={UserCog}>
            <p>Bạn có quyền:</p>
            <ul className="space-y-1">
              <li>• Truy cập và xem thông tin cá nhân</li>
              <li>• Yêu cầu sửa đổi hoặc xóa thông tin</li>
              <li>• Từ chối nhận email marketing</li>
              <li>• Khiếu nại nếu thấy vi phạm</li>
            </ul>
            <p className="mt-4 text-sm text-muted-foreground">
              Để thực hiện các quyền này, vui lòng liên hệ:{" "}
              <a
                href="mailto:privacy@busticket.vn"
                className="text-primary hover:underline"
              >
                privacy@busticket.vn
              </a>
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
