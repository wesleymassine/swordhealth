package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssignedTo  int64     `json:"assigned_to"`
	PerformedBy int64     `json:"performed_by"`
	PerformedAt time.Time `json:"performed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Claims struct {
	UserID int64    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
