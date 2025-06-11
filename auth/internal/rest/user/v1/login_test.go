package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Axel791/auth/internal/rest/user/v1/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubLogin struct {
	exec func(ctx context.Context, email, password string) (string, error)
}

func (s stubLogin) Execute(ctx context.Context, email, password string) (string, error) {
	if s.exec == nil {
		return "", errors.New("stub not configured")
	}
	return s.exec(ctx, email, password)
}

func TestLogin_InvalidBody(t *testing.T) {
	h := NewLogin(stubLogin{}) // до use-case кода не дойдёт

	req := httptest.NewRequest(http.MethodPost, "/users/login",
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
   ────────────────────────────────
   2. Ошибка use-case’а
   ────────────────────────────────
   NB: в ServeHTTP ДОЛЖЕН быть return после WriteErrorJSON,
   иначе сюда долетит 200 OK, и тест завалится.
*/

func TestLogin_UseCaseError(t *testing.T) {
	stub := stubLogin{
		exec: func(_ context.Context, email, _ string) (string, error) {
			if email != "john@doe.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			return "", fmt.Errorf("err")
		},
	}

	h := NewLogin(stub)

	body := `{"email":"john@doe.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/users/login",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized { // ToHTTPCode(err) → 401
		t.Fatalf("got %d, want %d (add return after WriteErrorJSON?)",
			rec.Code, http.StatusUnauthorized)
	}

	if !strings.Contains(rec.Body.String(), "invalid credentials") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

/*
   ────────────────────────────────
   3. Успешный сценарий
   ────────────────────────────────
*/

func TestLogin_Success(t *testing.T) {
	const (
		email    = "john@doe.com"
		password = "secret"
		token    = "jwt-123"
	)

	stub := stubLogin{
		exec: func(_ context.Context, e, p string) (string, error) {
			if e != email || p != password {
				t.Fatalf("wrong args: %s / %s", e, p)
			}
			return token, nil
		},
	}

	h := NewLogin(stub)

	reqBody := model.LoginRequest{Email: email, Password: password}
	raw, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users/login",
		bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusOK)
	}

	var got model.TokenResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if got.AccessToken != token {
		t.Fatalf("token mismatch: %q vs %q", got.AccessToken, token)
	}
}
