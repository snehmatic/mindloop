document.addEventListener("DOMContentLoaded", () => {
    // Select all task and habit forms
    const getItems = () => Array.from(document.querySelectorAll('form[action="/tasks/complete"], form[action="/habits/log"]'));
    let items = getItems();
    if (items.length === 0) return;

    let currentIndex = -1;

    // Restore index from sessionStorage if available
    const savedIndex = sessionStorage.getItem('vimNavIndex');
    if (savedIndex !== null) {
        currentIndex = parseInt(savedIndex, 10);
        // Ensure within bounds
        if (currentIndex >= items.length) {
            currentIndex = items.length - 1;
        }
        if (currentIndex >= 0) {
            items[currentIndex].classList.add('vim-focus');
            items[currentIndex].scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }

    const updateHighlight = (newIndex) => {
        if (currentIndex >= 0 && currentIndex < items.length) {
            items[currentIndex].classList.remove('vim-focus');
        }
        currentIndex = newIndex;
        if (currentIndex >= 0 && currentIndex < items.length) {
            items[currentIndex].classList.add('vim-focus');
            items[currentIndex].scrollIntoView({ behavior: 'smooth', block: 'center' });
            sessionStorage.setItem('vimNavIndex', currentIndex);
        }
    };

    document.addEventListener("keydown", (e) => {
        // Ignore if focus is in an input field
        const activeTag = document.activeElement.tagName.toLowerCase();
        if (['input', 'textarea', 'select'].includes(activeTag)) return;
        
        // Also ignore if Cmd/Ctrl is pressed
        if (e.ctrlKey || e.metaKey || e.altKey) return;

        items = getItems(); // Refresh items in case DOM changed (htmx)
        if (items.length === 0) return;

        if (e.key === 'j') {
            const next = Math.min(currentIndex + 1, items.length - 1);
            updateHighlight(next);
        } else if (e.key === 'k') {
            const prev = Math.max(currentIndex - 1, 0);
            updateHighlight(prev);
        } else if (e.key === 'x') {
            if (currentIndex >= 0 && currentIndex < items.length) {
                const item = items[currentIndex];
                const btn = item.querySelector('input[type="checkbox"], button[type="submit"]');
                if (btn) {
                    btn.click();
                    // Don't update index here; page will reload and restore it
                }
            }
        } else if (e.key === 'Escape') {
            if (currentIndex >= 0 && currentIndex < items.length) {
                items[currentIndex].classList.remove('vim-focus');
            }
            currentIndex = -1;
            sessionStorage.removeItem('vimNavIndex');
        }
    });
});
