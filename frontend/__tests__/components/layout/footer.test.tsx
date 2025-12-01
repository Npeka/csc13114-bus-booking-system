import { render, screen } from "@/lib/test-utils";
import { Footer } from "@/components/layout/footer";

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

describe("Footer component", () => {
  it("should render footer", () => {
    render(<Footer />);

    const footer = screen.getByRole("contentinfo");
    expect(footer).toBeInTheDocument();
  });

  it("should display company name", () => {
    render(<Footer />);

    expect(screen.getByText(/BusTicket\.vn/i)).toBeInTheDocument();
  });

  it("should display copyright information", () => {
    render(<Footer />);

    const currentYear = new Date().getFullYear();
    expect(
      screen.getByText(new RegExp(currentYear.toString())),
    ).toBeInTheDocument();
  });

  it("should render footer links", () => {
    render(<Footer />);

    // Check for common footer links (adjust based on actual footer content)
    const links = screen.getAllByRole("link");
    expect(links.length).toBeGreaterThan(0);
  });

  it("should have proper semantic HTML structure", () => {
    render(<Footer />);

    const footer = screen.getByRole("contentinfo");
    expect(footer.tagName).toBe("FOOTER");
  });
});
