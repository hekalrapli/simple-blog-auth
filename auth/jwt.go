package auth

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)



func Generate(userId int64,jwtSecret string) (string,error ){
	key := []byte(jwtSecret)
	expired := time.Now().Add(time.Minute * 15).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId" : userId,
		"exp" : expired,
	})
	s,err := t.SignedString(key)
	if err != nil {
		return "",err
		
	}
	return s,nil
}

func Verify(tokenString, jwtSecret string) (int64, error) {
	
	token,err := jwt.Parse(tokenString, func(token *jwt.Token) (any,error) {
		secretKey := []byte(jwtSecret)
		return secretKey,nil
	},jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, err
		
	}
	if token == nil || !token.Valid {
		return 0, errors.New("Token tidak valid")
	}
	if claims,ok := token.Claims.(jwt.MapClaims); ok {
		if val, ok := claims["userId"]; ok {
			f, ok := val.(float64)
			if ok {
				return int64(f),nil
			}
		}
	}
	return 0, errors.New("Claims tidak valid")
}