/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./web/**/*.templ'],
  plugins: [
    require('daisyui')
  ],
  daisyui: {
    themes: ['nord'],
  },
}
