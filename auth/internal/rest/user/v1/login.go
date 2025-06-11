package v1

import (
	"net/http"

	"github.com/Axel791/appkit"
	"github.com/Axel791/auth/internal/rest/user/v1/model"
	"github.com/Axel791/auth/internal/usecases/user"
)

// Login обрабатывает HTTP-запросы для входа пользователя
type Login struct {
	loginScenario user.Login
}

// NewLogin создаёт новый экземпляр Login с переданным сценарием входа
func NewLogin(loginScenario user.Login) *Login {
	return &Login{loginScenario: loginScenario}
}

// ServeHTTP обрабатывает POST-запросы на /login
func (h *Login) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var requestModel model.LoginRequest
	err := appkit.ReadFromBodyAndUnmarshalToModelJSON(request.Body, &requestModel)
	if err != nil {
		appkit.WriteErrorJSON(responseWriter, appkit.BadRequestError("invalid body"))

		return
	}

	token, err := h.loginScenario.Execute(
		request.Context(),
		requestModel.Email,
		requestModel.Password,
	)

	if err != nil {
		appkit.WriteErrorJSON(responseWriter, appkit.ToHTTPCode(err))

		return
	}

	appkit.WriteJSON(responseWriter, http.StatusOK, model.TokenResponse{AccessToken: token})
}
