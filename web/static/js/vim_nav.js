document.addEventListener('DOMContentLoaded', () => {
    const itemSelectors = [
        '.task-item',           
        '#intent-tasks-list .card.p-sm', 
        '.subtask-item',        
        '.habit-card',          
        '#intent-history-list > div[data-title]' 
    ];

    let items = [];
    let currentIndex = -1;

    function getItems() {
        const allItems = Array.from(document.querySelectorAll(itemSelectors.join(', ')));
        return allItems.filter(el => window.getComputedStyle(el).display !== 'none');
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
        const target = e.target;
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) return;
        if (e.metaKey || e.ctrlKey || e.altKey) return;

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
                const completeBtnSelectors = [
                    'form[action*="/complete"] button[type="submit"]', 
                    'button[title="Mark Complete"]', 
                    'form[action*="/habits/log"] button[type="submit"]', 
                    '.btn-success-outline' 
                ];
                
                for (let selector of completeBtnSelectors) {
                    const btn = item.querySelector(selector);
                    if (btn) {
                        btn.click();
                        setTimeout(() => {
                            items = getItems();
                            if (currentIndex >= items.length) currentIndex = Math.max(items.length - 1, 0);
                            updateHighlight();
                        }, 100);
                        break;
                    }
                }
            }
        }
    });
});
