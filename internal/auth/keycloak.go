package auth

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v12"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"net/http"
)

var (
	ErrAuthInvalidUserCredentials = errors.New("invalid user credentials")
	ErrEmailNotVerified           = errors.New("400 Bad Request: invalid_grant: Account is not fully set up")
)

type Auth interface {
	Register(ctx context.Context, rq domain.RegisterUser) (string, error)
	Login(ctx context.Context, rq domain.LoginUser) (string, error)
	Logout(ctx context.Context, token string) error
	UserExists(ctx context.Context, email string) bool
	GetUsersByEmail(ctx context.Context, filters domain.GetUserFilters) ([]*gocloak.User, error)
	SendVerifyEmail(ctx context.Context, userID string) error
	SendEmail(ctx context.Context, userID string) error
	GetIDFromToken(ctx context.Context, accessToken string) (string, error)
	Update(ctx context.Context, rq domain.RegisterUser, authID string) error
}

type Gocloak interface {
	LoginAdmin(ctx context.Context, username string, password string, realm string) (*gocloak.JWT, error)
	CreateUser(ctx context.Context, token string, realm string, user gocloak.User) (string, error)
	Login(ctx context.Context, clientID string, clientSecret string, realm string, username string, password string) (*gocloak.JWT, error)
	Logout(ctx context.Context, clientID string, clientSecret string, realm string, refreshToken string) error
	GetUsers(ctx context.Context, token string, realm string, params gocloak.GetUsersParams) ([]*gocloak.User, error)
	SendVerifyEmail(ctx context.Context, token string, userID string, realm string, params ...gocloak.SendVerificationMailParams) error
	ExecuteActionsEmail(ctx context.Context, token string, realm string, params gocloak.ExecuteActionsEmail) error
	GetUserInfo(ctx context.Context, accessToken string, realm string) (*gocloak.UserInfo, error)
	DecodeAccessToken(ctx context.Context, accessToken string, realm string) (*jwt.Token, *jwt.MapClaims, error)
	UpdateUser(ctx context.Context, token string, realm string, user gocloak.User) error
	GetUserByID(ctx context.Context, accessToken string, realm string, userID string) (*gocloak.User, error)
}

type KeycloakSettings struct {
	GoCloak      Gocloak
	ClientId     string
	ClientSecret string
	Realm        string
}

type auth struct {
	gocloak      Gocloak
	clientId     string
	clientSecret string
	realm        string
}

func NewAuth(settings KeycloakSettings) Auth {
	return &auth{
		gocloak:      settings.GoCloak,
		clientId:     settings.ClientId,
		clientSecret: settings.ClientSecret,
		realm:        settings.Realm,
	}
}

func (auth *auth) Register(ctx context.Context, rq domain.RegisterUser) (string, error) {
	token, err := auth.loginAdmin(ctx)
	if err != nil {
		return "", err
	}

	user := gocloak.User{
		FirstName:     gocloak.StringP(rq.Name),
		LastName:      gocloak.StringP(rq.LastName),
		Email:         gocloak.StringP(rq.Email),
		Enabled:       gocloak.BoolP(true),
		Username:      gocloak.StringP(rq.Email),
		EmailVerified: gocloak.BoolP(false),
		Credentials: &[]gocloak.CredentialRepresentation{
			{Value: &rq.Password, Temporary: gocloak.BoolP(false), Type: gocloak.StringP("password")}},
		RequiredActions: &[]string{"VERIFY_EMAIL"},
	}

	userID, err := auth.gocloak.CreateUser(ctx, token, "realm-test", user)
	if err != nil {
		logger.Error(err.Error())

		return "", err
	}

	return userID, nil
}

func (auth *auth) Update(ctx context.Context, rq domain.RegisterUser, authID string) error {
	token, err := auth.loginAdmin(ctx)
	if err != nil {
		return err
	}

	filters := domain.GetUserFilters{
		AuthID: authID,
	}
	users, err := auth.GetUsersByEmail(ctx, filters)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if len(users) == 0 {
		return errors.New("")
	}

	user := users[0]
	if rq.Name != "" {
		user.FirstName = &rq.Name
	}
	if rq.LastName != "" {
		user.LastName = &rq.LastName
	}
	if rq.Email != "" {
		user.Email = &rq.Email
		user.Username = &rq.Email
	}

	err = auth.gocloak.UpdateUser(ctx, token, auth.realm, *user)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (auth *auth) Login(ctx context.Context, rq domain.LoginUser) (string, error) {
	jwt, err := auth.gocloak.Login(ctx, auth.clientId, auth.clientSecret, auth.realm, rq.Email, rq.Password)
	if err != nil {
		logger.Error(err.Error())

		if apiError, ok := err.(*gocloak.APIError); ok {
			switch apiError.Code {
			case http.StatusUnauthorized:
				return "", ErrAuthInvalidUserCredentials
			case http.StatusBadRequest:
				return "", ErrEmailNotVerified
			default:
				return "", errors.New("internal error")
			}
		}

		return "", err
	}

	return jwt.AccessToken, nil
}

func (auth *auth) Logout(ctx context.Context, token string) error {
	err := auth.gocloak.Logout(ctx, auth.clientId, auth.clientSecret, auth.realm, token)
	if err != nil {
		return err
	}

	return nil
}

func (auth *auth) UserExists(ctx context.Context, email string) bool {
	filters := domain.GetUserFilters{
		Email: email,
	}
	users, _ := auth.GetUsersByEmail(ctx, filters)
	return len(users) > 0
}

// TODO get by more params
func (auth *auth) GetUsersByEmail(ctx context.Context, filters domain.GetUserFilters) ([]*gocloak.User, error) {
	token, err := auth.loginAdmin(ctx)
	if err != nil {
		logger.Error(err.Error())

		return []*gocloak.User{}, err
	}

	params := gocloak.GetUsersParams{}
	if filters.AuthID != "" {
		users, err := auth.gocloak.GetUserByID(ctx, token, auth.realm, filters.AuthID)
		if err != nil {
			logger.Error(err.Error())

			return []*gocloak.User{}, err
		}
		return []*gocloak.User{users}, nil
	}

	if filters.Email != "" {
		params.Email = &filters.Email
	}
	users, err := auth.gocloak.GetUsers(ctx, token, auth.realm, params)
	if err != nil {
		logger.Error(err.Error())

		return []*gocloak.User{}, err
	}

	return users, nil
}

func (auth *auth) SendVerifyEmail(ctx context.Context, userID string) error {
	token, err := auth.loginAdmin(ctx)
	if err != nil {
		logger.Error(err.Error())

		return err
	}

	err = auth.gocloak.SendVerifyEmail(ctx, token, userID, auth.realm)
	if err != nil {
		logger.Error(err.Error())

		return err
	}

	return nil
}

func (auth *auth) SendEmail(ctx context.Context, userID string) error {
	token, err := auth.loginAdmin(ctx)
	if err != nil {
		logger.Error(err.Error())

		return err
	}

	params := gocloak.ExecuteActionsEmail{
		ClientID:    &auth.clientId,
		RedirectURI: gocloak.StringP("localhost:8080/forgot"),
		UserID:      &userID,
		Actions:     &[]string{"UPDATE_PASSWORD"},
	}
	err = auth.gocloak.ExecuteActionsEmail(ctx, token, auth.realm, params)
	if err != nil {
		logger.Error(err.Error())

		return err
	}

	return nil
}

func (auth *auth) GetIDFromToken(ctx context.Context, accessToken string) (string, error) {
	token, _, err := auth.gocloak.DecodeAccessToken(ctx, accessToken, auth.realm)
	if err != nil {
		logger.Error(err.Error())

		return "", err
	}

	id, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
	logger.Any(id, ok)
	return id, nil
}

func (auth *auth) loginAdmin(ctx context.Context) (string, error) {
	token, err := auth.gocloak.LoginAdmin(ctx, "admin1", "admin1", "realm-test")
	if err != nil {
		logger.Error(err.Error())

		return "", err
	}

	return token.AccessToken, nil
}
