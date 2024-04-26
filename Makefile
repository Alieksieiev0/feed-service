# Build the application
build:
	@echo "Building..."
	
	@go build -o user-service cmd/user-service/main.go

# Run the application
run:
	@go run cmd/user-service/main.go

# Run the application in docker
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		docker-compose up; \
	fi


# Run the application in docker and rebuild it
docker-build-n-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		docker-compose up --build; \
	fi

# Shutdown docker
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests/... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f user-service


# Generate application proto files
proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air -c .air.toml; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi


.PHONY: build run docker-run docker-build-n-run docker-down test clean proto watch
