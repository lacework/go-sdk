
.PHONY: install
install:
	npm install

.PHONY: build
build:
	npm run asbuild


.PHONY: package
package: build
	zip -j hackathon.zip build/release.wasm spec.json
