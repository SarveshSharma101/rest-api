package middleware

import (
	"fmt"
	"net/http"
	"rest-api/rest-api/datamodels"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "my-secret"

func GenerateTokens(userId, role string) (*datamodels.Tokens, error) {
	jwtclaim := datamodels.Claim{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-rest-api-backend",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshclaim := datamodels.Claim{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-rest-api-backend",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaim).SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaim).SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &datamodels.Tokens{
		Jwt:     jwtToken,
		Refresh: refreshToken,
	}, nil
}

func ParseToken(toeknStr string) (*datamodels.Claim, error) {
	cliam := datamodels.Claim{}
	token, err := jwt.ParseWithClaims(
		toeknStr,
		&cliam,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*datamodels.Claim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}
	return claim, nil
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		fmt.Println("------------------")
		if header == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"err": "Authorization token not provided in header",
				},
			)
			return
		}

		auth := strings.Split(header, " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"err": "Authorisation header doesn't have bearer string",
				},
			)
			return
		}

		if _, err := ParseToken(auth[1]); err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"err": err.Error(),
				},
			)
			return
		}

		c.Next()
	}
}
