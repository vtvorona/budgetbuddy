/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{html,js}"],
  theme: {
    extend: {
      height: {
        'custom': 'calc(100svh - 64px - 64px)',
        'home': 'calc(100svh - 64px)',
      },
    },
  },
  plugins: [],
}

