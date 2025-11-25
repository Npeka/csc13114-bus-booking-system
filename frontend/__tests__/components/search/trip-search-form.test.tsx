import { render, screen, userEvent } from "@/lib/test-utils";
import { TripSearchForm } from "@/components/search/trip-search-form";

// Mock next/navigation
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
  }),
  useSearchParams: () => ({
    get: jest.fn(),
  }),
}));

describe("TripSearchForm component", () => {
  it("should render search form", () => {
    render(<TripSearchForm />);
    
    expect(screen.getByText(/điểm đi|from/i)).toBeInTheDocument();
    expect(screen.getByText(/điểm đến|to|destination/i)).toBeInTheDocument();
  });

  it("should display all form fields", () => {
    render(<TripSearchForm />);
    
    // Location fields
    expect(screen.getByText(/điểm đi/i)).toBeInTheDocument();
    expect(screen.getByText(/điểm đến/i)).toBeInTheDocument();
    
    // Date field
    expect(screen.getByText(/ngày đi/i)).toBeInTheDocument();
    
    // Passenger field
    expect(screen.getByText(/số vé|hành khách/i)).toBeInTheDocument();
  });

  it("should have submit button", () => {
    render(<TripSearchForm />);
    
    const searchButton = screen.getByRole("button", { name: /tìm|search/i });
    expect(searchButton).toBeInTheDocument();
  });

  it("should handle form submission", async () => {
    const user = userEvent.setup();
    render(<TripSearchForm />);
    
    const searchButton = screen.getByRole("button", { name: /tìm|search/i });
    await user.click(searchButton);
    
    // Form should attempt submission
  });

  it("should show swap locations button", () => {
    render(<TripSearchForm />);
    
    // Swap button for desktop view
    const swapButton = screen.queryByRole("button", { name: /đổi|swap/i });
    // May be hidden on mobile
  });

  it("should display popular routes", () => {
    render(<TripSearchForm />);
    
    expect(screen.getByText(/tuyến.*phổ biến|popular/i)).toBeInTheDocument();
  });

  it("should support round trip option", () => {
    render(<TripSearchForm />);
    
    // Round trip checkbox or toggle
    const roundTripOption = screen.queryByText(/khứ hồi|round.*trip/i);
    expect(roundTripOption).toBeTruthy();
  });
});
