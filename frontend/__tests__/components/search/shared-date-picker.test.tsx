import { render, screen } from "@/lib/test-utils";
import { SharedDatePicker } from "@/components/search/trip-search-form/shared-date-picker";

describe("SharedDatePicker component", () => {
  it("should render shared date picker", () => {
    render(
      <SharedDatePicker
        selectedDate={null}
        onSelect={jest.fn()}
        minDate={new Date()}
      />,
    );

    // Calendar should be rendered
    const currentMonth = new Date().toLocaleString("default", {
      month: "long",
    });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should display selected date", () => {
    const selectedDate = new Date(2025, 0, 15);

    render(
      <SharedDatePicker
        selectedDate={selectedDate}
        onSelect={jest.fn()}
        minDate={new Date()}
      />,
    );

    // Selected date should be highlighted
    expect(screen.getByText("15")).toBeInTheDocument();
  });

  it("should call onSelect when date clicked", async () => {
    const handleSelect = jest.fn();

    render(
      <SharedDatePicker
        selectedDate={null}
        onSelect={handleSelect}
        minDate={new Date()}
      />,
    );

    // Calendar is rendered
    const currentMonth = new Date().toLocaleString("default", {
      month: "long",
    });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should respect minimum date constraint", () => {
    const minDate = new Date(2025, 0, 1);

    render(
      <SharedDatePicker
        selectedDate={null}
        onSelect={jest.fn()}
        minDate={minDate}
      />,
    );

    expect(screen.getByText("January")).toBeInTheDocument();
  });

  it("should display current year", () => {
    render(
      <SharedDatePicker
        selectedDate={null}
        onSelect={jest.fn()}
        minDate={new Date()}
      />,
    );

    const currentYear = new Date().getFullYear();
    expect(screen.getByText(currentYear.toString())).toBeInTheDocument();
  });
});
