package jwt_user

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"userService/services/user-service/config"
)

func GetSignInToken() (token, refresh string, err error) {
	token, err = GetJWTToken()
	if err != nil {
		return "", "", err
	}

	refresh, err = GetRefreshToken()
	if err != nil {
		return "", "", err
	}
	return token, refresh, err
}

func GetRefreshToken() (string, error) {
	conf := config.UserConfig()
	jwtSecretToken := conf.Refresh

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 300).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecretToken))

	return tokenString, err
}

func GetJWTToken() (string, error) {
	conf := config.UserConfig()
	jwtSecretToken := conf.Jwt

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecretToken))

	return tokenString, err
}

func ForwardRefresh() (token, refresh string, err error) {
	token, err = GetJWTToken()
	if err != nil {
		return "", "", err
	}

	refresh, err = GetRefreshToken()
	if err != nil {
		return
	}
	return token, refresh, err
}

func CheckJWTToken(tokenString string) (error, jwt.MapClaims) {
	conf := config.UserConfig()
	jwtSecretToken := conf.Jwt

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretToken), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil, claims
	} else {
		return err, nil
	}
}

func CheckRefreshToken(tokenString string) (error, jwt.MapClaims) {
	conf := config.UserConfig()
	jwtSecretToken := conf.Refresh

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretToken), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil, claims
	} else {
		return err, nil
	}
}
