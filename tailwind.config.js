/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./cmd/web/**/*.templ'],
  plugins: [
    require('daisyui')
  ],
  daisyui: {
    themes: ['nord'],
  },
}
