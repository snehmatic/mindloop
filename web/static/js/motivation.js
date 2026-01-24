const quotes = [
    "The only way to do great work is to love what you do. – Steve Jobs",
    "It always seems impossible until it's done. – Nelson Mandela",
    "Don't watch the clock; do what it does. Keep going. – Sam Levenson",
    "The secret of getting ahead is getting started. – Mark Twain",
    "Quality is not an act, it is a habit. – Aristotle",
    "Focus on being productive instead of busy. – Tim Ferriss",
    "Simplicity is the ultimate sophistication. – Leonardo da Vinci",
    "What we fear of doing most is usually what we most need to do. – Ralph Waldo Emerson"
];

function showMotivation() {
    const quote = quotes[Math.floor(Math.random() * quotes.length)];
    const display = document.getElementById('motivation-display');
    if (display) {
        display.innerHTML = `<div class="card mb-md" style="border-left: 4px solid var(--primary); background: var(--primary-light); color: var(--primary-dark); font-weight: 500; font-style: italic; padding: 1rem;">"${quote}"</div>`;
        display.style.display = 'block';
    }
}

function toggleHistory() {
    const historySection = document.getElementById('history-section');
    const toggleBtn = document.getElementById('toggle-history-btn');
    if (historySection.style.display === 'none') {
        historySection.style.display = 'block';
        toggleBtn.textContent = 'Hide History';
    } else {
        historySection.style.display = 'none';
        toggleBtn.textContent = 'Show History';
    }
}
