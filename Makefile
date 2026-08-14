.PHONY: all frontend build test run clean

all: build

frontend:
	$(MAKE) -C frontend build

build: frontend
	go build -o listinator

# Depends on frontend, because the main package embeds frontend/dist and does
# not compile without it.
test: frontend
	go test ./...
	$(MAKE) -C frontend test

run: build
	LISTINATOR_SESSION_SECRET="secret" LISTINATOR_ADMIN_PASSWORD="secret" LISTINATOR_DATABASE_DIR=. ./listinator

clean:
	rm -rf listinator
	$(MAKE) -C frontend clean


