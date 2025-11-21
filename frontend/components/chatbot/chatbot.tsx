"use client";

import { useState, useRef, useEffect } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { MessageCircle, X, Send, Bot } from "lucide-react";
import { cn } from "@/lib/utils";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
  suggestions?: string[];
}

export function ChatBot() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([
    {
      id: "1",
      role: "assistant",
      content:
        "Xin chào! Tôi là trợ lý ảo của BusTicket.vn. Tôi có thể giúp bạn tìm chuyến xe, đặt vé hoặc trả lời các câu hỏi về dịch vụ. Bạn cần hỗ trợ gì?",
      timestamp: new Date(),
      suggestions: [
        "Tìm chuyến xe từ Hà Nội đi Đà Nẵng",
        "Giá vé bao nhiêu?",
        "Chính sách hoàn vé",
        "Liên hệ hỗ trợ",
      ],
    },
  ]);
  const [inputValue, setInputValue] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = async (message: string) => {
    if (!message.trim()) return;

    const userMessage: Message = {
      id: new Date().getTime().toString(),
      role: "user",
      content: message,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setInputValue("");
    setIsTyping(true);

    // Simulate AI response
    setTimeout(() => {
      const response = generateResponse(message);
      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: response.message,
        timestamp: new Date(),
        suggestions: response.suggestions,
      };
      setMessages((prev) => [...prev, assistantMessage]);
      setIsTyping(false);
    }, 1000);
  };

  const handleSuggestionClick = (suggestion: string) => {
    handleSendMessage(suggestion);
  };

  return (
    <>
      {/* Chat Button */}
      {!isOpen && (
        <Button
          onClick={() => setIsOpen(true)}
          className="shadow-elevated fixed right-6 bottom-6 z-50 h-14 w-14 rounded-full"
          size="icon"
        >
          <MessageCircle className="h-6 w-6" />
          <span className="sr-only">Mở chat</span>
        </Button>
      )}

      {/* Chat Window */}
      {isOpen && (
        <Card className="shadow-elevated fixed right-6 bottom-6 z-50 w-96 max-w-[calc(100vw-3rem)] py-0!">
          {/* Header */}
          <div className="flex items-center justify-between rounded-t-lg border-b p-4">
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full">
                <Bot className="h-6 w-6 text-primary" />
              </div>
              <div>
                <p className="font-semibold">Trợ lý ảo</p>
                <div className="flex items-center space-x-1">
                  <div className="bg-success h-2 w-2 animate-pulse rounded-full" />
                  <span className="text-xs">Đang online</span>
                </div>
              </div>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setIsOpen(false)}
              className="h-8 w-8"
            >
              <X className="h-5 w-5" />
            </Button>
          </div>

          {/* Messages */}
          <div className="h-96 space-y-4 overflow-y-auto p-4">
            {messages.map((message) => (
              <div key={message.id}>
                <div
                  className={cn(
                    "flex",
                    message.role === "user" ? "justify-end" : "justify-start",
                  )}
                >
                  <div
                    className={cn(
                      "flex max-w-[80%] space-x-2",
                      message.role === "user"
                        ? "flex-row-reverse space-x-reverse"
                        : "",
                    )}
                  >
                    {message.role !== "user" && (
                      <div
                        className={cn(
                          "flex h-8 w-8 shrink-0 items-center justify-center rounded-full border",
                        )}
                      >
                        <Bot className="h-4 w-4 text-primary" />
                      </div>
                    )}
                    <div>
                      <div
                        className={cn(
                          "rounded-lg p-3",
                          message.role === "user"
                            ? "rounded-br-xs bg-primary"
                            : "border",
                        )}
                      >
                        <p className="text-sm">{message.content}</p>
                      </div>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {message.timestamp.toLocaleTimeString("vi-VN", {
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Suggestions */}
                {message.role === "assistant" && message.suggestions && (
                  <div className="mt-2 flex flex-wrap gap-2 pl-10">
                    {message.suggestions.map((suggestion, index) => (
                      <Badge
                        key={index}
                        variant="secondary"
                        className="cursor-pointer transition-colors hover:bg-primary hover:text-white"
                        onClick={() => handleSuggestionClick(suggestion)}
                      >
                        {suggestion}
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            ))}

            {isTyping && (
              <div className="flex justify-start">
                <div className="flex space-x-2">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full border">
                    <Bot className="h-4 w-4 text-primary" />
                  </div>
                  <div className="rounded-lg border p-3">
                    <div className="flex space-x-1">
                      <div className="h-2 w-2 animate-bounce rounded-full bg-neutral-400" />
                      <div
                        className="h-2 w-2 animate-bounce rounded-full bg-neutral-400"
                        style={{ animationDelay: "0.2s" }}
                      />
                      <div
                        className="h-2 w-2 animate-bounce rounded-full bg-neutral-400"
                        style={{ animationDelay: "0.4s" }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* Input */}
          <div className="rounded-b-lg border-t p-4">
            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleSendMessage(inputValue);
              }}
              className="flex space-x-2"
            >
              <Input
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                placeholder="Nhập tin nhắn..."
                className="flex-1"
                disabled={isTyping}
              />
              <Button
                type="submit"
                size="icon"
                className="bg-primary text-white hover:bg-primary/90"
                disabled={!inputValue.trim() || isTyping}
              >
                <Send className="h-4 w-4" />
              </Button>
            </form>
          </div>
        </Card>
      )}
    </>
  );
}

// Simple response generator (in real app, this would call an AI API)
function generateResponse(message: string): {
  message: string;
  suggestions?: string[];
} {
  const lowerMessage = message.toLowerCase();

  if (
    lowerMessage.includes("tìm") ||
    lowerMessage.includes("chuyến") ||
    lowerMessage.includes("hà nội") ||
    lowerMessage.includes("đà nẵng")
  ) {
    return {
      message:
        "Để tìm chuyến xe, bạn có thể sử dụng form tìm kiếm trên trang chủ. Hoặc cho tôi biết:\n• Điểm đi\n• Điểm đến\n• Ngày khởi hành\n• Số hành khách\n\nTôi sẽ giúp bạn tìm chuyến xe phù hợp!",
      suggestions: ["Tìm xe từ HCM đi Đà Lạt", "Xe đi Sa Pa", "Xe giường nằm"],
    };
  }

  if (lowerMessage.includes("giá") || lowerMessage.includes("bao nhiêu")) {
    return {
      message:
        "Giá vé phụ thuộc vào:\n• Tuyến đường\n• Loại xe (ghế ngồi, giường nằm, limousine)\n• Nhà xe\n• Thời gian đặt vé\n\nGiá dao động từ 120.000đ - 500.000đ cho các tuyến phổ biến. Bạn muốn xem giá cụ thể cho tuyến nào?",
      suggestions: ["HCM - Đà Lạt", "Hà Nội - Đà Nẵng", "HCM - Nha Trang"],
    };
  }

  if (
    lowerMessage.includes("hoàn") ||
    lowerMessage.includes("hủy") ||
    lowerMessage.includes("chính sách")
  ) {
    return {
      message:
        "Chính sách hoàn/hủy vé:\n• Hủy trước 24h: hoàn 70% giá vé\n• Hủy từ 12-24h: hoàn 50%\n• Hủy dưới 12h: không hoàn\n\nLưu ý: Mỗi nhà xe có thể có chính sách khác nhau. Vui lòng kiểm tra kỹ khi đặt vé.",
      suggestions: ["Cách hủy vé", "Đổi chuyến", "Thời gian hoàn tiền"],
    };
  }

  if (
    lowerMessage.includes("liên hệ") ||
    lowerMessage.includes("hotline") ||
    lowerMessage.includes("hỗ trợ")
  ) {
    return {
      message:
        "Bạn có thể liên hệ với chúng tôi qua:\n📞 Hotline: 1900 989 901\n📧 Email: support@busticket.vn\n⏰ Thời gian: 24/7\n\nĐội ngũ của chúng tôi luôn sẵn sàng hỗ trợ bạn!",
      suggestions: ["Gửi email", "Gọi hotline", "FAQ"],
    };
  }

  return {
    message:
      "Cảm ơn bạn đã nhắn tin! Tôi có thể giúp bạn:\n• Tìm và đặt vé xe\n• Kiểm tra giá vé\n• Thông tin chính sách\n• Hỗ trợ và liên hệ\n\nBạn cần giúp gì?",
    suggestions: [
      "Tìm chuyến xe",
      "Xem giá vé",
      "Chính sách hoàn vé",
      "Liên hệ hỗ trợ",
    ],
  };
}
