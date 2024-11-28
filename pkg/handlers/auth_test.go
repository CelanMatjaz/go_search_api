package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/handlers"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	testcommon "github.com/CelanMatjaz/job_application_tracker_api/pkg/test_common"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestRegisterHandler(t *testing.T) {
	store := testcommon.CreateStore(t)
	authHandler := handlers.NewAuthHandler(store)
	path := "/api/auth/register"
	handler := authHandler.HandleRegister

	t.Run("check for correct body", func(t *testing.T) {
		body := types.RegisterBody{}
		res, req := newRequestAndRecorder(t, http.MethodPost, path, body)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.PasswordsDoNotMatch)

		t.Cleanup(func() {
			testcommon.ResetTables(store)
		})
	})

	body := types.RegisterBody{
		DisplayName:    "Display name",
		Email:          "test@test.test",
		Password:       "Password1!",
		PasswordVerify: "Password1!",
	}

	t.Run("test for matching passwords", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, body)
		err := handler(res, req)
		testcommon.AssertNotError(t, err, "")

		t.Cleanup(func() {
			testcommon.ResetTables(store)
		})
	})

	t.Run("test for mismatching passwords", func(t *testing.T) {
		newBody := body
		newBody.Email = "matching.password@test.test"
		newBody.Password = "Password1!!"

		res, req := newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.PasswordsDoNotMatch)

		t.Cleanup(func() {
			testcommon.ResetTables(store)
		})
	})

	t.Run("test for non existing account", func(t *testing.T) {
		newBody := body
		newBody.Email = "test1@test.test"

		t.Log("registering new account")
		res, req := newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err := handler(res, req)
		testcommon.AssertNotError(t, err, "")

		t.Log("registering new account with same body")
		res, req = newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err = handler(res, req)
		testcommon.AssertError(t, err, types.AccountAlreadyExists)

		t.Cleanup(func() {
			testcommon.ResetTables(store)
		})
	})

	t.Run("generic test", func(t *testing.T) {
		newBody := body

		res, req := newRequestAndRecorder(t, http.MethodPost, path, body)
		err := handler(res, req)
		testcommon.AssertNotError(t, err, "")

		account, exists, _ := store.GetAccountByEmail(newBody.Email)
		testcommon.Assert(t, exists, "account does not exist")
		testcommon.Assert(t, account.DisplayName == newBody.DisplayName, "account display name (%s) is not correct %s", account.DisplayName, newBody.DisplayName)
		testcommon.Assert(t, account.Email == newBody.Email, "account email (%s) is not correct %s", account.Email, newBody.Email)
		testcommon.Assert(t, !account.IsOauth, "account is oauth when it shouldn't be")

		t.Cleanup(func() {
			testcommon.ResetTables(store)
		})
	})
}

func TestLoginHandler(t *testing.T) {
	store := testcommon.CreateStore(t)
	authHandler := handlers.NewAuthHandler(store)
	path := "/api/auth/login"
	handler := authHandler.HandleLogin
	account, password := testcommon.SeedAccount(t, store)

	body := types.LoginBody{
		Email:    account.Email,
		Password: password,
	}

	t.Run("test with invalid body", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, "test")
		err := handler(res, req)
		testcommon.AssertError(t, err, types.InvalidJsonBody)
	})

	t.Run("test with valid body and correct creadentials", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, body)
		err := handler(res, req)
		testcommon.AssertNotError(t, err, "")
	})

	t.Run("test with valid body with wrong email", func(t *testing.T) {
		newBody := body
		newBody.Email = "email@email.com"
		res, req := newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.AccountDoesNotExist)
	})

	t.Run("test with valid body with wrong password", func(t *testing.T) {
		newBody := body
		newBody.Password = "wrong password"
		res, req := newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.InvalidPassword)
	})

	t.Run("test with valid body with wrong email and password", func(t *testing.T) {
		newBody := body
		newBody.Email = "email@email.com"
		res, req := newRequestAndRecorder(t, http.MethodPost, path, newBody)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.AccountDoesNotExist)
	})
}

func TestLogoutHandler(t *testing.T) {
	store := testcommon.CreateStore(t)
	authHandler := handlers.NewAuthHandler(store)
	path := "/api/auth/logout"
	handler := authHandler.HandleLogout

	t.Run("test without cookies", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, nil)
		err := handler(res, req)
		testcommon.Assert(t, err == nil, "unexpected error")
	})

	t.Run("test with cookies", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, nil)
		req.AddCookie(&http.Cookie{Name: "access_token"})
		req.AddCookie(&http.Cookie{Name: "refresh_token"})
		err := handler(res, req)
		testcommon.Assert(t, err == nil, "unexpected error")

		accessTokenCookie := getCookie(res, "access_token")
		testcommon.Assert(t, accessTokenCookie == nil || accessTokenCookie.Value == "", "access_token cookie not invalidated")
		refreshTokenCookie := getCookie(res, "refresh_token")
		testcommon.Assert(t, refreshTokenCookie == nil || refreshTokenCookie.Value == "", "refresh_token cookie not invalidated")
	})
}

func TestAuthCheckHandler(t *testing.T) {
	store := testcommon.CreateStore(t)
	authHandler := handlers.NewAuthHandler(store)
	path := "/api/auth/check"
	handler := authHandler.HandleAuthCheck

	t.Run("test with wrong account id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, nil)
		req = newRequestWithContextValues(req, middleware.AccountIdKey, 9999)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.AccountDoesNotExist)
	})

	account, _ := testcommon.SeedAccount(t, store)

	t.Run("test with correct account id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, nil)
		req = newRequestWithContextValues(req, middleware.AccountIdKey, account.Id)
		err := handler(res, req)
		testcommon.AssertNotError(t, err, "")
	})

	t.Run("test without account id", func(t *testing.T) {
		res, req := newRequestAndRecorder(t, http.MethodPost, path, nil)
		err := handler(res, req)
		testcommon.AssertError(t, err, types.Unauthenticated)
	})
}

func getCookie(res *httptest.ResponseRecorder, name string) *http.Cookie {
	result := res.Result()
	cookies := result.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

func cookieValid(res *httptest.ResponseRecorder, name string) bool {
	cookie := getCookie(res, name)
	return cookie != nil
}
