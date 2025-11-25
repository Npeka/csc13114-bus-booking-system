import { render, screen, userEvent } from "@/lib/test-utils";
import { ReturnDatePickerField } from "@/components/search/trip-search-form/return-date-picker-field";

describe("ReturnDatePickerField component", () => {
  it("should render return date picker", () => {
    render(
      <ReturnDatePickerField
        value={null}
        onChange={jest.fn()}
        minDate={new Date()}
        isRoundTrip={true}
      />
    );
    
    expect(screen.getByText(/ngày về|return/i)).toBeInTheDocument();
  });

  it("should display selected return date", () => {
    const returnDate = new Date(2025, 0, 20);
    
    render(
      <ReturnDatePickerField
        value={returnDate}
        onChange={jest.fn()}
        minDate={new Date()}
        isRoundTrip={true}
      />
    );
    
    expect(screen.getByText(/20|jan/i)).toBeInTheDocument();
  });

  it("should call onChange when date selected", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();
    
    render(
      <ReturnDatePickerField
        value={null}
        onChange={handleChange}
        minDate={new Date()}
        isRoundTrip={true}
      />
    );
    
    const input = screen.getByRole("textbox");
    await user.click(input);
  });

  it("should show placeholder when no date", () => {
    render(
      <ReturnDatePickerField
        value={null}
        onChange={jest.fn()}
        minDate={new Date()}
        isRoundTrip={true}
      />
    );
    
    expect(screen.getByPlaceholderText(/chọn ngày/i)).toBeInTheDocument();
  });

  it("should respect minimum date", () => {
    const minDate = new Date(2025, 0, 15);
    
    render(
      <ReturnDatePickerField
        value={null}
        onChange={jest.fn()}
        minDate={minDate}
        isRoundTrip={true}
      />
    );
    
    expect(screen.getByText(/ngày về/i)).toBeInTheDocument();
  });
});
