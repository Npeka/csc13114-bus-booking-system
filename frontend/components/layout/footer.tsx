import Image from "next/image";
import Link from "next/link";
import { Mail, Phone, MapPin } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t bg-linear-to-tl from-primary/30 via-primary/15 to-background">
      <div className="container py-12">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
          {/* Company Info */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2">
              <Image
                src="/favicon.png"
                alt="BusTicket.vn"
                width={52}
                height={52}
              />
              <span className="text-lg font-bold">
                BusTicket<span className="text-primary">.vn</span>
              </span>
            </div>
            <p className="text-sm text-neutral-400">
              Nền tảng đặt vé xe khách hàng đầu Việt Nam. Đặt vé nhanh, an toàn
              và tiện lợi.
            </p>
          </div>

          {/* Quick Links */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold tracking-wider uppercase">
              Liên kết nhanh
            </h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/about"
                  className="transition-colors hover:text-primary"
                >
                  Về chúng tôi
                </Link>
              </li>
              <li>
                <Link
                  href="/routes"
                  className="transition-colors hover:text-primary"
                >
                  Tuyến đường phổ biến
                </Link>
              </li>
              <li>
                <Link
                  href="/operators"
                  className="transition-colors hover:text-primary"
                >
                  Nhà xe
                </Link>
              </li>
              <li>
                <Link
                  href="/blog"
                  className="transition-colors hover:text-primary"
                >
                  Blog
                </Link>
              </li>
            </ul>
          </div>

          {/* Support */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold tracking-wider uppercase">
              Hỗ trợ
            </h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/faq"
                  className="transition-colors hover:text-primary"
                >
                  Câu hỏi thường gặp
                </Link>
              </li>
              <li>
                <Link
                  href="/booking-policy"
                  className="transition-colors hover:text-primary"
                >
                  Chính sách đặt vé
                </Link>
              </li>
              <li>
                <Link
                  href="/refund-policy"
                  className="transition-colors hover:text-primary"
                >
                  Chính sách hoàn tiền
                </Link>
              </li>
              <li>
                <Link
                  href="/terms"
                  className="transition-colors hover:text-primary"
                >
                  Điều khoản dịch vụ
                </Link>
              </li>
              <li>
                <Link
                  href="/privacy"
                  className="transition-colors hover:text-primary"
                >
                  Chính sách bảo mật
                </Link>
              </li>
            </ul>
          </div>

          {/* Contact */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold tracking-wider uppercase">
              Liên hệ
            </h3>
            <ul className="space-y-3 text-sm">
              <li className="flex items-start space-x-3">
                <Phone className="mt-0.5 h-5 w-5 shrink-0 text-primary" />
                <div>
                  <a
                    href="tel:1900989901"
                    className="transition-colors hover:text-primary"
                  >
                    1900 989 901
                  </a>
                </div>
              </li>
              <li className="flex items-start space-x-3">
                <Mail className="mt-0.5 h-5 w-5 shrink-0 text-primary" />
                <div>
                  <a
                    href="mailto:support@busticket.vn"
                    className="transition-colors hover:text-primary"
                  >
                    support@busticket.vn
                  </a>
                </div>
              </li>
              <li className="flex items-start space-x-3">
                <MapPin className="mt-0.5 h-5 w-5 shrink-0 text-primary" />
                <div>
                  <p>Quận 1, TP. Hồ Chí Minh, Việt Nam</p>
                </div>
              </li>
            </ul>
          </div>
        </div>

        {/* Payment Methods */}
        <div className="mt-12 border-t border-neutral-800 pt-8">
          <div className="flex flex-col items-center justify-between space-y-4 md:flex-row md:space-y-0">
            <div className="text-center md:text-left">
              <p className="text-xs text-neutral-400">
                © 2025 BusTicket.vn. All rights reserved.
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-xs">Phương thức thanh toán:</span>
              <div className="flex h-8 w-12 items-center justify-center rounded border border-neutral-700 text-xs font-bold">
                PayOS
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
