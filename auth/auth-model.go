package auth

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type SafeUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type Users struct {
	Users []string `json:"users"`
}

type Role int8

const (
	Admin Role = iota
	Gardener
)
