/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');

module.exports = {
  // These paths are just examples, customize them to match your project structure
  purge: [
    './public/**/*.html',
    './src/**/*.{vue,css}',
  ],
  content: [
    './public/**/*.html', 
    './src/**/*.{vue,css}'
  ],
  theme: {
    screens: {
      sm: '480px',
      md: '768px',
      lg: '976px',
      xl: '1440px',
    },
    colors: {
      gray: colors.gray,
      blue: colors.sky,
      red: colors.rose,
      pink: colors.fuchsia,
      teal: colors.teal,
      white: colors.white,
      black: colors.black,
      slate: colors.slate
    },
    fontFamily: {
      sans: ['sans-serif'],
      serif: ['serif'],
    },
    extend: {
      spacing: {
        '128': '32rem',
        '144': '36rem',
      },
      borderRadius: {
        '4xl': '2rem',
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'), 
    require('@tailwindcss/typography'),
  ]
}
