user_service() {
    set -a
    cp .env.local.example .env

    source .env

    if [ $1 ]; then
        USER_SERVICE_PORT=$1
    fi

    echo "Running User-Service on port ${USER_SERVICE_PORT}..."

    # Navigate to the correct directory and run the service
    cd ./user-management/cmd/user-service || exit
    go run main.go
}

task_service() {
    set -a
    cp .env.local.example .env

    source .env

    if [ $1 ]; then
        TASK_SERVICE_PORT=$1
    fi

    echo "Running Task-Service on port $(TASK_SERVICE_PORT)..."

    # Navigate to the correct directory and run the service
    cd ./task-management/cmd/task-service || exit
    go run main.go
}

notification_service() {
    set -a
    cp .env.local.example .env

    source .env

    if [ $1 ]; then
        NOTIFICATION_SERVICE_PORT=$1
    fi

    echo "Running Notification-Service on port $(NOTIFICATION_SERVICE_PORT)..."

    # Navigate to the correct directory and run the service
    cd ./user-notification/cmd/notification-service || exit
    go run main.go
}

migrate_up() {
    set -a
    cp .env.local.example .env

    source .env

	echo "Running migrations up..."

    # Navigate to the correct directory and run the service
    cd ./pkg/migrations|| exit
    go run cmd/main.go migrate-up
}

migrate_down() {
    set -a
    cp .env.local.example .env

    source .env

	echo "Reverting migrations..."

    # Navigate to the correct directory and run the service
    cd ./pkg/migrations|| exit
    go run cmd/main.go migrate-down
}