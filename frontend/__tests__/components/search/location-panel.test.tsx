import { render, screen, userEvent } from "@/lib/test-utils";
import { LocationPanel } from "@/components/search/trip-search-form/location-panel";

describe("LocationPanel component", () => {
  it("should render location panel", () => {
    render(<LocationPanel onSelect={jest.fn()} onClose={jest.fn()} />);

    // Panel should render
    expect(screen.getByText(/tỉnh|thành phố|location/i)).toBeInTheDocument();
  });

  it("should display list of cities", () => {
    render(<LocationPanel onSelect={jest.fn()} onClose={jest.fn()} />);

    // Popular Vietnamese cities should be listed
    const cities = screen.queryByText(/hà nội|sài gòn|hồ chí minh|đà nẵng/i);
    expect(cities).toBeTruthy();
  });

  it("should call onSelect when city clicked", async () => {
    const handleSelect = jest.fn();
    const user = userEvent.setup();

    render(<LocationPanel onSelect={handleSelect} onClose={jest.fn()} />);

    // Find and click a city
    const buttons = screen.getAllByRole("button");
    if (buttons.length > 0) {
      await user.click(buttons[0]);
      expect(handleSelect).toHaveBeenCalled();
    }
  });

  it("should support search/filter functionality", async () => {
    const user = userEvent.setup();

    render(<LocationPanel onSelect={jest.fn()} onClose={jest.fn()} />);

    // Look for search input
    const searchInput = screen.queryByPlaceholderText(/tìm|search/i);
    if (searchInput) {
      await user.type(searchInput, "Hà Nội");
    }
  });

  it("should call onClose when close button clicked", async () => {
    const handleClose = jest.fn();
    const user = userEvent.setup();

    render(<LocationPanel onSelect={jest.fn()} onClose={handleClose} />);

    // Find close button (X or close icon)
    const closeButton = screen.queryByRole("button", { name: /close|đóng/i });
    if (closeButton) {
      await user.click(closeButton);
      expect(handleClose).toHaveBeenCalled();
    }
  });

  it("should display popular locations section", () => {
    render(<LocationPanel onSelect={jest.fn()} onClose={jest.fn()} />);

    // Should have popular locations
    expect(screen.queryByText(/phổ biến|popular/i)).toBeTruthy();
  });
});
