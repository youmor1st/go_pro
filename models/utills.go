package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func GetUserByUsernameOrEmail(username, email string) (*User, error) {
	var user User

	query := `
        SELECT id, username, password, email FROM users
        WHERE username = $1 OR email = $2
        LIMIT 1
    `
	row := db.QueryRow(context.Background(), query, username, email)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Не найдено пользователя с таким именем или email
		}
		return nil, err
	}

	return &user, nil
}
func CreateUser(user *User) error {
	query := `
        INSERT INTO users (username, password, email)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	row := db.QueryRow(context.Background(), query, user.Username, user.Password, user.Email)
	err := row.Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User

	query := `
        SELECT id, username, password, email FROM users
        WHERE username = $1
        LIMIT 1
    `
	row := db.QueryRow(context.Background(), query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func UpdateUserProfile(updatedProfile *User) error {
	query := `
        UPDATE users
        SET username = $1, password = $2, email = $3
        WHERE id = $4
    `
	_, err := db.Exec(context.Background(), query, updatedProfile.Username, updatedProfile.Password, updatedProfile.Email, updatedProfile.ID)
	if err != nil {
		return err
	}

	return nil
}
func DeleteUserProfile(userID int) error {
	query := `
        DELETE FROM users
        WHERE id = $1
    `
	_, err := db.Exec(context.Background(), query, userID)
	if err != nil {
		return err
	}

	return nil
}
func CreateProduct(product *Product) error {
	query := `
		INSERT INTO products (name, description, price, stock_quantity, owner_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	row := db.QueryRow(context.Background(), query, product.Name, product.Description, product.Price, product.StockQuantity, product.OwnerID)
	err := row.Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func GetProducts(pageNumber, pageSize int, sortBy, filterBy string) ([]*Product, error) {
	query := `
        SELECT id, name, description, price, stock_quantity, created_at, owner_id
        FROM products
    `

	if filterBy != "" {
		query += fmt.Sprintf(" WHERE %s", filterBy)
	}

	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", sortBy)
	}

	offset := (pageNumber - 1) * pageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	// Выполняем запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.OwnerID)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetAllProducts() ([]*Product, error) {
	var products []*Product

	query := `
		SELECT id, name, description, price, stock_quantity, created_at, owner_id
		FROM products
	`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.OwnerID)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductByID(id int) (*Product, error) {
	var product Product

	query := `
        SELECT id, name, description, price, stock_quantity, created_at,owner_id  FROM products
        WHERE id = $1
        LIMIT 1
    `
	row := db.QueryRow(context.Background(), query, id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.OwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Продукт с указанным ID не найден
		}
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(productID int, updatedProduct *Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, stock_quantity = $4
        WHERE id = $5
    `
	_, err := db.Exec(context.Background(), query, updatedProduct.Name, updatedProduct.Description, updatedProduct.Price, updatedProduct.StockQuantity, productID)
	if err != nil {
		return err
	}

	return nil
}
func DeleteProduct(productID int) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`
	_, err := db.Exec(context.Background(), query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	return nil
}

func GetProductsByOwnerID(ownerID int) ([]*Product, error) {
	query := `
        SELECT id, name, description, price, stock_quantity, created_at, owner_id 
        FROM products 
        WHERE owner_id = $1
    `

	rows, err := db.Query(context.Background(), query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.OwnerID)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
func GetCartByUserID(userID int) ([]CartItem, error) {
	var cartItems []CartItem

	query := `
		SELECT id, user_id, product_id, quantity, created_at
		FROM cart_items
		WHERE user_id = $1
	`
	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cartItem CartItem
		err := rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func AddProductToCart(cartItem *CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(context.Background(), query, cartItem.UserID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func GetCartItemByProductID(cartID, productID int) (*CartItem, error) {
	var cartItem CartItem

	query := `
        SELECT id, user_id, product_id, quantity, created_at
        FROM cart_items
        WHERE cart_id = $1 AND product_id = $2
        LIMIT 1
    `
	row := db.QueryRow(context.Background(), query, cartID, productID)
	err := row.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity, &cartItem.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Продукт с указанным ID не найден в корзине
		}
		return nil, err
	}

	return &cartItem, nil
}

func UpdateCartItem(cartItem *CartItem) error {
	query := `
        UPDATE cart_items
        SET quantity = $1
        WHERE id = $2
    `
	_, err := db.Exec(context.Background(), query, cartItem.Quantity, cartItem.ID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCartByUserID(userID int, cart []CartItem) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	for _, item := range cart {
		_, err := tx.Exec(context.Background(), "INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)",
			userID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func GetOrdersByUserID(userID int) ([]*Order, error) {
	var orders []*Order

	query := `
        SELECT id, user_id, total_amount, status, created_at
        FROM orders
        WHERE user_id = $1
    `
	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrderByID(orderID int) (*Order, error) {
	var order Order

	query := `
        SELECT id, user_id, total_amount, status, created_at
        FROM orders
        WHERE id = $1
        LIMIT 1
    `
	row := db.QueryRow(context.Background(), query, orderID)
	err := row.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Заказ с указанным ID не найден
		}
		return nil, err
	}

	return &order, nil
}

func CreateOrder(order *Order) error {
	query := `
        INSERT INTO orders (user_id, total_amount, status)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	row := db.QueryRow(context.Background(), query, order.UserID, order.TotalAmount, order.Status)
	err := row.Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		fmt.Println("Error creating order:", err)
		return err
	}

	return nil
}

func CreateOrderItem(orderItem *OrderItem) error {
	query := `
        INSERT INTO order_items (order_id, product_id, quantity, price)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
	row := db.QueryRow(context.Background(), query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	err := row.Scan(&orderItem.ID, &orderItem.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrder(order *Order) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`
	_, err := db.Exec(context.Background(), query, order.Status, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}
	return nil
}

func DeleteOrder(orderID int) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`
	_, err := db.Exec(context.Background(), query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}
	return nil
}
