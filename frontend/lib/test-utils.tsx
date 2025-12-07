import React, { ReactElement } from "react";
import { render, RenderOptions } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/components/theme-provider";
import { User, UserStatus } from "@/lib/stores/auth-store";
import { Role } from "@/lib/auth/roles";

/**
 * Create a new QueryClient for testing with default options
 */
export function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false, // Disable retries in tests
        gcTime: 0, // Disable caching
      },
      mutations: {
        retry: false,
      },
    },
  });
}

/**
 * Wrapper with all providers (Theme + QueryClient)
 */
interface AllProvidersProps {
  children: React.ReactNode;
  queryClient?: QueryClient;
}

const AllTheProviders = ({
  children,
  queryClient = createTestQueryClient(),
}: AllProvidersProps) => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
        {children}
      </ThemeProvider>
    </QueryClientProvider>
  );
};

/**
 * Custom render function with all providers
 */
interface CustomRenderOptions extends Omit<RenderOptions, "wrapper"> {
  queryClient?: QueryClient;
}

export const customRender = (
  ui: ReactElement,
  options?: CustomRenderOptions,
) => {
  const { queryClient, ...renderOptions } = options || {};

  return render(ui, {
    wrapper: ({ children }) => (
      <AllTheProviders queryClient={queryClient}>{children}</AllTheProviders>
    ),
    ...renderOptions,
  });
};

// Re-export everything from React Testing Library
export * from "@testing-library/react";
export { default as userEvent } from "@testing-library/user-event";

// Export custom render as default render
export { customRender as render };

// ============================================================================
// Test Data Factories
// ============================================================================

/**
 * Create a mock user for testing
 */
export function createMockUser(overrides?: Partial<User>): User {
  return {
    id: "test-user-1",
    email: "test@example.com",
    phone: "+84123456789",
    full_name: "Test User",
    avatar: "https://example.com/avatar.jpg",
    role: Role.PASSENGER,
    status: UserStatus.Active,
    email_verified: true,
    phone_verified: true,
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    ...overrides,
  };
}

/**
 * Create a mock admin user
 */
export function createMockAdmin(overrides?: Partial<User>): User {
  return createMockUser({
    id: "admin-1",
    email: "admin@example.com",
    full_name: "Admin User",
    role: Role.ADMIN,
    ...overrides,
  });
}

/**
 * Create a mock access token
 */
export function createMockAccessToken(): string {
  return "mock-access-token-" + Math.random().toString(36).substring(7);
}

// ============================================================================
// Auth Store Mock Helpers
// ============================================================================

/**
 * Mock authenticated state for auth store
 */
export function mockAuthenticatedState(user?: Partial<User>) {
  return {
    user: createMockUser(user),
    accessToken: createMockAccessToken(),
    isAuthenticated: true,
    isLoading: false,
    error: null,
  };
}

/**
 * Mock unauthenticated state for auth store
 */
export function mockUnauthenticatedState() {
  return {
    user: null,
    accessToken: null,
    isAuthenticated: false,
    isLoading: false,
    error: null,
  };
}

/**
 * Mock loading state for auth store
 */
export function mockLoadingState() {
  return {
    user: null,
    accessToken: null,
    isAuthenticated: false,
    isLoading: true,
    error: null,
  };
}

/**
 * Mock error state for auth store
 */
export function mockErrorState(errorMessage: string = "Authentication error") {
  return {
    user: null,
    accessToken: null,
    isAuthenticated: false,
    isLoading: false,
    error: errorMessage,
  };
}

// ============================================================================
// Wait Helpers
// ============================================================================

/**
 * Wait for async operations with a timeout
 */
export const waitFor = async (
  callback: () => void | Promise<void>,
  timeout: number = 1000,
) => {
  const startTime = Date.now();

  while (Date.now() - startTime < timeout) {
    try {
      await callback();
      return;
    } catch {
      await new Promise((resolve) => setTimeout(resolve, 50));
    }
  }

  throw new Error(`waitFor timed out after ${timeout}ms`);
};
