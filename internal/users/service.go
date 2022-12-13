package users

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/leorodriguez/grupo-04/internal/auth"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"math/rand"
	"strings"
)

var (
	ErrInternal               = errors.New("internal error")
	ErrInvalidUserCredentials = errors.New("invalid user credentials")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrUserNotExists          = errors.New("user not exists")
	ErrEmailNotVerified       = errors.New("email not verified")
)

type UserDto struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	DNI      int    `json:"dni"`
	Phone    int    `json:"phone"`
	Email    string `json:"email"`
	CVU      string `json:"cvu"`
	Alias    string `json:"alias"`
}

type Service interface {
	Login(ctx context.Context, rq domain.LoginRequest) (domain.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	UpdateUser(ctx context.Context, accountDto domain.AccountDto, id int) error
	GetByID(ctx context.Context, id int) (domain.UserDB, error)
}

type service struct {
	auth       auth.Auth
	repository Repository
	aliasWords []string
}

func NewUsers(auth auth.Auth, repository Repository, aliasWords []string) Service {
	return &service{auth: auth, repository: repository, aliasWords: aliasWords}
}

func (s *service) Login(ctx context.Context, rq domain.LoginRequest) (domain.LoginResponse, error) {
	if exists := s.auth.UserExists(ctx, rq.Email); !exists {
		return domain.LoginResponse{}, ErrUserNotExists
	}

	lu := domain.LoginUser{
		Email:    rq.Email,
		Password: rq.Password,
	}

	token, err := s.auth.Login(ctx, lu)
	if err != nil {
		logger.Error(err.Error())

		switch err {
		case auth.ErrAuthInvalidUserCredentials:
			return domain.LoginResponse{}, ErrInvalidUserCredentials
		case auth.ErrEmailNotVerified:
			return domain.LoginResponse{}, ErrEmailNotVerified
		default:
			return domain.LoginResponse{}, ErrInternal
		}
	}

	return domain.LoginResponse{Token: token}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	return s.auth.Logout(ctx, token)
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	filters := domain.GetUserFilters{
		Email: email,
	}
	users, err := s.auth.GetUsersByEmail(ctx, filters)
	if err != nil {
		logger.Error(err.Error())

		return err
	}

	return s.auth.SendEmail(ctx, *users[0].ID)
}

func (s *service) UpdateUser(ctx context.Context, accountDto domain.AccountDto, id int) error {
	return s.repository.UpdateUser(ctx, accountDto, id)
}

func (s *service) GetByID(ctx context.Context, id int) (domain.UserDB, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *service) getNewCVU(ctx context.Context) string {
	var cvu string
	for true {
		cvu = fmt.Sprint(rand.Int63n(1000000000000000000))
		cvu += fmt.Sprint(rand.Int63n(10000))

		if exists := s.repository.CVUExist(ctx, cvu); !exists {
			break
		}
	}

	return cvu
}

func (s *service) getNewAlias(ctx context.Context) string {
	var aliasList []string
	var alias string
	listLen := len(s.aliasWords)

	for true {
		for i := 0; i < 3; i++ {
			y := rand.Intn(listLen)
			aliasList = append(aliasList, s.aliasWords[y])
		}

		alias = strings.Join(aliasList, ".")
		if exists := s.repository.AliasExist(ctx, alias); !exists {
			break
		}
	}

	return alias
}
