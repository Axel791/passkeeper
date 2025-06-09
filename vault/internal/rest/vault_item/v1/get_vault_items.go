package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Axel791/appkit"
	groupdomain "github.com/Axel791/vault/internal/domains/group"
	"github.com/Axel791/vault/internal/rest/vault_item/v1/model"
	"github.com/Axel791/vault/internal/services"
	"github.com/Axel791/vault/internal/usecases/vault_items"
	"github.com/go-chi/chi/v5"
)

// GetVaultItem обрабатывает HTTP-запросы на получение элементов хранилища для группы
type GetVaultItem struct {
	authService   services.Services
	getVaultItems vault_items.GetVaultItems
}

// NewGetVaultItem создаёт новый обработчик GetVaultItem
func NewGetVaultItem(authService services.Services, getVaultItems vault_items.GetVaultItems) *GetVaultItem {
	return &GetVaultItem{
		authService:   authService,
		getVaultItems: getVaultItems,
	}
}

// ServeHTTP обрабатывает GET-запросы на /vault/{groupID}
func (h *GetVaultItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "groupID")
	groupIDInt, err := strconv.ParseInt(groupIDStr, 10, 64)

	if err != nil || groupIDInt <= 0 {
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid group id"))

		return
	}

	groupID, err := groupdomain.NewGroupID(groupIDInt)
	if err != nil {
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid group id"))

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

	items, err := h.getVaultItems.Execute(r.Context(), userID, groupID)
	if err != nil {
		appkit.WriteErrorJSON(w, appkit.ToHTTPCode(err))

		return
	}

	resp := make([]model.VaultItemResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, model.VaultItemResponse{
			ID:            it.ID,
			GroupID:       it.GroupID,
			DataType:      it.DataType,
			EncryptedBlob: it.EncryptedBlob,
			CreatedAt:     it.CreatedAt,
			UpdatedAt:     it.UpdatedAt,
		})
	}

	appkit.WriteJSON(w, http.StatusOK, resp)
}
