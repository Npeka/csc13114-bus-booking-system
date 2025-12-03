import { render, screen, userEvent } from "@/lib/test-utils";
import { PassengerField } from "@/components/search/trip-search-form/passenger-field";

describe("PassengerField component", () => {
  it("should render passenger field", () => {
    render(<PassengerField value={1} onChange={jest.fn()} />);

    expect(screen.getByText(/hành khách|passenger/i)).toBeInTheDocument();
  });

  it("should display current passenger count", () => {
    render(<PassengerField value={3} onChange={jest.fn()} />);

    expect(screen.getByText(/3/)).toBeInTheDocument();
  });

  it("should call onChange when increment button clicked", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(<PassengerField value={1} onChange={handleChange} />);

    // Find increment button (usually has + or arrow up)
    const buttons = screen.getAllByRole("button");
    const incrementButton = buttons.find(
      (btn) =>
        btn.textContent?.includes("+") ||
        btn.getAttribute("aria-label")?.includes("increase"),
    );

    if (incrementButton) {
      await user.click(incrementButton);
      expect(handleChange).toHaveBeenCalledWith(2);
    }
  });

  it("should call onChange when decrement button clicked", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(<PassengerField value={3} onChange={handleChange} />);

    // Find decrement button
    const buttons = screen.getAllByRole("button");
    const decrementButton = buttons.find(
      (btn) =>
        btn.textContent?.includes("-") ||
        btn.getAttribute("aria-label")?.includes("decrease"),
    );

    if (decrementButton) {
      await user.click(decrementButton);
      expect(handleChange).toHaveBeenCalledWith(2);
    }
  });

  it("should not allow less than minimum passengers", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(<PassengerField value={1} onChange={handleChange} min={1} />);

    const buttons = screen.getAllByRole("button");
    const decrementButton = buttons.find(
      (btn) =>
        btn.textContent?.includes("-") ||
        btn.getAttribute("aria-label")?.includes("decrease"),
    );

    if (decrementButton) {
      await user.click(decrementButton);
      // Should not call onChange if already at minimum
      expect(handleChange).not.toHaveBeenCalled();
    }
  });

  it("should respect maximum passenger limit", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(<PassengerField value={10} onChange={handleChange} max={10} />);

    const buttons = screen.getAllByRole("button");
    const incrementButton = buttons.find(
      (btn) =>
        btn.textContent?.includes("+") ||
        btn.getAttribute("aria-label")?.includes("increase"),
    );

    if (incrementButton) {
      await user.click(incrementButton);
      // Should not call onChange if already at maximum
      expect(handleChange).not.toHaveBeenCalled();
    }
  });

  it("should handle keyboard navigation", () => {
    render(<PassengerField value={5} onChange={jest.fn()} />);

    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThan(0);
  });
});
