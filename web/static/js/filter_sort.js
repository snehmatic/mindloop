/**
 * Mindloop Client-Side Filtering & Sorting Controller
 * Automatically handles search, drop-down filtering, and sorting in the DOM.
 */

(function () {
    'use strict';

    console.log('[FilterSort] Script loaded and executing...');

    // Debounce helper
    function debounce(fn, delay) {
        let timeout;
        return function (...args) {
            clearTimeout(timeout);
            timeout = setTimeout(() => fn.apply(this, args), delay);
        };
    }

    // Initialize all filter-sort bars on the page
    function initAllFilterSortBars() {
        console.log('[FilterSort] initAllFilterSortBars called');
        const bars = document.querySelectorAll('.filter-sort-bar');
        console.log('[FilterSort] Found filter-sort-bars:', bars.length);
        bars.forEach(bar => {
            if (bar.dataset.initialized) {
                console.log('[FilterSort] Bar already initialized:', bar.id || bar.dataset.target);
                return;
            }
            setupFilterSortBar(bar);
            bar.dataset.initialized = 'true';
        });
    }

    function setupFilterSortBar(bar) {
        const targetSelector = bar.dataset.target;
        console.log('[FilterSort] Setting up filter-sort-bar targeting:', targetSelector);
        if (!targetSelector) return;

        const container = document.querySelector(targetSelector);
        if (!container) {
            console.warn('[FilterSort] Target container not found for selector:', targetSelector);
            return;
        }

        const searchInput = bar.querySelector('.search-input');
        const filterSelects = bar.querySelectorAll('.filter-select');
        const sortSelect = bar.querySelector('.sort-select');

        console.log('[FilterSort] Found elements for setup:', {
            hasSearchInput: !!searchInput,
            filterSelectsCount: filterSelects.length,
            hasSortSelect: !!sortSelect
        });

        // Store original DOM elements on the bar element to avoid stale closure references
        bar._originalElements = Array.from(container.children).filter(el => {
            return !el.classList.contains('empty-state') && el.id !== 'filter-empty-state';
        });
        console.log('[FilterSort] Initial original items count:', bar._originalElements.length);

        // Dynamically populate filter select choices from data attributes
        filterSelects.forEach(select => {
            const key = select.dataset.filterKey;
            if (select.dataset.dynamic === 'true' && key) {
                const uniqueValues = new Set();
                bar._originalElements.forEach(el => {
                    const val = el.dataset[key];
                    if (val) {
                        val.split(',').map(v => v.trim()).forEach(v => {
                            if (v) uniqueValues.add(v);
                        });
                    }
                });
                
                const firstOption = select.options[0];
                select.innerHTML = '';
                if (firstOption) select.appendChild(firstOption);
                
                Array.from(uniqueValues).sort().forEach(val => {
                    const opt = document.createElement('option');
                    opt.value = val;
                    opt.textContent = val;
                    select.appendChild(opt);
                });
                console.log(`[FilterSort] Dynamically populated filter key "${key}" with values:`, Array.from(uniqueValues));
            }
        });

        // Search handler
        if (searchInput) {
            console.log('[FilterSort] Registering input listener on search input');
            searchInput.addEventListener('input', debounce(() => {
                console.log('[FilterSort] Search input changed to:', searchInput.value);
                applyFiltersAndSort(container, bar);
            }, 150));
        }

        // Dropdown filters handler
        filterSelects.forEach(select => {
            select.addEventListener('change', () => {
                console.log(`[FilterSort] Filter select "${select.dataset.filterKey}" changed to:`, select.value);
                applyFiltersAndSort(container, bar);
            });
        });

        // Segmented controls handler (used in tasks.html)
        const segmentedBtns = bar.querySelectorAll('.segmented-btn');
        const barId = bar.id || 'default';
        let savedSegment = null;
        try {
            savedSegment = sessionStorage.getItem(barId + '-active-segment');
        } catch (e) {
            console.warn('[FilterSort] Failed to access sessionStorage:', e);
        }
        savedSegment = savedSegment || bar.dataset.activeSegment;
        
        if (savedSegment && barId !== 'default') {
            console.log(`[FilterSort] Restoring active segment "${savedSegment}" from storage`);
            bar.dataset.activeSegment = savedSegment;
            segmentedBtns.forEach(btn => {
                if (btn.dataset.filter === savedSegment) {
                    btn.classList.add('active');
                } else {
                    btn.classList.remove('active');
                }
            });
        }

        segmentedBtns.forEach(btn => {
            btn.addEventListener('click', (e) => {
                segmentedBtns.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');
                bar.dataset.activeSegment = btn.dataset.filter;
                console.log('[FilterSort] Segment clicked:', btn.dataset.filter);
                if (barId !== 'default') {
                    try {
                        sessionStorage.setItem(barId + '-active-segment', btn.dataset.filter);
                    } catch (e) {
                        console.warn('[FilterSort] Failed to write to sessionStorage:', e);
                    }
                }
                applyFiltersAndSort(container, bar);
            });
        });

        // Sort handler
        if (sortSelect) {
            sortSelect.addEventListener('change', () => {
                console.log('[FilterSort] Sort option changed to:', sortSelect.value);
                applyFiltersAndSort(container, bar);
            });
        }

        // Initial run
        applyFiltersAndSort(container, bar);
    }

    function applyFiltersAndSort(container, bar) {
        console.log('[FilterSort] applyFiltersAndSort running on container:', container.id || container.className);
        const searchInput = bar.querySelector('.search-input');
        const filterSelects = bar.querySelectorAll('.filter-select');
        const sortSelect = bar.querySelector('.sort-select');

        const searchQuery = searchInput ? searchInput.value.trim().toLowerCase() : '';
        const segmentFilter = bar.dataset.activeSegment || '';
        const originalElements = bar._originalElements || [];

        // Collect all filter selections
        const activeFilters = [];
        filterSelects.forEach(select => {
            const key = select.dataset.filterKey;
            const val = select.value;
            if (key && val) {
                activeFilters.push({ key, val });
            }
        });

        console.log('[FilterSort] Applying filters with query:', searchQuery, 'segment:', segmentFilter, 'selects:', activeFilters);

        // Collect matching items currently in container
        const items = Array.from(container.children).filter(el => {
            return !el.classList.contains('empty-state') && el.id !== 'filter-empty-state';
        });

        let visibleCount = 0;

        items.forEach(item => {
            let matchesSearch = true;
            let matchesSegment = true;
            let matchesSelects = true;

            // 1. Search filter
            if (searchQuery) {
                const searchableText = item.dataset.searchable || item.textContent;
                matchesSearch = searchableText.toLowerCase().includes(searchQuery);
            }

            // 2. Segmented filter (e.g. To-Do, Completed, All status filter on Tasks)
            if (segmentFilter && segmentFilter !== 'all') {
                const status = item.dataset.status || '';
                matchesSegment = (status === segmentFilter);
            }

            // 3. Dropdown filters
            for (let filter of activeFilters) {
                const itemVal = item.dataset[filter.key];
                
                // Custom checking for duration ranges
                if (filter.key === 'duration') {
                    if (!itemVal) {
                        matchesSelects = false;
                        break;
                    }
                    const mins = parseFloat(itemVal);
                    if (isNaN(mins)) {
                        matchesSelects = false;
                        break;
                    }
                    if (filter.val === 'short' && mins >= 25) {
                        matchesSelects = false;
                        break;
                    }
                    if (filter.val === 'pomodoro' && (mins < 25 || mins > 50)) {
                        matchesSelects = false;
                        break;
                    }
                    if (filter.val === 'long' && mins <= 50) {
                        matchesSelects = false;
                        break;
                    }
                    continue; // Skip standard equality checks
                }
                
                if (filter.val === 'standalone') {
                    if (itemVal && itemVal.trim() !== '') {
                        matchesSelects = false;
                        break;
                    }
                    continue; // Standalone matched successfully, skip standard checks
                } else if (!itemVal) {
                    matchesSelects = false;
                    break;
                }
                
                // Support multiple tags check if dataset value is comma-separated
                if (filter.key === 'labels' || (itemVal && itemVal.includes(','))) {
                    const tags = itemVal.split(',').map(t => t.trim().toLowerCase());
                    if (!tags.includes(filter.val.toLowerCase())) {
                        matchesSelects = false;
                        break;
                    }
                } else {
                    if (itemVal.toLowerCase() !== filter.val.toLowerCase()) {
                        matchesSelects = false;
                        break;
                    }
                }
            }

            // Show or hide based on all filters matching
            if (matchesSearch && matchesSegment && matchesSelects) {
                item.classList.remove('filter-item-hidden');
                visibleCount++;
            } else {
                item.classList.add('filter-item-hidden');
            }
        });

        console.log(`[FilterSort] Filter pass finished: ${visibleCount}/${items.length} items visible`);

        // 4. Sort logic
        if (sortSelect && sortSelect.value) {
            const sortVal = sortSelect.value;
            console.log('[FilterSort] Sorting items by:', sortVal);
            
            // Toggle drag-and-drop indicator class
            if (sortVal !== 'custom') {
                container.classList.add('filter-sort-bar-active-sort');
            } else {
                container.classList.remove('filter-sort-bar-active-sort');
            }

            if (sortVal === 'custom') {
                // Restore original order (append them back in their original sequence)
                originalElements.forEach(el => {
                    if (el.parentNode === container) {
                        container.appendChild(el);
                    }
                });
            } else {
                const [sortKey, sortDir] = sortVal.split('-');
                const sortedItems = items.sort((a, b) => {
                    let valA = a.dataset[sortKey] || '';
                    let valB = b.dataset[sortKey] || '';

                    // Try parsing as strict number first to avoid parsing year from ISO date strings
                    const isNumA = valA.trim() !== '' && !isNaN(valA);
                    const isNumB = valB.trim() !== '' && !isNaN(valB);
                    if (isNumA && isNumB) {
                        const numA = parseFloat(valA);
                        const numB = parseFloat(valB);
                        return sortDir === 'asc' ? numA - numB : numB - numA;
                    }

                    // Try parsing as dates
                    const dateA = Date.parse(valA);
                    const dateB = Date.parse(valB);
                    if (!isNaN(dateA) && !isNaN(dateB)) {
                        return sortDir === 'asc' ? dateA - dateB : dateB - dateA;
                    }

                    // Standard string comparison
                    valA = valA.toLowerCase();
                    valB = valB.toLowerCase();
                    if (valA < valB) return sortDir === 'asc' ? -1 : 1;
                    if (valA > valB) return sortDir === 'asc' ? 1 : -1;
                    return 0;
                });

                // Re-append items in sorted order
                sortedItems.forEach(el => container.appendChild(el));
            }
        }

        // Show/hide empty state
        const emptyState = container.querySelector('#filter-empty-state');
        if (emptyState) {
            if (visibleCount === 0 && items.length > 0) {
                console.log('[FilterSort] Showing empty state: zero visible items matching search/filters');
                emptyState.style.display = 'flex';
                
                // Configure empty state labels
                const titleEl = emptyState.querySelector('#filter-empty-title') || emptyState.querySelector('h3');
                const descEl = emptyState.querySelector('#filter-empty-desc') || emptyState.querySelector('p');
                const iconEl = emptyState.querySelector('#filter-empty-icon') || emptyState.querySelector('.empty-icon');
                
                if (titleEl) titleEl.textContent = 'No matching items found';
                if (descEl) descEl.textContent = 'Try adjusting your search query or filter options.';
                if (iconEl) {
                    iconEl.innerHTML = '<i data-lucide="search-x" style="width: 48px; height: 48px; color: var(--text-light);"></i>';
                    if (window.lucide) window.lucide.createIcons();
                }
            } else {
                emptyState.style.display = 'none';
            }
        }
    }

    // Watch for HTMX swaps to re-initialize and apply filters
    document.addEventListener('htmx:afterSwap', function (evt) {
        const targetContainer = evt.detail.target;
        console.log('[FilterSort] HTMX swap detected on target:', targetContainer);
        
        // Find if target container is inside an active filter-sort bar target
        const activeBars = document.querySelectorAll('.filter-sort-bar');
        activeBars.forEach(bar => {
            const targetSelector = bar.dataset.target;
            const container = document.querySelector(targetSelector);
            if (container && (container === targetContainer || container.contains(targetContainer))) {
                // Re-read original element list
                const originalElements = Array.from(container.children).filter(el => {
                    return !el.classList.contains('empty-state') && el.id !== 'filter-empty-state';
                });
                bar._originalElements = originalElements;
                console.log('[FilterSort] HTMX swap updated originalElements count to:', originalElements.length);
                applyFiltersAndSort(container, bar);
            }
        });
        
        // Re-init any new filter bars
        initAllFilterSortBars();
    });

    // Page Load Initializer
    document.addEventListener('DOMContentLoaded', () => {
        console.log('[FilterSort] DOMContentLoaded fired, triggering initialization');
        initAllFilterSortBars();
    });
})();
