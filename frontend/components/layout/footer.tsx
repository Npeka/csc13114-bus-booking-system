import Link from "next/link";
import { Facebook, Mail, Phone, MapPin } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t bg-neutral-900 text-neutral-300">
      <div className="container py-12">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
          {/* Company Info */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2">
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-primary">
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="h-5 w-5 text-white"
                >
                  <rect x="3" y="6" width="18" height="12" rx="2" />
                  <path d="M3 12h18" />
                  <path d="M8 6v6" />
                  <path d="M16 6v6" />
                </svg>
              </div>
              <span className="text-lg font-bold text-white">
                BusTicket<span className="text-brand-primary">.vn</span>
              </span>
            </div>
            <p className="text-sm text-neutral-400">
              Nền tảng đặt vé xe khách hàng đầu Việt Nam. Đặt vé nhanh, an toàn
              và tiện lợi.
            </p>
            <div className="flex space-x-4">
              <a
                href="https://facebook.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-neutral-400 hover:text-brand-primary transition-colors"
              >
                <Facebook className="h-5 w-5" />
                <span className="sr-only">Facebook</span>
              </a>
            </div>
          </div>

          {/* Quick Links */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold uppercase tracking-wider text-white">
              Liên kết nhanh
            </h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/about"
                  className="hover:text-brand-primary transition-colors"
                >
                  Về chúng tôi
                </Link>
              </li>
              <li>
                <Link
                  href="/routes"
                  className="hover:text-brand-primary transition-colors"
                >
                  Tuyến đường phổ biến
                </Link>
              </li>
              <li>
                <Link
                  href="/operators"
                  className="hover:text-brand-primary transition-colors"
                >
                  Nhà xe
                </Link>
              </li>
              <li>
                <Link
                  href="/blog"
                  className="hover:text-brand-primary transition-colors"
                >
                  Blog
                </Link>
              </li>
            </ul>
          </div>

          {/* Support */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold uppercase tracking-wider text-white">
              Hỗ trợ
            </h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/help"
                  className="hover:text-brand-primary transition-colors"
                >
                  Trung tâm trợ giúp
                </Link>
              </li>
              <li>
                <Link
                  href="/faq"
                  className="hover:text-brand-primary transition-colors"
                >
                  Câu hỏi thường gặp
                </Link>
              </li>
              <li>
                <Link
                  href="/booking-policy"
                  className="hover:text-brand-primary transition-colors"
                >
                  Chính sách đặt vé
                </Link>
              </li>
              <li>
                <Link
                  href="/refund-policy"
                  className="hover:text-brand-primary transition-colors"
                >
                  Chính sách hoàn tiền
                </Link>
              </li>
              <li>
                <Link
                  href="/terms"
                  className="hover:text-brand-primary transition-colors"
                >
                  Điều khoản dịch vụ
                </Link>
              </li>
              <li>
                <Link
                  href="/privacy"
                  className="hover:text-brand-primary transition-colors"
                >
                  Chính sách bảo mật
                </Link>
              </li>
            </ul>
          </div>

          {/* Contact */}
          <div className="space-y-4">
            <h3 className="text-sm font-semibold uppercase tracking-wider text-white">
              Liên hệ
            </h3>
            <ul className="space-y-3 text-sm">
              <li className="flex items-start space-x-3">
                <Phone className="h-5 w-5 mt-0.5 shrink-0 text-brand-primary" />
                <div>
                  <p className="font-medium text-white">Hotline</p>
                  <a
                    href="tel:1900989901"
                    className="hover:text-brand-primary transition-colors"
                  >
                    1900 989 901
                  </a>
                </div>
              </li>
              <li className="flex items-start space-x-3">
                <Mail className="h-5 w-5 mt-0.5 shrink-0 text-brand-primary" />
                <div>
                  <p className="font-medium text-white">Email</p>
                  <a
                    href="mailto:support@busticket.vn"
                    className="hover:text-brand-primary transition-colors"
                  >
                    support@busticket.vn
                  </a>
                </div>
              </li>
              <li className="flex items-start space-x-3">
                <MapPin className="h-5 w-5 mt-0.5 shrink-0 text-brand-primary" />
                <div>
                  <p className="font-medium text-white">Địa chỉ</p>
                  <p className="text-neutral-400">
                    Quận 1, TP. Hồ Chí Minh, Việt Nam
                  </p>
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
              <span className="text-xs text-neutral-400">
                Phương thức thanh toán:
              </span>
              <div className="flex items-center space-x-2">
                <div className="flex h-8 w-12 items-center justify-center rounded border border-neutral-700 bg-white text-xs font-bold">
                  MOMO
                </div>
                <div className="flex h-8 w-12 items-center justify-center rounded border border-neutral-700 bg-blue-600 text-xs font-bold text-white">
                  Zalo
                </div>
                <div className="flex h-8 w-12 items-center justify-center rounded border border-neutral-700 bg-white text-xs font-bold">
                  PayOS
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
