package auth

import (
	"context"
	"finance-crud-app/internal/types"
	"finance-crud-app/internal/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey string = "user_id"

func JWTAuthMiddleWare(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id")
			permissionDenied(w)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, UserKey, strconv.Itoa(userID))
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func CreateJWT(userID int) (string, error) {
	expiration := time.Minute * time.Duration(45)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Printf("signing error %v", tokenString)
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("access denied"))
}

func GetUserIDFromContext(ctx context.Context) string {
	valuee := ctx.Value("user_id")

	value, ok := ctx.Value(UserKey).(string)

	if !ok {
		log.Printf("user key from the get context function %v", valuee)
		return "-1"
	}

	if _, err := strconv.Atoi(value); err != nil {
		return "-1"
	}

	return value
}
