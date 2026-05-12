# Design Spec: Improve Name Gradient Readability

## Problem Statement
The current name gradient in the "Welcome back" greeting on the home page uses a very light color (`var(--primary-light)`) for the end stop. In light mode, this makes the last characters of the user's name almost invisible against the light background. In dark mode, it can also lack sufficient contrast depending on the specific colors used.

## Proposed Solution
Refine the gradient to use a more saturated "dark" version of the primary color for the end stop, ensuring high contrast in both light and dark modes. Simultaneously, ensure the "Welcome back," text remains standard high-contrast text.

## Technical Details

### CSS Changes (`web/static/css/style.css`)
Modify the `.text-gradient-name` class to use `var(--primary-dark)` instead of `var(--primary-light)`.

**Current:**
```css
.text-gradient-name {
    background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
    /* ... other properties ... */
}
```

**New:**
```css
.text-gradient-name {
    background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
    /* ... other properties ... */
}
```

### Contrast Ratios (Estimated)
*   **Light Mode**: Teal-600 (`#0d9488`) to Teal-700 (`#0f766e`). Both provide excellent contrast against the light background (`#f0fdfa`).
*   **Dark Mode**: Teal-400 (`#2dd4bf`) to Teal-500 (`#14b8a6`). Both provide high visibility against the dark background (`#050b14`).

## Success Criteria
- The user's name is clearly readable in both light and dark modes.
- The gradient effect is still visible but subtle and professional.
- The "Welcome back," text is high contrast (black in light mode, white in dark mode).
