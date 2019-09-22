package jwt

import (
	"encoding/json"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type jwtPayload struct {
	AppID int    `json:"uid"`
	Value string `json:"name"`
	jwt.StandardClaims
}

type jwtConfig struct {
	Secret string
	Expire time.Duration
}

func Config(secret string) jwtConfig {
	return jwtConfig{
		Secret: secret,
		Expire: 0,
	}
}

func ConfigWithExpire(secret string, expire time.Duration) jwtConfig {
	return jwtConfig{
		Secret: secret,
		Expire: expire,
	}
}

func (config jwtConfig) Encode(payload interface{}) string {
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	sc := jwt.StandardClaims{}
	if config.Expire != 0 {
		sc = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Expire).Unix(),
		}
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwtPayload{
		AppID:          1,
		Value:          string(b),
		StandardClaims: sc,
	})
	tokenstring, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		panic(err)
	}
	return tokenstring
}

func (config jwtConfig) Decode(tokenstring string) ([]byte, error) {
	payload := jwtPayload{}
	token, err := jwt.ParseWithClaims(tokenstring, &payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})
	if err != nil {
		return []byte{}, err
	}
	if !token.Valid {
		return []byte{}, errors.New("invalid token")
	}
	return []byte(payload.Value), nil
}
