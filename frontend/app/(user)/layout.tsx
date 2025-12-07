import { Header } from "@/components/layout/header";
import { Footer } from "@/components/layout/footer";
import { ChatBot } from "@/components/chatbot/chatbot";
import { HydrationGuard } from "@/components/auth/hydration-guard";

export default function UserLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <HydrationGuard>
        <Header />
      </HydrationGuard>
      <main className="flex-1">{children}</main>
      <Footer />
      <ChatBot />
    </>
  );
}
