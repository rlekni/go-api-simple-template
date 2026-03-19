package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func Auth0Middleware() func(http.Handler) http.Handler {
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")

	if auth0Domain == "" || auth0Audience == "" {
		// If not configured, we can't really validate.
		// For now, I'll return a middleware that does nothing or logs a warning.
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(auth0Audience, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := fmt.Sprintf("https://%s/", auth0Domain)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token, auth0Domain)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := jwtMiddleware.CheckJWT(w, r)
			if err != nil {
				return // Error response already sent by middleware
			}

			// Extract claims and put them in context
			userToken := r.Context().Value("user").(*jwt.Token)
			claims := userToken.Claims.(jwt.MapClaims)

			ctx := r.Context()
			// Auth0 standard claims: 'sub' is user id, 'email' might be there if scope allowed
			ctx = context.WithValue(ctx, "user_id", claims["sub"])
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, "email", email)
			} else if name, ok := claims["name"].(string); ok {
				// Fallback or handle based on how Auth0 is configured
				ctx = context.WithValue(ctx, "email", name)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getPemCert(token *jwt.Token, domain string) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", domain))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
