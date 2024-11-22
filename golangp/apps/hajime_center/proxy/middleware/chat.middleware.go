package middleware

import (
	"fmt"
	"hajime/golangp/apps/hajime_center/dify"
	"hajime/golangp/common/logging"
	"net/http"
)

func ChatMessageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := DeserializeUser(r)
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}
		isCreditsEnough := user.PreCheckBalance()

		difyClient, err := dify.GetDifyClient()
		if err != nil {
			logging.Warning("Auth Failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}

		Token, err := difyClient.GetUserToken(user.Role)
		if err != nil {
			logging.Warning("Token retrieval failed: " + err.Error())
			WriteErrorResponse(w, "401", err.Error(), http.StatusBadRequest)
			return
		}

		r.Header.Set("Authorization", "Bearer "+Token)

		if !isCreditsEnough {
			WriteErrorResponse(w, "200", "score not enough, you currently have "+fmt.Sprint(user.Balance)+" score", http.StatusBadRequest)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})

}
