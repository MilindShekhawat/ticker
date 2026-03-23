# Ticker - Design Specification

## 1: Design Philosophy

**Dark-first. Industrial. Exposed structure.**

Ticker's interface draws from technical schematics, engineering tools, and retro-industrial hardware. It looks like something that belongs on a workstation — not a consumer app. The aesthetic sits at the intersection of Marathon's sci-fi UI, Teenage Engineering's blocky hardware, and Factorio's "exposed internals" design language.

**Guiding Principles**

- **Function over form**: Every pixel earns its place. If it doesn't help the user manage work, it doesn't belong.
- **Density without clutter**: Show more information in less space. Use spacing and typography hierarchy — not boxes and borders — to separate concerns.
- **Quiet until needed**: The interface recedes. Actions, states, and data come forward only when relevant.
- **Structure is decoration**: Depth comes from wireframes, topology, exposed grids, and transparency — never from shadows or gradients. The skeleton of the UI *is* the visual language.
- **Keyboard-native**: Every interactive element is reachable and operable via keyboard. Focus states are always visible.

**Reference Board**

| Source                  | What we take                                                  |
| ----------------------- | ------------------------------------------------------------- |
| Marathon (game UI)      | Lime-on-black, serif display type, sharp rectangles, uppercase labels |
| Marathon (map website)  | Topographic wireframe textures, structural line patterns      |
| Teenage Engineering     | Blocky industrial forms, exposed-mechanism feel, bold function labeling |
| Factorio                | "Exposed internals" — showing inner structure gives identity and personality |
| Mux / Zed.dev           | Extended rule lines that bleed past containers                |
| Apple (blend modes)     | Transparent layering, edges as borders, color overlap         |
| C2MTL                   | CSS blend modes for accent moments                            |
| shadcn/ui               | Component API patterns, variant system (not visual style)     |

## 2: Color System

Colors use **Tailwind's built-in palette** — no custom LCH tokens. Zinc for neutrals, Lime for primary. Semantic and accent colors map directly to Tailwind color families.

### Neutral Base — Zinc

| Role        | Tailwind Class  | Usage                          |
| ----------- | --------------- | ------------------------------ |
| Base        | `zinc-950`      | Page root, deepest layer       |
| Surface     | `zinc-900`      | Cards, panels, sidebars        |
| Elevated    | `zinc-800`      | Hover fills, active rows       |
| Overlay     | `zinc-700`      | Dropdowns, tooltips, modals    |
| Invert      | `zinc-100`      | Inverted (light) surfaces      |

### Semantic Colors

| Role        | Tailwind Family | Usage                          |
| ----------- | --------------- | ------------------------------ |
| **Primary** | `lime`          | Buttons, active states, accent |
| **Success** | `green`         | Confirmations, positive states |
| **Danger**  | `red`           | Destructive actions, errors    |
| **Info**    | `cyan`          | Informational, links           |
| **Warning** | `yellow`        | Caution, attention needed      |

### Accent Colors

For tags, categories, and decorative differentiation:

| Role        | Tailwind Family | Usage                          |
| ----------- | --------------- | ------------------------------ |
| **Indigo**  | `indigo`        | Tags, categories               |
| **Pink**    | `pink`          | Tags, decorative               |
| **Purple**  | `purple`        | Tags, categories               |

### Text Colors

| Role        | Tailwind Class  | Usage                              |
| ----------- | --------------- | ---------------------------------- |
| Primary     | `zinc-100`      | Headings, labels, body             |
| Secondary   | `zinc-400`      | Descriptions, metadata             |
| Muted       | `zinc-500`      | Placeholders, hints                |
| Disabled    | `zinc-600`      | Disabled state text                |
| Invert      | `zinc-950`      | Text on light/primary backgrounds  |

### Borders

| Role        | Tailwind Class  | Usage                       |
| ----------- | --------------- | --------------------------- |
| Subtle      | `zinc-800`      | Hairlines, dividers         |
| Default     | `zinc-700`      | Component edges             |
| Strong      | `zinc-600`      | Emphasis, focus-adjacent    |
| Focus       | `lime-400`      | Focus rings, active borders |

### Color Blending

CSS `mix-blend-mode` is used for accent and decorative moments — not for core UI. Blend modes create depth through transparency and color interaction rather than shadows.

**Where to use blending:**

- Hero/landing areas — overlapping colored shapes with `mix-blend-mode: screen` or `multiply`
- Hover effects on accent elements — layered transparency instead of simple color swap
- Data visualization overlaps — when radar charts or status indicators overlap, blend rather than occlude
- Background textures — wireframe patterns at low opacity with `mix-blend-mode: lighten` on dark base

**Where NOT to use blending:**

- Core UI components (buttons, inputs, cards) — these stay solid and predictable
- Text — never blend text, always opaque and legible
- Interactive states that need to be unambiguous (selected, disabled, error)

## 3: Typography

### Font Stack

| Role           | Family                                        | Usage                                        |
| -------------- | --------------------------------------------- | -------------------------------------------- |
| Display        | Xanh Mono (400)                               | Major headings (h1, h2), hero text, branding |
| Primary        | JetBrains Mono (400, 500, 600, 700)           | Everything else — body, labels, UI, code     |

**Xanh Mono** is a serif monospace with distinctive character. It's used sparingly — only for display-level headings where its single weight (400) is not a limitation. Its serifs give major headings a typographic quality that contrasts with the utilitarian JetBrains Mono used everywhere else. This mirrors Marathon's serif display headings ("CRYO ARCHIVE", "TAU CETI IV").

**JetBrains Mono** is the workhorse. Body text, form inputs, buttons, badges, metadata, table cells, ticket IDs — everything renders in JetBrains Mono. This is a monospace-first interface. The tool feels like a terminal because it *is* typeset like one.

### Serif Mono Alternatives

If Xanh Mono's single weight becomes limiting for future needs:

| Font            | Weights         | Character                                    |
| --------------- | --------------- | -------------------------------------------- |
| Space Mono      | 400, 700        | Retro-techy, pairs well with the aesthetic   |
| Courier Prime   | 400, 700        | Clean serif mono, literary feel              |
| DM Mono         | 300, 400, 500   | Subtle serifs, geometric, versatile          |

### Type Scale

Fixed pixel values — no relative scaling. The interface is a tool, not a document.

| Token        | Size  | Font         | Usage                                |
| ------------ | ----- | ------------ | ------------------------------------ |
| `--text-xs`  | 11px  | JetBrains    | Badges, micro-labels, timestamps     |
| `--text-sm`  | 12px  | JetBrains    | Secondary text, table cells, hints   |
| `--text-base`| 14px  | JetBrains    | Body text, form inputs (default)     |
| `--text-md`  | 15px  | JetBrains    | Emphasized body                      |
| `--text-lg`  | 18px  | JetBrains    | Section headings, card titles        |
| `--text-xl`  | 22px  | Xanh Mono    | Page titles                          |
| `--text-2xl` | 28px  | Xanh Mono    | Hero headings (used sparingly)       |
| `--text-3xl` | 36px  | Xanh Mono    | Display text (reserved)              |

### Line Height & Tracking

| Token                | Value    | Usage                              |
| -------------------- | -------- | ---------------------------------- |
| `--leading-tight`    | 1.2      | Headings, compact labels           |
| `--leading-normal`   | 1.5      | Body text (default)                |
| `--leading-relaxed`  | 1.75     | Long-form content, descriptions    |
| `--tracking-tight`   | -0.02em  | Large headings                     |
| `--tracking-normal`  | 0em      | Body text                          |
| `--tracking-wide`    | 0.04em   | Badges, small mono text            |
| `--tracking-wider`   | 0.08em   | Uppercase labels, section headers  |

## 4: Visual Language — Depth & Structure

This is the core differentiator. Ticker doesn't use traditional depth cues (shadows, gradients, elevation). Instead, depth comes from **exposed structure** — wireframes, grid lines, transparency, and layered composition.

### 4.1: Extended Rule Lines

Inspired by Mux and Zed.dev — decorative border lines that extend beyond their container, creating a technical-drawing feel.

**Where to use:**
- Section dividers in the header and between content areas
- Extending from headings as horizontal rules
- Table header separators that bleed to viewport edges
- Active/selected states — a rule line extends to mark the current item

**Implementation:** Pseudo-elements (`::before`, `::after`) with `position: absolute` extending past parent bounds. Thin (1px) lines in `zinc-700` or `zinc-800`. Active states use `lime-400`.

### 4.2: Wireframe Textures

Inspired by Marathon's topographic map — fine structural line patterns used as background texture to create visual depth without filling with color.

**Where to use:**
- Empty states (no tickets, no projects) — a faint wireframe grid instead of a blank void
- Card backgrounds at very low opacity — gives "blueprint" feel
- Dashboard background — subtle topographic or grid lines underneath content
- Loading states — wireframe skeleton instead of pulsing blocks

**Implementation:** SVG patterns or CSS `repeating-linear-gradient` at 3-5% opacity. Line color: `zinc-800` or `zinc-700`. Never competes with content.

### 4.3: Exposed Internals (Factorio Principle)

Elements reveal their structure rather than hiding behind a polished surface. Components look like they have visible mechanisms.

**Where to use:**
- **Status indicators**: Instead of a simple colored dot, show a small schematic — a filled/empty gauge, a circuit-like pattern
- **Ticket cards**: Subtle internal grid lines or section marks visible at card level — like a printed circuit board
- **Icons**: Geometric, grid-aligned, showing "how it works" rather than abstracting it away
- **Progress states**: Wire-frame progress bars that fill with structure, not solid color

**Where NOT to use:**
- Text content — body text stays clean and unadorned
- Form inputs — functional areas stay crisp
- When it reduces legibility or adds cognitive load

### 4.4: Transparent Layering & Blend Modes

Inspired by Apple's transparent logo treatment and C2MTL's color overlap compositions.

**Where to use:**
- **Overlapping data viz**: Radar/spider charts where datasets overlap — use `mix-blend-mode: screen` so overlaps create new colors rather than occluding
- **Tag clusters**: When multiple tags sit near each other, their subtle backgrounds can blend
- **Hero/branding moments**: Landing page, project headers — layered colored shapes at low opacity
- **Focus/hover accents**: Instead of solid color swap, layer a transparent lime wash with `mix-blend-mode: lighten`

**CSS approach:**
```css
.blend-layer {
    mix-blend-mode: screen;   /* lighten blend — good on dark bg */
    opacity: 0.7;
}
.blend-multiply {
    mix-blend-mode: multiply; /* darken blend — good for overlapping fills */
}
```

### 4.5: Radar/Spider Charts

Polygon-based data visualization for project health, tag distribution, or status overview.

**Where to use:**
- Project dashboard — show ticket distribution across statuses as a radar chart
- Sprint/milestone health — visualize progress across multiple dimensions
- Member workload — compare assignment distribution

**Visual style:**
- Wireframe grid lines (concentric polygons) in `zinc-800`
- Data polygon with colored fill at 20-30% opacity + solid stroke
- Blend mode on fill so overlapping datasets create new color values
- Grid labels in JetBrains Mono `--text-xs`

### 4.6: Depth Hierarchy (Summary)

Instead of shadow-based elevation, Ticker uses a structural depth system:

| Depth Level | Technique                                              | Example                          |
| ----------- | ------------------------------------------------------ | -------------------------------- |
| 0 — Base    | Flat `zinc-950`, optional wireframe texture             | Page background                  |
| 1 — Surface | Solid `zinc-900` + `zinc-800` border                   | Cards, panels                    |
| 2 — Raised  | Extended rule lines + slightly lighter border           | Active cards, selected items     |
| 3 — Overlay | Solid `zinc-900` + `shadow-md` (exception: overlays need shadow for z-context) | Dropdowns, modals |
| 4 — Accent  | Transparent color layer with blend mode                 | Hover effects, data viz overlaps |

Shadows are *only* used at depth level 3 (overlays) where z-index stacking needs a visual cue. Everything else uses structure and borders.

## 5: Motion

Transitions are fast and functional. Nothing eases in slowly or bounces.

| Token               | Value                         | Usage                    |
| ------------------- | ----------------------------- | ------------------------ |
| `--duration-fast`   | 100ms                         | Color, opacity changes   |
| `--duration-normal` | 150ms                         | Default interactions     |
| `--duration-slow`   | 250ms                         | Layout shifts, reveals   |
| `--ease`            | cubic-bezier(0.16, 1, 0.3, 1)| Custom ease-out curve    |

No animations for decoration. Transitions exist only to smooth state changes and provide feedback.

## 6: Component Library

All components are custom-built using **CVA** (class-variance-authority) for type-safe variant composition. No external component library.

### Button

Uppercase text. Border-driven hover states. Sharp corners. Six variants.

| Variant       | Resting State                          | Hover State                           |
| ------------- | -------------------------------------- | ------------------------------------- |
| `primary`     | Lime bg, dark text, lime border        | Dark bg, lime text, lime border       |
| `white`       | White bg, dark text, white border      | Dark bg, lime text, lime border       |
| `outline`     | Dark bg, white text, white border      | Dark bg, lime text, lime border       |
| `black`       | Dark bg, white text, no visible border | Dark bg, lime text, lime border       |
| `invert`      | Dark bg, white text, no visible border | Lime bg, dark text, lime border       |
| `destructive` | Red bg, white text, red border         | Dark bg, red text, red border         |

Sizes: `xs` (h-4), `sm` (h-6), `md` (h-8, default), `lg` (h-10), `xl` (h-12)

All buttons show `disabled:opacity-40` and `disabled:pointer-events-none` when disabled.

### Input / Textarea

- Dark background (`zinc-900`), subtle border (`zinc-700`)
- Placeholder text in `zinc-500` (muted)
- Focus: lime border, no outline ring — the border *is* the focus indicator
- Error: red border, red error text below the field
- Sharp corners, no rounding
- Font: JetBrains Mono throughout

### Card

Composable compound component: `Card`, `CardHeader`, `CardTitle`, `CardDescription`, `CardContent`, `CardFooter`.

- Background: `zinc-900`
- Border: `zinc-800`
- Interactive variant: `hover:border-lime-400` on click
- Padding: `p-6` for header/content/footer sections
- Sharp corners
- Optional: faint wireframe grid texture on card background

### Badge

Rectangular, mono font, uppercase-adjacent tracking.

- Font: JetBrains Mono, 11px, weight 500
- No border radius — sharp rectangular badges
- Padding: 2px 8px
- Each semantic/accent color has a badge variant using subtle bg + text color + border

### Select (Dropdown)

Custom-built dropdown (no native `<select>`).

- Trigger styled as input field
- Dropdown panel rendered via portal
- Full keyboard navigation: Arrow keys, Enter, Escape
- ARIA attributes for accessibility
- Sharp-cornered panel with shadow (overlay exception — see depth level 3)

### Loader

Retro character-noise animation. Two variants:

- **Noise**: Random characters cycling rapidly in a grid
- **Scanline**: Blinking cursor-style animation

## 7: Layout Architecture

### AppFrame

The application shell wraps all authenticated pages.

```
┌─────────────────────────────────────────────────────────────┐
│ HEADER (h-14, sticky, z-40)                                 │
│ ┌────────┬──────────────────────────────┬────────┬────────┐ │
│ │ Logo   │ Section (page controls)      │Settings│Profile │ │
│ │ 112px  │ flexible                     │  56px  │  56px  │ │
│ └────────┴──────────────────────────────┴────────┴────────┘ │
├─────────────────────────────────────────────────────────────┤
│ MAIN CONTENT (px-4 py-4)                                    │
│                                                             │
│ Full width, scrollable                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

- **Logo**: "ticker" in Xanh Mono, lowercase, routes to dashboard
- **Section slot**: Receives page-specific controls (breadcrumbs, filters, view toggles)
- **Settings**: Pixel-art gear icon on lime background
- **Profile**: User avatar image
- **Extended rules**: Header bottom border extends full viewport width as a structural line

### Page Layouts

| Page             | Layout                                                |
| ---------------- | ----------------------------------------------------- |
| Login / Signup   | Centered card on `bg-base`, no AppFrame               |
| Dashboard        | AppFrame → responsive grid (1→2→3→4 cols)             |
| Table View       | AppFrame → header controls + table/card rows          |
| Kanban View      | AppFrame → horizontal scrolling columns               |
| Ticket Detail    | AppFrame → single-column detail layout                |
| Settings         | AppFrame → settings panels                            |

## 8: Interaction Patterns

### Hover States

All interactive elements use **border color transitions** as the primary hover signal, not background color changes. This keeps the dark base stable and avoids visual noise.

- Buttons: border transitions to lime (or red for destructive)
- Cards: border transitions to lime on interactive cards
- Links: text color transitions to lime
- Accent option: transparent lime wash via blend mode on hover (for special elements)

### Focus States

Global focus ring: `2px solid lime-400` with `2px offset`. Applied via `:focus-visible` — invisible for mouse users, visible for keyboard navigation.

### Selection

Text selection uses lime-tinted background with lime foreground text.

### Scrollbar

Custom webkit scrollbar (8px width):
- Track: `zinc-950`
- Thumb: `zinc-700`, sharp corners
- Thumb hover: `zinc-600`

## 9: Iconography

**No icon library.** All icons are inline SVGs — geometric, pixel-art inspired, built on a grid. Following the Factorio "exposed internals" principle: icons hint at inner workings rather than abstracting them away.

New icons should follow these rules:
- Build on a grid (400x400 or divisible equivalent)
- Use `polygon` or `rect` elements — no curves, no strokes
- Fill via `currentColor` so they inherit text color
- Keep detail minimal — these are functional glyphs, not illustrations
- Where possible, show "how it works" — a settings icon shows gears with visible teeth, a filter icon shows a literal funnel with internal lines

## 10: Responsive Behavior

The interface is **desktop-first** but adapts to smaller viewports:

- Project grid: 4 cols → 3 → 2 → 1 (via Tailwind responsive grid)
- Header: fixed layout, logo and controls stay visible
- Tables: horizontal scroll on narrow viewports
- Kanban: horizontal scroll is native to the pattern
- Forms: full-width on mobile, constrained on desktop
- Extended rule lines: clip to viewport on mobile, full bleed on desktop

Breakpoints follow Tailwind defaults: `sm` (640px), `md` (768px), `lg` (1024px), `xl` (1280px).

## 11: Accessibility

- **Color contrast**: All text-on-background combinations target WCAG AA minimum
- **Focus management**: `:focus-visible` on all interactive elements with lime outline
- **Keyboard navigation**: Custom components (Select, etc.) implement full arrow-key/enter/escape handling
- **ARIA attributes**: Custom widgets carry appropriate `role`, `aria-expanded`, `aria-selected`, etc.
- **No color-only signaling**: Status and severity always pair color with text labels
- **Blend modes**: Blended elements always have a fallback solid color for accessibility; never rely on blend output for meaning
- **Reduced motion**: Transitions are short enough (100-250ms) to be non-disruptive; a `prefers-reduced-motion` media query should be added for users who need it

## 12: Tech Stack

| Layer           | Technology                                     |
| --------------- | ---------------------------------------------- |
| Framework       | Preact 10.x                                    |
| Build           | Vite 7.x + @preact/preset-vite                 |
| Styling         | Tailwind CSS 4.x + custom CSS variables        |
| Variants        | class-variance-authority (CVA)                  |
| Routing         | preact-iso                                      |
| Fonts           | Google Fonts (Xanh Mono, JetBrains Mono)        |
| Icons           | Inline SVG (no library)                         |
| Component Lib   | Custom (no external UI library)                 |

## 13: Design Discipline

With this many visual techniques available, restraint is critical. Not every page needs every treatment.

### Usage Budget

| Technique              | Frequency      | Where                                        |
| ---------------------- | -------------- | -------------------------------------------- |
| Sharp corners          | **Always**     | Every element, no exceptions                 |
| Extended rule lines    | **Often**      | Headers, section breaks, active states       |
| Wireframe textures     | **Sometimes**  | Empty states, dashboard bg, loading skeletons |
| Exposed internals      | **Sometimes**  | Status indicators, icons, progress bars      |
| Blend modes            | **Rarely**     | Data viz overlaps, hero moments, hover accents |
| Radar/spider charts    | **Rarely**     | Dashboard analytics, project overview        |

**Rule of thumb:** If a page already uses extended rule lines and wireframe textures, it does *not* also need blend modes. Pick the one or two techniques that serve the content best. The base UI (zinc + lime + sharp corners + JetBrains Mono) should carry 80% of the experience on its own.
