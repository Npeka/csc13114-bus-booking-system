# Typography Scale

## Font Family

- **Display & Headings**: Inter (sans-serif)
  - Weight: 700 (bold)
  - Letter-spacing: -0.5px

- **Body & Small**: Inter (sans-serif)
  - Weight: 400 (regular) for body, 500 (medium) for labels
  - Letter-spacing: 0px

## Heading Scale

| Level | Size | Line Height | Weight | Usage                             |
| ----- | ---- | ----------- | ------ | --------------------------------- |
| H1    | 32px | 1.25        | 700    | Page titles, hero sections        |
| H2    | 28px | 1.33        | 700    | Section headers, dashboard titles |
| H3    | 24px | 1.33        | 700    | Subsection headers, card titles   |
| H4    | 20px | 1.4         | 600    | Minor headings, labels            |
| H5    | 16px | 1.5         | 600    | Tertiary headings                 |
| H6    | 14px | 1.5         | 600    | Smallest headings                 |

## Body Text

| Type       | Size | Line Height | Weight | Usage                              |
| ---------- | ---- | ----------- | ------ | ---------------------------------- |
| Body Large | 18px | 1.5         | 400    | Large body text, prominent content |
| Body       | 16px | 1.5         | 400    | Primary body text, default         |
| Body Small | 14px | 1.43        | 400    | Secondary text, descriptions       |
| Caption    | 12px | 1.33        | 400    | Small labels, metadata             |

## Visual Examples

### Heading Example

```
H1: Bảng điều khiển quản trị
H2: Đặt vé gần đây
H3: Bước 1: Chọn chuyến
Body: Vui lòng chọn một chuyến từ danh sách bên dưới
Caption: Cập nhật lần cuối: 23/11/2025
```

## Implementation in Code

### CSS Classes (Tailwind)

```tsx
// Headings
<h1 className="text-h1">Page Title</h1>
<h2 className="text-h2">Section Title</h2>

// Body text
<p className="text-base">Body text</p>
<p className="text-sm text-muted-foreground">Caption</p>
```

### Tailwind Config

```typescript
fontSize: {
  'h1': ['32px', { lineHeight: '1.25', fontWeight: '700' }],
  'h2': ['28px', { lineHeight: '1.33', fontWeight: '700' }],
  // ...
  'body': ['16px', { lineHeight: '1.5', fontWeight: '400' }],
  'caption': ['12px', { lineHeight: '1.33', fontWeight: '400' }],
}
```

## Spacing Between Text

- Heading + Paragraph: 12px margin-bottom
- Paragraph + Paragraph: 16px margin-bottom
- Heading + Section: 24px margin-bottom
- Section + Section: 32px margin-bottom

## Line Length

Optimal line length for readability: 50–75 characters

- Mobile: 40–50 characters (narrower)
- Desktop: 65–75 characters

Implemented via:

```tsx
<div className="max-w-2xl"> {/* ~65 characters per line */}
```

## Emphasis & Special Cases

### Links

- Underlined by default (accessibility)
- Color: primary
- Hover: darken primary

### Code/Monospace

- Font: Courier New / Monaco
- Size: 85% of body (14px for 16px body)
- Background: subtle background color
- Padding: 2px 4px

### Quotes

- Italic
- Left border: 4px primary
- Padding-left: 16px
- Color: muted-foreground

## Accessibility

1. **Minimum Font Size**: 14px (caption) on mobile, 12px on desktop
2. **Line Height**: Minimum 1.4 for legibility
3. **Contrast**: All text meets WCAG AA (4.5:1)
4. **Font**: Sans-serif for better web readability
5. **Font Scaling**: Support browser zoom up to 200%
