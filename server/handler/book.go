package handler

import (
	"api-fiber-gorm/common"
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func CreateBook(c *fiber.Ctx) error {

	db := database.DB

	validate := validator.New()

	Book := new(model.Book)

	if err := c.BodyParser(Book); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	if err := validate.Struct(Book); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	if err := db.Create(&Book).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "创建失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "创建成功", fiber.Map{"data": Book})
}

func GetBookList(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "20")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	db := database.DB

	var books []model.Book
	var totalCount int64

	// 查询总条数
	if err := db.Model(&model.Book{}).Count(&totalCount).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	if err := db.Preload("Category").Limit(intLimit).Offset(offset).Find(&books).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "获取成功", fiber.Map{
		"list":     books,
		"total":    totalCount,
		"page":     intPage,
		"pageSize": intLimit,
	})
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id 为空", nil)
	}

	db := database.DB

	var Book model.Book

	if err := db.Preload("Category").First(&Book, id).Error; err != nil {

		// 判断错误是否 orm 记录不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Response(c, fiber.StatusBadRequest, "未找到 id ", err.Error())
		}

		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "获取成功", Book)
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id 为空", nil)
	}

	db := database.DB

	var Book model.Book

	if err := db.First(&Book, id).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	if err := c.BodyParser(&Book); err != nil {
		return common.Response(c, fiber.StatusBadRequest, "更新失败", err.Error())
	}

	if err := db.Save(&Book).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "更新失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "更新成功", Book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return common.Response(c, fiber.StatusBadRequest, "id 为空", nil)
	}

	db := database.DB

	var Book model.Book

	if err := db.First(&Book, id).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "获取失败", err.Error())
	}

	if err := db.Delete(&Book).Error; err != nil {
		return common.Response(c, fiber.StatusBadRequest, "删除失败", err.Error())
	}

	return common.Response(c, fiber.StatusOK, "删除成功", nil)
}
