package v1

import (
	"net/http"
	"strings"

	"github.com/Axel791/appkit"
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	vaultdomain "github.com/Axel791/vault/internal/domains/vault_item"
	"github.com/Axel791/vault/internal/rest/vault_item/v1/model"
	"github.com/Axel791/vault/internal/services"
	"github.com/Axel791/vault/internal/usecases/vault_items"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
)

// UpdateVaultItem обрабатывает HTTP-запросы на обновление существующего элемента хранилища (vault item)
type UpdateVaultItem struct {
	authService     services.Services
	updateVaultItem vault_items.UpdateVaultItem
}

// NewUpdateVaultItem создаёт новый обработчик UpdateVaultItem
func NewUpdateVaultItem(authService services.Services, updateVaultItem vault_items.UpdateVaultItem) *UpdateVaultItem {
	return &UpdateVaultItem{
		authService:     authService,
		updateVaultItem: updateVaultItem,
	}
}

// ServeHTTP обрабатывает PATCH запросы на обновление хранилища
func (h *UpdateVaultItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var reqModel model.VaultItemUpdateRequest
	if err := appkit.ReadFromBodyAndUnmarshalToModelJSON(r.Body, &reqModel); err != nil {
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid body"))

		return
	}

	token := strings.TrimSpace(r.Header.Get("Authorization"))
	const bearer = "Bearer "
	if strings.HasPrefix(token, bearer) {
		token = token[len(bearer):]
	}
	if token == "" {
		appkit.WriteErrorJSON(w, appkit.ForbiddenError("missing auth token"))

		return
	}

	userID, err := h.authService.AuthenticateToken(r.Context(), token)
	if err != nil {
		appkit.WriteErrorJSON(w, appkit.ForbiddenError("invalid or expired token"))

		return
	}

	groupID, err := groupdomain.NewGroupID(reqModel.GroupID)
	if err != nil {
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid group id"))

		return
	}

	vaultID, err := vaultdomain.NewVaultID(reqModel.ID)
	if err != nil {
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid vault id"))

		return
	}

	input := dto.VaultItemUpdateInput{
		VaultID:       vaultID,
		GroupID:       groupID,
		DataType:      reqModel.DataType,
		EncryptedBlob: reqModel.EncryptedBlob,
	}

	if err = h.updateVaultItem.Execute(r.Context(), userID, input); err != nil {
		appkit.WriteErrorJSON(w, appkit.ToHTTPCode(err))
		return
	}

	appkit.WriteJSON(w, http.StatusOK, "vault item updated successfully")
}
