import { render, screen, userEvent } from "@/lib/test-utils";
import {
  Sheet,
  SheetTrigger,
  SheetContent,
  SheetHeader,
  SheetFooter,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";

describe("Sheet components", () => {
  describe("Basic Sheet", () => {
    it("should not show sheet initially", () => {
      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open Sheet</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetTitle>Sheet Title</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      expect(screen.queryByText("Sheet Title")).not.toBeInTheDocument();
    });

    it("should open sheet when trigger is clicked", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open Sheet</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetTitle>Sheet Title</SheetTitle>
            <div>Sheet Body</div>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open Sheet"));
      expect(screen.getByText("Sheet Title")).toBeInTheDocument();
      expect(screen.getByText("Sheet Body")).toBeInTheDocument();
    });
  });

  describe("Controlled Sheet", () => {
    it("should show sheet when open prop is true", () => {
      render(
        <Sheet open={true}>
          <SheetContent>
            <SheetTitle>Controlled Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      expect(screen.getByText("Controlled Sheet")).toBeInTheDocument();
    });

    it("should not show sheet when open prop is false", () => {
      render(
        <Sheet open={false}>
          <SheetContent>
            <SheetTitle>Hidden Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      expect(screen.queryByText("Hidden Sheet")).not.toBeInTheDocument();
    });

    it("should call onOpenChange when sheet state changes", async () => {
      const handleOpenChange = jest.fn();
      const user = userEvent.setup();

      render(
        <Sheet onOpenChange={handleOpenChange}>
          <SheetTrigger asChild>
            <Button>Toggle</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetTitle>Content</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Toggle"));
      expect(handleOpenChange).toHaveBeenCalledWith(true);
    });
  });

  describe("Sheet sides", () => {
    it("should support right side (default)", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent side="right">
            <SheetTitle>Right Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByText("Right Sheet")).toBeInTheDocument();
    });

    it("should support left side", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent side="left">
            <SheetTitle>Left Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByText("Left Sheet")).toBeInTheDocument();
    });

    it("should support top side", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent side="top">
            <SheetTitle>Top Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByText("Top Sheet")).toBeInTheDocument();
    });

    it("should support bottom side", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent side="bottom">
            <SheetTitle>Bottom Sheet</SheetTitle>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByText("Bottom Sheet")).toBeInTheDocument();
    });
  });

  describe("Sheet subcomponents", () => {
    it("should render SheetHeader", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader data-testid="sheet-header">
              <SheetTitle>Title</SheetTitle>
              <SheetDescription>Description</SheetDescription>
            </SheetHeader>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByTestId("sheet-header")).toBeInTheDocument();
      expect(screen.getByText("Title")).toBeInTheDocument();
      expect(screen.getByText("Description")).toBeInTheDocument();
    });

    it("should render SheetFooter", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetTitle>Title</SheetTitle>
            <SheetFooter data-testid="sheet-footer">
              <Button>Cancel</Button>
              <Button>Save</Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByTestId("sheet-footer")).toBeInTheDocument();
      expect(screen.getAllByRole("button")).toHaveLength(3); // Open + Cancel + Save
    });
  });

  describe("Complete Sheet", () => {
    it("should render complete sheet with all components", async () => {
      const user = userEvent.setup();

      render(
        <Sheet>
          <SheetTrigger asChild>
            <Button>Open Complete Sheet</Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Complete Sheet Title</SheetTitle>
              <SheetDescription>This is a description</SheetDescription>
            </SheetHeader>
            <div>Sheet body content</div>
            <SheetFooter>
              <Button variant="outline">Cancel</Button>
              <Button>Confirm</Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>,
      );

      await user.click(screen.getByText("Open Complete Sheet"));

      expect(screen.getByText("Complete Sheet Title")).toBeInTheDocument();
      expect(screen.getByText("This is a description")).toBeInTheDocument();
      expect(screen.getByText("Sheet body content")).toBeInTheDocument();
      expect(screen.getByText("Cancel")).toBeInTheDocument();
      expect(screen.getByText("Confirm")).toBeInTheDocument();
    });
  });
});
