package user

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type updateRequestBody struct {
	UserRoles []string `json:"roles"`
}

func NewCreateRequest(ctx context.Context, client http.Client, createInp CreateUserInput) *http.Request {
	body := model.PublicApiCreateUserInput{
		Name:        createInp.Username,
		FirstName:   createInp.FirstName,
		LastName:    createInp.LastName,
		Role:        createInp.UserRoles,
		ApiOnlyUser: createInp.ApiOnlyUser,
	}
	return client.NewPost(ctx, url.CreateUser(client.BaseUrl()), body)
}

func NewGenerateApiTokenRequest(ctx context.Context, client http.Client, userUid string) *http.Request {
	url := url.GenerateApiToken(client.BaseUrl(), userUid)
	return client.NewPost(ctx, url, nil)
}

func NewReadByUidRequest(ctx context.Context, client http.Client, uid string) *http.Request {
	url := url.ReadUserByUid(client.BaseUrl(), uid)
	return client.NewGet(ctx, url)
}

func NewReadByUsernameRequest(ctx context.Context, client http.Client, username string, apiOnlyUser bool) *http.Request {
	var readUrl string
	if apiOnlyUser {
		readUrl = url.ReadApiOnlyUserByUsername(client.BaseUrl())
	} else {
		readUrl = url.ReadUserByUsername(client.BaseUrl())
	}
	req := client.NewGet(ctx, readUrl)
	req.QueryParams.Add("q", fmt.Sprintf("name:%s", username))
	req.QueryParams.Add("limit", "1")
	req.QueryParams.Add("offset", "0")
	return req
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateInp UpdateUserInput) *http.Request {
	url := url.ReadUserByUid(client.BaseUrl(), updateInp.Uid)
	body := updateRequestBody{
		UserRoles: updateInp.UserRoles,
	}

	return client.NewPut(ctx, url, body)
}
