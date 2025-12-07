# BusTicket.vn Frontend

A production-ready Next.js 16 SPA for bus ticket booking with secure authentication, role-based authorization, and a comprehensive design system.

**Live Demo**: [https://csc13114-bus-booking-system.vercel.app](https://csc13114-bus-booking-system.vercel.app)

## Table of Contents

1. [Getting Started](#getting-started)
2. [Local Setup](#local-setup)
3. [Authentication & Authorization](#authentication--authorization)
4. [Architecture](#architecture)
5. [Development](#development)
6. [Testing](#testing)
7. [Design System](#design-system)
8. [Decisions & Tradeoffs](#decisions--tradeoffs)
9. [Deployment](#deployment)
10. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Prerequisites

- Node.js 16+ (18+ recommended)
- npm or yarn
- Firebase account with Google OAuth enabled
- Backend API running (see [API Setup](#api-setup))

### Quick Start

```bash
# Clone repo
git clone <repo-url>
cd frontend

# Install dependencies
npm install

# Setup environment variables
cp .env.example .env.local
# Edit .env.local with your Firebase credentials and API URL

# Run dev server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

---

## Local Setup

### Environment Variables

Create `.env.local` with:

```env
# Firebase
NEXT_PUBLIC_FIREBASE_API_KEY=your_firebase_api_key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_firebase_auth_domain
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_firebase_project_id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your_firebase_storage_bucket
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your_firebase_messaging_sender_id
NEXT_PUBLIC_FIREBASE_APP_ID=your_firebase_app_id

# Backend API
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Environment
NEXT_PUBLIC_APP_ENV=development
```

### API Setup

Backend services must be running:

```bash
# In /backend directory
docker-compose up

# Or manually start services:
cd user-service && make run &
cd booking-service && make run &
cd trip-service && make run &
cd payment-service && make run &
cd gateway-service && make run &
```

### Database Seeding

Backend handles database setup. Verify connectivity:

```bash
curl http://localhost:8080/api/v1/health
# Expected: {"status":"ok"}
```

---

## Authentication & Authorization

### Authentication Flow

We implement a secure, industry-standard token model:

#### 1. **Email + Password** (Firebase)

```
User → Signup/Login Form
  ↓
Firebase Authentication
  ↓
Backend validates & issues tokens
  ↓
Frontend stores tokens securely
```

#### 2. **Social Login** (Google OAuth)

```
User → "Login with Google"
  ↓
Google OAuth Dialog
  ↓
Firebase handles OAuth
  ↓
Backend validates ID token & issues session tokens
  ↓
Frontend stores tokens
```

### Token Management

**Security Model**:

- **Access Token** (short-lived, ~15 min):
  - Stored: In-memory via Zustand (fast, XSS-safe with httpOnly refresh)
  - Sent: In `Authorization: Bearer <token>` header
  - Expires: After 15 minutes
- **Refresh Token** (long-lived, ~30 days):
  - Stored: httpOnly, secure, sameSite cookies (XSS-proof)
  - Used: Automatically to refresh access token when expired
  - Expires: After 30 days

**Why This Model**:

- **httpOnly Cookies**: Immune to XSS attacks (JavaScript can't access)
- **In-Memory Access**: No cookie overhead per request, faster
- **Auto-Refresh**: Seamless experience without user re-login
- **Server Validation**: Refresh tokens never sent to client, only via HTTP

**Automatic Refresh Flow**:

```
1. API returns 401 Unauthorized
  ↓
2. Token Manager intercepts & detects refresh token exists
  ↓
3. Auto-refresh: POST /auth/refresh-token with refresh token
  ↓
4. Backend validates refresh token & issues new access token
  ↓
5. Store new access token in-memory
  ↓
6. Retry original request with new token
```

### Authorization (Role-Based Access Control)

**Role Model** (bit-flag encoding):

| Role      | Value | Permissions                                      |
| --------- | ----- | ------------------------------------------------ |
| Passenger | 1     | Book tickets, view own bookings                  |
| Admin     | 2     | User management, analytics, system config        |
| Operator  | 4     | Manage trips, seat availability, passenger lists |
| Support   | 8     | Customer support, refund processing              |

**Note**: Users can have multiple roles (bit-flag: e.g., role=3 means Admin + Passenger)

**Route Protection**:

```typescript
// Customer routes (any authenticated user)
<ProtectedRoute>
  <MyBookings />
</ProtectedRoute>

// Admin-only routes
<ProtectedRoute requiredRoles={[Role.ADMIN]}>
  <AdminDashboard />
</ProtectedRoute>

// Operator or Admin
<ProtectedRoute requiredRoles={[Role.OPERATOR, Role.ADMIN]}>
  <OperatorPanel />
</ProtectedRoute>
```

**UI Enforcement**:

- Admin users see "Quản trị" link in header
- Operators see "Điều hành" link in header
- Role badge displayed in dropdown menu
- Non-admin users redirected from `/admin/*` routes

**Server-Side Enforcement**:

- Backend validates user role for all API endpoints
- 403 Forbidden if insufficient permissions
- Refresh token invalidated if role changed

---

## Architecture

### High-Level Design

```
Frontend (Next.js)  ←→  Backend (Go Services)  ←→  PostgreSQL
      ↓
Zustand (Auth State) + TanStack Query (API Cache)
```

### Folder Structure

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed folder organization and data flow diagrams.

### Key Files

| Path                                  | Purpose                           |
| ------------------------------------- | --------------------------------- |
| `lib/auth/roles.ts`                   | Role utilities & bit-flag helpers |
| `lib/auth/useRole.ts`                 | Hook to check current user role   |
| `components/auth/protected-route.tsx` | Role-based route guard            |
| `components/auth/role-badge.tsx`      | Display user role UI              |
| `app/admin/page.tsx`                  | Admin dashboard (role-protected)  |
| `app/operator/dashboard/page.tsx`     | Operator panel (role-protected)   |

---

## Development

### Available Scripts

```bash
npm run dev          # Start dev server (http://localhost:3000)
npm run build        # Build for production
npm start            # Start production server
npm run lint         # Run ESLint
npm run format       # Format code with Prettier
npm test             # Run Jest tests
npm run test:watch   # Run tests in watch mode
npm run test:coverage # Generate coverage report
```

### Code Style

We enforce code quality with ESLint and Prettier:

```bash
# Linting (auto-fix where possible)
npm run lint

# Formatting
npm run format

# Pre-commit hook runs both automatically
git add .
git commit -m "my changes"  # Lint & format run automatically
```

### Development Workflow

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes following code style guidelines
3. Test locally: `npm test`
4. Commit (lint-staged runs automatically): `git commit -m "..."`
5. Push: `git push origin feature/my-feature`
6. Create Pull Request

---

## Testing

### Running Tests

```bash
# All tests
npm test

# Watch mode (re-run on file changes)
npm run test:watch

# Coverage report
npm run test:coverage
```

### Test Structure

```
__tests__/
├── lib/auth/
│   └── roles.test.ts           # Role utility tests
├── components/auth/
│   ├── role-badge.test.tsx     # RoleBadge component
│   └── protected-route.test.tsx # ProtectedRoute with role checking
└── components/ui/
    ├── button.test.tsx         # Button component
    └── badge.test.tsx          # Badge component
```

### Test Examples

```typescript
// Testing role utilities
import { isAdmin, hasRole, Role } from '@/lib/auth/roles';

test('should identify admin role', () => {
  expect(isAdmin(Role.ADMIN)).toBe(true);
  expect(isAdmin(Role.PASSENGER)).toBe(false);
});

// Testing role-protected components
test('ProtectedRoute should require admin role', () => {
  const { queryByText } = render(
    <ProtectedRoute requiredRoles={[Role.ADMIN]}>
      <div>Admin Only</div>
    </ProtectedRoute>
  );
  expect(queryByText('Admin Only')).not.toBeInTheDocument();
});
```

---

## Design System

### Colors

All colors use oklch color space for superior accessibility and consistency.

**Semantic Colors**:

- **Primary**: #E63946 (red) - Main actions, links
- **Success**: Green - Confirmed bookings
- **Error**: Red - Cancelled bookings, errors
- **Warning**: Yellow - Pending actions
- **Info**: Blue - Completed bookings

See [design/color-tokens.md](./design/color-tokens.md) for complete palette and usage guidelines.

### Typography

- **Headings**: Inter (sans-serif), weights 600–700
- **Body**: Inter (sans-serif), weight 400
- **Sizes**: h1 (32px) → h6 (14px), body (16px), caption (12px)

See [design/typography-scale.md](./design/typography-scale.md) for specifications.

### Components

17+ reusable UI components:

- Button (6 variants: default, outline, destructive, secondary, ghost, link)
- Card (with Header, Content, Footer, Title, Description)
- Badge (4 variants)
- Dialog, Tabs, Select, Checkbox, Input, Form
- And more...

See [design/component-specs.md](./design/component-specs.md) for usage and props.

---

## Decisions & Tradeoffs

### 1. Token Storage (Access in-memory, Refresh in httpOnly)

**Decision**: Access tokens in Zustand (in-memory), refresh tokens in httpOnly cookies

**Why**:

- Prevents XSS attacks (httpOnly cookies can't be stolen via JS)
- Fast access token reads (no cookie lookup)
- Auto-refresh handles expired tokens transparently

**Tradeoff**:

- Hard refresh (F5) loses access token but refresh token recovers it
- Page reload takes ~500ms to validate refresh token

**Alternatives Considered**:

- All in localStorage: Simpler but vulnerable to XSS
- All in httpOnly: Secure but slow (cookie sent with every request)

### 2. Bit-Flag Roles

**Decision**: Use bit-flag encoding (2^n) for roles, matching backend

**Why**:

- Efficient: Roles fit in single integer
- Flexible: User can have multiple roles (e.g., role=3 → Admin + Passenger)
- Consistent: Aligns with backend authorization model

**Implementation**:

```typescript
const isAdmin = (userRole & Role.ADMIN) === Role.ADMIN;
```

### 3. Protected Routes with Optional Role Check

**Decision**: Extend ProtectedRoute to accept `requiredRoles` prop

**Why**:

- Single component for auth + authz
- Easy to protect pages: `<ProtectedRoute requiredRoles={[Role.ADMIN]}>`
- Future-proof for multiple role requirements

### 4. Mock Data Strategy

**Decision**: Use mock data until backend API is available

**Why**:

- Faster frontend development
- Can test UI independently
- Easy feature flag to switch between mock and real data

**Current Status**: Mock data used for dashboard widgets

**Next Phase**: Wire to real API endpoints (see [/docs/NEXT_STEPS.md](./NEXT_STEPS.md))

---

## Deployment

### Vercel (Recommended)

Vercel is optimized for Next.js and provides seamless deployment:

1. **Connect GitHub**:
   - Go to [vercel.com](https://vercel.com)
   - Click "New Project"
   - Select GitHub repo
   - Vercel auto-detects Next.js

2. **Set Environment Variables**:
   - In Vercel dashboard: Project → Settings → Environment Variables
   - Add all variables from `.env.local`:
     ```
     NEXT_PUBLIC_FIREBASE_API_KEY=...
     NEXT_PUBLIC_API_BASE_URL=https://api.busticket.vn
     NEXT_PUBLIC_APP_ENV=production
     ```

3. **Deploy**:
   - Push to main branch → Automatic deploy
   - Or manual: `vercel --prod`

4. **Verify**:
   ```bash
   curl https://your-app.vercel.app/api/health
   ```

### Custom Deployment (Docker)

```dockerfile
FROM node:18-alpine

WORKDIR /app
COPY package*.json ./
RUN npm install --production

COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "start"]
```

```bash
docker build -t busticket-frontend .
docker run -p 3000:3000 \
  -e NEXT_PUBLIC_FIREBASE_API_KEY=... \
  -e NEXT_PUBLIC_API_BASE_URL=https://api.busticket.vn \
  busticket-frontend
```

---

## Troubleshooting

### Common Issues

#### "Firebase is not initialized"

**Cause**: Missing or invalid Firebase environment variables

**Fix**:

```bash
# Check .env.local
cat .env.local | grep FIREBASE

# Firebase credentials must not be blank
NEXT_PUBLIC_FIREBASE_API_KEY should be 39 chars
NEXT_PUBLIC_FIREBASE_PROJECT_ID should be non-empty
```

#### "401 Unauthorized" on API calls

**Cause**: Access token expired or missing

**Check**:

1. Open DevTools → Application → Cookies
2. Verify `refresh_token` cookie exists (httpOnly, secure)
3. Check network tab: Authorization header present?

**Fix**:

- Refresh page (triggers auto-refresh via refresh token)
- If persists: Clear cookies & login again

#### "Cannot POST /api/auth/refresh-token"

**Cause**: Backend not running or API URL wrong

**Fix**:

```bash
# Verify API is running
curl http://localhost:8080/api/v1/health

# Check .env.local
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Restart dev server
npm run dev
```

#### "Role badge not showing"

**Cause**: User role not loaded or is 0

**Fix**:

```bash
# Check in DevTools → Console
localStorage.getItem('auth-store')
# Should contain: {"user":{"role":2}}

# If empty, login again
```

#### Tests failing with "Cannot find module '@/lib/firebase'"

**Cause**: Jest module resolution not configured correctly

**Fix**:

```bash
# Clear Jest cache
npm test -- --clearCache

# Run tests again
npm test
```

---

## Project Status

### Completed (Assignment 1)

✅ Authentication (Email + Google OAuth)
✅ Authorization (Role-based access control)
✅ Layout & Design System (colors, typography, components)
✅ Dashboard pages (customer, admin, operator)
✅ Developer Tooling (ESLint, Prettier, Husky, Jest)
✅ Live Deployment (Vercel)

### In Progress / Planned

🚧 Real API Integration (backend endpoints)
🚧 Payment Processing (MoMo, ZaloPay)
🚧 E2E Testing (Cypress)
🚧 Multi-language i18n
🚧 Mobile Optimization

See [NEXT_STEPS.md](./NEXT_STEPS.md) for detailed roadmap.

---

## Contributing

1. Fork the repo
2. Create feature branch: `git checkout -b feature/xyz`
3. Commit changes: `git commit -m "feat: add xyz"`
4. Pre-commit hook runs lint & format automatically
5. Push: `git push origin feature/xyz`
6. Create Pull Request

---

## License

[Add license info]

---

## Support

For issues or questions:

1. Check [Troubleshooting](#troubleshooting) section
2. Review [ARCHITECTURE.md](./ARCHITECTURE.md) for system details
3. Check [design/](./design/) folder for design system
4. Open GitHub issue

---

## Additional Resources

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design & data flow
- [design/color-tokens.md](./design/color-tokens.md) - Color system
- [design/typography-scale.md](./design/typography-scale.md) - Typography
- [design/component-specs.md](./design/component-specs.md) - UI components
- [docs/DEPLOYMENT.md](./docs/DEPLOYMENT.md) - Deployment guide
- [NEXT_STEPS.md](./NEXT_STEPS.md) - Future roadmap
