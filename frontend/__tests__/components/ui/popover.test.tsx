import { render, screen, userEvent } from "@/lib/test-utils";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";

describe("Popover components", () => {
  describe("Basic Popover", () => {
    it("should not show popover initially", () => {
      render(
        <Popover>
          <PopoverTrigger asChild>
            <Button>Open</Button>
          </PopoverTrigger>
          <PopoverContent>Popover Content</PopoverContent>
        </Popover>,
      );

      expect(screen.queryByText("Popover Content")).not.toBeInTheDocument();
    });

    it("should open popover when trigger is clicked", async () => {
      const user = userEvent.setup();

      render(
        <Popover>
          <PopoverTrigger asChild>
            <Button>Open Popover</Button>
          </PopoverTrigger>
          <PopoverContent>Popover Content</PopoverContent>
        </Popover>,
      );

      await user.click(screen.getByText("Open Popover"));
      expect(screen.getByText("Popover Content")).toBeInTheDocument();
    });
  });

  describe("Controlled Popover", () => {
    it("should show popover when open prop is true", () => {
      render(
        <Popover open={true}>
          <PopoverTrigger asChild>
            <Button>Trigger</Button>
          </PopoverTrigger>
          <PopoverContent>Controlled Content</PopoverContent>
        </Popover>,
      );

      expect(screen.getByText("Controlled Content")).toBeInTheDocument();
    });

    it("should not show popover when open prop is false", () => {
      render(
        <Popover open={false}>
          <PopoverTrigger asChild>
            <Button>Trigger</Button>
          </PopoverTrigger>
          <PopoverContent>Hidden Content</PopoverContent>
        </Popover>,
      );

      expect(screen.queryByText("Hidden Content")).not.toBeInTheDocument();
    });

    it("should call onOpenChange when popover state changes", async () => {
      const handleOpenChange = jest.fn();
      const user = userEvent.setup();

      render(
        <Popover onOpenChange={handleOpenChange}>
          <PopoverTrigger asChild>
            <Button>Toggle</Button>
          </PopoverTrigger>
          <PopoverContent>Content</PopoverContent>
        </Popover>,
      );

      await user.click(screen.getByText("Toggle"));
      expect(handleOpenChange).toHaveBeenCalledWith(true);
    });
  });

  describe("PopoverContent", () => {
    it("should render popover content with custom className", async () => {
      const user = userEvent.setup();

      render(
        <Popover>
          <PopoverTrigger asChild>
            <Button>Open</Button>
          </PopoverTrigger>
          <PopoverContent
            className="custom-popover"
            data-testid="popover-content"
          >
            Content
          </PopoverContent>
        </Popover>,
      );

      await user.click(screen.getByText("Open"));
      const content = screen.getByTestId("popover-content");
      expect(content).toHaveClass("custom-popover");
    });

    it("should support different alignment options", async () => {
      const user = userEvent.setup();

      render(
        <Popover>
          <PopoverTrigger asChild>
            <Button>Open</Button>
          </PopoverTrigger>
          <PopoverContent align="start">Aligned Content</PopoverContent>
        </Popover>,
      );

      await user.click(screen.getByText("Open"));
      expect(screen.getByText("Aligned Content")).toBeInTheDocument();
    });
  });
});
