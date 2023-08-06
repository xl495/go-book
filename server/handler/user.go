package handler

import (
	"api-fiber-gorm/common"
	"github.com/go-playground/validator/v10"
	"strconv"

	"api-fiber-gorm/database"
	"api-fiber-gorm/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// GetUser 获取用户
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.Omit("password").Find(&user, id)
	if user.Username == "" {
		return common.Response(c, 404, "找不到该ID的用户", nil)
	}

	return common.Response(c, 200, "成功", fiber.Map{"data": user})
}

// CreateUser 创建新用户
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	validate := validator.New()

	db := database.DB

	user := new(model.User)

	if err := c.BodyParser(user); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "错误", err)
	}

	if err := validate.Struct(user); err != nil {
		return common.Response(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return common.Response(c, 500, "无法加密密码", err)
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		return common.Response(c, 500, "无法创建用户", err)
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return common.Response(c, 200, "成功", fiber.Map{"data": newUser})
}

// UpdateUser 更新用户
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return common.Response(c, 500, "请检查您的输入", err)
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return common.Response(c, 500, "无效的令牌ID", nil)
	}

	db := database.DB
	var user model.User

	db.First(&user, id)
	user.Names = uui.Names
	db.Save(&user)

	return common.Response(c, 200, "用户更新成功", fiber.Map{"data": user})
}

// DeleteUser 删除用户
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return common.Response(c, 500, "请检查您的输入", err)
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return common.Response(c, 500, "无效的令牌ID", nil)
	}

	if !validUser(id, pi.Password) {
		return common.Response(c, 500, "无效的用户", nil)
	}

	db := database.DB
	var user model.User

	db.First(&user, id)

	db.Delete(&user)
	return common.Response(c, 200, "用户删除成功", nil)
}
