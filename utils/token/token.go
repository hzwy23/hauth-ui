package token

import (
	"github.com/astaxie/beego/context"

	"github.com/hzwy23/hauth/token/hjwt"
)

func GetJwtClaims(ctx *context.Context) (*hjwt.JwtClaims, error) {
	cookie, err := ctx.Request.Cookie("Authorization")
	if err != nil {
		return nil, err
	}
	return hjwt.ParseJwt(cookie.Value)
}

func CheckoutToken(ctx *context.Context) bool {
	cookie, err := ctx.Request.Cookie("Authorization")
	if err != nil {
		return false
	}
	return hjwt.CheckToken(cookie.Value)
}
