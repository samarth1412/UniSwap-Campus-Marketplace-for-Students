/**
 * Display-only currency helpers.
 * Backend stores listing prices in USD.
 */
function safeNumber(value: number): number {
  return Number.isFinite(value) ? value : 0;
}

export function formatUsd(priceUsd: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2,
  }).format(safeNumber(priceUsd));
}
