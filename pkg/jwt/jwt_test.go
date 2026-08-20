package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	pkgjwt "pob/pkg/jwt"

	libjwt "github.com/golang-jwt/jwt/v5"
)

// newTestKeyPair はテスト用のRSA鍵ペアを生成する。
func newTestKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return priv, &priv.PublicKey
}

// signToken はテスト用トークンを秘密鍵で署名して返す。
func signToken(t *testing.T, claims libjwt.Claims, priv *rsa.PrivateKey) string {
	t.Helper()
	token := libjwt.NewWithClaims(libjwt.SigningMethodRS256, claims)
	s, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return s
}

// ── NewClaims ─────────────────────────────────────────────────────────────────

func TestNewClaims(t *testing.T) {
	t.Run("UserID が設定される", func(t *testing.T) {
		c := pkgjwt.NewClaims("user-123", time.Minute)
		if c.UserID != "user-123" {
			t.Errorf("UserID = %q, want %q", c.UserID, "user-123")
		}
	})

	t.Run("ExpiresAt が duration 後に設定される", func(t *testing.T) {
		before := time.Now()
		c := pkgjwt.NewClaims("u", 15*time.Minute)
		after := time.Now()

		exp := c.ExpiresAt.Time
		if exp.Before(before.Add(14 * time.Minute)) {
			t.Errorf("ExpiresAt too early: %v", exp)
		}
		if exp.After(after.Add(16 * time.Minute)) {
			t.Errorf("ExpiresAt too late: %v", exp)
		}
	})

	t.Run("IssuedAt が現在時刻付近に設定される", func(t *testing.T) {
		// jwt.NumericDate は Unix 秒精度のため、比較は秒単位で行う
		before := time.Now().Truncate(time.Second)
		c := pkgjwt.NewClaims("u", time.Minute)
		after := time.Now().Add(time.Second).Truncate(time.Second)

		iat := c.IssuedAt.Time
		if iat.Before(before) || iat.After(after) {
			t.Errorf("IssuedAt = %v, expected between %v and %v", iat, before, after)
		}
	})
}

// ── VerifyToken ───────────────────────────────────────────────────────────────

func TestVerifyToken(t *testing.T) {
	priv, pub := newTestKeyPair(t)

	t.Run("有効なトークン: Claims を返す", func(t *testing.T) {
		claims := pkgjwt.NewClaims("user-abc", 15*time.Minute)
		token := signToken(t, claims, priv)

		got, err := pkgjwt.VerifyToken(token, pub)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.UserID != "user-abc" {
			t.Errorf("UserID = %q, want %q", got.UserID, "user-abc")
		}
	})

	t.Run("有効なトークン: error が nil", func(t *testing.T) {
		token := signToken(t, pkgjwt.NewClaims("u", 5*time.Minute), priv)
		_, err := pkgjwt.VerifyToken(token, pub)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("期限切れトークン: error を返す", func(t *testing.T) {
		claims := pkgjwt.NewClaims("u", -1*time.Minute) // 過去に失効
		token := signToken(t, claims, priv)

		_, err := pkgjwt.VerifyToken(token, pub)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("別の鍵で検証: error を返す", func(t *testing.T) {
		_, otherPub := newTestKeyPair(t)
		token := signToken(t, pkgjwt.NewClaims("u", 5*time.Minute), priv)

		_, err := pkgjwt.VerifyToken(token, otherPub)
		if err == nil {
			t.Error("expected error for wrong public key, got nil")
		}
	})

	t.Run("改ざんされたトークン: error を返す", func(t *testing.T) {
		_, err := pkgjwt.VerifyToken("invalid.token.string", pub)
		if err == nil {
			t.Error("expected error for tampered token, got nil")
		}
	})

	t.Run("空文字列: error を返す", func(t *testing.T) {
		_, err := pkgjwt.VerifyToken("", pub)
		if err == nil {
			t.Error("expected error for empty token, got nil")
		}
	})

	t.Run("HMAC署名のトークン（アルゴリズム不一致）: error を返す", func(t *testing.T) {
		hmacClaims := libjwt.MapClaims{
			"user_id": "u",
			"exp":     time.Now().Add(5 * time.Minute).Unix(),
		}
		hmacToken := libjwt.NewWithClaims(libjwt.SigningMethodHS256, hmacClaims)
		tokenStr, err := hmacToken.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("failed to sign HMAC token: %v", err)
		}

		_, err = pkgjwt.VerifyToken(tokenStr, pub)
		if err == nil {
			t.Error("expected error for HMAC token with RSA verifier, got nil")
		}
	})
}

// ── ExtractBearerToken ────────────────────────────────────────────────────────

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"Bearer プレフィックスあり: 取り除く", "Bearer mytoken123", "mytoken123"},
		{"Bearer プレフィックスなし: そのまま返す", "mytoken123", "mytoken123"},
		{"空文字列: 空文字列を返す", "", ""},
		{"Bearer のみ: 空文字列を返す", "Bearer ", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkgjwt.ExtractBearerToken(tt.header)
			if got != tt.want {
				t.Errorf("ExtractBearerToken(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}
