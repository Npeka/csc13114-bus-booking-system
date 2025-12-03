import { render, screen, userEvent } from "@/lib/test-utils";
import { SeatMap } from "@/components/trips/seat-map";

// Mock seat data
const mockSeats = [
  { id: "A1", number: "A1", status: "available", price: 250000 },
  { id: "A2", number: "A2", status: "selected", price: 250000 },
  { id: "A3", number: "A3", status: "occupied", price: 250000 },
  { id: "B1", number: "B1", status: "available", price: 250000 },
];

describe("SeatMap component", () => {
  it("should render seat map", () => {
    render(
      <SeatMap seats={mockSeats} selectedSeats={[]} onSeatSelect={jest.fn()} />,
    );

    // Seat map should render
    expect(
      screen.getByText(/sơ đồ ghế|seat map|chọn ghế/i),
    ).toBeInTheDocument();
  });

  it("should display available seats", () => {
    render(
      <SeatMap seats={mockSeats} selectedSeats={[]} onSeatSelect={jest.fn()} />,
    );

    expect(screen.getByText("A1")).toBeInTheDocument();
    expect(screen.getByText("B1")).toBeInTheDocument();
  });

  it("should show selected seats", () => {
    render(
      <SeatMap
        seats={mockSeats}
        selectedSeats={["A2"]}
        onSeatSelect={jest.fn()}
      />,
    );

    expect(screen.getByText("A2")).toBeInTheDocument();
  });

  it("should call onSeatSelect when clicking available seat", async () => {
    const handleSeatSelect = jest.fn();
    const user = userEvent.setup();

    render(
      <SeatMap
        seats={mockSeats}
        selectedSeats={[]}
        onSeatSelect={handleSeatSelect}
      />,
    );

    const seatA1 = screen.getByText("A1");
    await user.click(seatA1);

    expect(handleSeatSelect).toHaveBeenCalledWith("A1");
  });

  it("should not allow selecting occupied seats", async () => {
    const handleSeatSelect = jest.fn();
    const user = userEvent.setup();

    render(
      <SeatMap
        seats={mockSeats}
        selectedSeats={[]}
        onSeatSelect={handleSeatSelect}
      />,
    );

    const occupiedSeat = screen.getByText("A3");
    await user.click(occupiedSeat);

    // Should not call handler for occupied seat
    expect(handleSeatSelect).not.toHaveBeenCalled();
  });

  it("should display seat legend", () => {
    render(
      <SeatMap seats={mockSeats} selectedSeats={[]} onSeatSelect={jest.fn()} />,
    );

    // Legend should show seat status indicators
    expect(screen.getByText(/trống|available/i)).toBeInTheDocument();
  });

  it("should show total price for selected seats", () => {
    render(
      <SeatMap
        seats={mockSeats}
        selectedSeats={["A1", "B1"]}
        onSeatSelect={jest.fn()}
      />,
    );

    // Should calculate total: 250,000 * 2 = 500,000
    expect(screen.getByText(/500[.,]000|tổng/i)).toBeInTheDocument();
  });
});
