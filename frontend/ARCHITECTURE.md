# Architecture & System Design

## Overview

BusTicket.vn is a Next.js 16 SPA featuring secure authentication, role-based authorization, and a responsive design system. This document outlines the system architecture, data flow, and design decisions.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Client (Browser)                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           React Components (TSX)                     │  │
│  │  - Pages (app/)                                     │  │
│  │  - Layout Components (header, footer)              │  │
│  │  - UI Components (button, card, badge, etc.)       │  │
│  │  - Auth Components (protected-route, role-badge)   │  │
│  └──────────────────────────────────────────────────────┘  │
│                           ↕                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │        State Management & Hooks                      │  │
│  │  - Zustand Auth Store (tokens, user data)          │  │
│  │  - React Hooks (useRole, useAuth)                  │  │
│  │  - TanStack Query (API caching & sync)            │  │
│  └──────────────────────────────────────────────────────┘  │
│                           ↕                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │      API Client (lib/api/)                           │  │
│  │  - Axios HTTP client                               │  │
│  │  - Auth interceptor (inject access token)          │  │
│  │  - 401 handler (trigger token refresh)             │  │
│  │  - Error handling & retry logic                    │  │
│  └──────────────────────────────────────────────────────┘  │
│                           ↕ HTTP                            │
└─────────────────────────────────────────────────────────────┘
                            │
                    ┌───────▼────────┐
                    │   Backend API  │
                    │  (Go Services) │
                    │  - User Service│
                    │  - Booking Svc │
                    │  - Trip Service│
                    │  - Payment Svc │
                    └────────────────┘
```

## Folder Structure

```
frontend/
├── app/                          # Next.js App Router pages
│   ├── layout.tsx               # Root layout (providers)
│   ├── page.tsx                 # Home page
│   ├── admin/
│   │   └── dashboard/
│   │       └── page.tsx         # Admin dashboard (role-protected)
│   ├── operator/
│   │   └── dashboard/
│   │       └── page.tsx         # Operator panel (role-protected)
│   ├── trips/
│   │   └── page.tsx             # Trip search & browsing
│   ├── my-bookings/
│   │   └── page.tsx             # User bookings (auth-protected)
│   ├── checkout/
│   │   └── page.tsx             # Checkout flow (auth-protected)
│   └── booking-confirmation/
│       └── page.tsx             # Confirmation page
│
├── components/                   # Reusable React components
│   ├── layout/
│   │   ├── header.tsx           # Top navigation (role-aware menu)
│   │   ├── footer.tsx           # Footer
│   │   └── chatbot.tsx          # AI chatbot widget
│   ├── auth/
│   │   ├── protected-route.tsx  # Role-based route guard
│   │   ├── auth-provider.tsx    # Session restoration
│   │   ├── role-badge.tsx       # Display user role
│   │   └── hydration-guard.tsx  # SSR/CSR sync
│   ├── ui/                       # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── badge.tsx
│   │   ├── tabs.tsx
│   │   ├── dialog.tsx
│   │   ├── form.tsx
│   │   └── ... (17+ components)
│   ├── search/
│   │   └── trip-search-form/
│   ├── trips/
│   │   ├── trip-card.tsx
│   │   ├── seat-map.tsx
│   │   └── trip-filters.tsx
│   └── theme/
│       └── mode-toggle.tsx      # Dark/light theme switcher
│
├── lib/                          # Utilities & libraries
│   ├── auth/
│   │   ├── roles.ts             # Role utilities (isAdmin, isOperator, etc.)
│   │   ├── useRole.ts           # Hook to read current user role
│   │   └── ... (auth-service, token-manager)
│   ├── api/
│   │   ├── auth-service.ts      # Auth API functions
│   │   ├── client.ts            # Axios client with interceptors
│   │   └── hooks/               # API hooks (useBookings, useTrips, etc.)
│   ├── stores/
│   │   └── auth-store.ts        # Zustand auth state (tokens, user)
│   ├── firebase.ts              # Firebase config & init
│   ├── utils.ts                 # General utilities
│   └── test-utils.tsx           # Jest testing helpers
│
├── __tests__/                    # Jest test files
│   ├── lib/auth/
│   │   └── roles.test.ts        # Role utility tests
│   ├── components/auth/
│   │   ├── role-badge.test.tsx
│   │   └── protected-route.test.tsx
│   └── components/ui/
│       ├── button.test.tsx
│       └── badge.test.tsx
│
├── design/                       # Design documentation
│   ├── color-tokens.md          # Color system (oklch, semantic colors)
│   ├── typography-scale.md      # Font sizes, weights, line heights
│   ├── component-specs.md       # UI component specifications
│   └── screenshots/             # Page screenshots by role
│
├── docs/                         # Additional documentation
│   └── DEPLOYMENT.md            # Deployment & environment setup
│
├── public/                       # Static assets
├── jest.config.ts               # Jest configuration
├── jest.setup.ts                # Jest setup (mocks, fixtures)
├── .lintstagedrc.json           # Pre-commit linting config
├── .eslintrc.mjs                # ESLint configuration
├── .prettierrc                  # Code formatter config
├── tsconfig.json                # TypeScript config
├── tailwind.config.ts           # Tailwind CSS config (design tokens)
├── next.config.ts               # Next.js config
└── package.json                 # Dependencies & scripts
```

## Data Flow

### User Authentication Flow

```
1. User clicks "Login"
   ↓
2. Opens auth dialog (phone OTP or Google OAuth)
   ↓
3. Sends credentials to Backend:
   - Phone + OTP → /api/v1/auth/verify-phone
   - Google Token → /api/v1/auth/firebase/auth
   ↓
4. Backend validates & returns:
   - accessToken (short-lived, JWT)
   - refreshToken (long-lived, httpOnly cookie)
   - User data (id, name, role, email)
   ↓
5. Frontend stores:
   - accessToken → Zustand auth store (in-memory)
   - refreshToken → httpOnly cookie (secure)
   - User → Zustand auth store
   ↓
6. Session persisted & restored on page reload
   - Zustand loads from localStorage
   - Checks refreshToken validity
   - Auto-refreshes if expired
```

### Authorization Flow

```
1. User tries to access protected route (e.g., /admin)
   ↓
2. ProtectedRoute component checks:
   - isAuthenticated? (access token exists)
   - hasRequiredRole? (user.role & ADMIN_ROLE === ADMIN_ROLE)
   ↓
3a. If authorized → Render protected content
   ↓
3b. If not authorized → Redirect to home
   ↓
4. Role displayed in header dropdown (RoleBadge)
   - Navigation shows role-specific links
   - Admin users see "Quản trị" link
   - Operators see "Điều hành" link
```

### API Request Flow

```
1. Component calls API via hook (e.g., useBookings())
   ↓
2. Axios client adds Authorization header:
   Authorization: Bearer <accessToken>
   ↓
3. Request sent to backend
   ↓
4a. Success (200) → Return data, cache in TanStack Query
   ↓
4b. Unauthorized (401) → Token expired
      - Trigger refresh token flow
      - Retry original request with new token
   ↓
4c. Error (other) → Show error toast, log to console
```

## State Management

### Authentication State (Zustand)

```typescript
interface AuthStore {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setUser(user: User): void;
  setAccessToken(token: string): void;
  logout(): void;
  restoreSession(): void;
}
```

**Persistence**: localStorage (auto-loaded on app start)

### Role State

```typescript
// Derived from auth.user.role (bitflag: 1=Passenger, 2=Admin, 4=Operator, 8=Support)
interface RoleState {
  role: number;
  isAdmin: boolean;
  isOperator: boolean;
  isPassenger: boolean;
  isSupport: boolean;
}
```

**Source**: Computed via `useRole()` hook from auth store

### API Data State (TanStack Query)

```typescript
// Bookings
const {
  data: bookings,
  isLoading,
  error,
} = useQuery({
  queryKey: ["bookings"],
  queryFn: () => api.getBookings(),
});

// Trips
const { data: trips } = useQuery({
  queryKey: ["trips", filters],
  queryFn: () => api.searchTrips(filters),
});
```

## Key Design Decisions

### 1. Token Storage Strategy

**Decision**: Access token in-memory, refresh token in httpOnly cookie

**Rationale**:

- **httpOnly Cookie**: Secure against XSS attacks (can't be stolen via JS)
- **In-Memory Access Token**: Fast access, no cookie lookup overhead
- **Trade-off**: Page reload loses access token, but refresh token recovers it

**Alternatives Considered**:

- All in localStorage: Simple but XSS-vulnerable
- All in sessionStorage: Lost on browser close, inconvenient
- All in httpOnly cookie: Secure but slower (cookie sent with every request)

### 2. Role-Based Access Control (Bit-Flag)

**Decision**: Align frontend role model with backend bit-flag encoding

**Rationale**:

- **Backend uses**: Passenger=1, Admin=2, Operator=4, Support=8 (powers of 2)
- **Frontend mirrors**: Same model for consistency, easy bitwise checks
- **Multiple roles**: User can have role=3 (Passenger + Admin)

**Implementation**:

```typescript
// Check if user is admin
const isAdmin = (userRole & Role.ADMIN) === Role.ADMIN;

// Check if user has any of multiple roles
const hasAccess = requiredRoles.some((role) => (userRole & role) === role);
```

### 3. Protected Routes with Role Checking

**Decision**: Extend ProtectedRoute component to accept optional `requiredRoles` prop

**Rationale**:

- Single component for both auth + authz checks
- Easy to protect new pages: `<ProtectedRoute requiredRoles={[Role.ADMIN]}>`
- Server-side protection via middleware (optional, can add in Phase 2)

**Usage**:

```tsx
export default function AdminDashboard() {
  return (
    <ProtectedRoute requiredRoles={[Role.ADMIN]}>
      <AdminContent />
    </ProtectedRoute>
  );
}
```

### 4. Dark Mode Support

**Decision**: Use next-themes for seamless theme switching

**Rationale**:

- **System preference**: Respects OS dark mode setting
- **User override**: Saves theme choice in localStorage
- **No flash**: Hydration guard prevents theme flicker on load

**Implementation**:

```tsx
<ThemeProvider attribute="class" defaultTheme="system" enableSystem>
  {/* App */}
</ThemeProvider>
```

### 5. Design Tokens via Tailwind

**Decision**: Define all colors, spacing, typography in tailwind.config.ts

**Rationale**:

- **Centralized**: Single source of truth for design system
- **Type-safe**: Intellisense in components
- **Scalable**: Easy to adjust brand colors globally
- **Dark mode**: Automatic CSS variable swapping

## Error Handling

### HTTP Errors

| Status | Handling                                 |
| ------ | ---------------------------------------- |
| 401    | Token refresh flow, retry request        |
| 403    | Show error toast, redirect to home       |
| 404    | Show 404 page                            |
| 500    | Show error toast, log to Sentry (future) |

### Form Validation

```typescript
// React Hook Form + Zod
const schema = z.object({
  email: z.string().email(),
  phone: z.string().min(9),
});

const {
  register,
  formState: { errors },
} = useForm({
  resolver: zodResolver(schema),
});
```

## Performance Optimizations

1. **Code Splitting**: Next.js automatic page splitting
2. **Image Optimization**: next/image component
3. **Bundle Size**: Tree-shaking, minification
4. **API Caching**: TanStack Query with stale-while-revalidate
5. **Dark Mode**: CSS-in-JS prevents flashing

## Security Measures

1. **HTTPS Only**: Refresh tokens in httpOnly, secure, sameSite cookies
2. **Token Expiry**: Access token expires after 15 mins (backend decision)
3. **CSRF Protection**: (depends on backend cookie-based CSRF token)
4. **XSS Prevention**: Sanitized HTML, Content Security Policy headers
5. **CORS**: Backend enforces allowed origins

## Testing Strategy

### Unit Tests (Jest + React Testing Library)

- **Utility functions**: Role checking, token validation
- **Components**: RoleBadge, ProtectedRoute, Button, Badge
- **Hooks**: useRole, useAuth

### Integration Tests (Playwright, future)

- Authentication flows (login, logout, token refresh)
- Role-based access (admin page redirects non-admins)
- Booking workflow (search → seat → checkout → confirm)

### E2E Tests (Cypress, future)

- Full user journeys in real browser
- Mobile responsiveness
- Dark mode switching

## Deployment & Environments

### Development

- Local dev server: `npm run dev`
- API: localhost backend or staging API
- Firebase: Dev project

### Production

- Vercel deployment
- API: Production backend
- Firebase: Production project
- Environment variables: Set in Vercel dashboard

## Future Enhancements

1. **Phase 2**: Real payment integration (MoMo, ZaloPay webhooks)
2. **Phase 3**: Advanced analytics, refund workflow
3. **Phase 4**: Mobile app port, multi-language i18n
4. **Phase 5**: Performance: Incremental Static Regeneration (ISR), edge caching
