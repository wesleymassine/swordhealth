# Variables
USER_SERVICE_PORT=8081
TASK_SERVICE_PORT=8082
NOTIFICATION_SERVICE_PORT=8083
DOCKER_COMPOSE_FILE=docker-compose.yaml

DB_URL="mysql://root:root@tcp(localhost:3306)/task"

.PHONY: up
up:
	@echo "Starting all services with Docker Compose..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop services with Docker Compose
.PHONY: down
down:
	@echo "Stopping all services..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# Clean Docker containers and images
.PHONY: clean
clean:
	@echo "Cleaning Docker containers and images..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all
	docker system prune -f

# Cleaning up unused Docker containers, images, and networks
.PHONY: prune
prune:
	@echo "Pruning unused Docker containers, images, and networks..."
	docker system prune -f

# Run migrations using docker exec
.PHONY: migrate
migrate:
	@echo "Running migrations..."
	cd ./pkg/migrations && go run cmd/main.go


# Run migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	cd ./pkg/migrations && go run cmd/main.go migrate-up

.PHONY: migrate-down
migrate-down:
	@echo "Reverting migrations..."
		cd ./pkg/migrations && go run cmd/main.go migrate-down

# Run services individually (for local testing without Docker Compose)

.PHONY: run-user-service
run-user-service:
	@echo "Running User-Service on port $(USER_SERVICE_PORT)..."
	cd ./user-management && cd ./cmd/user-service && go run main.go

.PHONY: run-task-service
run-task-service:
	@echo "Running Task-Service on port $(TASK_SERVICE_PORT)..."
	cd ./task-management && cd ./cmd/task-service && go run main.go

.PHONY: run-notification-service
run-notification-service:
	@echo "Running Notification-Service on port $(NOTIFICATION_SERVICE_PORT)..."
	cd ./user-notification && cd ./cmd/notification-service && go run main.go