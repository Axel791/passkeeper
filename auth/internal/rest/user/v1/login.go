package v1

import (
	"github.com/Axel791/appkit"
	"github.com/Axel791/auth/internal/rest/user/v1/model"
	"github.com/Axel791/auth/internal/usecases/user"
	"net/http"
)

type Login struct {
	loginScenario user.Login
}

func NewLogin(loginScenario user.Login) *Login {
	return &Login{loginScenario: loginScenario}
}

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
	}

	appkit.WriteJSON(responseWriter, http.StatusOK, model.TokenResponse{AccessToken: token})
}
