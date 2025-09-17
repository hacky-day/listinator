.PHONY: all build run clean

all: build

node_modules: package.json
	npm install

build: node_modules
	npm run build

run: node_modules
	npm run dev

clean:
	rm -rf node_modules dist
