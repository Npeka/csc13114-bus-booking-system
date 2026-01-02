import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Protected routes that require authentication
const protectedRoutes = [
  "/my-bookings",
  "/checkout",
  "/booking-confirmation",
  "/profile",
];

// Public guest routes that should NOT be protected
// These are explicitly allowed for unauthenticated users
const publicGuestRoutes = [
  "/checkout-guest",
  "/booking-confirmation-guest",
  "/booking-lookup",
  "/booking-details",
];

/**
 * Check if a pathname matches a protected route.
 * Uses exact prefix matching with path boundary awareness to avoid
 * matching /checkout-guest when checking for /checkout
 */
function isProtectedPath(pathname: string): boolean {
  return protectedRoutes.some((route) => {
    // Exact match
    if (pathname === route) return true;
    // Path prefix match (ensure we're matching a path segment, not a partial string)
    // e.g., /checkout/step1 should match, but /checkout-guest should not
    if (pathname.startsWith(route + "/")) return true;
    return false;
  });
}

/**
 * Check if a pathname is a public guest route
 */
function isPublicGuestPath(pathname: string): boolean {
  return publicGuestRoutes.some(
    (route) => pathname === route || pathname.startsWith(route + "/"),
  );
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Always allow public guest routes
  if (isPublicGuestPath(pathname)) {
    return NextResponse.next();
  }

  // Check if the route is protected
  const isProtected = isProtectedPath(pathname);

  // Check if user has refresh token cookie
  const refreshToken = request.cookies.get("refresh_token");
  const hasAuth = !!refreshToken;

  // If trying to access protected route without auth, redirect to home
  if (isProtected && !hasAuth) {
    const url = new URL("/", request.url);
    url.searchParams.set("redirect", pathname);
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

// Configure which routes the proxy should run on
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
};
