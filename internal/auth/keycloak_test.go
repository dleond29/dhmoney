package auth

import (
	"context"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"testing"

	"github.com/Nerzal/gocloak/v12"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type GocloakMock struct {
	mock.Mock
	gocloak.GoCloak
}

func (g *GocloakMock) LoginAdmin(ctx context.Context, username string, password string, realm string) (*gocloak.JWT, error) {
	args := g.Called(username, password, realm)
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

func (g *GocloakMock) CreateUser(ctx context.Context, token string, realm string, user gocloak.User) (string, error) {
	args := g.Called(token, realm, realm, user)
	return args.String(0), args.Error(1)
}

func (g *GocloakMock) Login(ctx context.Context, clientID string, clientSecret string, realm string, username string, password string) (*gocloak.JWT, error) {
	args := g.Called(clientID, clientSecret, realm, username, password)
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}

func (g *GocloakMock) Logout(ctx context.Context, clientID string, clientSecret string, realm string, refreshToken string) error {
	args := g.Called(clientID, clientSecret, realm, refreshToken)
	return args.Error(0)
}

func (g *GocloakMock) GetUsers(ctx context.Context, token string, realm string, params gocloak.GetUsersParams) ([]*gocloak.User, error) {
	args := g.Called(token, realm, params)
	return args.Get(0).([]*gocloak.User), args.Error(1)
}

func (g *GocloakMock) SendVerifyEmail(ctx context.Context, token string, userID string, realm string, params ...gocloak.SendVerificationMailParams) error {
	args := g.Called(token, userID, realm, params)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	var tests = []struct {
		name      string
		rq        domain.RegisterUser
		auth      Auth
		userID    string
		wantError bool
	}{
		{
			name: "register - ok",
			rq: domain.RegisterUser{
				Name:     "name",
				LastName: "lastname",
				Email:    "email",
				Password: "password",
			},
			auth: func() Auth {
				gMock := new(GocloakMock)
				gMock.On("LoginAdmin", "admin1", "admin1", "realm-test").
					Return(&gocloak.JWT{AccessToken: "accessToken"}, nil)

				user := gocloak.User{
					FirstName:     gocloak.StringP("name"),
					LastName:      gocloak.StringP("lastname"),
					Email:         gocloak.StringP("email"),
					Enabled:       gocloak.BoolP(true),
					Username:      gocloak.StringP("email"),
					EmailVerified: gocloak.BoolP(false),
					Credentials: &[]gocloak.CredentialRepresentation{
						{Value: gocloak.StringP("password"), Temporary: gocloak.BoolP(false), Type: gocloak.StringP("password")}},
					RequiredActions: &[]string{"VERIFY_EMAIL"},
				}

				gMock.On("CreateUser", "accessToken", "realm-test", "realm-test", user).
					Return("userID", nil)

				keycloakSettings := KeycloakSettings{
					GoCloak: gMock,
				}
				auth := NewAuth(keycloakSettings)
				return auth
			}(),
			userID:    "userID",
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID, err := tt.auth.Register(context.Background(), tt.rq)
			if tt.wantError {
				require.Error(t, err)
			}
			require.Equal(t, tt.userID, userID)
		})
	}
}

func TestLogin(t *testing.T) {
	var tests = []struct {
		name          string
		rq            domain.LoginUser
		auth          Auth
		tokenResponse string
		wantError     bool
	}{
		{
			name: "login - ok",
			rq: domain.LoginUser{
				Email:    "email@c.com",
				Password: "password",
			},
			auth: func() Auth {
				gMock := new(GocloakMock)

				gMock.On("Login", "clientID", "clientSecret", "realm-test", "email@c.com", "password").
					Return(&gocloak.JWT{RefreshToken: "accessToken"}, nil)

				keycloakSettings := KeycloakSettings{
					GoCloak:      gMock,
					ClientId:     "clientID",
					ClientSecret: "clientSecret",
					Realm:        "realm-test",
				}
				auth := NewAuth(keycloakSettings)
				return auth
			}(),
			tokenResponse: "accessToken",
			wantError:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			response, err := tt.auth.Login(context.Background(), tt.rq)
			if tt.wantError {
				require.Error(t, err)
			}
			require.Equal(t, tt.tokenResponse, response)
		})
	}
}

func TestLogout(t *testing.T) {
	var tests = []struct {
		name      string
		token     string
		auth      Auth
		wantError bool
	}{
		{
			name:  "logout - ok",
			token: "token",
			auth: func() Auth {
				gMock := new(GocloakMock)

				gMock.On("Logout", "clientID", "clientSecret", "realm-test", "token").
					Return(nil)

				keycloakSettings := KeycloakSettings{
					GoCloak:      gMock,
					ClientId:     "clientID",
					ClientSecret: "clientSecret",
					Realm:        "realm-test",
				}
				auth := NewAuth(keycloakSettings)
				return auth
			}(),
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.auth.Logout(context.Background(), tt.token)

			require.True(t, tt.wantError == (err != nil))
		})
	}
}

func TestGetUsersByEmail(t *testing.T) {
	var tests = []struct {
		name      string
		email     string
		auth      Auth
		usersOut  []*gocloak.User
		wantError bool
	}{
		{
			name:  "GetUsersByEmail - ok",
			email: "email@c.com",
			auth: func() Auth {
				gMock := new(GocloakMock)
				gMock.On("LoginAdmin", "admin1", "admin1", "realm-test").
					Return(&gocloak.JWT{AccessToken: "accessToken"}, nil)

				params := gocloak.GetUsersParams{Email: gocloak.StringP("email@c.com")}
				user := []*gocloak.User{
					{
						ID: gocloak.StringP("id1"),
					},
				}
				gMock.On("GetUsers", "accessToken", "realm-test", params).
					Return(user, nil)

				keycloakSettings := KeycloakSettings{
					GoCloak: gMock,
					Realm:   "realm-test",
				}
				auth := NewAuth(keycloakSettings)
				return auth
			}(),
			usersOut: []*gocloak.User{
				{
					ID: gocloak.StringP("id1"),
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			users, err := tt.auth.GetUsersByEmail(context.Background(), tt.email)
			if tt.wantError {
				require.Error(t, err)
			}
			require.Equal(t, tt.usersOut, users)
		})
	}
}

func TestSendVerifyEmail(t *testing.T) {
	var tests = []struct {
		name      string
		userID    string
		auth      Auth
		wantError bool
	}{
		{
			name:   "GetUsersByEmail - ok",
			userID: "userID",
			auth: func() Auth {
				gMock := new(GocloakMock)
				gMock.On("LoginAdmin", "admin1", "admin1", "realm-test").
					Return(&gocloak.JWT{AccessToken: "accessToken"}, nil)

				var params []gocloak.SendVerificationMailParams
				gMock.On("SendVerifyEmail", "accessToken", "userID", "realm-test", params).
					Return(nil)

				keycloakSettings := KeycloakSettings{
					GoCloak: gMock,
					Realm:   "realm-test",
				}
				auth := NewAuth(keycloakSettings)
				return auth
			}(),
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.auth.SendVerifyEmail(context.Background(), tt.userID)

			require.True(t, tt.wantError == (err != nil))
		})
	}
}
