package repository

import "testing"

func TestHash(t *testing.T) {
	t.Run("空でないハッシュ文字列を返す", func(t *testing.T) {
		h, err := Hash("password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h == "" {
			t.Error("expected non-empty hash")
		}
	})

	t.Run("同じ入力でも呼び出すたびに異なるハッシュが返る（bcryptのソルト）", func(t *testing.T) {
		h1, err := Hash("password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		h2, err := Hash("password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h1 == h2 {
			t.Error("expected different hashes for same input (bcrypt uses random salt)")
		}
	})

	t.Run("空文字列もハッシュできる", func(t *testing.T) {
		h, err := Hash("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h == "" {
			t.Error("expected non-empty hash")
		}
	})
}

func TestCompare(t *testing.T) {
	t.Run("正しいパスワード: true", func(t *testing.T) {
		h, err := Hash("correct_password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !Compare(h, "correct_password") {
			t.Error("expected Compare to return true for correct password")
		}
	})

	t.Run("誤ったパスワード: false", func(t *testing.T) {
		h, err := Hash("correct_password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if Compare(h, "wrong_password") {
			t.Error("expected Compare to return false for wrong password")
		}
	})

	t.Run("空文字列のパスワード: 正しい場合 true", func(t *testing.T) {
		h, err := Hash("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !Compare(h, "") {
			t.Error("expected Compare to return true for empty password")
		}
	})

	t.Run("空文字列のパスワード: 誤った場合 false", func(t *testing.T) {
		h, err := Hash("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if Compare(h, "not_empty") {
			t.Error("expected Compare to return false for non-empty against empty hash")
		}
	})

	t.Run("無効なハッシュ文字列: false", func(t *testing.T) {
		if Compare("not_a_valid_hash", "password") {
			t.Error("expected Compare to return false for invalid hash")
		}
	})
}
