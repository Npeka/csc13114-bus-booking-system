import { render, screen } from "@/lib/test-utils";
import { Skeleton } from "@/components/ui/skeleton";

describe("Skeleton component", () => {
  it("should render skeleton", () => {
    render(<Skeleton data-testid="skeleton" />);
    const skeleton = screen.getByTestId("skeleton");
    expect(skeleton).toBeInTheDocument();
  });

  it("should have skeleton class", () => {
    render(<Skeleton data-testid="skeleton" />);
    const skeleton = screen.getByTestId("skeleton");
    expect(skeleton).toHaveClass("animate-pulse");
  });

  it("should apply custom className", () => {
    render(<Skeleton className="custom-skeleton" data-testid="skeleton" />);
    const skeleton = screen.getByTestId("skeleton");
    expect(skeleton).toHaveClass("custom-skeleton");
  });

  it("should render as div by default", () => {
    render(<Skeleton data-testid="skeleton" />);
    const skeleton = screen.getByTestId("skeleton");
    expect(skeleton.tagName).toBe("DIV");
  });

  it("should render multiple skeletons", () => {
    render(
      <>
        <Skeleton data-testid="skeleton-1" />
        <Skeleton data-testid="skeleton-2" />
        <Skeleton data-testid="skeleton-3" />
      </>,
    );

    expect(screen.getAllByTestId(/skeleton-/)).toHaveLength(3);
  });
});
