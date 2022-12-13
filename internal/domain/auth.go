package domain

type RegisterUser struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string
	Password string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	DNI      int    `json:"dni"`
	Phone    int    `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ForgotRequest struct {
	Email string `json:"email"`
}

type GetUserFilters struct {
	AuthID string
	Email  string
}
