package users

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Nerzal/gocloak/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/mocks"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Register(ctx context.Context, rq domain.RegisterRequest) (UserDto, error) {
	args := r.Called(ctx, rq)
	return args.Get(0).(UserDto), args.Error(1)
}

func (r *repoMock) Login(ctx context.Context, rq domain.LoginRequest) (domain.LoginResponse, error) {
	args := r.Called(ctx, rq)
	return args.Get(0).(domain.LoginResponse), args.Error(1)
}

func (r *repoMock) Logout(ctx context.Context, token string) error {
	args := r.Called(ctx, token)
	return args.Error(0)
}

func (r *repoMock) ForgotPassword(ctx context.Context, email string) error {
	args := r.Called(ctx, email)
	return args.Error(0)
}

func TestService_Logout(t *testing.T) {
	var ctx = context.Background()
	token := "dfadfa3234234"

	testCases := []struct {
		name     string
		repoMock func(m *mock.Mock)
		authMock func(m *mock.Mock)
	}{
		{
			name: "Successfully logout",
			authMock: func(m *mock.Mock) {
				m.On("Logout", ctx, token).Return(nil).Once()
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			repo := NewRepository(db)

			authMock := &mocks.Auth{}
			testCase.authMock(&authMock.Mock)

			usersService := NewUsers(authMock, repo, []string{"Perro", "Gato"})

			err = usersService.Logout(ctx, token)

			assert.Equal(t, nil, err)
		})
	}
}

func Test_service_ForgotPassword(t *testing.T) {
	var ctx = context.Background()
	email := "digitalhouse@gmail.com"

	testCases := []struct {
		name          string
		repoMock      func(m *mock.Mock)
		domainMock    func(m *mock.Mock)
		expectedError error
	}{
		{
			name: "Successfully forgot password",
			domainMock: func(m *mock.Mock) {
				m.On("GetUsersByEmail", ctx, email).
					Return(
						[]*gocloak.User{}, errors.New("error"),
					).Once()
			},
			expectedError: errors.New("error"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, err := sqlmock.New()
			repo := NewRepository(db)

			domainMock := &mocks.Auth{}
			testCase.domainMock(&domainMock.Mock)

			usersService := NewUsers(domainMock, repo, []string{"Perro", "Gato"})

			err = usersService.ForgotPassword(ctx, email)

			assert.Equal(t, testCase.expectedError, err)

		})
	}
}
