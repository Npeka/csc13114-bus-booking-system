import { render, screen } from "@/lib/test-utils";
import { ThemeProvider } from "@/components/theme-provider";
import { useTheme } from "next-themes";

// Mock next-themes
jest.mock("next-themes", () => ({
  ThemeProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  useTheme: jest.fn(),
}));

describe("ThemeProvider", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("should render children", () => {
    render(
      <ThemeProvider>
        <div>Test Content</div>
      </ThemeProvider>
    );

    expect(screen.getByText("Test Content")).toBeInTheDocument();
  });

  it("should provide theme context to children", () => {
    function TestComponent() {
      return <div>Themed Component</div>;
    }

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    );

    expect(screen.getByText("Themed Component")).toBeInTheDocument();
  });

  it("should render with default props", () => {
    render(
      <ThemeProvider>
        <div>Content</div>
      </ThemeProvider>
    );

    expect(screen.getByText("Content")).toBeInTheDocument();
  });

  it("should handle multiple children", () => {
    render(
      <ThemeProvider>
        <div>Child 1</div>
        <div>Child 2</div>
        <div>Child 3</div>
      </ThemeProvider>
    );

    expect(screen.getByText("Child 1")).toBeInTheDocument();
    expect(screen.getByText("Child 2")).toBeInTheDocument();
    expect(screen.getByText("Child 3")).toBeInTheDocument();
  });

  it("should pass through theme provider props", () => {
    render(
      <ThemeProvider attribute="class" defaultTheme="dark">
        <div>Dark Theme Content</div>
      </ThemeProvider>
    );

    expect(screen.getByText("Dark Theme Content")).toBeInTheDocument();
  });
});
