package middleware

import (
	"errors"
	"fmt"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func DeserializeUser(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) < 2 || fields[0] != "Bearer" {
		return nil, errors.New("you are not logged in")
	}

	accessToken := fields[1]

	config, err := initializers.LoadEnv(".")
	if err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("the user belonging to this token no longer exists")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil
}
