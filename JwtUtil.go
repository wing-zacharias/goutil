package util

import (
	"github.com/golang-jwt/jwt"
	jwtMidWare "github.com/iris-contrib/middleware/jwt"
)

// JwtToken
//可自行加入验证内容:id,email,...
//Iat time.Now().Unix()
//Exp time.Now().Add(time.Minute * time.Duration(global.GlobalConfig.Common.AuthExpire)).Unix()
type JwtToken struct {
	LoginUser string `json:"loginUser"`
	Iss       string `json:"iss"`
	Iat       int64  `json:"iat"`
	Exp       int64  `json:"exp"`
}

func GenJwtToken(jwtToken *JwtToken, secret string) (string, error) {
	token := jwtMidWare.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"loginUser": jwtToken.LoginUser,
		"iss":       jwtToken.Iss,
		"iat":       jwtToken.Iat,
		"exp":       jwtToken.Exp,
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(token string, secret string) (jt *JwtToken, err error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	jwtToken := &JwtToken{
		LoginUser: claim.Claims.(jwt.MapClaims)["loginUser"].(string),
		Iss:       claim.Claims.(jwt.MapClaims)["iss"].(string),
		Iat:       claim.Claims.(jwt.MapClaims)["iat"].(int64),
		Exp:       claim.Claims.(jwt.MapClaims)["exp"].(int64),
	}
	return jwtToken, nil
}
