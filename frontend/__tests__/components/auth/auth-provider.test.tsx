import { render } from "@/lib/test-utils";
import { AuthProvider } from "@/components/auth/auth-provider";
import { restoreSession } from "@/lib/auth/session";

// Mock the session module
jest.mock("@/lib/auth/session");

describe("AuthProvider component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("should render children", () => {
    const { container } = render(
      <AuthProvider>
        <div data-testid="child">Child Content</div>
      </AuthProvider>
    );

    expect(container.querySelector('[data-testid="child"]')).toBeInTheDocument();
  });

  it("should call restoreSession on mount", () => {
    render(
      <AuthProvider>
        <div>Content</div>
      </AuthProvider>
    );

    expect(restoreSession).toHaveBeenCalledTimes(1);
  });

  it("should only call restoreSession once", () => {
    const { rerender } = render(
      <AuthProvider>
        <div>Content</div>
      </AuthProvider>
    );

    expect(restoreSession).toHaveBeenCalledTimes(1);

    // Rerender should not call restoreSession again
    rerender(
      <AuthProvider>
        <div>Updated Content</div>
      </AuthProvider>
    );

    expect(restoreSession).toHaveBeenCalledTimes(1);
  });

  it("should render children even if restoreSession throws", () => {
    (restoreSession as jest.Mock).mockImplementation(() => {
      throw new Error("Session restore failed");
    });

    const { container } = render(
      <AuthProvider>
        <div data-testid="child">Child Content</div>
      </AuthProvider>
    );

    // Children should still render
    expect(container.querySelector('[data-testid="child"]')).toBeInTheDocument();
  });
});
