package auth

import "testing"

func TestHashAndVerify(t *testing.T) {
	h, err := HashPassword("correct horse battery")
	if err != nil {
		t.Fatal(err)
	}
	if !VerifyPassword(h, "correct horse battery") {
		t.Fatal("valid password rejected")
	}
	if VerifyPassword(h, "wrong password!!") {
		t.Fatal("wrong password accepted")
	}
	h2, _ := HashPassword("correct horse battery")
	if h == h2 {
		t.Fatal("hashes must be salted")
	}
}

func TestPasswordPolicy(t *testing.T) {
	if _, err := HashPassword("short"); err != ErrPasswordPolicy {
		t.Fatalf("want policy error, got %v", err)
	}
}

func TestVerifyRejectsGarbage(t *testing.T) {
	for _, s := range []string{"", "$argon2id$", "$bcrypt$v=1$x$y$z", "$argon2id$v=19$m=x$a$b"} {
		if VerifyPassword(s, "whatever12345") {
			t.Fatalf("accepted %q", s)
		}
	}
}

func TestHasRole(t *testing.T) {
	if !hasRole("admin", []string{"admin", "editor"}) || hasRole("editor", []string{"admin"}) {
		t.Fatal("role check wrong")
	}
}
