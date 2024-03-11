package auth

import (
	"context"
	"fmt"
	"url-shortener/external/auth/dto"
	"url-shortener/lib"
	"url-shortener/src/helper"
)

type authService struct {
	HttpClient lib.HttpClient
}

func NewAuthService(httpClient lib.HttpClient) IAuthService {
	return authService{
		HttpClient: httpClient,
	}
}

type IAuthService interface {
	Authenticate(ctx context.Context, headers map[string]string, payload map[string]interface{}, serviceConfig map[string]interface{}) (*dto.AuthenticateResponse, error)
}

func (a authService) Authenticate(ctx context.Context, headers map[string]string, payload map[string]interface{}, serviceConfig map[string]interface{}) (*dto.AuthenticateResponse, error) {

	url := ""

	for _, v := range payload {
		url = fmt.Sprintf("https://AUTH-SERVICE/%s", v)
	}

	result, err := a.HttpClient.Get(url, headers, nil)

	if err != nil {
		return nil, err
	}

	var authResponse dto.AuthenticateResponse

	helper.ConvertMapToStruct(result, &authResponse)

	return &authResponse, nil
}
