// Package jwt provides JWT token generation and validation.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/secure"
)

// TypeBearerToken value "Bearer" for the token type field.
const TypeBearerToken string = "Bearer"

// ErrJwtSignMethod is an error when jwt signing method is wrong.
var ErrJwtSignMethod = errors.New("unexpected signing method")

// Config is a config for jwt module.
type Config struct {
	Secret     string `yaml:"secret"     env:"JWT_SECRET"`
	Encryption string `yaml:"encryption" env:"JWT_ENCRYPTION"`
	Expire     int    `yaml:"expire"     env:"JWT_EXPIRE"     env-default:"1"` // in hours
}

// Provider implements jwt token generation and validation.
type Provider struct {
	cfg         *Config
	secretBytes []byte
}

// New is a constructor for [Provider].
func New(cfg *Config) *Provider {
	return &Provider{
		cfg:         cfg,
		secretBytes: []byte(cfg.Secret),
	}
}

// CreateAccessToken generates new JWT.
func (j *Provider) CreateAccessToken(userID int64) (string, error) {
	key := []byte(j.cfg.Secret)

	jwtClaims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.cfg.Expire))),
		"sub": userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := accessToken.SignedString(key)
	if err != nil {
		return "", errs.Wrap(err, "sign token")
	}

	if j.cfg.Encryption != "" {
		return secure.Encrypt(signedToken, j.cfg.Encryption)
	}

	return signedToken, nil
}

// ParseAccessToken tries to parse JWT from incoming string.
func (j *Provider) ParseAccessToken(accessTokenString string) (*jwt.Token, error) {
	var err error

	if j.cfg.Encryption != "" {
		accessTokenString, err = secure.Decrypt(accessTokenString, j.cfg.Encryption)
		if err != nil {
			return nil, errs.Wrap(err, "decrypt token")
		}
	}

	accessToken, err := jwt.Parse(
		accessTokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errs.Wrap(ErrJwtSignMethod, "verify token signature")
			}
			return j.secretBytes, nil
		},
	)
	if err != nil {
		return nil, errs.Wrap(err, "parse token")
	}

	return accessToken, nil
}
