package jwt

import (
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/config"
)

const (
	Dir            = "jwt"
	PrivateKeyFile = "private.key"
	PublicKeyFile  = "public.key"
)

type Payload struct {
	ID          string   `json:"id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

type Jwt interface {
	Decode(token string) (*Payload, error)
}

type jwtService struct {
	cfg       *config.Config
	publicKey *rsa.PublicKey
}

type Claims struct {
	jwt.RegisteredClaims
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       []string `json:"scope,omitempty"`
}

func NewJWT(cfg *config.Config) (Jwt, error) {
	publicKey, err := loadPublicKey(cfg)
	if err != nil {
		return nil, err
	}

	return &jwtService{
		cfg:       cfg,
		publicKey: publicKey,
	}, nil
}

func (j *jwtService) Decode(token string) (*Payload, error) {
	claims := &Claims{}

	result, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return false, errors.ErrInvalidSigningMethod
			}
			return j.publicKey, nil
		})

	if err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, errors.ErrInvalidToken
	}

	return &Payload{
		ID:          claims.ID,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		Scope:       claims.Scope,
	}, nil
}

func loadPrivateKey(cfg *config.Config) (*rsa.PrivateKey, error) {
	filePath := filepath.Join(cfg.CertPath, Dir, PrivateKeyFile)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func loadPublicKey(cfg *config.Config) (*rsa.PublicKey, error) {
	filePath := filepath.Join(cfg.CertPath, Dir, PublicKeyFile)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
