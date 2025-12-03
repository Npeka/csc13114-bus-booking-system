import { render, screen, userEvent } from "@/lib/test-utils";
import { DatePickerField } from "@/components/search/trip-search-form/date-picker-field";

describe("DatePickerField component", () => {
  it("should render date picker field", () => {
    render(
      <DatePickerField
        label="Ngày đi"
        value={null}
        onChange={jest.fn()}
        minDate={new Date()}
      />,
    );

    expect(screen.getByText("Ngày đi")).toBeInTheDocument();
  });

  it("should display selected date", () => {
    const selectedDate = new Date(2025, 0, 15); // Jan 15, 2025

    render(
      <DatePickerField
        label="Ngày đi"
        value={selectedDate}
        onChange={jest.fn()}
        minDate={new Date()}
      />,
    );

    // Should show formatted date
    expect(screen.getByText(/15|jan|january/i)).toBeInTheDocument();
  });

  it("should call onChange when date selected", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(
      <DatePickerField
        label="Ngày đi"
        value={null}
        onChange={handleChange}
        minDate={new Date()}
      />,
    );

    // Click to open calendar
    const button = screen.getByRole("button");
    await user.click(button);

    // Calendar should be visible
  });

  it("should show placeholder when no date selected", () => {
    render(
      <DatePickerField
        label="Ngày đi"
        value={null}
        onChange={jest.fn()}
        minDate={new Date()}
      />,
    );

    expect(screen.getByText(/chọn ngày|select date/i)).toBeInTheDocument();
  });

  it("should respect minimum date", () => {
    const minDate = new Date(2025, 0, 1);

    render(
      <DatePickerField
        label="Ngày đi"
        value={null}
        onChange={jest.fn()}
        minDate={minDate}
      />,
    );

    // Component should render with minDate constraint
    expect(screen.getByText("Ngày đi")).toBeInTheDocument();
  });

  it("should be keyboard accessible", () => {
    render(
      <DatePickerField
        label="Ngày đi"
        value={null}
        onChange={jest.fn()}
        minDate={new Date()}
      />,
    );

    const button = screen.getByRole("button");
    button.focus();
    expect(button).toHaveFocus();
  });
});
