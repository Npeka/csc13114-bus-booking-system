# Component Specifications

## Overview

This document outlines the reusable UI components used throughout the BusTicket.vn application. All components follow the design system tokens (colors, typography, spacing) defined in `color-tokens.md` and `typography-scale.md`.

## Core Components

### Button

**Props**:

```typescript
interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?:
    | "default"
    | "destructive"
    | "outline"
    | "secondary"
    | "ghost"
    | "link";
  size?: "default" | "sm" | "lg" | "icon";
  asChild?: boolean;
}
```

**Variants**:
| Variant | Background | Use Case |
|---------|------------|----------|
| default | Primary | Main actions, CTAs |
| destructive | Error | Delete, cancel operations |
| outline | Transparent + border | Secondary actions |
| secondary | Secondary | Alternative actions |
| ghost | None | Subtle actions, toolbar buttons |
| link | None, underlined text | Links within content |

**Sizes**:
| Size | Height | Padding | Font Size |
|------|--------|---------|-----------|
| sm | 32px | 8px 12px | 14px |
| default | 40px | 12px 16px | 16px |
| lg | 48px | 16px 20px | 18px |
| icon | 40px | 8px | - |

**Examples**:

```tsx
<Button>Confirm Booking</Button>
<Button variant="outline">Cancel</Button>
<Button size="sm">Save</Button>
<Button variant="destructive" size="lg">Delete Account</Button>
```

---

### Card

**Props**:

```typescript
interface CardProps extends HTMLAttributes<HTMLDivElement> {}
```

**Sub-components**: `CardHeader`, `CardContent`, `CardFooter`, `CardTitle`, `CardDescription`

**Usage**:

```tsx
<Card>
  <CardHeader>
    <CardTitle>Booking Details</CardTitle>
    <CardDescription>Confirm your booking</CardDescription>
  </CardHeader>
  <CardContent>{/* Content */}</CardContent>
  <CardFooter>{/* Actions */}</CardFooter>
</Card>
```

**Styling**:

- Background: card color
- Border: 1px solid muted
- Border-radius: 8px
- Padding: 24px (header), 16px (content/footer)
- Shadow: 0 1px 3px rgba(0,0,0,0.1)

---

### Badge

**Props**:

```typescript
interface BadgeProps extends HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "secondary" | "destructive" | "outline";
}
```

**Variants**:
| Variant | Style | Use Case |
|---------|-------|----------|
| default | Primary background, white text | Active status |
| secondary | Secondary background | Secondary status |
| destructive | Error background | Error/warning status |
| outline | Transparent, border | Neutral status |

**Examples**:

```tsx
<Badge>Confirmed</Badge>
<Badge variant="secondary" className="bg-success/10 text-success">Approved</Badge>
<Badge variant="destructive">Cancelled</Badge>
```

---

### Dialog / Modal

**Props**:

```typescript
interface DialogProps {
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}
```

**Sub-components**: `DialogTrigger`, `DialogContent`, `DialogHeader`, `DialogTitle`, `DialogDescription`, `DialogFooter`

**Usage**:

```tsx
<Dialog>
  <DialogTrigger asChild>
    <Button>Open Dialog</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Confirm Action</DialogTitle>
    </DialogHeader>
    {/* Content */}
    <DialogFooter>
      <Button>Cancel</Button>
      <Button>Confirm</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

---

### Form Components

#### Input

```tsx
<Input placeholder="Enter email" type="email" />
```

#### Select

```tsx
<Select>
  <SelectTrigger>
    <SelectValue placeholder="Choose option" />
  </SelectTrigger>
  <SelectContent>
    <SelectItem value="opt1">Option 1</SelectItem>
    <SelectItem value="opt2">Option 2</SelectItem>
  </SelectContent>
</Select>
```

#### Checkbox

```tsx
<Checkbox id="terms" />
<label htmlFor="terms">I agree to terms</label>
```

---

### Tabs

**Props**:

```typescript
interface TabsProps {
  value?: string;
  onValueChange?: (value: string) => void;
  defaultValue?: string;
}
```

**Usage**:

```tsx
<Tabs defaultValue="upcoming">
  <TabsList>
    <TabsTrigger value="upcoming">Upcoming</TabsTrigger>
    <TabsTrigger value="past">Past</TabsTrigger>
  </TabsList>
  <TabsContent value="upcoming">Content for upcoming</TabsContent>
  <TabsContent value="past">Content for past</TabsContent>
</Tabs>
```

---

### Dropdown Menu

**Usage**:

```tsx
<DropdownMenu>
  <DropdownMenuTrigger asChild>
    <Button variant="ghost" size="icon">
      <MenuIcon />
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent>
    <DropdownMenuItem>Profile</DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem>Logout</DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

---

## Spacing Scale

All spacing uses a 4px base unit:

| Token | Value | Use Case           |
| ----- | ----- | ------------------ |
| xs    | 4px   | Micro spacing      |
| sm    | 8px   | Small gaps         |
| md    | 16px  | Default spacing    |
| lg    | 24px  | Larger sections    |
| xl    | 32px  | Major sections     |
| 2xl   | 48px  | Page-level spacing |

**Implementation**:

```tsx
<div className="space-y-md">
  {" "}
  {/* 16px gap */}
  <div>Item 1</div>
  <div>Item 2</div>
</div>
```

---

## Border & Radius

- **Default border radius**: 8px (`rounded-lg`)
- **Subtle radius**: 4px (`rounded`)
- **Pill shape**: 9999px (`rounded-full`)
- **Border width**: 1px (default)

---

## Shadow System

- **None**: No shadow (flat design)
- **sm**: `0 1px 2px 0 rgba(0, 0, 0, 0.05)`
- **md**: `0 4px 6px -1px rgba(0, 0, 0, 0.1)`
- **lg**: `0 10px 15px -3px rgba(0, 0, 0, 0.1)`
- **xl**: `0 20px 25px -5px rgba(0, 0, 0, 0.1)`

**Usage**:

```tsx
<Card className="shadow-md">
```

---

## Responsive Behavior

All components are mobile-first and responsive:

```tsx
<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
  {/* Stacks on mobile, 2 cols on tablet, 3 cols on desktop */}
</div>
```

Breakpoints:

- Mobile: < 640px
- sm: 640px
- md: 768px
- lg: 1024px
- xl: 1280px
- 2xl: 1536px

---

## Custom Components

### RoleBadge

**Props**:

```typescript
interface RoleBadgeProps {
  userRole: number;
  className?: string;
}
```

**Usage**:

```tsx
<RoleBadge userRole={user.role} />
// Displays: "Quản trị viên" (Admin), "Nhà điều hành" (Operator), etc.
```

---

## Accessibility

All components implement:

- Keyboard navigation (Tab, Enter, Esc)
- ARIA labels for screen readers
- Focus management
- Color-independent status indication
- Minimum touch target size: 48px × 48px
