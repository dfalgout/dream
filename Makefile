.PHONY: watch css-watch js-watch build

watch:
	@echo "Watching for changes..."
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

css-watch:
	@bun tailwindcss -i input.css -o assets/css/output.css --watch

js-watch:
	@bun esbuild --bundle javascript/*.js --outdir=assets/js --watch

build:
	@bun tailwindcss -i input.css -o assets/css/output.css --minify
	@bun esbuild --bundle javascript/*.js --outdir=assets/js --minify