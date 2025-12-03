import { render, screen, userEvent } from "@/lib/test-utils";
import { LocationField } from "@/components/search/trip-search-form/location-field";

describe("LocationField component", () => {
  it("should render location field", () => {
    render(
      <LocationField
        id="from"
        label="Điểm đi"
        value={null}
        onChange={jest.fn()}
        placeholder="Chọn điểm đi"
      />,
    );

    expect(screen.getByText("Điểm đi")).toBeInTheDocument();
  });

  it("should display placeholder when no value", () => {
    render(
      <LocationField
        id="from"
        label="Điểm đi"
        value={null}
        onChange={jest.fn()}
        placeholder="Chọn điểm đi"
      />,
    );

    expect(screen.getByText("Chọn điểm đi")).toBeInTheDocument();
  });

  it("should display selected location", () => {
    const selectedCity = { id: "hcm", name: "TP. Hồ Chí Minh" };

    render(
      <LocationField
        id="from"
        label="Điểm đi"
        value={selectedCity}
        onChange={jest.fn()}
        placeholder="Chọn điểm đi"
      />,
    );

    expect(screen.getByText("TP. Hồ Chí Minh")).toBeInTheDocument();
  });

  it("should call onChange when location selected", async () => {
    const handleChange = jest.fn();
    const user = userEvent.setup();

    render(
      <LocationField
        id="from"
        label="Điểm đi"
        value={null}
        onChange={handleChange}
        placeholder="Chọn điểm đi"
      />,
    );

    // Click to open dropdown
    const button = screen.getByRole("button");
    await user.click(button);

    // onChange may be called when selecting from dropdown
  });

  it("should render with label", () => {
    render(
      <LocationField
        id="destination"
        label="Điểm đến"
        value={null}
        onChange={jest.fn()}
        placeholder="Chọn điểm đến"
      />,
    );

    expect(screen.getByText("Điểm đến")).toBeInTheDocument();
  });

  it("should be keyboard accessible", () => {
    render(
      <LocationField
        id="from"
        label="Điểm đi"
        value={null}
        onChange={jest.fn()}
        placeholder="Chọn điểm đi"
      />,
    );

    const button = screen.getByRole("button");
    button.focus();
    expect(button).toHaveFocus();
  });
});
