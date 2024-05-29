/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./cmd/web/**/*.templ'],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui')
  ],
  daisyui: {
    themes: ['nord'],
  },
}

