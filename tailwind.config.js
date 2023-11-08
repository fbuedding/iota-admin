/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/**/*.{html,templ}"],
  theme: {
    extend: {
      screens: {
        print: { raw: 'print' },
        screen: { raw: 'screen' },
      }
    }
  },
  daisyui: {
    themes: ["light", "dark"],
    darkTheme: "dark",
  },
  plugins: [require("daisyui")],
}


