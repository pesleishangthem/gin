package middleware

type User struct {
	Sub       string   `json:"sub"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	Email     string   `json:"email"`
}
