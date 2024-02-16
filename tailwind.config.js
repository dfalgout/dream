/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./**/*.templ"],
    theme: {
        extend: {},
    },
    plugins: [
        require("@tailwindcss/forms"),
        require("@tailwindcss/typography"),
        require("daisyui"),
    ],
    daisyui: {
        logs: false,
    }
};