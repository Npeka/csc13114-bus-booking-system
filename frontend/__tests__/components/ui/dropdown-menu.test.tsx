import { render, screen, userEvent } from "@/lib/test-utils";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";

describe("DropdownMenu components", () => {
  describe("Basic DropdownMenu", () => {
    it("should not show menu initially", () => {
      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Open Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem>Item 1</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      expect(screen.queryByText("Item 1")).not.toBeInTheDocument();
    });

    it("should open menu when trigger is clicked", async () => {
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Open Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem>Menu Item</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Open Menu"));
      expect(screen.getByText("Menu Item")).toBeInTheDocument();
    });

    it("should handle menu item clicks", async () => {
      const handleClick = jest.fn();
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onSelect={handleClick}>Action</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Menu"));
      await user.click(screen.getByText("Action"));

      expect(handleClick).toHaveBeenCalled();
    });
  });

  describe("DropdownMenuLabel", () => {
    it("should render menu label", async () => {
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuItem>Profile</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Menu"));
      expect(screen.getByText("My Account")).toBeInTheDocument();
    });
  });

  describe("DropdownMenuSeparator", () => {
    it("should render separator", async () => {
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem>Item 1</DropdownMenuItem>
            <DropdownMenuSeparator data-testid="separator" />
            <DropdownMenuItem>Item 2</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Menu"));
      expect(screen.getByTestId("separator")).toBeInTheDocument();
    });
  });

  describe("DropdownMenuCheckboxItem", () => {
    it("should handle checkbox state", async () => {
      const handleCheckedChange = jest.fn();
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuCheckboxItem
              checked={false}
              onCheckedChange={handleCheckedChange}
            >
              Option
            </DropdownMenuCheckboxItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Menu"));
      await user.click(screen.getByText("Option"));

      expect(handleCheckedChange).toHaveBeenCalledWith(true);
    });
  });

  describe("DropdownMenuRadioGroup", () => {
    it("should handle radio selection", async () => {
      const handleValueChange = jest.fn();
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Menu</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuRadioGroup
              value="option1"
              onValueChange={handleValueChange}
            >
              <DropdownMenuRadioItem value="option1">
                Option 1
              </DropdownMenuRadioItem>
              <DropdownMenuRadioItem value="option2">
                Option 2
              </DropdownMenuRadioItem>
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Menu"));
      await user.click(screen.getByText("Option 2"));

      expect(handleValueChange).toHaveBeenCalledWith("option2");
    });
  });

  describe("Complete DropdownMenu", () => {
    it("should render complete menu with all components", async () => {
      const user = userEvent.setup();

      render(
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button>Options</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>Profile</DropdownMenuItem>
            <DropdownMenuItem>Settings</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>Logout</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      );

      await user.click(screen.getByText("Options"));

      expect(screen.getByText("Account")).toBeInTheDocument();
      expect(screen.getByText("Profile")).toBeInTheDocument();
      expect(screen.getByText("Settings")).toBeInTheDocument();
      expect(screen.getByText("Logout")).toBeInTheDocument();
    });
  });
});
