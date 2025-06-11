package v1

import (
	"net/http"
	"strings"

	"github.com/Axel791/appkit"
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	"github.com/Axel791/vault/internal/rest/vault_item/v1/model"
	"github.com/Axel791/vault/internal/services"
	"github.com/Axel791/vault/internal/usecases/vault_items"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
)

// CreateVaultItemV1 обрабатывает HTTP-запросы на создание нового элемента хранилища (vault item)
type CreateVaultItemV1 struct {
	authService     services.Services
	createVaultItem vault_items.CreateVaultItem
}

// NewCreateVaultItemV1 создаёт новый обработчик CreateVaultItemV1
func NewCreateVaultItemV1(authService services.Services, createVaultItem vault_items.CreateVaultItem) *CreateVaultItemV1 {
	return &CreateVaultItemV1{
		authService:     authService,
		createVaultItem: createVaultItem,
	}
}

// ServeHTTP обрабатывает POST-запросы на создание элемента хранилища
func (h *CreateVaultItemV1) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var reqModel model.VaultItemCreateRequest
	if err := appkit.ReadFromBodyAndUnmarshalToModelJSON(req.Body, &reqModel); err != nil {
		appkit.WriteErrorJSON(resp, appkit.BadRequestError("invalid body"))
		return
	}

	token := strings.TrimSpace(req.Header.Get("Authorization"))
	const bearer = "Bearer "
	if strings.HasPrefix(token, bearer) {
		token = token[len(bearer):]
	}

	if token == "" {
		appkit.WriteErrorJSON(resp, appkit.ForbiddenError("missing auth token"))

		return
	}

	userID, err := h.authService.AuthenticateToken(req.Context(), token)
	if err != nil {
		appkit.WriteErrorJSON(resp, appkit.ForbiddenError("invalid or expired token"))

		return
	}

	groupID, err := groupdomain.NewGroupID(reqModel.GroupID)
	if err != nil {
		appkit.WriteErrorJSON(resp, appkit.BadRequestError("invalid group id"))
	}

	input := dto.VaultItemInput{
		GroupID:       groupID,
		DataType:      reqModel.DataType,
		EncryptedBlob: reqModel.EncryptedBlob,
	}

	if err = h.createVaultItem.Execute(req.Context(), userID, input); err != nil {
		appkit.WriteErrorJSON(resp, appkit.ToHTTPCode(err))

		return
	}

	appkit.WriteJSON(resp, http.StatusCreated, "vault item created successfully")
}
