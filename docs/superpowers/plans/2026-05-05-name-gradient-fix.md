# Name Gradient Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve name gradient readability on the home page.

**Architecture:** Modify CSS variables in `style.css` to use a darker teal for the gradient end stop, ensuring contrast in both light and dark modes.

**Tech Stack:** Go, HTML Templates, Vanilla CSS.

---

### Task 1: Refine Name Gradient Style

**Files:**
- Modify: `web/static/css/style.css`

- [ ] **Step 1: Locate and update the `.text-gradient-name` class**

```css
/* Find this block in web/static/css/style.css */
.text-gradient-name {
    background: linear-gradient(135deg, var(--primary) 0%, var(--primary-light) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    color: var(--primary);
}

/* Change to: */
.text-gradient-name {
    background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    color: var(--primary);
}
```

- [ ] **Step 2: Verify the "Welcome back," text is correctly targeted**

In `web/templates/home.html`, ensure the greeting is outside the gradient span:
```html
Welcome back, <span class="text-gradient-name">{{ .UserName }}</span>
```

- [ ] **Step 3: Commit the changes**

```bash
git add web/static/css/style.css
git commit -m "style: refine name gradient for better readability"
```

---

### Task 2: Verification

- [ ] **Step 1: Build and run the server**

Run: `make run-server`
Expected: Server starts on port 8765.

- [ ] **Step 2: Visual check in browser**

Open: `http://localhost:8765`
Check:
- Name is clearly readable in Light Mode.
- Toggle to Dark Mode and verify name is still clearly readable.
- "Welcome back," should be high-contrast text.
