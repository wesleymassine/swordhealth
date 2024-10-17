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

| Method | Endpoint                 | Description                 |
|--------|--------------------------|-----------------------------|
| POST   | `/users/register`        | Register a new user         |
| POST   | `/users/login`           | Login a user                |
| GET    | `/users/profile/:id`     | Get current user profile    |
| PUT    | `/users/update/:id`      | Update user profile         |
| DELETE | `/users/delete/:id`      | Delete a user by ID         |
| GET    | `/healthcheck`           | Healt check                 |            

### Requests and Responses

#### 1. Register a New User
- **POST** `/users/register`
  
**Request Body**:
```json
{
  "name": "Manager",
  "email": "manager@gmail.com",
  "password": "manager",
  "role": "manager"
}
```
**Response**:
```json
{
    "username": "Manager",
    "email": "manager@gmail.com",
    "role": "Manager"
}
```

#### 2. User Login
- **POST** `/users/login`

**Request Body**:
```json
{
    "email": "manager@gmail.com",
    "password": "manager"
}
```
**Response**:
```json
{
  "token": "JWT-TOKEN-HERE",
}
```

---

## Task-Service

### Routes

| Method | Endpoint                | Description                 |
|--------|-------------------------|-----------------------------|
| POST   | `/tasks`                | Create a new task           |
| GET    | `/tasks`                | List all tasks              |
| GET    | `/tasks/:task_id`       | Update a task by ID         |
| PUT    | `/tasks/:task_id/status`| PUT task task               |
| GET    | `/healthcheck`          | Healt check                 |            

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
    "id": 1,
    "title": "Fix the bug in the server-side code.",
    "description": "The server is down and needs immediate attention",
    "status": "pending",
    "assigned_to": 2,
    "performed_by": 0,
    "performed_at": "0001-01-01T00:00:00Z",
    "created_at": "2024-10-17T04:26:02.071374177+01:00"
}
```

#### 2. Update a Task
- **PUT** `/tasks/:task_id/status`

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
| GET    | `/healthcheck`          | Healt check                 |            

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

