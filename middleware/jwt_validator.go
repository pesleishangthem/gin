package middleware

import (
	"fmt"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pesleishangthem/gin/config"
)

type Validator struct {
	jwks *keyfunc.JWKS
}

func NewValidator(cfg config.ServerConfig) (*Validator, error) {

	jwksURL := fmt.Sprintf("%s/protocol/openid-connect/certs", cfg.GetKeycloakRealmURL())

	fmt.Printf("\nKeycloakRealmURL: %s\n", cfg.GetKeycloakRealmURL())

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}

	return &Validator{jwks: jwks}, nil
}

func (v *Validator) Validate(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, v.jwks.Keyfunc)
}
