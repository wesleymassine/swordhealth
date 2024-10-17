package usecase

import (
	"context"

	"github.com/wesleymassine/swordhealth/user-management/domain"
)

type UserUsecase struct {
	UserRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{UserRepo: repo}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	r, err := uc.UserRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		Role:     r.Role,
	}, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int64) (*domain.UserResponse, error) {
	r, err := uc.UserRepo.GetUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        r.ID,
		Username:  r.Username,
		Email:     r.Email,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}, nil

}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.UserRepo.GetUserByEmail(ctx, email)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	if err := uc.UserRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return uc.GetUserByID(ctx, user.ID)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	return uc.UserRepo.DeleteUser(ctx, id)
}
