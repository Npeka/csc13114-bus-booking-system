"use client";

import { useState, useRef, useEffect } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { MessageCircle, X, Send, Bot } from "lucide-react";
import { cn } from "@/lib/utils";
import { sendChatMessage, ChatbotTripData } from "@/lib/api/chatbot-service";
import { ChatbotTripList } from "./chatbot-trip-card";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
  suggestions?: string[];
  trips?: ChatbotTripData[];
}

// Helper function to extract trips from various response data structures
function extractTripsFromData(data: unknown): ChatbotTripData[] | undefined {
  if (!data || typeof data !== "object") return undefined;

  const dataObj = data as Record<string, unknown>;

  // Check if data contains trips directly
  if (dataObj.trips && Array.isArray(dataObj.trips)) {
    return dataObj.trips as ChatbotTripData[];
  }

  // Check for nested data.data structure (API wrapper)
  if (dataObj.data && typeof dataObj.data === "object") {
    const nestedData = dataObj.data as Record<string, unknown>;
    if (nestedData.trips && Array.isArray(nestedData.trips)) {
      return nestedData.trips as ChatbotTripData[];
    }
    // Some APIs return trips array directly in data.data
    if (Array.isArray(nestedData)) {
      return nestedData as ChatbotTripData[];
    }
  }

  return undefined;
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

    try {
      // Call real chatbot API
      const response = await sendChatMessage({
        message,
        history: messages
          .filter((m) => m.id !== "1") // Exclude initial greeting
          .map((m) => ({
            role: m.role,
            content: m.content,
          })),
      });

      // Extract trips from response data if available
      const trips = response.data?.trips || extractTripsFromData(response.data);

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content: response.message,
        timestamp: new Date(),
        suggestions: response.suggestions || [],
        trips: trips,
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      console.error("Chatbot error:", error);

      // Show error message to user
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: "assistant",
        content:
          "Xin lỗi, đã có lỗi xảy ra. Vui lòng thử lại sau hoặc liên hệ hỗ trợ.",
        timestamp: new Date(),
        suggestions: ["Tìm chuyến xe", "Liên hệ hỗ trợ"],
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsTyping(false);
    }
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
          className="fixed right-6 bottom-6 z-50 h-14 w-14 rounded-full text-white shadow-elevated"
          size="icon"
        >
          <MessageCircle className="h-6 w-6" />
          <span className="sr-only">Mở chat</span>
        </Button>
      )}

      {/* Chat Window */}
      {isOpen && (
        <Card className="fixed right-6 bottom-6 z-50 w-96 max-w-[calc(100vw-3rem)] py-0! shadow-elevated">
          {/* Header */}
          <div className="flex items-center justify-between rounded-t-[12px] border-b bg-primary p-4 text-white">
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full">
                <Bot className="h-6 w-6" />
              </div>
              <div>
                <p className="font-semibold">Trợ lý ảo</p>
                <div className="flex items-center space-x-1">
                  <div className="h-2 w-2 animate-pulse rounded-full bg-success" />
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
                            ? "rounded-br-xs bg-primary text-white"
                            : "border",
                        )}
                      >
                        <p className="text-sm">{message.content}</p>
                      </div>
                      {/* Trip Cards */}
                      {message.role === "assistant" &&
                        message.trips &&
                        message.trips.length > 0 && (
                          <div className="mt-2">
                            <ChatbotTripList trips={message.trips} />
                          </div>
                        )}
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
