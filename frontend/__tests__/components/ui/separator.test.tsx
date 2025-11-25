import { render, screen } from "@/lib/test-utils";
import { Separator } from "@/components/ui/separator";

describe("Separator component", () => {
  it("should render separator", () => {
    render(<Separator data-testid="separator" />);
    const separator = screen.getByTestId("separator");
    expect(separator).toBeInTheDocument();
  });

  it("should have proper aria role", () => {
    render(<Separator />);
    const separator = screen.getByRole("separator");
    expect(separator).toBeInTheDocument();
  });

  it("should handle horizontal orientation by default", () => {
    render(<Separator data-testid="separator" />);
    const separator = screen.getByTestId("separator");
    expect(separator).toHaveAttribute("data-orientation", "horizontal");
  });

  it("should handle vertical orientation", () => {
    render(<Separator orientation="vertical" data-testid="separator" />);
    const separator = screen.getByTestId("separator");
    expect(separator).toHaveAttribute("data-orientation", "vertical");
  });

  it("should apply custom className", () => {
    render(<Separator className="custom-separator" data-testid="separator" />);
    const separator = screen.getByTestId("separator");
    expect(separator).toHaveClass("custom-separator");
  });

  it("should handle decorative separators", () => {
    render(<Separator decorative data-testid="separator" />);
    const separator = screen.getByTestId("separator");
    // Decorative separators should not have role in accessibility tree
    expect(separator).toHaveAttribute("role", "none");
  });
});
