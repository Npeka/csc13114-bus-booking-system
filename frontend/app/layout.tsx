import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/components/auth/auth-provider";
import { ThemeProvider } from "@/components/theme-provider";
import { QueryProvider } from "@/components/providers/query-provider";
import { Toaster } from "sonner";

export const metadata: Metadata = {
  title: {
    default: "BusTicket.vn - Đặt vé xe khách trực tuyến",
    template: "%s | BusTicket.vn",
  },
  description:
    "Đặt vé xe khách trực tuyến nhanh chóng, an toàn và tiện lợi. Hàng trăm tuyến đường khắp Việt Nam với giá cả cạnh tranh.",
  keywords: [
    "đặt vé xe khách",
    "vé xe online",
    "đặt vé xe buýt",
    "bus ticket vietnam",
    "đặt vé trực tuyến",
  ],
  authors: [{ name: "BusTicket.vn Team" }],
  creator: "BusTicket.vn",
  publisher: "BusTicket.vn",
  metadataBase: new URL("https://busticket.vn"),
  openGraph: {
    type: "website",
    locale: "vi_VN",
    url: "https://busticket.vn",
    title: "BusTicket.vn - Đặt vé xe khách trực tuyến",
    description:
      "Đặt vé xe khách trực tuyến nhanh chóng, an toàn và tiện lợi. Hàng trăm tuyến đường khắp Việt Nam.",
    siteName: "BusTicket.vn",
  },
  robots: {
    index: true,
    follow: true,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="vi" suppressHydrationWarning>
      <body
        className="flex min-h-screen flex-col antialiased"
        suppressHydrationWarning
      >
        <QueryProvider>
          <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
            <AuthProvider>{children}</AuthProvider>
          </ThemeProvider>
        </QueryProvider>
        <Toaster position="top-center" />
      </body>
    </html>
  );
}
