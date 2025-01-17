build:
	mkdir -p bin
	go build -o bin .

build-musl:
	mkdir -p bin
	go build -tags "musl" -o bin .

pull:
	@while true; do \
		echo "Running git pull..."; \
		git pull; \
		echo "Waiting for 30 seconds..."; \
		sleep 30; \
	done