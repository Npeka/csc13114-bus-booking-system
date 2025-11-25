import { render, screen } from "@/lib/test-utils";
import { TripCard } from "@/components/trips/trip-card";

// Mock trip data
const mockTrip = {
  id: "trip-1",
  operator_id: "op-1",
  route_id: "route-1",
  bus_id: "bus-1",
  departure_time: "2025-01-15T08:00:00Z",
  arrival_time: "2025-01-15T12:00:00Z",
  price: 250000,
  available_seats: 15,
  total_seats: 40,
  status: "active" as const,
  operator: {
    id: "op-1",
    name: "Phương Trang",
    logo_url: "/logos/phuong-trang.png",
  },
  route: {
    id: "route-1",
    origin: "TP. Hồ Chí Minh",
    destination: "Đà Lạt",
    distance_km: 300,
  },
};

describe("TripCard component", () => {
  it("should render trip card", () => {
    render(<TripCard trip={mockTrip} />);
    
    expect(screen.getByText(mockTrip.operator.name)).toBeInTheDocument();
  });

  it("should display operator name", () => {
    render(<TripCard trip={mockTrip} />);
    
    expect(screen.getByText("Phương Trang")).toBeInTheDocument();
  });

  it("should display route information", () => {
    render(<TripCard trip={mockTrip} />);
    
    expect(screen.getByText(/TP\. Hồ Chí Minh/i)).toBeInTheDocument();
    expect(screen.getByText(/Đà Lạt/i)).toBeInTheDocument();
  });

  it("should display price", () => {
    render(<TripCard trip={mockTrip} />);
    
    // Price formatted: 250.000₫
    expect(screen.getByText(/250[.,]000/)).toBeInTheDocument();
  });

  it("should display available seats", () => {
    render(<TripCard trip={mockTrip} />);
    
    expect(screen.getByText(/15/)).toBeInTheDocument();
  });

  it("should display departure and arrival times", () => {
    render(<TripCard trip={mockTrip} />);
    
    // Times should be formatted and displayed
    expect(screen.getByText(/08:00|8:00/)).toBeInTheDocument();
    expect(screen.getByText(/12:00/)).toBeInTheDocument();
  });

  it("should show sold out status when no seats available", () => {
    const soldOutTrip = {
      ...mockTrip,
      available_seats: 0,
    };
    
    render(<TripCard trip={soldOutTrip} />);
    
    expect(screen.getByText(/hết chỗ|sold out/i)).toBeInTheDocument();
  });

  it("should handle missing operator logo", () => {
    const tripWithoutLogo = {
      ...mockTrip,
      operator: {
        ...mockTrip.operator,
        logo_url: undefined,
      },
    };
    
    render(<TripCard trip={tripWithoutLogo} />);
    
    expect(screen.getByText("Phương Trang")).toBeInTheDocument();
  });
});
