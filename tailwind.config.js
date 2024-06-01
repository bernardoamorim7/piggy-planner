/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./cmd/web/**/*.templ'],
  theme: {
    extend: {
      // backgroundImage: theme => ({
      //   'piggy-pattern': "url('/assets/imgs/favicon.png')",
      // })
    },
  },
  plugins: [
    require('daisyui')
  ],
  daisyui: {
    themes: ['nord'],
  },
}

