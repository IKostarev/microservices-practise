package pass_utils

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/argon2"
	"strings"
)

type PasswordConfig struct {
	Time    uint32 `envconfig:"PASS_TIME" required:"true" default:"1"`
	Memory  uint32 `envconfig:"PASS_MEMORY" required:"true" default:"65536"`
	Threads uint8  `envconfig:"PASS_THREADS" required:"true" default:"4"`
	KeyLen  uint32 `envconfig:"PASS_KEY_LEN" required:"true" default:"32"`
}

type PasswordUtils struct {
	cfg *PasswordConfig
}

func NewPasswordUtils(cfg *PasswordConfig) *PasswordUtils {
	return &PasswordUtils{
		cfg: cfg,
	}
}

// GeneratePassword создает пароль на основе библиотеки golang.org/x/crypto/argon2
func (u *PasswordUtils) GeneratePassword(ctx context.Context, password string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GeneratePassword")
	defer span.Finish()

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, u.cfg.Time, u.cfg.Memory, u.cfg.Threads, u.cfg.KeyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, u.cfg.Memory, u.cfg.Time, u.cfg.Threads, b64Salt, b64Hash)
	return full, nil
}

// ComparePassword сравниваем пароль и переданный хэш пароля на основе библиотеки golang.org/x/crypto/argon2
func (u *PasswordUtils) ComparePassword(ctx context.Context, password, hash string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ComparePassword")
	defer span.Finish()

	parts := strings.Split(hash, "$")

	c := &PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
