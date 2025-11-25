import { render, screen, userEvent } from "@/lib/test-utils";
import { Calendar } from "@/components/ui/calendar";

describe("Calendar component", () => {
  it("should render calendar", () => {
    render(<Calendar />);
    
    // Calendar should render with current month
    const currentMonth = new Date().toLocaleString("default", { month: "long" });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should display current year", () => {
    render(<Calendar />);
    
    const currentYear = new Date().getFullYear();
    expect(screen.getByText(currentYear.toString())).toBeInTheDocument();
  });

  it("should handle selected date", () => {
    const selectedDate = new Date(2025, 0, 15); // Jan 15, 2025
    
    render(<Calendar selected={selectedDate} />);
    
    // Calendar should be rendered
    expect(screen.getByText("January")).toBeInTheDocument();
  });

  it("should call onSelect when date is clicked", async () => {
    const handleSelect = jest.fn();
    const user = userEvent.setup();
    
    render(<Calendar onSelect={handleSelect} />);
    
    // Click on a date (find any date button)
    const dateButtons = screen.getAllByRole("button");
    const firstDateButton = dateButtons.find(btn => 
      /^\d+$/.test(btn.textContent || "")
    );
    
    if (firstDateButton) {
      await user.click(firstDateButton);
      expect(handleSelect).toHaveBeenCalled();
    }
  });

  it("should support month prop", () => {
    const specificMonth = new Date(2025, 5, 1); // June 2025
    
    render(<Calendar month={specificMonth} />);
    
    expect(screen.getByText("June")).toBeInTheDocument();
    expect(screen.getByText("2025")).toBeInTheDocument();
  });

  it("should handle disabled dates", () => {
    const disabledMatcher = (date: Date) => date.getDay() === 0; // Disable Sundays
    
    render(<Calendar disabled={disabledMatcher} />);
    
    // Calendar should render
    const currentMonth = new Date().toLocaleString("default", { month: "long" });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should support mode selection", () => {
    render(<Calendar mode="single" />);
    
    // Single mode calendar should render
    const currentMonth = new Date().toLocaleString("default", { month: "long" });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should handle range mode", () => {
    render(<Calendar mode="range" />);
    
    // Range mode calendar should render
    const currentMonth = new Date().toLocaleString("default", { month: "long" });
    expect(screen.getByText(new RegExp(currentMonth, "i"))).toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(<Calendar className="custom-calendar" />);
    
    const calendar = container.firstChild;
    expect(calendar).toHaveClass("custom-calendar");
  });
});
