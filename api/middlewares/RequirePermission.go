package middlewares

import (
	"context"
	"net/http"

	"bitpos.com/m/repos"
	"bitpos.com/m/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value("userID").(int64)
	return userID, ok
}

/* ================= FIXED PERMISSION MIDDLEWARE ================= */

func RequirePermission(userAccessRepo *repos.UserAccessRepo, permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			userID, ok := GetUserID(r.Context())
			if !ok {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			has, err := userAccessRepo.HasPermission(userID, permission)
			if err != nil {
				utils.ErrorResponse(w, http.StatusInternalServerError, "Permission Check Failed")
				return
			}

			if !has {
				utils.ErrorResponse(w, http.StatusForbidden, "Insufficient Permissions")
				return
			}

			next(w, r)
		}
	}
}