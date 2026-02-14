#!/bin/bash

# Restore missing CSS classes and fix regressions

cat <<EOT >> web/static/css/style.css

/* --- Restored Styles --- */

/* Progress Bars (Habits) */
.progress-container {
    background-color: var(--bg-surface-alt); /* Updated variable */
    border-radius: var(--radius-full);
    height: 0.6rem;
    overflow: hidden;
    margin-top: 0.75rem;
}

.progress-bar {
    background: var(--gradient-primary);
    height: 100%;
    border-radius: var(--radius-full);
    transition: width 0.5s ease-out;
}

/* About Hero (Visuals) */
.about-hero {
    text-align: center;
    padding: 4rem 1.5rem;
    margin-bottom: 3rem;
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--border);
    /* Restored background image */
    background-image: linear-gradient(rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.8)), url('/static/images/about-graded.png');
    background-size: cover;
    background-position: center;
}

/* Home Grid (Responsive) */
.home-grid {
    grid-template-columns: 2fr 1fr;
}

@media (max-width: 768px) {
    .home-grid {
        grid-template-columns: 1fr;
    }
}
EOT

echo "Restored CSS styles."
