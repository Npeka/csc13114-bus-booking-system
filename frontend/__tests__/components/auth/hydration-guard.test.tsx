import { render, screen } from "@/lib/test-utils";
import { HydrationGuard } from "@/components/auth/hydration-guard";
import { useAuthStore } from "@/lib/stores/auth-store";

// Mock persist functions
const mockHasHydrated = jest.fn(() => true);
const mockOnFinishHydration = jest.fn(() => jest.fn());

jest.mock("@/lib/stores/auth-store", () => {
  const actual = jest.requireActual("@/lib/stores/auth-store");
  return {
    ...actual,
    useAuthStore: Object.assign(actual.useAuthStore, {
      persist: {
        hasHydrated: mockHasHydrated,
        onFinishHydration: mockOnFinishHydration,
      },
    }),
  };
});

describe("HydrationGuard component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockHasHydrated.mockReturnValue(true);
  });

  it("should render children after hydration", () => {
    render(
      <HydrationGuard>
        <div>Hydrated Content</div>
      </HydrationGuard>
    );

    expect(screen.getByText("Hydrated Content")).toBeInTheDocument();
  });

  it("should show loading splash when not hydrated", () => {
    mockHasHydrated.mockReturnValue(false);

    render(
      <HydrationGuard>
        <div>Content</div>
      </HydrationGuard>
    );

    expect(screen.queryByText("Content")).not.toBeInTheDocument();
  });

  it("should handle multiple children when hydrated", () => {
    render(
      <HydrationGuard>
        <div>First Child</div>
        <div>Second Child</div>
        <div>Third Child</div>
      </HydrationGuard>
    );

    expect(screen.getByText("First Child")).toBeInTheDocument();
    expect(screen.getByText("Second Child")).toBeInTheDocument();
    expect(screen.getByText("Third Child")).toBeInTheDocument();
  });
});
