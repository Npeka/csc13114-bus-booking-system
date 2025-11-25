import { render, screen, userEvent } from "@/lib/test-utils";
import { SwapLocationsButton } from "@/components/search/trip-search-form/swap-locations-button";

describe("SwapLocationsButton component", () => {
  it("should render swap button", () => {
    render(<SwapLocationsButton onSwap={jest.fn()} />);
    
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
  });

  it("should call onSwap when clicked", async () => {
    const handleSwap = jest.fn();
    const user = userEvent.setup();
    
    render(<SwapLocationsButton onSwap={handleSwap} />);
    
    const button = screen.getByRole("button");
    await user.click(button);
    
    expect(handleSwap).toHaveBeenCalledTimes(1);
  });

  it("should have accessible label", () => {
    render(<SwapLocationsButton onSwap={jest.fn()} />);
    
    const button = screen.getByRole("button");
    expect(button).toHaveAttribute("aria-label");
  });

  it("should be keyboard accessible", async () => {
    const handleSwap = jest.fn();
    const user = userEvent.setup();
    
    render(<SwapLocationsButton onSwap={handleSwap} />);
    
    const button = screen.getByRole("button");
    button.focus();
    await user.keyboard("{Enter}");
    
    expect(handleSwap).toHaveBeenCalled();
  });

  it("should handle multiple clicks", async () => {
    const handleSwap = jest.fn();
    const user = userEvent.setup();
    
    render(<SwapLocationsButton onSwap={handleSwap} />);
    
    const button = screen.getByRole("button");
    await user.click(button);
    await user.click(button);
    await user.click(button);
    
    expect(handleSwap).toHaveBeenCalledTimes(3);
  });

  it("should render icon", () => {
    const { container } = render(<SwapLocationsButton onSwap={jest.fn()} />);
    
    // Check for SVG icon
    const svg = container.querySelector("svg");
    expect(svg).toBeInTheDocument();
  });
});
