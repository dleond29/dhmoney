package domain

type User struct {
	ID     int
	IDAuth string
	DNI    int
	Phone  int
	CVU    string
	Alias  string
}

type UserDto struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	DNI      int    `json:"dni"`
	Phone    int    `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	CVU      string `json:"cvu"`
	Alias    string `json:"alias"`
}

type UserInfo struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	DNI      int    `json:"dni"`
	Phone    int    `json:"phone"`
}

type UserDB struct {
	ID    int
	DNI   int
	Phone int
}
