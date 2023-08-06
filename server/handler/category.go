package handler

import (
	"api-fiber-gorm/common"
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateCategory(c *fiber.Ctx) error {

	db := database.DB

	validate := validator.New()

	category := new(model.Category)

	if err := c.BodyParser(category); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	if err := validate.Struct(category); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	if err := db.Create(&category).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "创建成功", fiber.Map{"data": category})
}

func GetCategoryList(c *fiber.Ctx) error {
	db := database.DB

	var categories []model.Category

	if err := db.Preload("Books").Find(&categories).Error; err != nil {
		return common.Response(c, fiber.StatusInternalServerError, "获取失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "获取成功", categories)
}

func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id为空", nil)
	}

	db := database.DB

	var category model.Category

	if err := db.Preload("Books").First(&category, id).Error; err != nil {
		// 判断错误是否 orm 记录不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Response(c, fiber.StatusBadRequest, "未找到id", err.Error())
		}
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "获取成功", category)
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id 为空", nil)
	}

	db := database.DB

	var category model.Category

	if err := db.First(&category, id).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	if err := c.BodyParser(&category); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "更新失败", err.Error())
	}

	if err := db.Save(&category).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "更新失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "更新成功", category)
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id 为空", nil)
	}

	db := database.DB

	var category model.Category

	if err := db.First(&category, id).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "此分类不存在!", nil)
	}

	// 查找 book 是否该分类下的书籍
	var book model.Book
	if err := db.Where("category_id = ?", category.ID).First(&book).Error; errors.Is(err, gorm.ErrRecordNotFound) == false {
		return common.Response(c, fiber.StatusBadRequest, "分类下存在书籍,请先删除全部书籍!", nil)
	}

	if err := db.Delete(&category).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "删除失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "删除成功", nil)
}
