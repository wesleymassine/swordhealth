# Variables
USER_SERVICE_PORT=8081
TASK_SERVICE_PORT=8082
NOTIFICATION_SERVICE_PORT=8083
DOCKER_COMPOSE_FILE=docker-compose.yaml

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
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all
	docker system prune -f

# Cleaning up unused Docker containers, images, and networks
.PHONY: prune
prune:
	@echo "Pruning unused Docker containers, images, and networks..."
	docker system prune -f

# Run migrations
.PHONY: migrate-up
migrate-up:
	bash -c "source ./scripts.sh && migrate_up"

.PHONY: migrate-down
migrate-down:
	bash -c "source ./scripts.sh && migrate_down"	

# Run services individually (for local testing without Docker Compose)

.PHONY: run-user-service
run-user-service:
	bash -c "source ./scripts.sh && user_service $(USER_SERVICE_PORT)"

.PHONY: run-task-service
run-task-service:
	bash -c "source ./scripts.sh && task_service $(TASK_SERVICE_PORT)"

.PHONY: run-notification-service
run-notification-service:
	bash -c "source ./scripts.sh && notification_service $(NOTIFICATION_SERVICE_PORT)"
