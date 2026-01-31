package jwtguard

import (
	"net"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGuard struct {
	Secret string `json:"secret"`
}

func init() {
	caddy.RegisterModule(JWTGuard{})
}

func (JWTGuard) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.jwt_guard",
		New: func() caddy.Module { return new(JWTGuard) },
	}
}

func (j *JWTGuard) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "JWT missing", http.StatusUnauthorized)
		return nil
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid JWT", http.StatusUnauthorized)
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)

	if !claims.VerifyExpiresAt(time.Now(), true) {
		http.Error(w, "JWT expired", http.StatusUnauthorized)
		return nil
	}

	if claims["ip"] != ip {
		http.Error(w, "IP mismatch", http.StatusForbidden)
		return nil
	}

	return next.ServeHTTP(w, r)
}

func (j *JWTGuard) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&j.Secret) {
			return d.ArgErr()
		}
	}
	return nil
}