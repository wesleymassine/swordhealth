# swordhealth
Sword Health - Technical backend challenge 

# swordhealth Architecture for User, Task, and Notification Services

This repository contains three services: **User-Service**, **Task-Service**, and **Notification-Service**. Each service is built with Go, utilizing clean architecture principles, and is ready for production use with Docker, Docker Compose, and RabbitMQ integration.

## Services Overview

### User-Service
The **User-Service** manages user registration, login, profile updates, and role-based access control.

### Task-Service
The **Task-Service** handles task creation, updates, and assignment. It also notifies managers when tasks are completed or updated.

### Notification-Service
The **Notification-Service** is responsible for sending notifications to managers whenever a task is updated by a technician.

## Table of Contents
- [Routes](#routes)
- [Requests and Responses](#requests-and-responses)
- [Running Locally](#running-locally)
- [Environment Variables](#environment-variables)

---

## User-Service

### Routes

| Method | Endpoint                | Description                 |
|--------|-------------------------|-----------------------------|
| POST   | `/users/register`        | Register a new user         |
| POST   | `/users/login`           | Login a user                |
| GET    | `/users/me`              | Get current user profile    |
| PUT    | `/users/update`          | Update user profile         |
| DELETE | `/users/delete/:id`      | Delete a user by ID         |

### Requests and Responses

#### 1. Register a New User
- **POST** `/users/register`
  
**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "role": "technician"
}
```
**Response**:
```json
{
  "message": "User registered successfully"
}
```

#### 2. User Login
- **POST** `/users/login`

**Request Body**:
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```
**Response**:
```json
{
  "token": "JWT-TOKEN-HERE",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "technician"
  }
}
```

---

## Task-Service

### Routes

| Method | Endpoint                | Description                 |
|--------|-------------------------|-----------------------------|
| POST   | `/tasks`                | Create a new task           |
| GET    | `/tasks`                | List all tasks              |
| GET    | `/tasks/:id`            | Get task details by ID      |
| PUT    | `/tasks/update/:id`     | Update a task by ID         |
| DELETE | `/tasks/delete/:id`     | Delete a task by ID         |

### Requests and Responses

#### 1. Create a Task
- **POST** `/tasks`

**Request Body**:
```json
{
  "title": "Fix server issue",
  "description": "Fix the bug in the server-side code.",
  "status": "pending",
  "assigned_to": 2
}
```
**Response**:
```json
{
  "message": "Task created successfully"
}
```

#### 2. Update a Task
- **PUT** `/tasks/update/:id`

**Request Body**:
```json
{
  "status": "completed"
}
```
**Response**:
```json
{
  "message": "Task updated successfully"
}
```

---

## Notification-Service

### Routes

| Method | Endpoint                | Description                 |
|--------|-------------------------|-----------------------------|
| GET    | `/notifications`        | List all notifications      |
| POST   | `/notifications/send`   | Send a notification manually|

### Requests and Responses

#### 1. List All Notifications
- **GET** `/notifications`

**Response**:
```json
[
  {
    "id": 1,
    "task_id": 1,
    "message": "Technician 1 completed task 'Fix server issue' on 2024-10-08.",
    "status": "sent",
    "sent_at": "2024-10-08T12:30:00Z"
  }
]
```

#### 2. Send a Notification Manually
- **POST** `/notifications/send`

**Request Body**:
```json
{
  "task_id": 1,
  "message": "Technician X completed task Y."
}
```
**Response**:
```json
{
  "message": "Notification sent successfully"
}
```

---

## Running Locally

To run all services (User-Service, Task-Service, Notification-Service) together locally:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo/swordhealth.git
   cd swordhealth
   ```

2. **Run all services with Docker Compose**:
   ```bash
   docker-compose up
   ```

This will start all three services:
- **User-Service** at `http://localhost:8081`
- **Task-Service** at `http://localhost:8082`
- **Notification-Service** at `http://localhost:8083`

Additionally, the services will communicate through RabbitMQ, which will be automatically set up based on the previous RabbitMQ configuration.

---

## Environment Variables

### User-Service Environment Variables

| Variable         | Description                      |
|------------------|----------------------------------|
| `DB_DSN`         | Database connection string       |
| `JWT_SECRET`     | Secret key for JWT authentication|
| `PORT`           | Service port (default: `8081`)   |

### Task-Service Environment Variables

| Variable         | Description                      |
|------------------|----------------------------------|
| `DB_DSN`         | Database connection string       |
| `RABBITMQ_URL`   | RabbitMQ connection URL          |
| `JWT_SECRET`     | Secret key for JWT authentication|
| `PORT`           | Service port (default: `8082`)   |

### Notification-Service Environment Variables

| Variable         | Description                      |
|------------------|----------------------------------|
| `DB_DSN`         | Database connection string       |
| `RABBITMQ_URL`   | RabbitMQ connection URL          |
| `PORT`           | Service port (default: `8083`)   |

## What I'd do differently or additional points

