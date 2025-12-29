import { escapeHtml } from './escapeHtml';

export function highlightText(
    text: string,
    kind: 'info' | 'success' | 'error' = 'info'
) {
    const colour =
        kind === 'success'
            ? 'text-success'
            : kind === 'error'
            ? 'text-error'
            : 'text-info';
    return `<span class="${colour} font-semibold">${escapeHtml(text)}</span>`;
}
