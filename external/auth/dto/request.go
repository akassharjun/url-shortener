package dto

func MakeAuthenticateRequest(userId string, auth string) (map[string]string, AuthenticatePayload) {
	headers := map[string]string{}

	payload := AuthenticatePayload{
		UserId: userId,
		Auth:   auth,
	}

	return headers, payload
}
