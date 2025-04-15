package model

type UserPage struct {
	Count  int           `json:"count"`
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Items  []UserDetails `json:"items"`
}

type UserDetails struct {
	Uid         string   `json:"uid"`
	Username    string   `json:"name"`
	Roles       []string `json:"roles"`
	ApiOnlyUser bool     `json:"apiOnlyUser"`
}

type PublicApiCreateUserInput struct {
	Name        string  `json:"name"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Role        string  `json:"role"`
	ApiOnlyUser bool    `json:"apiOnlyUser"`
}
