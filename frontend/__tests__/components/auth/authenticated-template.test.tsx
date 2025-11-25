import { render, screen, waitFor } from "@/lib/test-utils";
import { AuthenticatedTemplate } from "@/components/auth/authenticated-template";
import { useAuthStore } from "@/lib/stores/auth-store";
import { hasValidSession } from "@/lib/api/auth-service";
import { createMockUser } from "@/lib/test-utils";
import { act } from "react";

// Mock the auth service
jest.mock("@/lib/api/auth-service", () => ({
  hasValidSession: jest.fn(),
}));

describe("AuthenticatedTemplate component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (hasValidSession as jest.Mock).mockResolvedValue(true);
    act(() => {
      useAuthStore.getState().logout();
    });
  });

  it("should render children when authenticated", () => {
    const mockUser = createMockUser();

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    render(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.getByText("Authenticated Content")).toBeInTheDocument();
  });

  it("should not render children when not authenticated", () => {
    render(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.queryByText("Authenticated Content")).not.toBeInTheDocument();
  });

  it("should update when authentication state changes", () => {
    const { rerender } = render(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.queryByText("Authenticated Content")).not.toBeInTheDocument();

    const mockUser = createMockUser();
    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    rerender(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.getByText("Authenticated Content")).toBeInTheDocument();
  });

  it("should hide content when user logs out", () => {
    const mockUser = createMockUser();

    act(() => {
      useAuthStore.getState().login(mockUser, "token");
    });

    const { rerender } = render(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.getByText("Authenticated Content")).toBeInTheDocument();

    act(() => {
      useAuthStore.getState().logout();
    });

    rerender(
      <AuthenticatedTemplate>
        <div>Authenticated Content</div>
      </AuthenticatedTemplate>
    );

    expect(screen.queryByText("Authenticated Content")).not.toBeInTheDocument();
  });

  it("should render nothing (null) when not authenticated", () => {
    const { container } = render(
      <AuthenticatedTemplate>
        <div>Content</div>
      </AuthenticatedTemplate>
    );

    expect(container.firstChild).toBeNull();
  });
});
