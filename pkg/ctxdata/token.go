package ctxdata

import "github.com/golang-jwt/jwt"

const Identify = "zero-im"

func GetJwtToken(key string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims[Identify] = uid
	claims["iat"] = iat
	claims["exp"] = iat + seconds

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(key))
}
