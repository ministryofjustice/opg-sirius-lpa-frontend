export default function donorHeader() {
    const btn = document.getElementById('panel-calendar-btn');
    const panel = document.getElementById('panel-calendar');
    if (!btn || !panel) return;

    btn.addEventListener('click', () => {
        const isOpen = btn.getAttribute('aria-expanded') === 'true';
        btn.setAttribute('aria-expanded', String(!isOpen));
        panel.hidden = isOpen;
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && !panel.hidden) {
            panel.hidden = true;
            btn.setAttribute('aria-expanded', 'false');
            btn.focus();
        }
    });
}