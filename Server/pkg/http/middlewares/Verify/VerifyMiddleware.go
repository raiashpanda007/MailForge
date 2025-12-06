package httpmiddleware

import (
	"context"
	"net/http"

	"github.com/raiashpanda007/MailForge/pkg/http/controllers/auth"
	"github.com/raiashpanda007/MailForge/pkg/utils"
)

func VerifyMiddleware(tokenProvider auth.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ClientSideToken, err := req.Cookie("Access-Token")
			if err != nil {
				utils.WriteJson(res, http.StatusUnauthorized, utils.Data{Message: "PLEASE LOGIN", Data: utils.GeneralError(err, "PLEASE LOGIN")})
				return
			}
			verifiedUser, err := tokenProvider.VerifyToken(ClientSideToken.Value)
			if err != nil {
				utils.WriteJson(res, http.StatusUnauthorized, utils.Data{Message: "INVALID TOKEN PLEASE LOGIN AGAIN", Data: utils.GeneralError(err, "INVALID TOKEN")})
				return
			}
			ctx := context.WithValue(req.Context(), "USER", verifiedUser)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
