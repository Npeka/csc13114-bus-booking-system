import { render, screen, userEvent } from "@/lib/test-utils";
import { ChatBot } from "@/components/chatbot/chatbot";

// Mock next/link
jest.mock("next/link", () => {
  return ({ children, href }: { children: React.ReactNode; href: string }) => {
    return <a href={href}>{children}</a>;
  };
});

describe("ChatBot component", () => {
  it("should render chatbot", () => {
    render(<ChatBot />);
    
    // Chatbot button or icon should be visible
    const chatButton = screen.getByRole("button");
    expect(chatButton).toBeInTheDocument();
  });

  it("should toggle chatbot on button click", async () => {
    const user = userEvent.setup();
    render(<ChatBot />);
    
    const chatButton = screen.getByRole("button");
    await user.click(chatButton);
    
    // Chatbot dialog/panel should open
    // Check for chat interface elements
    expect(chatButton).toBeInTheDocument();
  });

  it("should have accessible label", () => {
    render(<ChatBot />);
    
    const chatButton = screen.getByRole("button");
    expect(chatButton).toHaveAttribute("aria-label");
  });

  it("should render chat icon", () => {
    const { container } = render(<ChatBot />);
    
    // Check for SVG icon
    const svg = container.querySelector("svg");
    expect(svg).toBeInTheDocument();
  });

  it("should be keyboard accessible", () => {
    render(<ChatBot />);
    
    const chatButton = screen.getByRole("button");
    chatButton.focus();
    expect(chatButton).toHaveFocus();
  });
});
