# Design System Inspired by Resend

## 1. Visual Theme & Atmosphere

Resend's website is a dark, cinematic canvas that treats email infrastructure like a luxury product. The entire page is draped in pure black (`#000000`) with text that glows in near-white (`#f0f0f0`), creating a theater-like experience where content performs on a void stage. This isn't the typical developer-tool darkness — it's the controlled darkness of a photography gallery, where every element is lit with intention and nothing competes for attention.

The typography system is the star of the show. Three carefully chosen typefaces create a hierarchy that feels both editorial and technical: Domaine Display (a Klim Type Foundry serif) appears at massive 96px for hero headlines with barely-there line-height (1.00) and negative tracking (-0.96px), creating display text that feels like a magazine cover. ABC Favorit (by Dinamo) handles section headings with an even more aggressive letter-spacing (-2.8px at 56px), giving a compressed, engineered quality to mid-tier text. Inter takes over for body and UI, providing the clean readability that lets the display fonts shine. Commit Mono rounds out the family for code blocks.

What makes Resend distinctive is its icy, blue-tinted border system. Instead of neutral gray borders, Resend uses `rgba(214, 235, 253, 0.19)` — a frosty, slightly blue-tinted line at 19% opacity that gives every container and divider a cold, crystalline quality against the black background. Combined with pill-shaped buttons (9999px radius), multi-color accent system (orange, green, blue, yellow, red — each with its own CSS variable scale), and OpenType stylistic sets (`"ss01"`, `"ss03"`, `"ss04"`, `"ss11"`), the result is a design system that feels premium, precise, and quietly confident.

**Key Characteristics:**
- Pure black background with near-white (`#f0f0f0`) text — theatrical, gallery-like darkness
- Three-font hierarchy: Domaine Display (serif hero), ABC Favorit (geometric sections), Inter (body/UI)
- Icy blue-tinted borders: `rgba(214, 235, 253, 0.19)` — every border has a cold, crystalline shimmer
- Multi-color accent system: orange, green, blue, yellow, red — each with numbered CSS variable scales
- Pill-shaped buttons and tags (9999px radius) with transparent backgrounds
- OpenType stylistic sets (`"ss01"`, `"ss03"`, `"ss04"`, `"ss11"`) on display fonts
- Commit Mono for code — monospace as a design element, not an afterthought
- Whisper-level shadows using blue-tinted ring: `rgba(176, 199, 217, 0.145) 0px 0px 0px 1px`

## 2. Color Palette & Roles

### Primary
- **Void Black** (`#000000`): Page background, the defining canvas color
- **Near White** (`#f0f0f0`): Primary text, button text, high-contrast elements
- **Pure White** (`#ffffff`): Maximum emphasis text, link highlights

### Accent Scale — Orange
- **Orange 4** (`#ff5900`): at 22% opacity — subtle warm glow
- **Orange 10** (`#ff801f`): primary orange accent — warm, energetic
- **Orange 11** (`#ffa057`): lighter orange for secondary use

### Accent Scale — Green
- **Green 3** (`#22ff99`): at 12% opacity — faint emerald wash
- **Green 4** (`#11ff99`): at 18% opacity — success indicator glow
- **Green 10** (`#22ff99`): bright green

### Accent Scale — Blue
- **Blue 4** (`#0075ff`): at 34% opacity — medium blue accent
- **Blue 5** (`#0081fd`): at 42% opacity — stronger blue
- **Blue 10** (`#3b9eff`): bright blue — links, interactive elements

### Accent Scale — Other
- **Yellow 9** (`#ffc53d`): warm gold for warnings or highlights
- **Red 5** (`#ff2047`): at 34% opacity — error states, destructive actions
- **Red 10** (`#ff2047`): bright red

### Neutral Scale
- **Silver** (`#a1a4a5`): Secondary text, muted links, descriptions
- **Dark Gray** (`#464a4d`): Tertiary text, de-emphasized content

### Borders & Shadows
- **Frost Border** (`rgba(214, 235, 253, 0.19)`): The signature — icy blue-tinted borders
- **Frost Border Alt** (`rgba(217, 237, 254, 0.145)`): Lighter variant
- **Ring Shadow** (`rgba(176, 199, 217, 0.145) 0px 0px 0px 1px`): Blue-tinted shadow-as-border
- **Focus Ring** (`rgb(0, 0, 0) 0px 0px 0px 8px`): Heavy black focus ring
- **White Hover** (`rgba(255, 255, 255, 0.28)`): Button hover state on dark

## 3. Typography Rules

### Font Families
- **Display Serif**: `domaine` (Domaine Display) — hero headlines only
- **Display Sans**: `aBCFavorit` (ABC Favorit), fallback `Inter` — section headings, nav
- **Body / UI**: `Inter` — body text, buttons, forms, all UI
- **Monospace**: `commitMono`, fallback `ui-monospace` — code blocks

### Hierarchy

| Role | Font | Size | Weight | Line Height | Letter Spacing |
|------|------|------|--------|-------------|----------------|
| Display Hero | domaine | 96px | 400 | 1.00 | -0.96px |
| Section Heading | aBCFavorit | 56px | 400 | 1.20 | -2.8px |
| Sub-heading | aBCFavorit | 20px | 400 | 1.30 | normal |
| Feature Title | Inter | 24px | 500 | 1.50 | normal |
| Body Large | Inter | 18px | 400 | 1.50 | normal |
| Body | Inter | 16px | 400 | 1.50 | normal |
| Nav Link | aBCFavorit | 14px | 500 | 1.43 | +0.35px |
| Button / UI | Inter | 14px | 500 | 1.43 | normal |
| Caption | Inter | 14px | 400 | 1.60 | normal |
| Small | Inter | 12px | 500 | 1.33 | normal |
| Code | commitMono | 14–16px | 400 | 1.50 | normal |

## 4. Component Stylings

### Buttons
**Primary Transparent Pill**: transparent bg, `#f0f0f0` text, frost border, 5px 12px padding, 9999px radius. Hover: `rgba(255,255,255,0.28)` bg.
**White Solid Pill**: `#ffffff` bg, `#000000` text, 9999px radius. High-contrast CTA.
**Ghost Button**: transparent, `#f0f0f0` text, 4px radius, no border. Hover: subtle tint.
**Danger Button**: `rgba(255,32,71,0.34)` bg, `#ff2047` text, frost border.

### Cards & Containers
- Background: transparent or `rgba(255,255,255,0.03)`
- Border: `1px solid rgba(214, 235, 253, 0.19)`
- Radius: 16px (standard), 24px (large sections)
- Shadow: `rgba(176, 199, 217, 0.145) 0px 0px 0px 1px`

### Forms
- Input bg: transparent or very dark
- Input border: frost border `rgba(214, 235, 253, 0.19)`
- Input radius: 8px
- Focus: blue ring `rgba(0,117,255,0.34)`
- Error border: `rgba(255,32,71,0.34)`
- Label: Inter 14px weight 500, `#a1a4a5`
- Error message: Inter 13px, `#ff2047`

### Navigation
- Sticky, `rgba(0,0,0,0.85)` + `backdrop-filter: blur(12px)`
- Bottom border: frost border
- Nav links: aBCFavorit 14px 500 +0.35px tracking
- Pill CTAs right-aligned

### Badges
- Pill shape (9999px), 2px 10px padding, 12px Inter weight 500
- Orange: `rgba(255,89,0,0.22)` bg, `#ffa057` text
- Green: `rgba(34,255,153,0.12)` bg, `#22ff99` text
- Blue: `rgba(0,117,255,0.34)` bg, `#3b9eff` text
- Yellow: `rgba(255,197,61,0.20)` bg, `#ffc53d` text
- Red: `rgba(255,32,71,0.34)` bg, `#ff2047` text

## 5. Layout Principles

### Page Background
- `#000000` — pure void black. Non-negotiable.

### Spacing System (base 8px)
- Scale: 4, 6, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96px

### Border Radius Scale
- Sharp: 4px (inputs, ghost buttons)
- Subtle: 6px (menu panels)
- Standard: 8px (tabs, form inputs)
- Card: 16px (feature cards)
- Section: 24px (large panels)
- Pill: 9999px (CTAs, badges, tags)

## 6. CSS Classes Available (from app.css)

```
Buttons:
  .btn          → base — Inter 14px 500, pill 9999px, 5px 12px
  .btn-primary  → transparent bg, frost border, #f0f0f0 text
  .btn-white    → white bg, black text
  .btn-ghost    → transparent, no border, 4px radius
  .btn-danger   → red tinted
  .btn-lg       → 10px 20px, 16px font
  .btn-sm       → 3px 10px, 12px font

Forms:
  .form-group   → flex column, gap 6px, margin-bottom 20px
  .form-label   → Inter 14px 500, #a1a4a5
  .form-input   → dark bg, frost border, 8px radius, 10px 14px padding
  .form-error   → 13px, #ff2047
  .input-error  → on .form-input when invalid

Cards:
  .card         → transparent bg, frost border, 16px radius, ring shadow

Badges:
  .badge        → pill, 2px 10px
  .badge-orange → warm orange
  .badge-green  → emerald green
  .badge-blue   → bright blue
  .badge-yellow → warm gold
  .badge-red    → danger red
```

## 7. Do's and Don'ts

### Do
- Use `#000000` as page background — always
- Apply frost borders (`rgba(214, 235, 253, 0.19)`) for ALL structural lines
- Use pill radius (9999px) for CTAs and badges
- Use transparent buttons with frost borders on dark backgrounds
- Keep Inter for all body/UI text
- Use multi-color accents at low opacity for backgrounds, full for text

### Don't
- Don't use gray neutral borders — always use the frost blue tint
- Don't make opaque colored backgrounds — use tinted opacity variants
- Don't use traditional box-shadows on dark — use frost borders for depth
- Don't use positive letter-spacing on headings (only nav links get +0.35px)
- Don't lighten the page background above `#000000`

## 8. Agent Prompt Guide

### Quick Color Reference
- Page BG: `#000000` (pure black)
- Primary text: `#f0f0f0`
- Secondary text: `#a1a4a5`
- Tertiary text: `#464a4d`
- Border: `rgba(214, 235, 253, 0.19)` (frost)
- Orange: `#ff801f`
- Green: `#22ff99`
- Blue: `#3b9eff`
- Yellow: `#ffc53d`
- Red: `#ff2047`
- White hover: `rgba(255, 255, 255, 0.28)`
