import { render, screen, waitFor } from "@/lib/test-utils";
import { AuthGuard } from "@/components/auth/auth-guard";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/auth-service";
import { Role } from "@/lib/auth/roles";
import { createMockUser, createMockAdmin } from "@/lib/test-utils";
import { routerMocks } from "@/jest.setup";
import { act } from "react";

// Mock the auth service
jest.mock("@/lib/api/auth-service");

describe("AuthGuard component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    act(() => {
      useAuthStore.getState().logout();
    });
    (hasValidSession as jest.Mock).mockResolvedValue(false);
  });

  it("should show loading state initially", () => {
    act(() => {
      useAuthStore.getState().setLoading(true);
    });

    render(
      <AuthGuard>
        <div>Protected Content</div>
      </AuthGuard>
    );

    expect(screen.getByText(/Đang kiểm tra quyền truy cập/)).toBeInTheDocument();
  });

  it("should render children when authenticated", async () => {
    const mockUser = createMockUser();
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    render(
      <AuthGuard>
        <div>Protected Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(screen.getByText("Protected Content")).toBeInTheDocument();
    });
  });

  it("should redirect to home when not authenticated", async () => {
    (hasValidSession as jest.Mock).mockResolvedValue(false);

    render(
      <AuthGuard>
        <div>Protected Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(routerMocks.mockPush).toHaveBeenCalledWith("/");
    });

    expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
  });

  it("should render children when user has required role", async () => {
    const mockAdmin = createMockAdmin();
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().login(mockAdmin, "token");
    });

    render(
      <AuthGuard requiredRole={Role.ADMIN}>
        <div>Admin Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(screen.getByText("Admin Content")).toBeInTheDocument();
    });
  });

  it("should redirect when user lacks required role", async () => {
    const mockUser = createMockUser({ role: Role.PASSENGER });
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    render(
      <AuthGuard requiredRole={Role.ADMIN}>
        <div>Admin Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(routerMocks.mockPush).toHaveBeenCalledWith("/");
    });

    expect(screen.queryByText("Admin Content")).not.toBeInTheDocument();
  });

  it("should not render when authenticated but lacks required role", async () => {
    const mockUser = createMockUser({ role: Role.PASSENGER });
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    render(
      <AuthGuard requiredRole={Role.OPERATOR}>
        <div>Operator Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(screen.queryByText("Operator Content")).not.toBeInTheDocument();
    });
  });

  it("should handle session validation async correctly", async () => {
    const mockUser = createMockUser();
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    render(
      <AuthGuard>
        <div>Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(hasValidSession).toHaveBeenCalled();
    });

    expect(screen.getByText("Content")).toBeInTheDocument();
  });

  it("should not check auth while loading", () => {
    act(() => {
      useAuthStore.getState().setLoading(true);
    });

    render(
      <AuthGuard>
        <div>Content</div>
      </AuthGuard>
    );

    expect(hasValidSession).not.toHaveBeenCalled();
  });

  it("should allow multi-role users with required role", async () => {
    const multiRole = Role.ADMIN | Role.PASSENGER;
    const mockUser = createMockUser({ role: multiRole });
    (hasValidSession as jest.Mock).mockResolvedValue(true);

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    render(
      <AuthGuard requiredRole={Role.ADMIN}>
        <div>Admin Content</div>
      </AuthGuard>
    );

    await waitFor(() => {
      expect(screen.getByText("Admin Content")).toBeInTheDocument();
    });
  });
});
