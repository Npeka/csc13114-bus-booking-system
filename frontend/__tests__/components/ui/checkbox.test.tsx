import { render, screen } from "@/lib/test-utils";
import { Checkbox } from "@/components/ui/checkbox";
import React from "react";

describe("Checkbox component", () => {
  it("should render checkbox", () => {
    render(<Checkbox />);
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).toBeInTheDocument();
  });

  it("should handle checked state", () => {
    render(<Checkbox checked={true} />);
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).toBeChecked();
  });

  it("should handle unchecked state", () => {
    render(<Checkbox checked={false} />);
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).not.toBeChecked();
  });

  it("should handle disabled state", () => {
    render(<Checkbox disabled />);
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).toBeDisabled();
  });

  it("should call onCheckedChange when clicked", () => {
    const handleChange = jest.fn();
    render(<Checkbox onCheckedChange={handleChange} />);
    const checkbox = screen.getByRole("checkbox");

    checkbox.click();
    expect(handleChange).toHaveBeenCalledWith(true);
  });

  it("should apply custom className", () => {
    render(<Checkbox className="custom-checkbox" data-testid="checkbox" />);
    const checkbox = screen.getByTestId("checkbox");
    expect(checkbox).toHaveClass("custom-checkbox");
  });

  it("should forward ref correctly", () => {
    const ref = React.createRef<HTMLButtonElement>();
    render(<Checkbox ref={ref} />);
    expect(ref.current).toBeInstanceOf(HTMLButtonElement);
  });

  it("should have proper aria attributes", () => {
    render(<Checkbox aria-label="Accept terms" />);
    const checkbox = screen.getByLabelText("Accept terms");
    expect(checkbox).toBeInTheDocument();
  });
});
