package errors

type ApiProvider struct {
	Error
}

func (a ApiProvider) GetError(message string, status int) ApiProvider {
	return ApiProvider{
		Error{
			Code:       "LLM_PROVIDER_ERROR",
			Message:    message,
			Details:    GetDetails(),
			StatusCode: status,
		},
	}
}
