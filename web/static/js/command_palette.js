document.addEventListener('DOMContentLoaded', () => {
    const palette = document.getElementById('command-palette');
    if (!palette) return;
    const input = document.getElementById('cmd-input');
    const results = document.getElementById('cmd-results');
    
    const commands = [
        { name: 'Home', url: '/', icon: 'home' },
        { name: 'Intent', url: '/intent', icon: 'crosshair' },
        { name: 'Focus', url: '/focus', icon: 'zap' },
        { name: 'Tasks', url: '/tasks', icon: 'list-todo' },
        { name: 'Habits', url: '/habits', icon: 'check-square' },
        { name: 'Journal', url: '/journal', icon: 'book' },
        { name: 'Notes', url: '/notes', icon: 'file-text' },
        { name: 'Void', url: '/void', icon: 'ghost' },
        { name: 'Summary', url: '/summary', icon: 'bar-chart-2' },
        { name: 'About', url: '/about', icon: 'info' },
        { name: 'Settings', url: '/settings', icon: 'settings' },
        { name: 'Toggle Theme', action: 'toggleTheme', icon: 'moon' }
    ];

    let selectedIndex = 0;

    function openPalette() {
        palette.classList.add('active');
        input.value = '';
        renderResults('');
        setTimeout(() => input.focus(), 50);
    }

    function closePalette() {
        palette.classList.remove('active');
    }

    function executeCommand(cmd) {
        closePalette();
        if (cmd.url) {
            window.location.href = cmd.url;
        } else if (cmd.action === 'toggleTheme') {
            const toggle = document.getElementById('theme-toggle');
            if (toggle) toggle.click();
        }
    }

    function renderResults(query) {
        query = query.toLowerCase();
        const filtered = commands.filter(c => c.name.toLowerCase().includes(query));
        results.innerHTML = '';
        
        if (filtered.length === 0) {
            results.innerHTML = '<div class="cmd-empty">No matching commands</div>';
            return;
        }

        filtered.forEach((cmd, idx) => {
            const item = document.createElement('div');
            item.className = 'cmd-item' + (idx === selectedIndex ? ' selected' : '');
            item.innerHTML = `<i data-lucide="${cmd.icon}"></i> <span>${cmd.name}</span>`;
            item.addEventListener('click', () => executeCommand(cmd));
            item.addEventListener('mouseenter', () => {
                selectedIndex = idx;
                updateSelection();
            });
            results.appendChild(item);
        });
        
        if (window.lucide) {
            window.lucide.createIcons({ root: results });
        }
    }

    function updateSelection() {
        const items = results.querySelectorAll('.cmd-item');
        items.forEach((item, idx) => {
            if (idx === selectedIndex) {
                item.classList.add('selected');
                item.scrollIntoView({ block: 'nearest' });
            } else {
                item.classList.remove('selected');
            }
        });
    }

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
            e.preventDefault();
            if (palette.classList.contains('active')) {
                closePalette();
            } else {
                openPalette();
            }
        } else if (e.key === 'Escape' && palette.classList.contains('active')) {
            closePalette();
        }
    });

    input.addEventListener('keydown', (e) => {
        const query = input.value.toLowerCase();
        const filtered = commands.filter(c => c.name.toLowerCase().includes(query));
        
        if (e.key === 'ArrowDown') {
            e.preventDefault();
            selectedIndex = (selectedIndex + 1) % filtered.length;
            updateSelection();
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            selectedIndex = (selectedIndex - 1 + filtered.length) % filtered.length;
            updateSelection();
        } else if (e.key === 'Enter') {
            e.preventDefault();
            if (filtered.length > 0) {
                executeCommand(filtered[selectedIndex]);
            }
        }
    });

    input.addEventListener('input', (e) => {
        selectedIndex = 0;
        renderResults(e.target.value);
    });

    palette.addEventListener('click', (e) => {
        if (e.target === palette) {
            closePalette();
        }
    });
});
