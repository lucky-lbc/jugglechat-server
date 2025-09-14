package apis

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/juggleim/commons/ctxs"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/responses"
	utils "github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/admins/services"
	"github.com/juggleim/jugglechat-server/storages/dbs"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"
)

func Validate(ctx *gin.Context) {
	session := fmt.Sprintf("admin_%s", utils.GenerateUUIDShort11())
	ctx.Header(Header_RequestId, session)
	ctx.Set(string(ctxs.CtxKey_Session), session)

	urlPath := ctx.Request.URL.Path
	if strings.HasSuffix(urlPath, "/login") || strings.HasSuffix(urlPath, "/apps/create") || strings.HasSuffix(urlPath, "/apps/active") {
		return
	}
	authStr := ctx.Request.Header.Get(Header_Authorization)
	account, roleType, err := validateAuthorization(authStr)
	if err != nil {
		responses.AdminErrorHttpResp(ctx, errs.AdminErrorCode_AuthFail)
		ctx.Abort()
		return
	}
	//check account
	code := services.CheckAccountState(account)
	if code != errs.AdminErrorCode_Success {
		responses.AdminErrorHttpResp(ctx, code)
		ctx.Abort()
		return
	}
	ctx.Set(string(ctxs.CtxKey_Account), account)
	ctx.Set(string(ctxs.CtxKey_RoleType), roleType)
}

func GetLoginedAccount(ctx *gin.Context) string {
	if account, ok := ctx.Value(ctxs.CtxKey_Account).(string); ok {
		return account
	}
	return ""
}

func GetAccountRoleType(ctx *gin.Context) dbs.RoleType {
	if rt, ok := ctx.Value(ctxs.CtxKey_RoleType).(dbs.RoleType); ok {
		return rt
	}
	return dbs.RoleType_NormalAdmin
}

var jwtkey = []byte("jug9le1m")

type Claims struct {
	Account  string
	RoleType int
	jwt.RegisteredClaims
}

func TestAu() {
	fmt.Println(generateAuthorization("admin1", dbs.RoleType_SuperAdmin))
}

func generateAuthorization(account string, role dbs.RoleType) (string, error) {
	expireTime := time.Now().Add(time.Hour)
	claims := &Claims{
		Account:  account,
		RoleType: int(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expireTime,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			Issuer:  "jugglechat",
			Subject: "juggle",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateAuthorization(authorization string) (string, dbs.RoleType, error) {
	token, claims, err := parseToken(authorization)
	if err != nil || !token.Valid {
		return "", dbs.RoleType_NormalAdmin, fmt.Errorf("auth fail")
	}
	return claims.Account, dbs.RoleType(claims.RoleType), nil
}

func parseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
