document.addEventListener('DOMContentLoaded', () => {
    // Selectors for navigable items across different pages
    const itemSelectors = [
        '.task-item',           // Tasks page
        '#intent-tasks-list .card.p-sm', // Intent inline tasks
        '.subtask-item',        // Subtasks
        '.habit-card',          // Habit cards
        '#intent-history-list > div[data-title]' // Intent history
    ];

    let items = [];
    let currentIndex = -1;

    // Refresh the list of navigable items based on what's visible
    function getItems() {
        const allItems = Array.from(document.querySelectorAll(itemSelectors.join(', ')));
        // Filter out items that are display: none
        return allItems.filter(el => {
            return window.getComputedStyle(el).display !== 'none';
        });
    }

    function updateHighlight() {
        items.forEach((item, index) => {
            if (index === currentIndex) {
                item.classList.add('vim-focused');
                item.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            } else {
                item.classList.remove('vim-focused');
            }
        });
    }

    document.addEventListener('keydown', (e) => {
        // Disable when typing in input, textarea, or contenteditable
        const target = e.target;
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
            return;
        }

        // Only handle j, k, x without modifiers (except shift for some keys if we want, but typically lowercase)
        if (e.metaKey || e.ctrlKey || e.altKey) {
            return;
        }

        const key = e.key.toLowerCase();
        
        if (key === 'j' || key === 'k' || key === 'x') {
            items = getItems();
            if (items.length === 0) return;
        }

        if (key === 'j') {
            currentIndex = Math.min(currentIndex + 1, items.length - 1);
            updateHighlight();
            e.preventDefault();
        } else if (key === 'k') {
            currentIndex = Math.max(currentIndex - 1, 0);
            updateHighlight();
            e.preventDefault();
        } else if (key === 'x') {
            if (currentIndex >= 0 && currentIndex < items.length) {
                const item = items[currentIndex];
                
                // Try to find a completion button inside the item
                // Target generic success buttons or specific forms
                const completeBtnSelectors = [
                    'form[action*="/complete"] button[type="submit"]', // Standard task complete
                    'button[title="Mark Complete"]', // Subtask complete
                    'form[action*="/habits/log"] button[type="submit"]', // Habit log
                    '.btn-success-outline' // Generic success outline button
                ];
                
                for (let selector of completeBtnSelectors) {
                    const btn = item.querySelector(selector);
                    if (btn) {
                        btn.click();
                        // Re-fetch items after a short delay to adjust selection if item disappears/changes
                        setTimeout(() => {
                            items = getItems();
                            if (currentIndex >= items.length) {
                                currentIndex = Math.max(items.length - 1, 0);
                            }
                            updateHighlight();
                        }, 100);
                        break;
                    }
                }
            }
        }
    });
});
