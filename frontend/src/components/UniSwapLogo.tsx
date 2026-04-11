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
        <linearGradient id="uniswap-grad" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stopColor="#4f46e5" />
          <stop offset="1" stopColor="#e11d48" />
        </linearGradient>
      </defs>
      <circle cx="32" cy="32" r="28" fill="url(#uniswap-grad)" opacity="0.14" />
      <path
        d="M17 31.2c0-10.6 8.2-19.2 18.3-19.2 7.4 0 13.8 4.3 16.9 10.6L46 39.6h-9.2l-13.1-8.4z"
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
