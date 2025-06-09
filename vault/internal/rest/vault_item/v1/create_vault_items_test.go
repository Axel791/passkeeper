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
   ────────────────────────────────────────────
   Мини-заглушки интерфейсов
   ────────────────────────────────────────────
*/

// stubAuthService реализует services.Services.
type stubAuthService struct {
	auth func(ctx context.Context, token string) (userdomain.UserID, error)
}

func (s stubAuthService) AuthenticateToken(ctx context.Context, token string) (userdomain.UserID, error) {
	if s.auth == nil {
		return userdomain.UserID{}, errors.New("stub not configured")
	}
	return s.auth(ctx, token)
}

// stubCreateVaultItem реализует vault_items.CreateVaultItem.
type stubCreateVaultItem struct {
	exec func(ctx context.Context, uid userdomain.UserID, in dto.VaultItemInput) error
}

func (s stubCreateVaultItem) Execute(ctx context.Context, uid userdomain.UserID, in dto.VaultItemInput) error {
	if s.exec == nil {
		return errors.New("stub not configured")
	}
	return s.exec(ctx, uid, in)
}

/*
   ────────────────────────────────────────────
   1. Невалидное тело запроса
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_InvalidBody(t *testing.T) {
	h := &CreateVaultItemV1{} // ни auth, ни use-case не понадобятся

	req := httptest.NewRequest(http.MethodPost, "/vault", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────────────────
   2. Отсутствует заголовок Authorization
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_MissingToken(t *testing.T) {
	body := `{"groupId":"abc","dataType":"text","encryptedBlob":"blob"}`
	h := &CreateVaultItemV1{}

	req := httptest.NewRequest(http.MethodPost, "/vault", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusForbidden)
	}
	if !strings.Contains(rec.Body.String(), "missing auth token") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────────────────
   3. Токен не валиден / истёк
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_BadToken(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, errors.New("bad token")
		},
	}

	body := `{"groupId":"abc","dataType":"text","encryptedBlob":"blob"}`
	h := &CreateVaultItemV1{
		authService: authStub,
	}

	req := httptest.NewRequest(http.MethodPost, "/vault", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer bad-token")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusForbidden)
	}
	if !strings.Contains(rec.Body.String(), "invalid or expired token") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────────────────
   4. Некорректный groupID
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_InvalidGroupID(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, nil
		},
	}

	// groupId явно не UUID → groupdomain.NewGroupID вернёт ошибку
	reqBody := model.VaultItemCreateRequest{
		GroupID:       1,
		DataType:      2,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &CreateVaultItemV1{
		authService: authStub,
	}

	req := httptest.NewRequest(http.MethodPost, "/vault", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer good-token")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "invalid group id") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────────────────
   5. Ошибка use-case’a (пример: конфликт)
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_UseCaseError(t *testing.T) {
	authStub := stubAuthService{
		auth: func(_ context.Context, _ string) (userdomain.UserID, error) {
			return userdomain.UserID{}, nil
		},
	}

	useCaseStub := stubCreateVaultItem{
		exec: func(_ context.Context, _ userdomain.UserID, _ dto.VaultItemInput) error {
			return appkit.New(appkit.Conflict, "err conflict")
		},
	}

	reqBody := model.VaultItemCreateRequest{
		GroupID:       1,
		DataType:      3,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &CreateVaultItemV1{
		authService:     authStub,
		createVaultItem: useCaseStub,
	}

	req := httptest.NewRequest(http.MethodPost, "/vault", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer good-token")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusConflict)
	}
	if !strings.Contains(rec.Body.String(), "duplicate") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────────────────
   6. Успешное создание
   ────────────────────────────────────────────
*/

func TestCreateVaultItem_Success(t *testing.T) {
	const (
		groupID = 1
		token   = "good-token"
	)

	authStub := stubAuthService{
		auth: func(_ context.Context, tk string) (userdomain.UserID, error) {
			if tk != token {
				t.Fatalf("unexpected token: %s", tk)
			}
			return userdomain.UserID{}, nil
		},
	}

	useCaseStub := stubCreateVaultItem{
		exec: func(_ context.Context, _ userdomain.UserID, in dto.VaultItemInput) error {
			if in.GroupID.ToInt64() != groupID {
				t.Fatalf("wrong groupID: %s", in.GroupID)
			}
			if in.DataType != 1 {
				t.Fatalf("wrong datatype: %s", in.DataType)
			}
			return nil
		},
	}

	reqBody := model.VaultItemCreateRequest{
		GroupID:       groupID,
		DataType:      1,
		EncryptedBlob: []byte("blob"),
	}
	raw, _ := json.Marshal(reqBody)

	h := &CreateVaultItemV1{
		authService:     authStub,
		createVaultItem: useCaseStub,
	}

	req := httptest.NewRequest(http.MethodPost, "/vault", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusCreated)
	}
	if !strings.Contains(rec.Body.String(), "vault item created successfully") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
