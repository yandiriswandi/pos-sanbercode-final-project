package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

func CreateCart(c *gin.Context) {
	var cart models.Cart

	// Ambil userID dari context (dari JWT via middleware)

	// Validasi input
	if err := c.ShouldBindJSON(&cart); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		if cart.ProductID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Product can't be empty",
			})
			return
		}

		if cart.Quantity == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: " can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", cart.ProductID).Scan(&exists)
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
				Message: "product not found",
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

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Status:  "failed",
			Code:    401,
			Message: "User ID not found in context",
		})
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Status:  "failed",
			Code:    401,
			Message: "User ID invalid",
		})
		return
	}
	cart.UserID = userID

	// Insert ke database
	query := `
			INSERT INTO cart (user_id, product_id, quantity, price, subtotal, note,address)
			VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id
			`
	err := config.DB.QueryRow(query, cart.UserID, cart.ProductID, cart.Quantity, cart.Price, cart.Subtotal, cart.Note, cart.Address).Scan(&cart.ID)
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
		Data:    cart,
	})
}

func UpdateCart(c *gin.Context) {
	var cart models.Cart
	id := c.Param("id") // ambil id dari URL parameter

	if err := c.ShouldBindJSON(&cart); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if cart.UserID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "user can't be empty",
			})
			return
		}

		if cart.ProductID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Product can't be empty",
			})
			return
		}

		if cart.Quantity == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: " can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", cart.ProductID).Scan(&exists)
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
				Message: "product not found",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var existsUser bool
		errUser := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE id = $1)", cart.UserID).Scan(&existsUser)
		if errUser != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: errUser.Error(),
			})
			return
		}

		if !existsUser {
			c.JSON(http.StatusNotFound, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "user not found",
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

	query := `UPDATE cart SET user_id = $1, product_id = $2, quantity = $3, price = $4, subtotal = $5, note = $6 ,address=$7 WHERE id = $8`
	res, err := config.DB.Exec(query, cart.UserID, cart.ProductID, cart.Quantity, cart.Price, cart.Subtotal, cart.Note, cart.Address, id)
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
			Message: "cart not found",
		})
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes update data",
		Data:    cart,
	})
}

func DeleteCart(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM cart WHERE id = $1`
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
			Message: "cart not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "success delete data",
	})
}

func GetCarts(c *gin.Context) {
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

	// Ambil user_id dari context (dari middleware JWT)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.Response{
			Status:  "failed",
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}
	userID := userIDInterface.(int)

	// Hitung total data untuk user tersebut
	var totalData int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM cart WHERE user_id = $1`, userID).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Ambil data cart untuk user tersebut
	rows, err := config.DB.Query(`
		SELECT
			cart.id,
			cart.user_id,
			users.name AS user_name,
			cart.product_id,
			product.name AS product_name,
			cart.quantity,
			cart.price,
			cart.address,
			cart.subtotal,
			cart.note
		FROM cart
		INNER JOIN users ON cart.user_id = users.id
		INNER JOIN product ON cart.product_id = product.id
		WHERE cart.user_id = $1
		ORDER BY cart.id
		LIMIT $2 OFFSET $3
	`, userID, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var carts []models.Cart
	for rows.Next() {
		var cart models.Cart
		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.UserName,
			&cart.ProductID,
			&cart.ProductName,
			&cart.Quantity,
			&cart.Price,
			&cart.Address,
			&cart.Subtotal,
			&cart.Note,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		carts = append(carts, cart)
	}

	c.JSON(http.StatusOK, models.SuccessList{
		Status:    "success",
		Code:      200,
		Message:   "Success get data",
		Data:      carts,
		Page:      page,
		Size:      size,
		TotalData: totalData,
	})
}

func GetCart(c *gin.Context) {
	id := c.Param("id")

	var cart models.Cart
	query := `
		SELECT
			cart.id,
			cart.user_id,
			users.name AS user_name,
			cart.product_id,
			product.name AS product_name,
			cart.quantity,
			cart.price,
			cart.address,
			cart.subtotal,
			cart.note
		FROM cart
		INNER JOIN users ON cart.user_id = users.id
		INNER JOIN product ON cart.product_id = product.id
		WHERE cart.id = $1
	`

	err := config.DB.QueryRow(query, id).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.UserName,
		&cart.ProductID,
		&cart.ProductName,
		&cart.Quantity,
		&cart.Price,
		&cart.Address,
		&cart.Subtotal,
		&cart.Note,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Status:  "failed",
				Code:    404,
				Message: "data not found",
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
		Data:    cart,
	})
}
