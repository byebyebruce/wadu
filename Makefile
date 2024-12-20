build:
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin .

pull:
	@while true; do \
		echo "Running git pull..."; \
		git pull; \
		echo "Waiting for 30 seconds..."; \
		sleep 30; \
	done