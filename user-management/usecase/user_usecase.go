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

func (uc *UserUsecase) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return uc.UserRepo.Create(ctx, user)
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return uc.UserRepo.GetUserByID(ctx, id)
}

func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.UserRepo.GetUserByEmail(ctx, email)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return uc.UserRepo.UpdateUser(ctx, user)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	return uc.UserRepo.DeleteUser(ctx, id)
}
