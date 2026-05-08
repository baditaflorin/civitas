/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#17202a",
        paper: "#f6f4ee",
        panel: "#fffdfa",
        line: "#d8d2c3",
        civic: "#0b5cad",
        signal: "#bb4d00",
        field: "#2f6f4e",
      },
      boxShadow: {
        soft: "0 18px 60px rgba(23, 32, 42, 0.10)",
      },
    },
  },
  plugins: [],
};
