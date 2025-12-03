import { render, screen, userEvent } from "@/lib/test-utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  SelectGroup,
  SelectLabel,
} from "@/components/ui/select";

describe("Select components", () => {
  describe("Basic Select", () => {
    it("should render select trigger", () => {
      render(
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="Select option" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="option1">Option 1</SelectItem>
          </SelectContent>
        </Select>,
      );

      expect(screen.getByText("Select option")).toBeInTheDocument();
    });

    it("should display placeholder", () => {
      render(
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="Choose an option" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1">Item 1</SelectItem>
          </SelectContent>
        </Select>,
      );

      expect(screen.getByText("Choose an option")).toBeInTheDocument();
    });

    it("should open dropdown when trigger is clicked", async () => {
      const user = userEvent.setup();

      render(
        <Select>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="option1">Option 1</SelectItem>
            <SelectItem value="option2">Option 2</SelectItem>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      await user.click(trigger);

      expect(screen.getByText("Option 1")).toBeInTheDocument();
      expect(screen.getByText("Option 2")).toBeInTheDocument();
    });

    it("should select an option when clicked", async () => {
      const user = userEvent.setup();
      const handleValueChange = jest.fn();

      render(
        <Select onValueChange={handleValueChange}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="opt1">Option 1</SelectItem>
            <SelectItem value="opt2">Option 2</SelectItem>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      await user.click(trigger);
      await user.click(screen.getByText("Option 1"));

      expect(handleValueChange).toHaveBeenCalledWith("opt1");
    });
  });

  describe("Controlled Select", () => {
    it("should display selected value", () => {
      render(
        <Select value="option1">
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="option1">First Option</SelectItem>
            <SelectItem value="option2">Second Option</SelectItem>
          </SelectContent>
        </Select>,
      );

      expect(screen.getByText("First Option")).toBeInTheDocument();
    });

    it("should handle defaultValue", () => {
      render(
        <Select defaultValue="default">
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="default">Default Option</SelectItem>
            <SelectItem value="other">Other Option</SelectItem>
          </SelectContent>
        </Select>,
      );

      expect(screen.getByText("Default Option")).toBeInTheDocument();
    });
  });

  describe("SelectGroup with SelectLabel", () => {
    it("should render grouped options with labels", async () => {
      const user = userEvent.setup();

      render(
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="Select fruit" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Fruits</SelectLabel>
              <SelectItem value="apple">Apple</SelectItem>
              <SelectItem value="banana">Banana</SelectItem>
            </SelectGroup>
            <SelectGroup>
              <SelectLabel>Vegetables</SelectLabel>
              <SelectItem value="carrot">Carrot</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      await user.click(trigger);

      expect(screen.getByText("Fruits")).toBeInTheDocument();
      expect(screen.getByText("Vegetables")).toBeInTheDocument();
      expect(screen.getByText("Apple")).toBeInTheDocument();
      expect(screen.getByText("Carrot")).toBeInTheDocument();
    });
  });

  describe("Select with disabled state", () => {
    it("should handle disabled select", () => {
      render(
        <Select disabled>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1">Item 1</SelectItem>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      expect(trigger).toBeDisabled();
    });

    it("should handle disabled items", async () => {
      const user = userEvent.setup();

      render(
        <Select>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="enabled">Enabled Item</SelectItem>
            <SelectItem value="disabled" disabled>
              Disabled Item
            </SelectItem>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      await user.click(trigger);

      const disabledItem = screen.getByText("Disabled Item");
      expect(disabledItem.closest('[role="option"]')).toHaveAttribute(
        "data-disabled",
      );
    });
  });

  describe("SelectTrigger variants", () => {
    it("should support different sizes", () => {
      render(
        <Select>
          <SelectTrigger size="sm">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1">Item</SelectItem>
          </SelectContent>
        </Select>,
      );

      const trigger = screen.getByRole("combobox");
      expect(trigger).toHaveAttribute("data-size", "sm");
    });
  });

  describe("Multiple selects on same page", () => {
    it("should handle multiple independent selects", async () => {
      const user = userEvent.setup();

      render(
        <>
          <Select>
            <SelectTrigger aria-label="First select">
              <SelectValue placeholder="First" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1a">First A</SelectItem>
            </SelectContent>
          </Select>
          <Select>
            <SelectTrigger aria-label="Second select">
              <SelectValue placeholder="Second" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="2a">Second A</SelectItem>
            </SelectContent>
          </Select>
        </>,
      );

      const triggers = screen.getAllByRole("combobox");
      expect(triggers).toHaveLength(2);

      await user.click(triggers[0]);
      expect(screen.getByText("First A")).toBeInTheDocument();
    });
  });
});
