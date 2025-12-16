gen:
  tailwindcss --minify --input ./internal/frontend/static/input.css --output ./internal/frontend/static/output.css
  rustywind --write --output-css-file ./internal/frontend/static/output.css .
  go tool templ generate
