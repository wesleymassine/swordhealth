# Variables
USER_SERVICE_PORT=8081
TASK_SERVICE_PORT=8082
NOTIFICATION_SERVICE_PORT=8083
DOCKER_COMPOSE_FILE=docker-compose.yml

# Run services individually (for local testing without Docker Compose)

.PHONY: run-user-service
run-user-service:
	@echo "Running User-Service on port $(USER_SERVICE_PORT)..."
	cd ./cmd/user-service && go run main.go
