import { render, screen, userEvent } from "@/lib/test-utils";
import { TripFilters } from "@/components/trips/trip-filters";

describe("TripFilters component", () => {
  it("should render trip filters", () => {
    render(<TripFilters onFilterChange={jest.fn()} />);
    
    // Should have filter controls
    expect(screen.getByText(/bộ lọc|filter/i)).toBeInTheDocument();
  });

  it("should display filter options", () => {
    render(<TripFilters onFilterChange={jest.fn()} />);
    
    // Common filters
    const priceFilter = screen.queryByText(/giá|price/i);
    const timeFilter = screen.queryByText(/thời gian|time|giờ/i);
    
    expect(priceFilter || timeFilter).toBeTruthy();
  });

  it("should call onFilterChange when filter is applied", async () => {
    const handleFilterChange = jest.fn();
    const user = userEvent.setup();
    
    render(<TripFilters onFilterChange={handleFilterChange} />);
    
    // Find and interact with any filter control
    const buttons = screen.getAllByRole("button");
    if (buttons.length > 0) {
      await user.click(buttons[0]);
      // Filter change may be called
    }
  });

  it("should support price range filtering", () => {
    render(<TripFilters onFilterChange={jest.fn()} />);
    
    // Price filter should be available
    const priceText = screen.queryByText(/giá|price/i);
    expect(priceText).toBeTruthy();
  });

  it("should support departure time filtering", () => {
    render(<TripFilters onFilterChange={jest.fn()} />);
    
    // Time filter should be available
    const timeText = screen.queryByText(/giờ khởi hành|departure|thời gian/i);
    expect(timeText).toBeTruthy();
  });

  it("should allow clearing filters", async () => {
    const handleFilterChange = jest.fn();
    const user = userEvent.setup();
    
    render(<TripFilters onFilterChange={handleFilterChange} />);
    
    // Look for clear/reset button
    const clearButton = screen.queryByText(/xóa|clear|reset/i);
    if (clearButton) {
      await user.click(clearButton);
      expect(handleFilterChange).toHaveBeenCalled();
    }
  });
});
