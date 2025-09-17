.PHONY: all frontend build run clean

all: build

frontend:
	$(MAKE) -C frontend build

build: frontend
	go build -o listinator

run: build
	LISTINATOR_SESSION_SECRET="secret" LISTINATOR_ADMIN_PASSWORD="secret" LISTINATOR_DATABASE_DIR=. ./listinator

clean:
	rm -rf listinator
	$(MAKE) -C frontend clean


