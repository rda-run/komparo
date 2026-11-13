package config

import (
	"strings"
	"testing"
)

func TestRedactConnectionString(t *testing.T) {
	for _, in := range []string{
		"postgres://user:s3cret@host:5432/db?sslmode=disable",
		"postgresql://user:s3cret@host:5432/db",
		"host=h port=5432 user=u dbname=d password=s3cret sslmode=disable",
		`host=h password='s3cret with space' sslmode=disable`,
		"postgres://user@host/db?password=s3cret&sslmode=disable",
	} {
		out := RedactConnectionString(in)
		if strings.Contains(out, "s3cret") || strings.Contains(out, "space") {
			t.Errorf("password leaked: %q -> %q", in, out)
		}
		if strings.Contains(in, "sslmode") && !strings.Contains(out, "sslmode=disable") {
			t.Errorf("redaction ate other parameters: %q -> %q", in, out)
		}
	}

	if in := "postgres://host:5432/db"; RedactConnectionString(in) != in {
		t.Errorf("URL without password changed: %q", RedactConnectionString(in))
	}
}
