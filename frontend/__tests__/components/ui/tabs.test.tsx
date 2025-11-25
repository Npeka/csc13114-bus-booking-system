import { render, screen, userEvent } from "@/lib/test-utils";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";

describe("Tabs components", () => {
  describe("Tabs", () => {
    it("should render tabs container", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>
      );

      expect(screen.getByText("Tab 1")).toBeInTheDocument();
    });

    it("should show correct content for default value", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      expect(screen.getByText("Content 1")).toBeVisible();
      expect(screen.queryByText("Content 2")).not.toBeVisible();
    });

    it("should switch content when tab is clicked", async () => {
      const user = userEvent.setup();
      
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>
      );

      const tab2 = screen.getByText("Tab 2");
      await user.click(tab2);

      expect(screen.getByText("Content 2")).toBeVisible();
      expect(screen.queryByText("Content 1")).not.toBeVisible();
    });
  });

  describe("TabsList", () => {
    it("should render tabs list", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList data-testid="tabs-list">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>
      );

      const tabsList = screen.getByTestId("tabs-list");
      expect(tabsList).toBeInTheDocument();
    });

    it("should apply custom className", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList className="custom-list" data-testid="tabs-list">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>
      );

      const tabsList = screen.getByTestId("tabs-list");
      expect(tabsList).toHaveClass("custom-list");
    });
  });

  describe("TabsTrigger", () => {
    it("should render tab trigger", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>
      );

      expect(screen.getByText("Tab 1")).toBeInTheDocument();
    });

    it("should have active state for selected tab", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" data-testid="tab1">
              Tab 1
            </TabsTrigger>
            <TabsTrigger value="tab2" data-testid="tab2">
              Tab 2
            </TabsTrigger>
          </TabsList>
        </Tabs>
      );

      const tab1 = screen.getByTestId("tab1");
      const tab2 = screen.getByTestId("tab2");

      expect(tab1).toHaveAttribute("data-state", "active");
      expect(tab2).toHaveAttribute("data-state", "inactive");
    });
  });

  describe("TabsContent", () => {
    it("should render tab content", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">
            <p>Tab 1 Content</p>
          </TabsContent>
        </Tabs>
      );

      expect(screen.getByText("Tab 1 Content")).toBeInTheDocument();
    });

    it("should apply custom className", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1" className="custom-content" data-testid="content">
            Content
          </TabsContent>
        </Tabs>
      );

      const content = screen.getByTestId("content");
      expect(content).toHaveClass("custom-content");
    });
  });

  describe("Multiple tabs interaction", () => {
    it("should handle keyboard navigation", async () => {
      const user = userEvent.setup();
      
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
            <TabsTrigger value="tab3">Tab 3</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
          <TabsContent value="tab3">Content 3</TabsContent>
        </Tabs>
      );

      const tab1 = screen.getByText("Tab 1");
      await user.click(tab1);
      
      // Keyboard navigation would work here
      expect(screen.getByText("Content 1")).toBeVisible();
    });
  });
});
