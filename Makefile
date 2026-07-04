.PHONY: scalar deps

# dependencies
deps:
	pnpm ci

# copy scalar standalone
scalar: deps
	cp node_modules/@scalar/api-reference/dist/browser/standalone.js ./standalone.js