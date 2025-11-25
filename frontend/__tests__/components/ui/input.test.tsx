import { render, screen, userEvent } from "@/lib/test-utils";
import { Input } from "@/components/ui/input";
import React from "react";

describe("Input component", () => {
  it("should render input element", () => {
    render(<Input />);
    const input = screen.getByRole("textbox");
    expect(input).toBeInTheDocument();
  });

  it("should accept and display value", async () => {
    const user = userEvent.setup();
    render(<Input />);
    const input = screen.getByRole("textbox");

    await user.type(input, "test value");
    expect(input).toHaveValue("test value");
  });

  it("should handle placeholder", () => {
    render(<Input placeholder="Enter text" />);
    const input = screen.getByPlaceholderText("Enter text");
    expect(input).toBeInTheDocument();
  });

  it("should handle disabled state", () => {
    render(<Input disabled />);
    const input = screen.getByRole("textbox");
    expect(input).toBeDisabled();
  });

  it("should handle onChange event", async () => {
    const user = userEvent.setup();
    const handleChange = jest.fn();
    render(<Input onChange={handleChange} />);
    const input = screen.getByRole("textbox");

    await user.type(input, "a");
    expect(handleChange).toHaveBeenCalled();
  });

  it("should apply custom className", () => {
    render(<Input className="custom-class" data-testid="input" />);
    const input = screen.getByTestId("input");
    expect(input).toHaveClass("custom-class");
  });

  it("should handle different input types", () => {
    const { rerender } = render(<Input type="email" data-testid="input" />);
    let input = screen.getByTestId("input");
    expect(input).toHaveAttribute("type", "email");

    rerender(<Input type="password" data-testid="input" />);
    input = screen.getByTestId("input");
    expect(input).toHaveAttribute("type", "password");

    rerender(<Input type="number" data-testid="input" />);
    input = screen.getByTestId("input");
    expect(input).toHaveAttribute("type", "number");
  });

  it("should forward ref correctly", () => {
    const ref = React.createRef<HTMLInputElement>();
    render(<Input ref={ref} />);
    expect(ref.current).toBeInstanceOf(HTMLInputElement);
  });

  it("should handle required attribute", () => {
    render(<Input required data-testid="input" />);
    const input = screen.getByTestId("input");
    expect(input).toBeRequired();
  });

  it("should handle readOnly attribute", () => {
    render(<Input readOnly data-testid="input" />);
    const input = screen.getByTestId("input");
    expect(input).toHaveAttribute("readOnly");
  });
});
