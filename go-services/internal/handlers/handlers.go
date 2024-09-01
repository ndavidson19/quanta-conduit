package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"conduit/internal/api"
	"conduit/internal/db" // Import the db generated code
	"conduit/internal/store"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

// Implement the ServerInterface
var _ api.ServerInterface = (*Handler)(nil)

// Users

func (h *Handler) GetAccount(ctx echo.Context, accountId int) error {
	account, err := h.store.GetAccount(ctx.Request().Context(), int32(accountId))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Account not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get account")
	}

	id := int(account.ID)
	userId := int(account.UserID.Int32)
	accountIdInt := int(account.AccountID)
	balance := account.Balance
	isSystemAccount := account.IsSystemAccount.Bool
	accountType := account.AccountType.String
	createdAt := account.CreatedAt.Time
	updatedAt := account.UpdatedAt.Time

	response := api.Account{
		Id:              &id,
		UserId:          &userId,
		AccountId:       &accountIdInt,
		Balance:         &balance,
		AccountType:     &accountType,
		IsSystemAccount: &isSystemAccount,
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) ListUsers(ctx echo.Context, params api.ListUsersParams) error {
	users, err := h.store.ListUsers(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list users")
	}

	response := make([]api.User, len(users))
	for i, user := range users {
		id := int(user.ID)
		emailVerified := user.EmailVerified.Bool
		createdAt := user.CreatedAt.Time
		updatedAt := user.UpdatedAt.Time
		provider := user.Provider.String
		response[i] = api.User{
			Id:            &id,
			Email:         &user.Email,
			Provider:      &provider,
			EmailVerified: &emailVerified,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateUser(ctx echo.Context) error {
	var req api.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	// Prepare the CreateUserParams
	createParams := db.CreateUserParams{
		Email:           req.Email,
		HashedPassword:  sql.NullString{String: string(hashedPassword), Valid: true},
		Provider:        sql.NullString{String: "", Valid: false}, // Set default as empty
		EmailVerified:   sql.NullBool{Bool: false, Valid: true},   // Default to false
		TwoFactorSecret: sql.NullString{String: "", Valid: false}, // Default to empty
	}

	// If provider is provided, set it
	if req.Provider != nil {
		createParams.Provider = sql.NullString{String: *req.Provider, Valid: true}
	}

	user, err := h.store.CreateUser(ctx.Request().Context(), createParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation
			return echo.NewHTTPError(http.StatusConflict, "User with this email already exists")
		}
		log.Printf("Failed to create user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	// Prepare the response
	id := int(user.ID)
	emailVerified := user.EmailVerified.Bool
	createdAt := user.CreatedAt.Time
	updatedAt := user.UpdatedAt.Time
	response := api.User{
		Id:            &id,
		Email:         &user.Email,
		Provider:      &user.Provider.String,
		EmailVerified: &emailVerified,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (h *Handler) GetUser(ctx echo.Context, userId int) error {
	user, err := h.store.GetUser(ctx.Request().Context(), int32(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	id := int(user.ID)
	emailVerified := user.EmailVerified.Bool
	createdAt := user.CreatedAt.Time
	updatedAt := user.UpdatedAt.Time
	response := api.User{
		Id:            &id,
		Email:         &user.Email,
		Provider:      &user.Provider.String,
		EmailVerified: &emailVerified,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteUser(ctx echo.Context, userId int) error {
	err := h.store.DeleteUser(ctx.Request().Context(), int32(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Accounts

func (h *Handler) ListAccounts(ctx echo.Context, params api.ListAccountsParams) error {
	var userID sql.NullInt32
	if params.UserId != nil {
		userID = sql.NullInt32{Int32: int32(*params.UserId), Valid: true}
	}

	accounts, err := h.store.ListAccountsByUserID(ctx.Request().Context(), userID)
	if err != nil {
		log.Printf("Failed to list accounts: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list accounts: %v", err))
	}

	response := make([]api.Account, len(accounts))
	for i, account := range accounts {
		id := int(account.ID)
		userId := int(account.UserID.Int32)
		accountId := int(account.AccountID)
		isSystemAccount := account.IsSystemAccount.Bool
		createdAt := account.CreatedAt.Time
		updatedAt := account.UpdatedAt.Time

		response[i] = api.Account{
			Id:              &id,
			UserId:          &userId,
			AccountId:       &accountId,
			Balance:         &account.Balance,
			AccountType:     &account.AccountType.String,
			IsSystemAccount: &isSystemAccount,
			CreatedAt:       &createdAt,
			UpdatedAt:       &updatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateAccount(ctx echo.Context) error {
	var req api.CreateAccountRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	account, err := h.store.CreateAccount(ctx.Request().Context(), db.CreateAccountParams{
		UserID:      sql.NullInt32{Int32: int32(req.UserId), Valid: true},
		AccountID:   int64(req.UserId), // You might want to generate this differently
		AccountType: sql.NullString{String: req.AccountType, Valid: true},
		Balance:     "0", // Start with zero balance
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation
			return echo.NewHTTPError(http.StatusConflict, "Account already exists for this user")
		}
		log.Printf("Failed to create account: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create account")
	}
	id := int(account.ID)
	userId := int(account.UserID.Int32)
	accountId := int(account.AccountID)
	isSystemAccount := account.IsSystemAccount.Bool
	createdAt := account.CreatedAt.Time
	updatedAt := account.UpdatedAt.Time

	response := api.Account{
		Id:              &id,
		UserId:          &userId,
		AccountId:       &accountId,
		Balance:         &account.Balance,
		AccountType:     &account.AccountType.String,
		IsSystemAccount: &isSystemAccount,
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// Orders

func (h *Handler) ListOrders(ctx echo.Context, params api.ListOrdersParams) error {
	var accountID sql.NullInt64
	var status sql.NullString

	if params.AccountId != nil {
		accountID = sql.NullInt64{Int64: int64(*params.AccountId), Valid: true}
	}
	if params.Status != nil {
		status = sql.NullString{String: string(*params.Status), Valid: true}
	}

	orders, err := h.store.ListOrders(ctx.Request().Context(), store.ListOrdersParams{
		AccountID: accountID,
		Status:    status,
		Limit:     int32(*params.Limit),
		Offset:    int32(*params.Offset),
	})
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list orders: %v", err))
	}

	response := make([]api.Order, len(orders))
	for i, order := range orders {
		id := int(order.ID)
		accountId := int(order.AccountID)
		assetId := int(order.AssetID)
		orderType := api.OrderOrderType(order.OrderType)
		side := api.OrderSide(order.Side)
		status := api.OrderStatus(order.Status)
		createdAt := order.CreatedAt.Time
		updatedAt := order.UpdatedAt.Time

		response[i] = api.Order{
			Id:        &id,
			AccountId: &accountId,
			AssetId:   &assetId,
			OrderType: &orderType,
			Side:      &side,
			Quantity:  &order.Quantity,
			Price:     &order.Price.String,
			Status:    &status,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateOrder(ctx echo.Context) error {
	var req api.CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		log.Printf("Failed to bind create order request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
	}

	// Add request validation
	if req.AccountId == 0 || req.AssetId == 0 || req.Quantity == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	// Check if the asset exists
	_, err := h.store.GetAsset(ctx.Request().Context(), int32(req.AssetId))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Asset not found")
		}
		log.Printf("Failed to get asset: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check asset")
	}

	order, err := h.store.CreateOrder(ctx.Request().Context(), db.CreateOrderParams{
		AccountID: int64(req.AccountId),
		AssetID:   int32(req.AssetId),
		OrderType: string(req.OrderType),
		Side:      string(req.Side),
		Quantity:  string(req.Quantity),
		Price:     sql.NullString{String: *req.Price, Valid: req.Price != nil},
		Status:    "pending",
	})
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create order")
	}

	id := int(order.ID)
	accountId := int(order.AccountID)
	assetId := int(order.AssetID)
	orderType := api.OrderOrderType(order.OrderType)
	side := api.OrderSide(order.Side)
	status := api.OrderStatus(order.Status)
	createdAt := order.CreatedAt.Time
	updatedAt := order.UpdatedAt.Time

	response := api.Order{
		Id:        &id,
		AccountId: &accountId,
		AssetId:   &assetId,
		OrderType: &orderType,
		Side:      &side,
		Quantity:  &order.Quantity,
		Price:     &order.Price.String,
		Status:    &status,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (h *Handler) CancelOrder(ctx echo.Context, orderId int) error {
	_, err := h.store.CancelOrder(ctx.Request().Context(), int32(orderId))
	if err != nil {
		log.Printf("Failed to cancel order: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to cancel order: %v", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}

// Portfolios

func (h *Handler) CreatePortfolio(ctx echo.Context) error {
	var req api.CreatePortfolioRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	portfolio, err := h.store.CreatePortfolio(ctx.Request().Context(), db.CreatePortfolioParams{
		AccountID:  int64(req.AccountId),
		TotalValue: "0", // Start with zero total value
	})
	if err != nil {
		log.Printf("Failed to create portfolio: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create portfolio")
	}

	id := int(portfolio.ID)
	accountId := int(portfolio.AccountID)
	createdAt := portfolio.CreatedAt.Time
	updatedAt := portfolio.UpdatedAt.Time

	response := api.Portfolio{
		Id:         &id,
		AccountId:  &accountId,
		TotalValue: &portfolio.TotalValue,
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// Add other handlers as needed...
func (h *Handler) GetOrder(ctx echo.Context, orderId int) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) GetPortfolioSummary(ctx echo.Context, accountId int) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) GetPortfolio(ctx echo.Context, accountId int) error {

	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) ListAccountTransactions(ctx echo.Context, accountId int, params api.ListAccountTransactionsParams) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) ListAssets(ctx echo.Context, params api.ListAssetsParams) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) ListPortfolioHoldings(ctx echo.Context, accountId int) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) ListPortfolios(ctx echo.Context, params api.ListPortfoliosParams) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) UpdateUser(ctx echo.Context, userId int) error {
	// Implementation here
	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}
