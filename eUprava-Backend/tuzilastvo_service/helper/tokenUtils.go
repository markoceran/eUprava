package helper

import (
	"errors"
	"github.com/cristalhq/jwt/v4"
	"log"
	"net/http"
	"os"
	"strings"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

var verifier, _ = jwt.NewVerifierHS(jwt.HS256, jwtKey)

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return token, nil
}

func ExtractUserType(r *http.Request) (string, error) {
	claims := ExtractClaims(r)
	role, ok := claims["rola"]
	if !ok {
		return "", errors.New("role claim not found or not a string")
	}
	return role, nil
}

func ExtractClaims(r *http.Request) map[string]string {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return nil
	}

	bearerToken := strings.Split(bearer, "Bearer ")
	if len(bearerToken) != 2 {
		return nil
	}

	tokenString := bearerToken[1]
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil
	}

	var claims map[string]string

	err = jwt.ParseClaims(token.Bytes(), verifier, &claims)
	if err != nil {
		log.Println(err)
	}

	return claims
}
