package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

func CreateCategory(c *gin.Context) {
	var category models.Category
	// Validasi input
	if err := c.ShouldBindJSON(&category); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if category.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Name can't empty",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    400,
			Message: validationError,
		})
		return
	}

	query := `INSERT INTO category (name, code, description)
          VALUES ($1, $2, $3) RETURNING id`

	err := config.DB.QueryRow(query,
		category.Name,
		category.Code,
		category.Description,
	).Scan(&category.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Code:    200,
		Status:  "success",
		Message: "success to add data",
		Data:    category,
	})
}

func UpdateCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id") // ambil id dari URL parameter

	// Validasi input
	if err := c.ShouldBindJSON(&category); err != nil {
		validationError := err.Error()

		if category.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Name can't empty",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    400,
			Message: validationError,
		})
		return
	}

	query := `UPDATE category SET name = $1, description = $2,code = $3  WHERE id = $4`
	res, err := config.DB.Exec(query, category.Name, category.Description, category.Code, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "success updated data",
		"data":    category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM category WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    404,
			Message: "category not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Code:    200,
		Message: "success deleted data",
	})
}

func GetCategoryList(c *gin.Context) {
	// Ambil query parameter page dan size, default kalau kosong
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	// Hitung total data
	var totalData int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM category`).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Ambil data category
	rows, err := config.DB.Query(`SELECT id, name, description,code FROM category ORDER BY id LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var categorys []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Code); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		categorys = append(categorys, category)
	}

	c.JSON(http.StatusOK, models.SuccessList{
		Status:    "success",
		Code:      200,
		Message:   "Success get data",
		Data:      categorys,
		Page:      page,
		Size:      size,
		TotalData: totalData,
	})
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	query := `SELECT id, name, description,code FROM category WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "category not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "Success get data",
		Data:    category,
	})
}
