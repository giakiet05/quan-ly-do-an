/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{svelte,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'var(--primary-color)',
          hover: 'var(--primary-color-hover)',
        },
        text: {
          DEFAULT: 'var(--text-color)',
          secondary: 'var(--text-color-secondary)',
        },
        background: 'var(--background)',
        surface: 'var(--surface-color)',
        border: 'var(--border-color)',
      },
      fontFamily: {
        primary: ['Lexend', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
