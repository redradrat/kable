mix = require("laravel-mix")

mix.js("public/js/*.js", "dist/public/js")
    .sass("public/css/main.scss", "dist/public/css")
    .copy("public/img/*.png", "dist/public/img")
    .copy("views/*.pug", "dist/views")
    .copy("package.dist.json", "dist/package.json")
