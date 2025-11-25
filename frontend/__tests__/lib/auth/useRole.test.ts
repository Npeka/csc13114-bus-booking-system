import { renderHook } from "@testing-library/react";
import { useRole } from "@/lib/auth/useRole";
import { useAuthStore } from "@/lib/stores/auth-store";
import { Role } from "@/lib/auth/roles";
import { createMockUser, createMockAdmin, createMockOperator } from "@/lib/test-utils";
import { act } from "react";

describe("useRole hook", () => {
  beforeEach(() => {
    // Reset auth store before each test
    act(() => {
      useAuthStore.getState().logout();
    });
  });

  it("should return default values when no user is authenticated", () => {
    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(0);
    expect(result.current.isAdmin).toBe(false);
    expect(result.current.isOperator).toBe(false);
    expect(result.current.isPassenger).toBe(false);
    expect(result.current.isSupport).toBe(false);
  });

  it("should return correct role for passenger user", () => {
    const mockUser = createMockUser({ role: Role.PASSENGER });
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(Role.PASSENGER);
    expect(result.current.isPassenger).toBe(true);
    expect(result.current.isAdmin).toBe(false);
    expect(result.current.isOperator).toBe(false);
    expect(result.current.isSupport).toBe(false);
  });

  it("should return correct role for admin user", () => {
    const mockUser = createMockAdmin();
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(Role.ADMIN);
    expect(result.current.isAdmin).toBe(true);
    expect(result.current.isPassenger).toBe(false);
    expect(result.current.isOperator).toBe(false);
    expect(result.current.isSupport).toBe(false);
  });

  it("should return correct role for operator user", () => {
    const mockUser = createMockOperator();
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(Role.OPERATOR);
    expect(result.current.isOperator).toBe(true);
    expect(result.current.isAdmin).toBe(false);
    expect(result.current.isPassenger).toBe(false);
    expect(result.current.isSupport).toBe(false);
  });

  it("should return correct role for support user", () => {
    const mockUser = createMockUser({ role: Role.SUPPORT });
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(Role.SUPPORT);
    expect(result.current.isSupport).toBe(true);
    expect(result.current.isAdmin).toBe(false);
    expect(result.current.isPassenger).toBe(false);
    expect(result.current.isOperator).toBe(false);
  });

  it("should handle multi-role users (bitwise combination)", () => {
    // User with both Admin and Passenger roles
    const multiRoleValue = Role.ADMIN | Role.PASSENGER; // 3
    const mockUser = createMockUser({ role: multiRoleValue });
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.role).toBe(multiRoleValue);
    expect(result.current.isAdmin).toBe(true);
    expect(result.current.isPassenger).toBe(true);
    expect(result.current.isOperator).toBe(false);
    expect(result.current.isSupport).toBe(false);
  });

  it("should update when user changes", () => {
    const passengerUser = createMockUser({ role: Role.PASSENGER });
    
    act(() => {
      useAuthStore.getState().setUser(passengerUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.isPassenger).toBe(true);
    expect(result.current.isAdmin).toBe(false);

    // Simulate user change to admin
    const adminUser = createMockAdmin();
    
    act(() => {
      useAuthStore.getState().setUser(adminUser);
    });

    expect(result.current.isAdmin).toBe(true);
    expect(result.current.isPassenger).toBe(false);
  });

  it("should handle user logout (null user)", () => {
    const mockUser = createMockAdmin();
    
    act(() => {
      useAuthStore.getState().setUser(mockUser);
    });

    const { result } = renderHook(() => useRole());

    expect(result.current.isAdmin).toBe(true);

    // Simulate logout
    act(() => {
      useAuthStore.getState().logout();
    });

    expect(result.current.role).toBe(0);
    expect(result.current.isAdmin).toBe(false);
  });
});
