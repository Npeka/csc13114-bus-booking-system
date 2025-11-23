# Color Tokens & Design System

## Color Palette

Our design system uses oklch color space for superior perceptual uniformity and accessibility.

### Primary Colors

- **Primary**: `oklch(64% 0.22 29)` - Brand red (#E63946)
  - Primary action buttons, links, active states
  - Accessible contrast ratio: 4.5:1 on white

- **Secondary**: `oklch(50% 0.18 200)` - Brand blue
  - Secondary actions, hover states
  - Supporting UI elements

### Semantic Colors

#### Success (Positive)

- **Success**: `oklch(67% 0.25 146)`
- Use cases: Confirmed bookings, successful actions, positive states

#### Error (Negative)

- **Error**: `oklch(59% 0.25 31)`
- Use cases: Cancelled bookings, errors, destructive actions
- Accessibility: Strong visual distinction for colorblind users

#### Warning (Caution)

- **Warning**: `oklch(65% 0.23 61)`
- Use cases: Pending actions, warnings, caution states

#### Info (Neutral)

- **Info**: `oklch(64% 0.21 257)`
- Use cases: Completed bookings, informational states

### Neutral Colors

- **Foreground**: Text, primary content
- **Muted Foreground**: Secondary text, captions, hints
- **Background**: Page background
- **Muted**: Subtle backgrounds, borders
- **Card**: Card backgrounds

## Semantic Color Usage

### Booking Status Colors

| Status    | Color   | Usage                     |
| --------- | ------- | ------------------------- |
| Confirmed | Success | Approved, active bookings |
| Pending   | Warning | Awaiting confirmation     |
| Completed | Info    | Past bookings             |
| Cancelled | Error   | Cancelled bookings        |
| Refunded  | Info    | Refund issued             |

### UI Element Colors

| Element             | Color     | Notes               |
| ------------------- | --------- | ------------------- |
| Buttons (Primary)   | Primary   | Main CTAs           |
| Buttons (Secondary) | Secondary | Alternative actions |
| Badges              | Primary   | Default status      |
| Links               | Primary   | Hyperlinks          |
| Borders             | Muted     | Subtle separation   |
| Focus Ring          | Primary   | Keyboard focus      |

## Dark Mode Support

All colors support both light and dark mode via CSS variables.

```css
@light {
  --background: oklch(98% 0 0);
  --foreground: oklch(20% 0 0);
}

@dark {
  --background: oklch(15% 0 0);
  --foreground: oklch(95% 0 0);
}
```

## Accessibility Considerations

1. **Color Contrast**: All colors meet WCAG AA standards (4.5:1 min)
2. **Colorblind Safe**: Semantic meaning conveyed with icons + color
3. **Dark Mode**: Adequate contrast in both themes
4. **Perception Uniformity**: oklch ensures consistent perceived brightness

## Implementation

Colors are defined in Tailwind CSS config:

```typescript
colors: {
  primary: 'oklch(64% 0.22 29)',
  secondary: 'oklch(50% 0.18 200)',
  success: 'oklch(67% 0.25 146)',
  error: 'oklch(59% 0.25 31)',
  // ...
}
```

Usage in components:

```tsx
<Button className="bg-primary">Confirm</Button>
<Badge className="bg-success/10 text-success">Approved</Badge>
```
