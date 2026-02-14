#!/bin/bash

# Update style.css to fix alignment

# 1. Fix theme toggle button padding
# Find .btn-sm and update padding to be equal for better icon centering
sed -i '' 's/padding: 0.4rem 0.8rem;/padding: 0.4rem;/g' web/static/css/style.css

# 2. Ensure nav items are centered vertically
# .nav-link already has align-items: center, but let's make sure the icon inside doesn't have weird margins
# We will append a specific rule for the theme toggle button if needed, but adjusting btn-sm might be enough.
# Let's also check .nav-links alignment in style.css

echo "Updated style.css"
