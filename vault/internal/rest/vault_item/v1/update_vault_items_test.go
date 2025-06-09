package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	userdomain "github.com/Axel791/vault/internal/domains/user"
	"github.com/Axel791/vault/internal/rest/vault_item/v1/model"
	"github.com/Axel791/vault/internal/usecases/vault_items/dto"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Axel791/appkit"
)

/*
   ──────────────────────────────────────────
   Заглушки интерфейсов
   ──────────────────────────────────────────
*/

type stubUpdateVaultItem struct {
	exec func(ctx context.Context, uid userdomain.UserID, in dto.VaultItemUpdateInput) error
}

func (s stubUpdateVaultItem) Execute(ctx context.Context, uid userdomain.UserID, in dto.VaultItemUpdateInput) error {
	if s.exec == nil {
		return errors.New("stub not configured")
	}
	return s.exec(ctx, uid, in)
}

/*
──────────────────────────────────────────
1. Невалидный JSON
──────────────────────────────────────────
*/
func TestUpdateVaultItem_InvalidBody(t *testing.T) {
	h := &UpdateVaultItem{}

	req := httptest.NewRequest(http.MethodPatch, "/vault", strings.NewReader("oops"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
2. Нет токена
──────────────────────────────────────────
*/
func TestUpdateVaultItem_MissingToken(t *testing.T) {
	body := `{"id":1,"groupId":1,"dataType":2,"encryptedBlob":"ZGI="}`
	h := &UpdateVaultItem{}

	req := httptest.NewRequest(http.MethodPatch, "/vault", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got %d, want 403", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "missing auth token") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
3. Просроченный / неверный токен
──────────────────────────────────────────
*/
func TestUpdateVaultItem_BadToken(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, errors.New("token error")
		},
	}

	body := `{"id":1,"groupId":1,"dataType":2,"encryptedBlob":"ZGI="}`
	h := &UpdateVaultItem{authService: authStub}

	req := httptest.NewRequest(http.MethodPatch, "/vault", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer bad")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got %d, want 403", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid or expired token") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
4. Неверный groupID
──────────────────────────────────────────
*/
func TestUpdateVaultItem_InvalidGroupID(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, nil
		},
	}

	reqBody := model.VaultItemUpdateRequest{
		ID:            1,
		GroupID:       0, // вызовет ошибку groupdomain.NewGroupID
		DataType:      2,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &UpdateVaultItem{authService: authStub}

	req := httptest.NewRequest(http.MethodPatch, "/vault", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer ok")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid group id") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
5. Неверный vaultID
──────────────────────────────────────────
*/
func TestUpdateVaultItem_InvalidVaultID(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, nil
		},
	}

	reqBody := model.VaultItemUpdateRequest{
		ID:            0, // vaultdomain.NewVaultID вернёт ошибку
		GroupID:       1,
		DataType:      2,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &UpdateVaultItem{authService: authStub}

	req := httptest.NewRequest(http.MethodPatch, "/vault", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer ok")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid vault id") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
6. Use-case возвращает бизнес-ошибку
──────────────────────────────────────────
*/
func TestUpdateVaultItem_UseCaseError(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, nil
		},
	}

	useCaseStub := stubUpdateVaultItem{
		exec: func(_ context.Context, _ userdomain.UserID, _ dto.VaultItemUpdateInput) error {
			return appkit.New(appkit.Conflict, "conflict")
		},
	}

	reqBody := model.VaultItemUpdateRequest{
		ID:            1,
		GroupID:       1,
		DataType:      2,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &UpdateVaultItem{
		authService:     authStub,
		updateVaultItem: useCaseStub,
	}

	req := httptest.NewRequest(http.MethodPatch, "/vault", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer ok")
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("got %d, want 409", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "conflict") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
──────────────────────────────────────────
7. Успешное обновление
──────────────────────────────────────────
*/
func TestUpdateVaultItem_Success(t *testing.T) {
	const (
		token   = "ok"
		vaultID = 1
		groupID = 1
	)

	authStub := stubAuthService{
		auth: func(_ context.Context, tk string) (userdomain.UserID, error) {
			if tk != token {
				t.Fatalf("token mismatch: %s", tk)
			}
			return userdomain.UserID{}, nil
		},
	}

	useCaseStub := stubUpdateVaultItem{
		exec: func(_ context.Context, _ userdomain.UserID, in dto.VaultItemUpdateInput) error {
			if in.VaultID.ToInt64() != vaultID || in.GroupID.ToInt64() != groupID {
				t.Fatalf("wrong IDs: %v / %v", in.VaultID, in.GroupID)
			}
			return nil
		},
	}

	reqBody := model.VaultItemUpdateRequest{
		ID:            vaultID,
		GroupID:       groupID,
		DataType:      2,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &UpdateVaultItem{
		authService:     authStub,
		updateVaultItem: useCaseStub,
	}

	req := httptest.NewRequest(http.MethodPatch, "/vault", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "vault item updated successfully") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
