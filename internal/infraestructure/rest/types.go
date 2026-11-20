package rest

type accessDTO struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
}
