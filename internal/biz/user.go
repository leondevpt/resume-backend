package biz

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{}
}

type UserRepo interface {
	GetUserByID()
	GetUserByEmail(email string)
	SaveUser() error
}
