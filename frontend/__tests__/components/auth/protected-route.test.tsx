import { render, screen } from "@/lib/test-utils";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { useAuthStore } from "@/lib/stores/auth-store";
import { Role } from "@/lib/auth/roles";

// Mock the auth store
jest.mock("@/lib/stores/auth-store");
jest.mock("@/lib/api/auth-service", () => ({
  hasValidSession: jest.fn(),
}));

describe("ProtectedRoute component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("should show loading state initially", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
      user: null,
    });

    render(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>,
    );

    expect(screen.getByText(/Đang kiểm tra/)).toBeInTheDocument();
  });

  it("should render children when authenticated", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: {
        id: "123",
        role: Role.PASSENGER,
        full_name: "Test User",
      },
    });

    render(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>,
    );

    expect(screen.getByText("Protected Content")).toBeInTheDocument();
  });

  it("should not render children when not authenticated", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      user: null,
    });

    render(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>,
    );

    expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
  });

  it("should check required roles", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: {
        id: "123",
        role: Role.ADMIN,
        full_name: "Admin User",
      },
    });

    render(
      <ProtectedRoute requiredRoles={[Role.ADMIN]}>
        <div>Admin Content</div>
      </ProtectedRoute>,
    );

    expect(screen.getByText("Admin Content")).toBeInTheDocument();
  });

  it("should not render children when user lacks required role", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: {
        id: "123",
        role: Role.PASSENGER,
        full_name: "Test User",
      },
    });

    render(
      <ProtectedRoute requiredRoles={[Role.ADMIN]}>
        <div>Admin Content</div>
      </ProtectedRoute>,
    );

    expect(screen.queryByText("Admin Content")).not.toBeInTheDocument();
  });

  it("should allow access if user has one of multiple required roles", () => {
    (useAuthStore as jest.Mock).mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      user: {
        id: "123",
        role: Role.OPERATOR,
        full_name: "Operator User",
      },
    });

    render(
      <ProtectedRoute requiredRoles={[Role.ADMIN, Role.OPERATOR]}>
        <div>Privileged Content</div>
      </ProtectedRoute>,
    );

    expect(screen.getByText("Privileged Content")).toBeInTheDocument();
  });
});
