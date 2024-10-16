package utils

import "testing"

func TestVerifyAccessToken(t *testing.T) {
	t.Cleanup(func() {
		JwtAuth = JwtData{}
	})

	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "issuer")
	InitJwt()

	providedAccountId := 1
	token, err := CreateAccessToken(providedAccountId)
	if err != nil {
		t.Fatalf("Failed creating access token")
	}

	accountId, ok := VerifyAccessToken(token)
	if !ok {
		t.Fatalf("Failed to verify access token")
	}

	if accountId != providedAccountId {
		t.Fatalf("Access token's account id does not match the provided account id")
	}
}

func TestVerifyRefreshToken(t *testing.T) {
	t.Cleanup(func() {
		JwtAuth = JwtData{}
	})

	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "issuer")
	InitJwt()

	providedAccountId := 1
	providedVersion := 1
	token, err := CreateRefreshToken(providedAccountId, providedVersion)
	if err != nil {
		t.Fatalf("Failed creating refresh token")
	}

	accountId, version, ok := VerifyRefreshToken(token)
	if !ok {
		t.Fatalf("Failed to verify access token")
	}

	if accountId != providedAccountId {
		t.Fatalf("Refresh token's account id does not match the provided account id")
	}

	if version != providedVersion {
		t.Fatalf("Refresh token's version does not match the provided version")
	}
}
