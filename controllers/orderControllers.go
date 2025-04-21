package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/config"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

func CreateOrder(c *gin.Context) {
	var order models.Order

	// Ambil userID dari context (dari JWT via middleware)

	// Validasi input
	if err := c.ShouldBindJSON(&order); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		if order.ProductID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Product can't be empty",
			})
			return
		}

		if order.Quantity == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: " can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", order.ProductID).Scan(&exists)
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
	order.UserID = userID

	query := `
	INSERT INTO "orders" 
	(user_id, product_id, quantity, price, subtotal, note, address, status, code)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, transaction_date`

	err := config.DB.QueryRow(query,
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.Price,
		order.Subtotal,
		order.Note,
		order.Address,
		order.Status,
		order.Code,
	).Scan(&order.ID, &order.TransactionDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	if order.CoID != 0 {
		_, err := config.DB.Exec("DELETE FROM cart WHERE id = $1", order.CoID)
		if err != nil {
			// Kalau error, kirim warning di response tapi jangan gagalkan order
			c.JSON(http.StatusOK, models.SuccessAddUpdate{
				Status:  "success_with_warning",
				Code:    200,
				Message: "Order added, but failed to delete cart",
				Data:    order,
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes add data",
		Data:    order,
	})
}

func UpdateOrder(c *gin.Context) {
	var order models.Order
	id := c.Param("id") // ambil id dari URL parameter

	if err := c.ShouldBindJSON(&order); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if order.UserID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "user can't be empty",
			})
			return
		}

		if order.ProductID == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: "Product can't be empty",
			})
			return
		}

		if order.Quantity == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Status:  "failed",
				Code:    400,
				Message: " can't be empty",
			})
			return
		}

		// Cek apakah KategoriID ada di table category
		var exists bool
		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)", order.ProductID).Scan(&exists)
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
		errUser := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE id = $1)", order.UserID).Scan(&existsUser)
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

	query := `
	UPDATE orders
	SET user_id = $1,
		product_id = $2,
		quantity = $3,
		price = $4,
		subtotal = $5,
		note = $6,
		address = $7,
		status = $8,
		code = $9
	WHERE id = $10
`
	res, err := config.DB.Exec(query,
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.Price,
		order.Subtotal,
		order.Note,
		order.Address,
		order.Status,
		order.Code,
		id,
	)
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
			Message: "order not found",
		})
	}

	c.JSON(http.StatusOK, models.SuccessAddUpdate{
		Status:  "success",
		Code:    200,
		Message: "succes update data",
		Data:    order,
	})
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM orders WHERE id = $1`
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
			Message: "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "success delete data",
	})
}

func GetOrders(c *gin.Context) {
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// Ambil data order untuk user tersebut
	rows, err := config.DB.Query(`
		SELECT
			orders.id,
			orders.user_id,
			users.name AS user_name,
			orders.product_id,
			product.name AS product_name,
			orders.quantity,
			orders.price,
			orders.address,
			orders.subtotal,
			orders.note,
			orders.status,
			orders.code,
			orders.transaction_date
		FROM orders
		INNER JOIN users ON orders.user_id = users.id
		INNER JOIN product ON orders.product_id = product.id
		WHERE orders.user_id = $1
		ORDER BY orders.id
		LIMIT $2 OFFSET $3
	`, userID, size, offset)

	if err != nil {
		// c.JSON(http.StatusInternalServerError, models.Response{
		// 	Status:  "failed",
		// 	Code:    500,
		// 	Message: err.Error(),
		// })
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.UserName,
			&order.ProductID,
			&order.ProductName,
			&order.Quantity,
			&order.Price,
			&order.Address,
			&order.Subtotal,
			&order.Note,
			&order.Status,
			&order.Code,
			&order.TransactionDate,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "failed",
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, models.SuccessList{
		Status:    "success",
		Code:      200,
		Message:   "Success get data",
		Data:      orders,
		Page:      page,
		Size:      size,
		TotalData: totalData,
	})
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	query := `
		SELECT
			order.id,
			order.user_id,
			users.name AS user_name,
			order.product_id,
			product.name AS product_name,
			order.quantity,
			order.price,
			order.address,
			order.subtotal,
			order.note
		FROM orders
		INNER JOIN users ON order.user_id = users.id
		INNER JOIN product ON order.product_id = product.id
		WHERE order.id = $1
	`

	err := config.DB.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.UserName,
		&order.ProductID,
		&order.ProductName,
		&order.Quantity,
		&order.Price,
		&order.Address,
		&order.Subtotal,
		&order.Note,
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
		Data:    order,
	})
}
