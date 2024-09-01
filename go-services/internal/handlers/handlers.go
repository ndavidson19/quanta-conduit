package handlers

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    "conduit/internal/api"
    "conduit/internal/store"
    "conduit/internal/db"  // Import the db generated code
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
    // Note: We're not using limit and offset here as the store method doesn't support pagination
    // You might want to implement pagination in your store method if needed
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
        response[i] = api.User{
            Id:            &id,
            Email:         &user.Email,
            Provider:      &user.Provider.String,
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

    user, err := h.store.CreateUser(ctx.Request().Context(), db.CreateUserParams{
        Email:    req.Email,
        Provider: sql.NullString{String: *req.Provider, Valid: req.Provider != nil},
    })
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
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
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) CreateAccount(ctx echo.Context) error {
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

// Orders

func (h *Handler) ListOrders(ctx echo.Context, params api.ListOrdersParams) error {
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) CreateOrder(ctx echo.Context) error {
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) CancelOrder(ctx echo.Context, orderId int) error {
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}

// Portfolios

func (h *Handler) CreatePortfolio(ctx echo.Context) error {
    // Implementation here
    return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
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

