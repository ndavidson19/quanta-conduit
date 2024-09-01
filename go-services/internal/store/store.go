package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"conduit/internal/db" // Import the generated db code

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	PostgresURL string
}

type Store struct {
	*db.Queries
	db *sql.DB
}

func NewStore(config Config) (*Store, error) {
	log.Printf("Attempting to connect to PostgreSQL with URL: %s", config.PostgresURL)

	var sqlDB *sql.DB
	var err error

	// Retry logic for database connection
	for i := 0; i < 5; i++ {
		sqlDB, err = sql.Open("postgres", config.PostgresURL)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database connection: %v. Retrying in 5 seconds...", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Set connection pool settings
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)

		// Try to ping the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err == nil {
			log.Printf("Successfully connected to the database on attempt %d", i+1)
			break
		}
		log.Printf("Attempt %d: Failed to ping database: %v. Retrying in 5 seconds...", i+1, err)
		sqlDB.Close() // Close the connection before retrying
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
	}

	// Test query to verify connection and permissions
	var testResult int
	err = sqlDB.QueryRow("SELECT 1").Scan(&testResult)
	if err != nil {
		return nil, fmt.Errorf("failed to execute test query: %w", err)
	}

	log.Printf("Successfully executed test query. Result: %d", testResult)

	store := &Store{
		Queries: db.New(sqlDB),
		db:      sqlDB,
	}

	// Ensure special accounts exist
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//if err := store.EnsureSpecialAccountsExist(ctx); err != nil {
	//    return nil, fmt.Errorf("failed to ensure special accounts exist: %w", err)
	//}
	// Initialize schema
	if err := store.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	// Verify schema
	if err := store.VerifySchema(); err != nil {
		log.Fatalf("Schema verification failed: %v", err)
	}

	log.Println("Successfully initialized store and ensured special accounts exist")

	return store, nil
}

// EnsureSpecialAccountsExist checks and creates special system accounts if they don't exist
func (s *Store) EnsureSpecialAccountsExist(ctx context.Context) error {
	specialAccounts := []struct {
		accountID   int64
		accountType string
	}{
		{1, "deposit"},
		{2, "withdrawal"},
	}

	for _, sa := range specialAccounts {
		account, err := s.GetAccountByAccountID(ctx, sa.accountID)
		if err == sql.ErrNoRows {
			// Account doesn't exist, create it
			_, err = s.CreateAccount(ctx, db.CreateAccountParams{
				AccountID:       sa.accountID,
				Balance:         "0", // Assuming Balance is a string in the generated code
				AccountType:     sql.NullString{String: sa.accountType, Valid: true},
				IsSystemAccount: sql.NullBool{Bool: true, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to create special account %d: %w", sa.accountID, err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check for special account %d: %w", sa.accountID, err)
		} else if !account.IsSystemAccount.Bool {
			return fmt.Errorf("account %d exists but is not a system account", sa.accountID)
		}
	}

	return nil
}

// ExecTx executes a function within a database transaction
func (s *Store) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// PlaceOrder places a new order and updates the portfolio
func (s *Store) PlaceOrder(ctx context.Context, params db.CreateOrderParams) error {
	return s.ExecTx(ctx, func(q *db.Queries) error {
		// Create the order
		_, err := q.CreateOrder(ctx, params)
		if err != nil {
			return err
		}

		// Note: Additional logic for updating portfolio and account balance
		// needs to be implemented based on your specific requirements

		return nil
	})
}

// ExecuteOrder marks an order as executed and updates the portfolio
func (s *Store) ExecuteOrder(ctx context.Context, orderID int32, executedPrice float64) error {
	return s.ExecTx(ctx, func(q *db.Queries) error {
		// Get the original order
		order, err := q.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		// Create executed order record
		_, err = q.CreateExecutedOrder(ctx, db.CreateExecutedOrderParams{
			OrderID:          orderID,
			ExecutedQuantity: order.Quantity,
			ExecutedPrice:    fmt.Sprintf("%.2f", executedPrice), // Assuming ExecutedPrice is a string in the generated code
		})
		if err != nil {
			return err
		}

		// Update order status
		_, err = q.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
			ID:     orderID,
			Status: "executed",
		})
		if err != nil {
			return err
		}

		// Note: Additional logic for updating portfolio
		// needs to be implemented based on your specific requirements

		return nil
	})
}

type ListAccountsParams struct {
	UserID sql.NullInt32
	Limit  int32
	Offset int32
}

type ListOrdersParams struct {
	AccountID sql.NullInt64
	Status    sql.NullString
	Limit     int32
	Offset    int32
}

func (s *Store) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]db.Account, error) {
	query := `
        SELECT id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at
        FROM accounts
        WHERE ($1::integer IS NULL OR user_id = $1)
        ORDER BY id
        LIMIT $2 OFFSET $3
    `
	rows, err := s.db.QueryContext(ctx, query, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []db.Account
	for rows.Next() {
		var account db.Account
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.AccountID,
			&account.Balance,
			&account.PortfolioValue,
			&account.AccountTier,
			&account.AccountType,
			&account.IsSystemAccount,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *Store) ListOrders(ctx context.Context, arg ListOrdersParams) ([]db.Order, error) {
	query := `
        SELECT id, account_id, asset_id, order_type, side, quantity, price, status, created_at, updated_at
        FROM orders
        WHERE ($1::bigint IS NULL OR account_id = $1)
          AND ($2::text IS NULL OR status = $2)
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4
    `
	rows, err := s.db.QueryContext(ctx, query, arg.AccountID, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []db.Order
	for rows.Next() {
		var order db.Order
		if err := rows.Scan(
			&order.ID,
			&order.AccountID,
			&order.AssetID,
			&order.OrderType,
			&order.Side,
			&order.Quantity,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetPortfolioSummary retrieves a summary of the portfolio for a given account
// Note: This method needs to be implemented based on your specific requirements
func (s *Store) GetPortfolioSummary(ctx context.Context, accountID int64) (interface{}, error) {
	// Implement the logic to get portfolio summary
	// This might involve multiple queries or a custom SQL query
	// Return an appropriate struct or map with the portfolio summary data
	return nil, fmt.Errorf("GetPortfolioSummary not implemented")
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	user, err := s.Queries.CreateUser(ctx, arg)
	if err != nil {
		log.Printf("Database error in CreateUser: %v", err) // Add this log
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *Store) VerifySchema() error {
	_, err := s.db.Exec(`
        SELECT id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at 
        FROM users 
        LIMIT 1
    `)
	if err != nil {
		return fmt.Errorf("users table schema verification failed: %w", err)
	}
	return nil
}

func (s *Store) InitializeSchema() error {
	// Check if the schema is already initialized
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if schema exists: %w", err)
	}

	if exists {
		log.Println("Database schema already initialized, skipping initialization")
		return nil
	}

	log.Println("Initializing database schema...")

	schemaSQL, err := ioutil.ReadFile("schema/postgres_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split the schema into individual statements
	statements := strings.Split(string(schemaSQL), ";")

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err = s.db.Exec(stmt)
		if err != nil {
			// Log the error but continue with other statements
			log.Printf("Error executing statement: %v\nStatement: %s", err, stmt)
		}
	}

	log.Println("Database schema initialization completed")
	return nil
}
