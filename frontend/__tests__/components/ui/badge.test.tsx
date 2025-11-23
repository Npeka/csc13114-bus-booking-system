import { render, screen } from "@/lib/test-utils";
import { Badge } from "@/components/ui/badge";

describe("Badge component", () => {
  it("should render badge with text", () => {
    render(<Badge>Test Badge</Badge>);
    expect(screen.getByText("Test Badge")).toBeInTheDocument();
  });

  it("should render badge with variant", () => {
    render(<Badge variant="outline">Outline Badge</Badge>);
    const badge = screen.getByText("Outline Badge");
    expect(badge).toHaveClass("border");
  });

  it("should render badge with default variant", () => {
    render(<Badge>Default Badge</Badge>);
    const badge = screen.getByText("Default Badge");
    expect(badge).toHaveClass("bg-primary");
  });
});
