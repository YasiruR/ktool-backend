package http

import (
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"net/http"
	"strings"
)

func authenticateMiddleware(next http.Handler) http.Handler {
	ctx := traceable_context.WithUUID(uuid.New())
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//user validation by token header
		token := req.Header.Get("Authorization")
		if len(strings.Split(token, "Bearer")) < 2 {
			log.Logger.ErrorContext(ctx, "token format is invalid", token)
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
		if !ok {
			log.Logger.DebugContext(ctx, "invalid user", token)
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(res, req)
	})
}
