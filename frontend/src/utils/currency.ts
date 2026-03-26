/**
 * Display-only currency helpers.
 * Backend stores listing prices in INR; the UI shows USD.
 */
const INR_TO_USD = 83;

function safeNumber(value: number): number {
  return Number.isFinite(value) ? value : 0;
}

export function convertINRtoUSD(priceInr: number): number {
  return safeNumber(priceInr) / INR_TO_USD;
}

export function convertUSDToINR(priceUsd: number): number {
  return safeNumber(priceUsd) * INR_TO_USD;
}

export function formatUsdFromInr(inr: number): string {
  const usd = convertINRtoUSD(inr);
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2,
  }).format(usd);
}
