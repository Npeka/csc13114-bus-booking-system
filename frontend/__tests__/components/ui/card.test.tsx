import { render, screen } from "@/lib/test-utils";
import {
  Card,
  CardHeader,
  CardFooter,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";

describe("Card components", () => {
  describe("Card", () => {
    it("should render card", () => {
      render(<Card data-testid="card">Card content</Card>);
      const card = screen.getByTestId("card");
      expect(card).toBeInTheDocument();
    });

    it("should apply custom className", () => {
      render(
        <Card className="custom-card" data-testid="card">
          Content
        </Card>
      );
      const card = screen.getByTestId("card");
      expect(card).toHaveClass("custom-card");
    });

    it("should render children", () => {
      render(
        <Card>
          <div>Child content</div>
        </Card>
      );
      expect(screen.getByText("Child content")).toBeInTheDocument();
    });
  });

  describe("CardHeader", () => {
    it("should render card header", () => {
      render(<CardHeader data-testid="header">Header</CardHeader>);
      const header = screen.getByTestId("header");
      expect(header).toBeInTheDocument();
    });

    it("should apply custom className", () => {
      render(
        <CardHeader className="custom-header" data-testid="header">
          Header
        </CardHeader>
      );
      const header = screen.getByTestId("header");
      expect(header).toHaveClass("custom-header");
    });
  });

  describe("CardTitle", () => {
    it("should render card title", () => {
      render(<CardTitle>Title Text</CardTitle>);
      expect(screen.getByText("Title Text")).toBeInTheDocument();
    });

    it("should render as h3 by default", () => {
      render(<CardTitle data-testid="title">Title</CardTitle>);
      const title = screen.getByTestId("title");
      expect(title.tagName).toBe("H3");
    });
  });

  describe("CardDescription", () => {
    it("should render card description", () => {
      render(<CardDescription>Description text</CardDescription>);
      expect(screen.getByText("Description text")).toBeInTheDocument();
    });

    it("should render as p element", () => {
      render(<CardDescription data-testid="desc">Description</CardDescription>);
      const desc = screen.getByTestId("desc");
      expect(desc.tagName).toBe("P");
    });
  });

  describe("CardContent", () => {
    it("should render card content", () => {
      render(<CardContent data-testid="content">Content</CardContent>);
      const content = screen.getByTestId("content");
      expect(content).toBeInTheDocument();
    });

    it("should render children", () => {
      render(
        <CardContent>
          <p>Paragraph 1</p>
          <p>Paragraph 2</p>
        </CardContent>
      );
      expect(screen.getByText("Paragraph 1")).toBeInTheDocument();
      expect(screen.getByText("Paragraph 2")).toBeInTheDocument();
    });
  });

  describe("CardFooter", () => {
    it("should render card footer", () => {
      render(<CardFooter data-testid="footer">Footer</CardFooter>);
      const footer = screen.getByTestId("footer");
      expect(footer).toBeInTheDocument();
    });

    it("should apply custom className", () => {
      render(
        <CardFooter className="custom-footer" data-testid="footer">
          Footer
        </CardFooter>
      );
      const footer = screen.getByTestId("footer");
      expect(footer).toHaveClass("custom-footer");
    });
  });

  describe("Complete Card", () => {
    it("should render complete card with all components", () => {
      render(
        <Card>
          <CardHeader>
            <CardTitle>Card Title</CardTitle>
            <CardDescription>Card Description</CardDescription>
          </CardHeader>
          <CardContent>
            <p>Card Content</p>
          </CardContent>
          <CardFooter>Card Footer</CardFooter>
        </Card>
      );

      expect(screen.getByText("Card Title")).toBeInTheDocument();
      expect(screen.getByText("Card Description")).toBeInTheDocument();
      expect(screen.getByText("Card Content")).toBeInTheDocument();
      expect(screen.getByText("Card Footer")).toBeInTheDocument();
    });
  });
});
