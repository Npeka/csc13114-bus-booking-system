import { render, screen } from "@/lib/test-utils";
import { Button } from "@/components/ui/button";

describe("Button component", () => {
  it("should render button with text", () => {
    render(<Button>Click me</Button>);
    expect(
      screen.getByRole("button", { name: /click me/i }),
    ).toBeInTheDocument();
  });

  it("should render button with variant", () => {
    render(<Button variant="outline">Outline Button</Button>);
    const button = screen.getByRole("button", { name: /outline button/i });
    expect(button).toHaveClass("border");
  });

  it("should render button with size", () => {
    render(<Button size="sm">Small Button</Button>);
    const button = screen.getByRole("button", { name: /small button/i });
    expect(button).toHaveClass("h-8");
  });

  it("should handle disabled state", () => {
    render(<Button disabled>Disabled</Button>);
    const button = screen.getByRole("button");
    expect(button).toBeDisabled();
  });

  it("should handle click events", () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click</Button>);
    const button = screen.getByRole("button");
    button.click();
    expect(handleClick).toHaveBeenCalled();
  });
});
