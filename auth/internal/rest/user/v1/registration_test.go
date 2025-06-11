// register_handler_test.go
package v1 // ← замените, если хэндлер лежит в другом пакете

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Axel791/auth/internal/rest/user/v1/model"
)

type stubRegistration struct {
	exec func(ctx context.Context, email, password string) error
}

func (s stubRegistration) Execute(ctx context.Context, email, password string) error {
	if s.exec == nil {
		return errors.New("stub not configured")
	}
	return s.exec(ctx, email, password)
}

/*
   ─────────────────────────────────────────────
   1. Невалидный JSON
   ─────────────────────────────────────────────
*/

func TestRegister_InvalidBody(t *testing.T) {
	h := NewRegister(stubRegistration{}) // до use-case кода не дойдёт

	req := httptest.NewRequest(http.MethodPost, "/users/registration",
		strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusBadRequest)
	}

	want := `{"error":"invalid body"}`
	if strings.TrimSpace(rec.Body.String()) != want {
		t.Fatalf("body mismatch: %s", rec.Body.String())
	}
}

/*
   ─────────────────────────────────────────────
   2. Ошибка use-case’a (пример: e-mail уже занят)
   ─────────────────────────────────────────────
*/

var errEmailExists = errors.New("email already exists") // имитация бизнес-ошибки

func TestRegister_UseCaseError(t *testing.T) {
	stub := stubRegistration{
		exec: func(_ context.Context, email, _ string) error {
			if email != "john@doe.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			return errEmailExists
		},
	}

	h := NewRegister(stub)

	body := `{"email":"john@doe.com","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/users/registration",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	// Предполагаем, что appkit.ToHTTPCode(errEmailExists) → 409 Conflict.
	// Если маппите иначе — поменяйте число.
	if rec.Code != http.StatusConflict {
		t.Fatalf("got %d, want %d (проверьте ToHTTPCode/return)", rec.Code, http.StatusConflict)
	}

	if !strings.Contains(rec.Body.String(), "email already exists") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ─────────────────────────────────────────────
   3. Успех (пользователь создан)
   ─────────────────────────────────────────────
*/

func TestRegister_Success(t *testing.T) {
	const (
		email    = "john@doe.com"
		password = "secret"
	)

	stub := stubRegistration{
		exec: func(_ context.Context, e, p string) error {
			if e != email || p != password {
				t.Fatalf("wrong args: %s / %s", e, p)
			}
			return nil
		},
	}

	h := NewRegister(stub)

	reqBody := model.RegistrationRequest{Email: email, Password: password}
	raw, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/registration",
		bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusCreated)
	}

	// appkit.WriteJSON пишет строку как JSON: "user created successfully"
	if !strings.Contains(rec.Body.String(), "user created successfully") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
