.PHONY: all build clean
all: build

build:
	$(MAKE) -C core build
	$(MAKE) -C typifier build

clean:
	$(MAKE) -C core clean
	$(MAKE) -C typifier clean

run-typifier:
	$(MAKE) -C typifier run

run-core:
	$(MAKE) -C core run

run-frontend:
	$(MAKE) -C core/frontend run
