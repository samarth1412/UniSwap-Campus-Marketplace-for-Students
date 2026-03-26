import type { HTMLAttributes } from 'react';

type UniSwapLogoProps = HTMLAttributes<SVGSVGElement>;

export function UniSwapLogo({ className, ...rest }: UniSwapLogoProps) {
  return (
    <svg
      viewBox="0 0 64 64"
      role="img"
      aria-label="UniSwap"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      {...rest}
    >
      <defs>
        <linearGradient id="uniswapp-grad" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stopColor="#4f46e5" />
          <stop offset="1" stopColor="#e11d48" />
        </linearGradient>
      </defs>
      <circle cx="32" cy="32" r="28" fill="url(#uniswapp-grad)" opacity="0.14" />
      <path
        d="M17 31.2c0-10.6 8.2-19.2 18.3-19.2 7.4 0 13.8 4.3 16.9 10.6L46 39.6l-9.2 0L23.7 31.2z"
        fill="none"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinejoin="round"
      />
      <path
        d="M47 32.8c0 10.6-8.2 19.2-18.3 19.2-7.4 0-13.8-4.3-16.9-10.6L18 24.4h9.2l13.1 8.4z"
        fill="none"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinejoin="round"
        opacity="0.9"
      />
    </svg>
  );
}

/** Simple swap-arrows mark for the navbar (~28px). */
export function UniSwapLogo({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      width={28}
      height={28}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden
    >
      <path
        d="M7 16V4M7 4L3 8M7 4L11 8"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M17 8V20M17 20L21 16M17 20L13 16"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
