package accounts

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/leorodriguez/grupo-04/internal/auth"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/internal/transactions"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"math/rand"
	"strings"
)

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrTokenExpired       = errors.New("expired token")
)

type Service interface {
	Register(ctx context.Context, rq domain.RegisterRequest) (*users.UserDto, error)
	GetAccountInfo(ctx context.Context, id int, token string) (domain.AccountInfo, error)
	GetUserInfo(ctx context.Context, id int) (domain.UserInfo, error)
	GetTransactionsLastFive(ctx context.Context, id int, token string) ([]domain.TransactionInfo, error)
	IsAuthorized(ctx context.Context, id int, isUserID bool, token string) (bool, error)
	UpdateAccount(ctx context.Context, rq domain.RegisterRequest, id int) (*users.UserDto, error)
	UpdateAlias(ctx context.Context, accountID int, alias string) error
}

type service struct {
	auth                   auth.Auth
	usersService           users.Service
	accountsRepository     Repository
	transactionsRepository transactions.Repository
	aliasWords             []string
}

func NewService(usersService users.Service, accountsRepository Repository, transactionsRepository transactions.Repository, auth auth.Auth,
	aliasWords []string) Service {
	return &service{
		usersService:           usersService,
		accountsRepository:     accountsRepository,
		transactionsRepository: transactionsRepository,
		auth:                   auth,
		aliasWords:             aliasWords,
	}
}

func (s *service) Register(ctx context.Context, rq domain.RegisterRequest) (*users.UserDto, error) {
	if exists := s.auth.UserExists(ctx, rq.Email); exists {
		return &users.UserDto{}, users.ErrEmailAlreadyRegistered
	}

	ru := domain.RegisterUser{
		Name:     rq.Name,
		LastName: rq.LastName,
		Email:    rq.Email,
		Password: rq.Password,
	}

	userIDAuth, err := s.auth.Register(ctx, ru)
	if err != nil {
		logger.Error(err.Error())
		return &users.UserDto{}, users.ErrInternal
	}

	cvu := s.getNewCVU(ctx)
	alias := s.getNewAlias(ctx)

	accountToSave := domain.AccountDto{
		AuthID: userIDAuth,
		DNI:    rq.DNI,
		Phone:  rq.Phone,
		CVU:    cvu,
		Alias:  alias,
	}

	user, err := s.accountsRepository.SaveAccount(ctx, accountToSave)
	if err != nil {
		logger.Error(err.Error())
		return &users.UserDto{}, users.ErrInternal
	}

	err = s.auth.SendVerifyEmail(ctx, userIDAuth)
	if err != nil {
		logger.Error(err.Error())
	}

	return user, nil
}

func (s *service) GetAccountInfo(ctx context.Context, id int, token string) (domain.AccountInfo, error) {
	account, err := s.accountsRepository.GetAccountByID(ctx, id)
	if err != nil {
		return domain.AccountInfo{}, err
	}

	accountInfo := domain.AccountInfo{
		AccountID: account.ID,
		UserID:    account.User.ID,
		CVU:       account.CVU,
		Alias:     account.Alias,
		Balance:   account.Balance,
	}
	return accountInfo, nil
}

func (s *service) GetUserInfo(ctx context.Context, userID int) (domain.UserInfo, error) {
	user, err := s.usersService.GetByID(ctx, userID)
	if err != nil {
		return domain.UserInfo{}, err
	}

	account, err := s.accountsRepository.GetAccountByUserID(ctx, userID)
	if err != nil {
		return domain.UserInfo{}, err
	}

	filters := domain.GetUserFilters{
		AuthID: account.AuthID,
	}
	authUsers, err := s.auth.GetUsersByEmail(ctx, filters)
	if err != nil {
		return domain.UserInfo{}, err
	}
	if len(authUsers) == 0 {
		// todo error
		return domain.UserInfo{}, nil
	}

	accountInfo := domain.UserInfo{
		Name:     *authUsers[0].FirstName,
		LastName: *authUsers[0].LastName,
		Email:    *authUsers[0].Email,
		DNI:      user.DNI,
		Phone:    user.Phone,
	}

	return accountInfo, nil
}

func (s *service) IsAuthorized(ctx context.Context, id int, isUserID bool, token string) (bool, error) {
	authID, err := s.auth.GetIDFromToken(ctx, token)
	if err != nil {
		switch err.Error() {
		case "could not decode accessToken with custom claims: Token is expired":
			return false, ErrTokenExpired
		default:
			return false, err
		}
	}

	isAuthorized, err := s.accountsRepository.IsAuthorized(ctx, id, isUserID, authID)
	if err != nil {
		return false, err
	}

	return isAuthorized, nil
}

func (s *service) GetTransactionsLastFive(ctx context.Context, id int, token string) ([]domain.TransactionInfo, error) {
	trx, err := s.transactionsRepository.GetAllByIDLimit(ctx, id, 5)
	if err != nil {
		return []domain.TransactionInfo{}, err
	}

	return trx, nil
}

func (s *service) UpdateAccount(ctx context.Context, rq domain.RegisterRequest, id int) (*users.UserDto, error) {
	account, err := s.accountsRepository.GetAccountByID(ctx, id)
	if err != nil {
		return &users.UserDto{}, err
	}

	if rq.Email != "" {
		if exists := s.auth.UserExists(ctx, rq.Email); exists {
			return &users.UserDto{}, users.ErrEmailAlreadyRegistered
		}
	}

	if rq.Name != "" || rq.LastName != "" || rq.Email != "" {
		ru := domain.RegisterUser{}
		if rq.Name != "" {
			ru.Name = rq.Name
		}
		if rq.LastName != "" {
			ru.LastName = rq.LastName
		}
		if rq.Email != "" {
			ru.Email = rq.Email
		}

		err = s.auth.Update(ctx, ru, account.AuthID)
		if err != nil {
			logger.Error(err.Error())
			return &users.UserDto{}, users.ErrInternal
		}
	}

	if rq.DNI != 0 || rq.Phone != 0 {
		accountToSave := domain.AccountDto{}
		if rq.DNI != 0 {
			accountToSave.DNI = rq.DNI
		}
		if rq.Phone != 0 {
			accountToSave.Phone = rq.Phone
		}

		err := s.usersService.UpdateUser(ctx, accountToSave, account.User.ID)
		if err != nil {
			logger.Error(err.Error())
			return &users.UserDto{}, users.ErrInternal
		}
	}
	return &users.UserDto{}, nil
}

func (s *service) UpdateAlias(ctx context.Context, accountID int, alias string) error {
	exists := s.accountsRepository.AliasExist(ctx, alias)
	if exists {
		return ErrAliasAlreadyExists
	}

	err := s.accountsRepository.UpdateAlias(ctx, accountID, alias)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) getNewCVU(ctx context.Context) string {
	var cvu string
	for true {
		cvu = fmt.Sprint(rand.Int63n(1000000000000000000))
		cvu += fmt.Sprint(rand.Int63n(10000))

		if exists := s.accountsRepository.CVUExist(ctx, cvu); !exists {
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
		if exists := s.accountsRepository.AliasExist(ctx, alias); !exists {
			break
		}
	}

	return alias
}
