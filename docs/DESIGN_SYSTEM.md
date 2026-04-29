# 🎨 GoondVR Design System

## Color System

### Primary Gradients
```css
--gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%)
  Usage: Primary buttons, active states, header text
  
--gradient-success: linear-gradient(135deg, #11998e 0%, #38ef7d 100%)
  Usage: Recording badges, success notifications
  
--gradient-danger: linear-gradient(135deg, #ee0979 0%, #ff6a00 100%)
  Usage: Delete buttons, error states, paused badges
  
--gradient-warning: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)
  Usage: Warning states, alerts
```

### Background
```css
Body: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%)
Glass: rgba(255, 255, 255, 0.05) with backdrop-blur(10px)
```

### Borders & Overlays
```css
--glass-border: rgba(255, 255, 255, 0.1)
--glass-bg: rgba(255, 255, 255, 0.05)
```

---

## Typography

### Font Stack
```css
font-family: 'Noto Sans TC', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif
```

### Sizes
- **Huge Header**: 3rem (48px) - Main title
- **Large Header**: 1.5rem (24px) - Section headers
- **Body**: 1rem (16px) - Regular text
- **Small**: 0.875rem (14px) - Meta information
- **Tiny**: 0.75rem (12px) - Labels, badges

### Weights
- **Heavy**: 700 - Headers, important text
- **Bold**: 600 - Buttons, labels
- **Medium**: 500 - Body text
- **Regular**: 400 - Secondary text

---

## Spacing System

### Scale (rem)
```
0.25rem = 4px   - Tiny gaps
0.5rem  = 8px   - Small gaps
0.75rem = 12px  - Medium gaps
1rem    = 16px  - Standard gap
1.5rem  = 24px  - Large gap
2rem    = 32px  - Extra large gap
3rem    = 48px  - Huge gap
```

### Component Padding
- **Container**: 2rem (32px)
- **Card**: 1rem (16px)
- **Button**: 0.5rem 1rem (8px 16px)
- **Input**: 0.75rem (12px)

---

## Border Radius

```css
Small:       8px  - Badges, small buttons
Medium:      12px - Inputs, thumbnails
Large:       16px - Cards, containers
Extra Large: 24px - Modals, main container
Circle:      50%  - Icons, dots
```

---

## Shadows

### Elevation System
```css
--shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.15)
  Usage: Buttons, small cards

--shadow-md: 0 4px 16px rgba(0, 0, 0, 0.2)
  Usage: Cards, dropdowns

--shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.3)
  Usage: Modals, elevated containers

--shadow-glow: 0 0 20px rgba(102, 126, 234, 0.3)
  Usage: Active states, focus rings
```

---

## Animation System

### Timing Functions
```css
Standard: cubic-bezier(0.4, 0, 0.2, 1)  - Most transitions
Ease-out: cubic-bezier(0, 0, 0.2, 1)    - Entrances
Ease-in:  cubic-bezier(0.4, 0, 1, 1)    - Exits
```

### Durations
```css
Fast:     0.15s - Hover states
Standard: 0.3s  - Most transitions
Slow:     0.5s  - Page loads, fades
Very Slow: 2-3s - Ambient animations
```

### Keyframe Animations

#### Glow (Header)
```css
@keyframes glow {
  from { filter: drop-shadow(0 0 10px rgba(102, 126, 234, 0.3)); }
  to   { filter: drop-shadow(0 0 20px rgba(118, 75, 162, 0.5)); }
}
Duration: 3s, infinite alternate
```

#### Pulse (Badges)
```css
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50%      { opacity: 0.8; }
}
Duration: 2s, infinite
```

#### Blink (Recording Dot)
```css
@keyframes blink {
  0%, 100% { opacity: 1; }
  50%      { opacity: 0.5; }
}
Duration: 2s, infinite
```

#### Shimmer (Loading)
```css
@keyframes shimmer {
  0%   { background-position: -1000px 0; }
  100% { background-position: 1000px 0; }
}
Duration: 2s, infinite
```

#### Slide In (Toast)
```css
@keyframes slideIn {
  from { transform: translateX(400px); opacity: 0; }
  to   { transform: translateX(0); opacity: 1; }
}
Duration: 0.3s
```

---

## Component Patterns

### Glassmorphism Card
```css
background: rgba(255, 255, 255, 0.05);
backdrop-filter: blur(10px);
border: 1px solid rgba(255, 255, 255, 0.1);
border-radius: 16px;
box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
```

### Gradient Button
```css
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
border: none;
border-radius: 12px;
box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

&:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}
```

### Focus Ring
```css
&:focus {
  outline: none;
  border-color: rgba(102, 126, 234, 0.5);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}
```

### Hover Lift
```css
transition: transform 0.3s ease, box-shadow 0.3s ease;

&:hover {
  transform: translateY(-8px);
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}
```

---

## Interactive States

### Button States
```
Default:  Normal appearance
Hover:    Lift + enhanced shadow
Active:   Press down + ripple effect
Disabled: Reduced opacity + no pointer
Loading:  Text change + disabled state
```

### Card States
```
Default:  Normal appearance
Hover:    Lift + glow + border highlight
Active:   Pressed state
Selected: Persistent glow
```

### Input States
```
Default:  Subtle background
Focus:    Glow ring + brighter background
Error:    Red border + error message
Success:  Green border + success icon
Disabled: Reduced opacity
```

---

## Responsive Breakpoints

```css
Mobile:  max-width: 768px
Tablet:  769px - 1024px
Desktop: 1025px+
```

### Mobile Adjustments
- Reduced padding (1rem instead of 2rem)
- Single column grid
- Smaller header (2rem instead of 3rem)
- Simplified animations
- Touch-friendly targets (min 44px)

---

## Accessibility

### Contrast Ratios
- **Body text**: 4.5:1 minimum
- **Large text**: 3:1 minimum
- **Interactive elements**: Clear focus indicators

### Focus Management
- Visible focus rings on all interactive elements
- Logical tab order
- Skip links for navigation

### Motion
- Respects `prefers-reduced-motion`
- Essential animations only
- No flashing content

---

## Icon System

### Icon Library
Using Tocas UI icon set with semantic names:
- `is-plus-icon` - Add actions
- `is-gear-icon` - Settings
- `is-trash-icon` - Delete
- `is-play-icon` - Resume
- `is-pause-icon` - Pause
- `is-eye-icon` - Viewers
- `is-clock-icon` - Time
- `is-hard-drive-icon` - Storage
- `is-tower-broadcast-icon` - Live

### Icon Sizes
- Small: 16px
- Medium: 20px
- Large: 24px

---

## Grid System

### Channel Grid
```css
display: grid;
grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
gap: 1.5rem;
```

### Form Grid
```css
display: grid;
grid-template-columns: repeat(2, 1fr);
gap: 1rem;
```

---

## Best Practices

### Performance
1. Use `transform` and `opacity` for animations
2. Enable hardware acceleration with `will-change`
3. Minimize repaints and reflows
4. Use CSS containment where appropriate

### Maintainability
1. Use CSS variables for theming
2. Follow BEM naming convention
3. Keep specificity low
4. Document complex selectors

### Consistency
1. Use design tokens from this system
2. Maintain consistent spacing
3. Follow animation patterns
4. Use semantic color names

---

## Component Library

### Buttons
- Primary: Gradient background
- Secondary: Outlined with hover fill
- Negative: Red gradient
- Icon: Icon-only compact

### Cards
- Channel card: Full details
- Grid card: Compact view
- Collapsed card: Single line

### Badges
- Recording: Green gradient
- Offline: Gray
- Paused: Red gradient

### Inputs
- Text input: Single line
- Textarea: Multi-line
- Select: Dropdown
- Checkbox: Toggle
- Radio: Single choice

### Modals
- Settings: Large form
- Create: Channel creation
- Confirm: Destructive actions

---

**Design System Version 1.0**
*Last Updated: 2026-04-29*
