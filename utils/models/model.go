package model

type (
	CreateUserReq struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Mobile    string `json:"mobile"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	GenericRes struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Meta    interface{} `json:"meta,omitempty"`
	}
	BotReq struct {
		Input string `json:"input"`
	}
)
