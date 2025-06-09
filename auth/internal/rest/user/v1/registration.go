package v1

import (
	"net/http"

	"github.com/Axel791/appkit"
	"github.com/Axel791/auth/internal/rest/user/v1/model"
	"github.com/Axel791/auth/internal/usecases/user"
)

// Register обрабатывает HTTP-запросы для регистрации нового пользователя
type Register struct {
	registration user.RegistrationUseCase
}

// NewRegister создаёт новый экземпляр Register с переданным сценарием регистрации
func NewRegister(registration user.RegistrationUseCase) *Register {
	return &Register{
		registration: registration,
	}
}

// ServeHTTP обрабатывает POST-запросы на /register
func (h *Register) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var requestModel model.RegistrationRequest
	err := appkit.ReadFromBodyAndUnmarshalToModelJSON(request.Body, &requestModel)
	if err != nil {
		appkit.WriteErrorJSON(responseWriter, appkit.BadRequestError("invalid body"))

		return
	}

	err = h.registration.Execute(request.Context(), requestModel.Email, requestModel.Password)
	if err != nil {
		appkit.WriteErrorJSON(responseWriter, appkit.ToHTTPCode(err))

		return
	}

	appkit.WriteJSON(responseWriter, http.StatusCreated, "user created successfully")
}
