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

async function showMotivation() {
    const display = document.getElementById('motivation-display');
    if (!display) return;

    display.style.display = 'block';
    display.innerHTML = `<div class="text-center text-muted">Thinking...</div>`;

    try {
        const response = await fetch('/api/quote');
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();

        // Check for rate limit wait time
        let notice = '';
        if (data.retry_after && data.retry_after > 0) {
            notice = `
                <div style="font-size: 0.8rem; color: var(--text-muted); background-color: var(--secondary-light); border-top: 1px solid var(--border); padding: 0.5rem; margin-top: 1rem; border-radius: 4px;">
                    <span style="font-size: 1.2em; vertical-align: middle;">⏳</span> 
                    Slow down! New quote in <span id="retry-countdown">${data.retry_after}</span>s
                </div>`;

            // Start countdown
            startCountdown(data.retry_after);
        }

        display.innerHTML = `
            <div class="card mb-md" style="background: var(--bg-surface); backdrop-filter: blur(12px); border: 1px solid var(--border); color: var(--text-main); padding: 1.5rem; position: relative; overflow: hidden;">
                 <div style="position: absolute; top: 0; left: 0; width: 4px; height: 100%; background: var(--gradient-primary);"></div>
                <div style="font-weight: 500; font-style: italic; font-size: 1.1rem; margin-bottom: 0.75rem; line-height: 1.6;">"${data.q}"</div>
                <div style="text-align: right; font-size: 0.9rem; color: var(--text-muted); font-weight: 600;">— ${data.a}</div>
                ${notice}
            </div>`;
    } catch (error) {
        // Fallback to local quotes
        const quote = quotes[Math.floor(Math.random() * quotes.length)];
        display.innerHTML = `
            <div class="card mb-md" style="border-left: 4px solid var(--secondary); background: var(--bg-surface); color: var(--text-main); padding: 1rem;">
                <div style="font-weight: 500; font-style: italic; margin-bottom: 0.5rem;">"${quote}"</div>
                <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 1rem; border-top: 1px solid var(--border); padding-top: 0.5rem;">
                    <span style="font-size: 1.2em; vertical-align: middle;">📜</span> Showing a classic quote.
                </div>
            </div>`;
    }
}

let countdownInterval;

function startCountdown(seconds) {
    if (countdownInterval) clearInterval(countdownInterval);

    let left = seconds;
    countdownInterval = setInterval(() => {
        left--;
        const el = document.getElementById('retry-countdown');
        if (el) {
            el.textContent = left;
        }

        if (left <= 0) {
            clearInterval(countdownInterval);
            if (el && el.parentElement) {
                el.parentElement.innerHTML = '<span style="font-size: 1.2em; vertical-align: middle;">✅</span> Ready for a new quote!';
                el.parentElement.style.color = 'var(--success)';
                el.parentElement.style.backgroundColor = 'var(--bg-surface)';
                el.parentElement.style.borderColor = 'var(--success)';
            }
        }
    }, 1000);
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
