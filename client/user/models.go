package user

import (
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/model"
)

type CreateUserInput struct {
	Username    string
	FirstName   *string
	LastName    *string
	UserRoles   string
	ApiOnlyUser bool
}

type UpdateUserInput struct {
	Uid       string
	UserRoles []string
}

type DeleteUserInput struct {
	Uid string
}

type DeleteUserOutput struct{}

type RevokeApiTokenOutput struct{}

// CreateUser endpoint returns a user-tenant association for whatever reason
type UserTenantAssociation struct {
	Uid    string      `json:"uid"`
	Source Association `json:"source"`
}

type CreateUserOutput = model.UserDetails

type UpdateUserOutput = model.UserDetails

type ReadUserOutput = model.UserDetails

type Association struct {
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Uid       string `json:"uid"`
}

type ReadByUsernameInput struct {
	Name        string `json:"name"`
	ApiOnlyUser bool   `json:"apiOnlyUser"`
}

type GenerateApiTokenInput struct {
	Name string `json:"name"`
}

type RevokeApiTokenInput struct {
	Name        string `json:"name"`
	ApiOnlyUser bool   `json:"apiOnlyUser"`
}

type RevokeOAuthTokenInput struct {
	ApiTokenId string `json:"apiTokenId"`
}

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

type ApiTokenResponse struct {
	ApiToken string `json:"apiToken"`
}

func NewCreateUserInput(username string, userRoles string, apiOnlyUser bool, firstName *string, lastName *string) *CreateUserInput {
	return &CreateUserInput{
		Username:    username,
		UserRoles:   userRoles,
		ApiOnlyUser: apiOnlyUser,
		FirstName:   firstName,
		LastName:    lastName,
	}
}

func NewReadByUsernameInput(name string, apiOnlyUser bool) *ReadByUsernameInput {
	return &ReadByUsernameInput{
		Name:        name,
		ApiOnlyUser: apiOnlyUser,
	}
}

func NewGenerateApiTokenInput(name string) *GenerateApiTokenInput {
	return &GenerateApiTokenInput{
		Name: name,
	}
}

func NewRevokeApiTokenInput(name string) *RevokeApiTokenInput {
	return &RevokeApiTokenInput{
		Name: name,
	}
}

func NewRevokeOAuthTokenInput(apiTokenId string) *RevokeOAuthTokenInput {
	return &RevokeOAuthTokenInput{
		ApiTokenId: apiTokenId,
	}
}

func NewReadByUidInput(uid string) *ReadByUidInput {
	return &ReadByUidInput{
		Uid: uid,
	}
}

func NewUpdateByUidInput(uid string, userRoles []string) *UpdateUserInput {
	return &UpdateUserInput{
		Uid:       uid,
		UserRoles: userRoles,
	}
}
