import { render, screen, userEvent } from "@/lib/test-utils";
import { useForm } from "react-hook-form";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

// Test component that uses the form
function TestForm({ onSubmit = jest.fn() }) {
  const form = useForm({
    defaultValues: {
      email: "",
      username: "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="email"
          rules={{ required: "Email is required" }}
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="email@example.com" {...field} />
              </FormControl>
              <FormDescription>Enter your email address</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}

describe("Form components", () => {
  describe("Form with FormField", () => {
    it("should render form fields", () => {
      render(<TestForm />);
      
      expect(screen.getByLabelText("Email")).toBeInTheDocument();
      expect(screen.getByLabelText("Username")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /submit/i })).toBeInTheDocument();
    });

    it("should render FormLabel", () => {
      render(<TestForm />);
      
      expect(screen.getByText("Email")).toBeInTheDocument();
      expect(screen.getByText("Username")).toBeInTheDocument();
    });

    it("should render FormDescription", () => {
      render(<TestForm />);
      
      expect(screen.getByText("Enter your email address")).toBeInTheDocument();
    });

    it("should accept input values", async () => {
      const user = userEvent.setup();
      render(<TestForm />);
      
      const emailInput = screen.getByLabelText("Email");
      await user.type(emailInput, "test@example.com");
      
      expect(emailInput).toHaveValue("test@example.com");
    });

    it("should handle form submission", async () => {
      const handleSubmit = jest.fn();
      const user = userEvent.setup();
      
      render(<TestForm onSubmit={handleSubmit} />);
      
      const emailInput = screen.getByLabelText("Email");
      const usernameInput = screen.getByLabelText("Username");
      
      await user.type(emailInput, "test@example.com");
      await user.type(usernameInput, "testuser");
      await user.click(screen.getByRole("button", { name: /submit/i }));
      
      expect(handleSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          email: "test@example.com",
          username: "testuser",
        }),
        expect.anything()
      );
    });

    it("should display validation errors", async () => {
      const user = userEvent.setup();
      render(<TestForm />);
      
      // Submit without filling required field
      await user.click(screen.getByRole("button", { name: /submit/i }));
      
      // Wait for validation error to appear
      const errorMessage = await screen.findByText("Email is required");
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe("FormItem", () => {
    it("should render form item wrapper", () => {
      render(<TestForm />);
      
      // FormItem wraps each field, verify structure exists
      expect(screen.getByLabelText("Email").closest("div")).toBeInTheDocument();
    });
  });

  describe("FormControl", () => {
    it("should render form control correctly", () => {
      render(<TestForm />);
      
      const emailInput = screen.getByLabelText("Email");
      expect(emailInput).toBeInTheDocument();
      expect(emailInput).toHaveAttribute("type", "text");
    });
  });

  describe("Integration with react-hook-form", () => {
    it("should integrate with react-hook-form validation", async () => {
      const user = userEvent.setup();
      
      function ValidationForm() {
        const form = useForm({
          defaultValues: { email: "" },
        });

        return (
          <Form {...form}>
            <form onSubmit={form.handleSubmit(() => {})}>
              <FormField
                control={form.control}
                name="email"
                rules={{
                  required: "Required",
                  pattern: {
                    value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                    message: "Invalid email",
                  },
                }}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">Submit</Button>
            </form>
          </Form>
        );
      }

      render(<ValidationForm />);
      
      const emailInput = screen.getByLabelText("Email");
      await user.type(emailInput, "invalid-email");
      await user.click(screen.getByRole("button", { name: /submit/i }));
      
      const errorMessage = await screen.findByText("Invalid email");
      expect(errorMessage).toBeInTheDocument();
    });
  });
});
