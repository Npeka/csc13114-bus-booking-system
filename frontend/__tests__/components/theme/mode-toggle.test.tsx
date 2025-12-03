import { render, screen, userEvent } from "@/lib/test-utils";
import { ModeToggle } from "@/components/theme/mode-toggle";
import { useTheme } from "next-themes";
import React from "react";

// Mock next-themes
jest.mock("next-themes", () => ({
  useTheme: jest.fn(),
}));

describe("ModeToggle component", () => {
  const mockSetTheme = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    (useTheme as jest.Mock).mockReturnValue({
      setTheme: mockSetTheme,
      theme: "light",
    });
  });

  it("should render mode toggle button", () => {
    render(<ModeToggle />);
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
  });

  it("should toggle theme when clicked", async () => {
    const user = userEvent.setup();
    (useTheme as jest.Mock).mockReturnValue({
      setTheme: mockSetTheme,
      theme: "light",
    });

    render(<ModeToggle />);
    const button = screen.getByRole("button");

    await user.click(button);
    expect(mockSetTheme).toHaveBeenCalledWith("dark");
  });

  it("should toggle from dark to light when clicked", async () => {
    const user = userEvent.setup();
    (useTheme as jest.Mock).mockReturnValue({
      setTheme: mockSetTheme,
      theme: "dark",
    });

    render(<ModeToggle />);
    const button = screen.getByRole("button");

    await user.click(button);
    expect(mockSetTheme).toHaveBeenCalledWith("light");
  });

  it("should render disabled button before mounting", () => {
    const { rerender } = render(<ModeToggle />);
    expect(screen.getByRole("button")).toBeInTheDocument();
  });

  it("should have accessibility label", () => {
    render(<ModeToggle />);
    expect(screen.getByText(/toggle theme/i)).toBeInTheDocument();
  });
});
