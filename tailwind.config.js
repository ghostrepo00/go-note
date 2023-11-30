/** @type {import('tailwindcss').Config} */
module.exports = {
  //content: ["./internal/handler/web/template/**/*.{html,js}"],
  content: ["./web/template/**/*.html"],
  theme: {    
    colors: {
      'base': {
        '50': '#f6f6f9',
        '100': '#ebebf3',
        '200': '#d3d4e4',
        '300': '#adafcc',
        '400': '#8084b0',
        '500': '#606597',
        '600': '#4c4f7d',
        '700': '#3e4066',
        '800': '#363856',
        '900': '#313249',
        '950': '#0e0e15',
      }
    },
    textColor: {
      'base': {
        '50': '#f5f4ef',
        '100': '#e8e6d9',
        '200': '#d2cfb6',
        '300': '#b9b18d',
        '400': '#a3956c',
        '500': '#92825d',
        '600': '#7d6a4f',
        '700': '#665442',
        '800': '#59473b',
        '900': '#4e3f36',
        '950': '#15100e',
      }
    },
    extend: {
    },
  },
  plugins: [],
}

