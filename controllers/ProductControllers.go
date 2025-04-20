package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	// Validasi input
	if err := c.ShouldBindJSON(&product); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if product.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "name can't be empty",
			})
			return
		}
		if product.Code == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "code can't be empty",
			})
			return
		}

		if product.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "category can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE id = $1)", product.CategoryID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}

		if !exists {
			c.JSON(http.StatusNotFound, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "category not found",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    404,
			Message: validationError,
		})
		return
	}

	query := `INSERT INTO product (name, code, category_id,stock,image, description) VALUES ($1, $2, $3, $4,$5,$6) RETURNING id`
	err := config.DB.QueryRow(query, product.Name, product.Code, product.CategoryID, product.Stock, product.Image, product.Description).Scan(&product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes add data",
		Data:    product,
	})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id") // ambil id dari URL parameter

	if err := c.ShouldBindJSON(&product); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if product.Name == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "name can't be empty",
			})
			return
		}
		if product.Code == "" {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "code can't be empty",
			})
			return
		}

		if product.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "category can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE id = $1)", product.CategoryID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}

		if !exists {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "category not found",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, models.Response{
			Status:  "failed",
			Code:    404,
			Message: validationError,
		})
		return
	}

	query := `UPDATE product SET name = $1, code = $2, category_id = $3, stock = $4, image = $5, description = $6 WHERE id = $7`
	res, err := config.DB.Exec(query, product.Name, product.Code, product.CategoryID, product.Stock, product.Image, product.Description, id)
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
		c.JSON(http.StatusNotFound, models.Response{
			Status:  "failed",
			Code:    404,
			Message: "product not found",
		})
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes update data",
		Data:    product,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM book WHERE id = $1`
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
		c.JSON(http.StatusNotFound, models.Response{
			Status:  "failed",
			Code:    404,
			Message: "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "success add data",
	})
}

func GetProducts(c *gin.Context) {
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM product`).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Ambil data book dengan JOIN ke category
	rows, err := config.DB.Query(`
		SELECT
			product.id,
			product.name,
			product.code,
			product.category_id,
			category.name AS category_name,
			product.stock,
			product.description
		FROM product
		INNER JOIN category ON product.category_id = category.id
		ORDER BY product.id
		LIMIT $1 OFFSET $2
	`, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal mengambil data book",
		})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Code,
			&product.CategoryID,
			&product.CategoryName,
			&product.Stock,
			&product.Description,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal membaca data book",
			})
			return
		}
		products = append(products, product)
	}

	// Cek error setelah rows.Next()
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, models.SuccessList{
		Status:    "success",
		Code:      200,
		Message:   "Success get data",
		Data:      products,
		Page:      page,
		Size:      size,
		TotalData: totalData,
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	query := `SELECT id, name, category_id,code,stock,description FROM book WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.CategoryID, &product.Stock, &product.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "data not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal mengambil detail book",
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "Success get data",
		Data:    product,
	})
}
