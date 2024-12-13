package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"hajime/golangp/apps/hajime_center/constants"
	"hajime/golangp/apps/hajime_center/initializers"
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/common/logging"
	"hajime/golangp/common/mail_server"
	"hajime/golangp/common/utils"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	DB *gorm.DB
	cs *CreditSystem
}

func SendEmail(email string, data *mail_server.EmailData, emailTemp string) error {
	config, err := initializers.LoadEnv(".")
	if err != nil {
		return err
	}
	emailConfig := mail_server.EmailConfig{
		EmailFrom: config.EmailFrom,
		SMTPPass:  config.SMTPPass,
		SMTPUser:  config.SMTPUser,
		SMTPHost:  config.SMTPHost,
		SMTPPort:  config.SMTPPort,
	}

	templatesPath := "golangp/apps/hajime_center/templates"
	defaultTemplatePath := "/srv/HajimeCenter/templates"
	mail_server.SendEmail(&emailConfig, email, data, emailTemp, templatesPath, defaultTemplatePath)
	return nil
}

func NewAuthController(DB *gorm.DB, creditSystem *CreditSystem) AuthController {
	return AuthController{DB, creditSystem}
}

// SignUpUser SignUp User
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:              payload.Name,
		Email:             strings.ToLower(payload.Email),
		Password:          hashedPassword,
		Role:              constants.RoleEditor,
		Verified:          false,
		Photo:             "test",
		Provider:          "local",
		Balance:           constants.GiftedPoints,
		FromCode:          payload.FromCode,
		UserMaxCodeAmount: constants.RoleEditorMaxCodeAmount,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists, " +
			"try use forget password to reset it."})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}
	// Generate Verification Code
	code := randstr.String(6)

	verificationCode := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verificationCode
	ac.DB.Save(newUser)

	config, _ := initializers.LoadEnv(".")

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := mail_server.EmailData{
		URL:              config.ClientOrigin + "/verifyemail/" + code,
		VerificationCode: code,
		FirstName:        firstName,
		Subject:          "Your account verification code",
	}

	err = SendEmail(newUser.Email, &emailData, "verificationCode.html")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}

	message := "We sent an email with a verification code to " + newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "BAD_REQUEST", "message": err.Error()})
		return
	}

	var user models.User
	if payload.Email != "" {
		// Login with email
		result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "account_not_found", "message": "Invalid email or password"})
			return
		}
	} else if payload.Name != "" {
		// Login with name
		result := ac.DB.First(&user, "name = ?", payload.Name)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "accout_not_found", "message": "Invalid username or password"})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "BAD_REQUEST", "message": "Email or username is required"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, gin.H{"result": "fail", "code": "UNVERIFIED_ACCOUNT", "message": "Please verify your email"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "INVALID_CREDENTIALS", "message": "Invalid email or password"})
		return
	}

	now := time.Now()
	if user.LoginTime == nil || !user.IsSameDay(*user.LoginTime, now) {
		// 更新余额
		if err := user.UpdateBalance(constants.DailySignInPoints, "DailySignInPoints"); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"result": "fail", "code": "UPDATE_BALANCE_ERROR", "message": err.Error()})
			return
		}

		// 更新登录时间
		if err := user.UpdateLoginTime(&now); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"result": "fail", "code": "UPDATE_LOGINTIME_ERROR", "message": err.Error()})
			return
		}
	}

	config, _ := initializers.LoadEnv(".")

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "TOKEN_GENERATION_ERROR", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"result": "fail", "code": "TOKEN_GENERATION_ERROR", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"result": "success", "data": gin.H{"access_token": accessToken, "refresh_token": refreshToken}})
}

// RefreshAccessToken Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := initializers.LoadEnv(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	config, _ := initializers.LoadEnv(".")

	ctx.SetCookie("access_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Domain, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// VerifyEmail [...] Verify Email
func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, err := initializers.LoadEnv(".")
	if err != nil {
		logging.Danger("Could not load openai-config", err)
	}

	// Generate Verification Code
	resetToken := randstr.String(6)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := mail_server.EmailData{
		URL:              config.ClientOrigin + "/resetpassword/" + resetToken,
		VerificationCode: resetToken,
		FirstName:        firstName,
		Subject:          "Your password reset token (valid for 10min)",
	}

	err = SendEmail(user.Email, &emailData, "verificationCode.html")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")
	config, _ := initializers.LoadEnv(".")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(payload.Password)

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.Verified = true
	updatedUser.PasswordResetToken = ""
	ac.DB.Save(&updatedUser)

	ctx.SetCookie("token", "", -1, "/", config.Domain, false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}

func (ac *AuthController) PasswordChange(ctx *gin.Context) {
	var payload *models.PasswordChangeInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.NewPassword != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "New password and confirmation password do not match"})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)

	if err := utils.VerifyPassword(currentUser.Password, payload.CurrentPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Current password is incorrect"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	currentUser.Password = hashedPassword
	ac.DB.Save(&currentUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password updated successfully"})
}

// AddUser 处理添加新用户
func (ac *AuthController) AddUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	if currentUser.Role != constants.RoleAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Only Admin can add users"})
		return
	}

	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var maxCodeAmount int

	if payload.Role == constants.RoleAdmin {
		maxCodeAmount = constants.RoleAdminMaxCodeAmount
	} else if payload.Role == constants.RoleEditor {
		maxCodeAmount = constants.RoleEditorMaxCodeAmount
	} else {
		maxCodeAmount = constants.RoleUserMaxCodeAmount
	}

	now := time.Now()
	newUser := models.User{
		Name:              payload.Name,
		Email:             strings.ToLower(payload.Email),
		Password:          hashedPassword,
		Role:              payload.Role, // 允许管理员设置用户角色
		Verified:          true,         // 直接设置为已验证
		Photo:             "test",
		Provider:          "local",
		FromCode:          payload.FromCode,
		UserMaxCodeAmount: maxCodeAmount,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	result := ac.DB.Create(&newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
			return
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User added successfully"})
}

// DeleteUser 处理删除用户
func (ac *AuthController) DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	if currentUser.Role != constants.RoleAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Only Admin can update credits"})
		return
	}

	userId := ctx.Param("userId")

	var user models.User
	result := ac.DB.First(&user, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	ac.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}

func (ac *AuthController) UpdateUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	if currentUser.Role != constants.RoleAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Only Admin can update users"})
		return
	}

	var payload *models.UpdateUserInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userId := ctx.Param("userId")
	var user models.User
	result := ac.DB.First(&user, "id = ?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	if payload.Name != "" {
		user.Name = payload.Name
	}
	if payload.Email != "" {
		user.Email = strings.ToLower(payload.Email)
	}
	if payload.Role != "" {
		user.Role = payload.Role
	}

	ac.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully"})
}

// GetAllUsers 获取所有用户
func (ac *AuthController) GetAllUsers(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	if currentUser.Role != constants.RoleAdmin {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Only Admin can update credits"})
		return
	}

	var users []models.User
	result := ac.DB.Find(&users)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to retrieve users"})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponse := models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Photo:     user.Photo,
			Role:      user.Role,
			Verified:  user.Verified,
			Balance:   user.Balance,
			Provider:  user.Provider,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		userResponses = append(userResponses, userResponse)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"users": userResponses}})
}

func (ac *AuthController) GetMe(ctx *gin.Context) {
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Photo:     currentUser.Photo,
		Role:      currentUser.Role,
		Verified:  currentUser.Verified,
		Balance:   currentUser.Balance,
		Provider:  currentUser.Provider,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) LoginWithWallet(ctx *gin.Context) {
	// 获取当前用户
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// 从请求中获取 SignLoginModel
	var form models.SignLoginModel

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	// Check if currentUser.Address exists
	if currentUser.Address == "" {
		err := currentUser.UpdateBalance(constants.WalletLinkPoints, "WalletLinkPoints")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	// 将 WalletAddress 转为小写
	form.WalletAddress = strings.ToLower(form.WalletAddress)

	// 调用 LoginWithSign 方法
	loginResponse, err := currentUser.LoginWithSign(form, currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 返回登录响应
	ctx.JSON(http.StatusOK, loginResponse)
}

func (ac *AuthController) UnlinkWallet(ctx *gin.Context) {
	// 获取当前用户
	currentUser, ok := ctx.MustGet("currentUser").(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User not found"})
		return
	}

	// 调用 UnlinkWallet 方法
	err := currentUser.UnlinkWallet(currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Wallet unlinked successfully"})
}
