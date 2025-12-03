import { render, screen } from "@/lib/test-utils";
import { Label } from "@/components/ui/label";

describe("Label component", () => {
  it("should render label", () => {
    render(<Label>Label Text</Label>);
    expect(screen.getByText("Label Text")).toBeInTheDocument();
  });

  it("should render as label element", () => {
    render(<Label data-testid="label">Text</Label>);
    const label = screen.getByTestId("label");
    expect(label.tagName).toBe("LABEL");
  });

  it("should handle htmlFor attribute", () => {
    render(<Label htmlFor="input-id">Label</Label>);
    const label = screen.getByText("Label");
    expect(label).toHaveAttribute("for", "input-id");
  });

  it("should apply custom className", () => {
    render(
      <Label className="custom-label" data-testid="label">
        Text
      </Label>,
    );
    const label = screen.getByTestId("label");
    expect(label).toHaveClass("custom-label");
  });

  it("should render children correctly", () => {
    render(
      <Label>
        <span>Required</span> Field
      </Label>,
    );
    expect(screen.getByText("Required")).toBeInTheDocument();
    expect(screen.getByText(/Field/)).toBeInTheDocument();
  });

  it("should associate with input via htmlFor", () => {
    render(
      <>
        <Label htmlFor="test-input">Username</Label>
        <input id="test-input" />
      </>,
    );

    const label = screen.getByText("Username");
    const input = screen.getByRole("textbox");

    expect(label).toHaveAttribute("for", "test-input");
    expect(input).toHaveAttribute("id", "test-input");
  });
});
