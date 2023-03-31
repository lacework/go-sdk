
.PHONY: install
install:
	npm install

.PHONY: build
build:
	npm run asbuild

.PHONY: package
package: build
	cp build/release.wasm build/hackathon.wasm
	zip -j hackathon.zip build/hackathon.wasm spec.json
