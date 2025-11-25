import { render, screen, userEvent } from "@/lib/test-utils";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

describe("Dialog components", () => {
  describe("Dialog with trigger", () => {
    it("should not show dialog content initially", () => {
      render(
        <Dialog>
          <DialogTrigger asChild>
            <Button>Open Dialog</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Dialog Title</DialogTitle>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      expect(screen.queryByText("Dialog Title")).not.toBeInTheDocument();
    });

    it("should open dialog when trigger is clicked", async () => {
      const user = userEvent.setup();
      
      render(
        <Dialog>
          <DialogTrigger asChild>
            <Button>Open Dialog</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Dialog Title</DialogTitle>
              <DialogDescription>Dialog Description</DialogDescription>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      const trigger = screen.getByText("Open Dialog");
      await user.click(trigger);

      expect(screen.getByText("Dialog Title")).toBeInTheDocument();
      expect(screen.getByText("Dialog Description")).toBeInTheDocument();
    });
  });

  describe("Controlled Dialog", () => {
    it("should show dialog when open prop is true", () => {
      render(
        <Dialog open={true}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Controlled Dialog</DialogTitle>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      expect(screen.getByText("Controlled Dialog")).toBeInTheDocument();
    });

    it("should not show dialog when open prop is false", () => {
      render(
        <Dialog open={false}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Hidden Dialog</DialogTitle>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      expect(screen.queryByText("Hidden Dialog")).not.toBeInTheDocument();
    });

    it("should call onOpenChange when dialog state changes", async () => {
      const handleOpenChange = jest.fn();
      const user = userEvent.setup();

      render(
        <Dialog onOpenChange={handleOpenChange}>
          <DialogTrigger asChild>
            <Button>Open</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogTitle>Test</DialogTitle>
          </DialogContent>
        </Dialog>
      );

      await user.click(screen.getByText("Open"));
      expect(handleOpenChange).toHaveBeenCalledWith(true);
    });
  });

  describe("DialogHeader", () => {
    it("should render dialog header", () => {
      render(
        <Dialog open={true}>
          <DialogContent>
            <DialogHeader data-testid="dialog-header">
              <DialogTitle>Title</DialogTitle>
            </DialogHeader>
          </DialogContent>
        </Dialog>
      );

      expect(screen.getByTestId("dialog-header")).toBeInTheDocument();
    });
  });

  describe("DialogFooter", () => {
    it("should render dialog footer", () => {
      render(
        <Dialog open={true}>
          <DialogContent>
            <DialogFooter data-testid="dialog-footer">
              <Button>Cancel</Button>
              <Button>Confirm</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      );

      expect(screen.getByTestId("dialog-footer")).toBeInTheDocument();
      expect(screen.getByText("Cancel")).toBeInTheDocument();
      expect(screen.getByText("Confirm")).toBeInTheDocument();
    });
  });

  describe("Complete Dialog", () => {
    it("should render complete dialog with all components", async () => {
      const user = userEvent.setup();

      render(
        <Dialog>
          <DialogTrigger asChild>
            <Button>Open Complete Dialog</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Complete Dialog Title</DialogTitle>
              <DialogDescription>This is a description</DialogDescription>
            </DialogHeader>
            <div>Dialog body content</div>
            <DialogFooter>
              <Button variant="outline">Cancel</Button>
              <Button>Confirm</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      );

      await user.click(screen.getByText("Open Complete Dialog"));

      expect(screen.getByText("Complete Dialog Title")).toBeInTheDocument();
      expect(screen.getByText("This is a description")).toBeInTheDocument();
      expect(screen.getByText("Dialog body content")).toBeInTheDocument();
      expect(screen.getByText("Cancel")).toBeInTheDocument();
      expect(screen.getByText("Confirm")).toBeInTheDocument();
    });
  });
});
