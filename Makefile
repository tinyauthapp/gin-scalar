.PHONY: scalar deps test vet

# dependencies
deps:
	pnpm ci

# copy scalar standalone
scalar: deps
	cp node_modules/@scalar/api-reference/dist/browser/standalone.js ./standalone.js

# test
test:
	go test -race ./...

# vet
vet:
	go vet ./...