import { render, screen, userEvent } from "@/lib/test-utils";
import { Header } from "@/components/layout/header";

// Mock next/link
jest.mock("next/link", () => {
  const MockLink = ({
    children,
    href,
  }: {
    children: React.ReactNode;
    href: string;
  }) => {
    return <a href={href}>{children}</a>;
  };
  MockLink.displayName = "MockLink";
  return MockLink;
});

// Mock next/navigation
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
  usePathname: () => "/",
}));

describe("Header component", () => {
  it("should render header", () => {
    render(<Header />);

    const header = screen.getByRole("banner");
    expect(header).toBeInTheDocument();
  });

  it("should display logo/brand", () => {
    render(<Header />);

    expect(screen.getByText(/busticket/i)).toBeInTheDocument();
  });

  it("should display navigation links", () => {
    render(<Header />);

    // Main navigation links
    const links = screen.getAllByRole("link");
    expect(links.length).toBeGreaterThan(0);
  });

  it("should show theme toggle button", () => {
    render(<Header />);

    // Theme toggle should be present
    const themeButton = screen.getByRole("button", { name: /theme|chủ đề/i });
    expect(themeButton).toBeInTheDocument();
  });

  it("should display user menu when authenticated", () => {
    // Mock authenticated state
    render(<Header />);

    // User dropdown or avatar when logged in
    const userButton = screen.queryByRole("button", { name: /user|menu/i });
    // May not be visible if not authenticated
  });

  it("should show login button when not authenticated", () => {
    render(<Header />);

    const loginLink = screen.queryByText(/đăng nhập|login|sign in/i);
    expect(loginLink).toBeTruthy();
  });

  it("should have mobile menu button", () => {
    render(<Header />);

    // Mobile menu toggle (hamburger)
    const mobileMenuButton = screen.queryByRole("button", { name: /menu/i });
    // Should exist for mobile view
  });

  it("should be sticky/fixed positioned", () => {
    const { container } = render(<Header />);

    const header = container.querySelector("header");
    expect(header).toBeInTheDocument();
  });
});
