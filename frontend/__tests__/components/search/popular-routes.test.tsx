import { render, screen, userEvent } from "@/lib/test-utils";
import { PopularRoutes } from "@/components/search/trip-search-form/popular-routes";

describe("PopularRoutes component", () => {
  it("should render popular routes section", () => {
    render(<PopularRoutes onRouteSelect={jest.fn()} />);

    const heading = screen.getByText(/tuyến phổ biến|popular routes/i);
    expect(heading).toBeInTheDocument();
  });

  it("should display route cards", () => {
    render(<PopularRoutes onRouteSelect={jest.fn()} />);

    // Should have multiple route options
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThan(0);
  });

  it("should call onRouteSelect when route is clicked", async () => {
    const handleRouteSelect = jest.fn();
    const user = userEvent.setup();

    render(<PopularRoutes onRouteSelect={handleRouteSelect} />);

    const buttons = screen.getAllByRole("button");
    if (buttons.length > 0) {
      await user.click(buttons[0]);
      expect(handleRouteSelect).toHaveBeenCalled();
    }
  });

  it("should display route information", () => {
    render(<PopularRoutes onRouteSelect={jest.fn()} />);

    // Common routes in Vietnam
    const routes = screen.queryByText(
      /hà nội|sài gòn|đà nẵng|nha trang|đà lạt/i,
    );
    expect(routes).toBeTruthy();
  });

  it("should be keyboard accessible", async () => {
    const user = userEvent.setup();
    render(<PopularRoutes onRouteSelect={jest.fn()} />);

    const buttons = screen.getAllByRole("button");
    if (buttons.length > 0) {
      buttons[0].focus();
      expect(buttons[0]).toHaveFocus();
    }
  });
});
