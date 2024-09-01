// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for AssetAssetType.
const (
	AssetAssetTypeCrypto AssetAssetType = "crypto"
	AssetAssetTypeForex  AssetAssetType = "forex"
	AssetAssetTypeStock  AssetAssetType = "stock"
)

// Defines values for CreateOrderRequestOrderType.
const (
	CreateOrderRequestOrderTypeLimit  CreateOrderRequestOrderType = "limit"
	CreateOrderRequestOrderTypeMarket CreateOrderRequestOrderType = "market"
)

// Defines values for CreateOrderRequestSide.
const (
	CreateOrderRequestSideBuy  CreateOrderRequestSide = "buy"
	CreateOrderRequestSideSell CreateOrderRequestSide = "sell"
)

// Defines values for OrderOrderType.
const (
	OrderOrderTypeLimit  OrderOrderType = "limit"
	OrderOrderTypeMarket OrderOrderType = "market"
)

// Defines values for OrderSide.
const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// Defines values for OrderStatus.
const (
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusExecuted  OrderStatus = "executed"
	OrderStatusPending   OrderStatus = "pending"
)

// Defines values for ListAssetsParamsAssetType.
const (
	ListAssetsParamsAssetTypeCrypto ListAssetsParamsAssetType = "crypto"
	ListAssetsParamsAssetTypeForex  ListAssetsParamsAssetType = "forex"
	ListAssetsParamsAssetTypeStock  ListAssetsParamsAssetType = "stock"
)

// Defines values for ListOrdersParamsStatus.
const (
	ListOrdersParamsStatusCancelled ListOrdersParamsStatus = "cancelled"
	ListOrdersParamsStatusExecuted  ListOrdersParamsStatus = "executed"
	ListOrdersParamsStatusPending   ListOrdersParamsStatus = "pending"
)

// Account defines model for Account.
type Account struct {
	AccountId       *int       `json:"accountId,omitempty"`
	AccountType     *string    `json:"accountType,omitempty"`
	Balance         *string    `json:"balance,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	Id              *int       `json:"id,omitempty"`
	IsSystemAccount *bool      `json:"isSystemAccount,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	UserId          *int       `json:"userId,omitempty"`
}

// Asset defines model for Asset.
type Asset struct {
	AssetType *AssetAssetType `json:"assetType,omitempty"`
	CreatedAt *time.Time      `json:"createdAt,omitempty"`
	Id        *int            `json:"id,omitempty"`
	Name      *string         `json:"name,omitempty"`
	Symbol    *string         `json:"symbol,omitempty"`
	UpdatedAt *time.Time      `json:"updatedAt,omitempty"`
}

// AssetAssetType defines model for Asset.AssetType.
type AssetAssetType string

// CreateAccountRequest defines model for CreateAccountRequest.
type CreateAccountRequest struct {
	AccountType string `json:"accountType"`
	UserId      int    `json:"userId"`
}

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	AccountId int                         `json:"accountId"`
	AssetId   int                         `json:"assetId"`
	OrderType CreateOrderRequestOrderType `json:"orderType"`
	Price     *string                     `json:"price,omitempty"`
	Quantity  string                      `json:"quantity"`
	Side      CreateOrderRequestSide      `json:"side"`
}

// CreateOrderRequestOrderType defines model for CreateOrderRequest.OrderType.
type CreateOrderRequestOrderType string

// CreateOrderRequestSide defines model for CreateOrderRequest.Side.
type CreateOrderRequestSide string

// CreatePortfolioRequest defines model for CreatePortfolioRequest.
type CreatePortfolioRequest struct {
	AccountId int `json:"accountId"`
}

// CreateUserRequest defines model for CreateUserRequest.
type CreateUserRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Provider *string `json:"provider,omitempty"`
}

// Order defines model for Order.
type Order struct {
	AccountId *int            `json:"accountId,omitempty"`
	AssetId   *int            `json:"assetId,omitempty"`
	CreatedAt *time.Time      `json:"createdAt,omitempty"`
	Id        *int            `json:"id,omitempty"`
	OrderType *OrderOrderType `json:"orderType,omitempty"`
	Price     *string         `json:"price,omitempty"`
	Quantity  *string         `json:"quantity,omitempty"`
	Side      *OrderSide      `json:"side,omitempty"`
	Status    *OrderStatus    `json:"status,omitempty"`
	UpdatedAt *time.Time      `json:"updatedAt,omitempty"`
}

// OrderOrderType defines model for Order.OrderType.
type OrderOrderType string

// OrderSide defines model for Order.Side.
type OrderSide string

// OrderStatus defines model for Order.Status.
type OrderStatus string

// Portfolio defines model for Portfolio.
type Portfolio struct {
	AccountId  *int       `json:"accountId,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	Id         *int       `json:"id,omitempty"`
	TotalValue *string    `json:"totalValue,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

// PortfolioHolding defines model for PortfolioHolding.
type PortfolioHolding struct {
	AssetId         *int       `json:"assetId,omitempty"`
	AverageBuyPrice *string    `json:"averageBuyPrice,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	Id              *int       `json:"id,omitempty"`
	PortfolioId     *int       `json:"portfolioId,omitempty"`
	Quantity        *string    `json:"quantity,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

// Transaction defines model for Transaction.
type Transaction struct {
	Amount          *string    `json:"amount,omitempty"`
	CreditAccountId *int       `json:"creditAccountId,omitempty"`
	DebitAccountId  *int       `json:"debitAccountId,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Id              *string    `json:"id,omitempty"`
	Timestamp       *time.Time `json:"timestamp,omitempty"`
}

// UpdateUserRequest defines model for UpdateUserRequest.
type UpdateUserRequest struct {
	Email         *string `json:"email,omitempty"`
	EmailVerified *bool   `json:"emailVerified,omitempty"`
	Password      *string `json:"password,omitempty"`
	Provider      *string `json:"provider,omitempty"`
}

// User defines model for User.
type User struct {
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	Email         *string    `json:"email,omitempty"`
	EmailVerified *bool      `json:"emailVerified,omitempty"`
	Id            *int       `json:"id,omitempty"`
	Provider      *string    `json:"provider,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

// ListAccountsParams defines parameters for ListAccounts.
type ListAccountsParams struct {
	UserId *int `form:"userId,omitempty" json:"userId,omitempty"`
	Limit  *int `form:"limit,omitempty" json:"limit,omitempty"`
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListAccountTransactionsParams defines parameters for ListAccountTransactions.
type ListAccountTransactionsParams struct {
	Limit  *int `form:"limit,omitempty" json:"limit,omitempty"`
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListAssetsParams defines parameters for ListAssets.
type ListAssetsParams struct {
	AssetType *ListAssetsParamsAssetType `form:"assetType,omitempty" json:"assetType,omitempty"`
	Limit     *int                       `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int                       `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListAssetsParamsAssetType defines parameters for ListAssets.
type ListAssetsParamsAssetType string

// ListOrdersParams defines parameters for ListOrders.
type ListOrdersParams struct {
	AccountId *int                    `form:"accountId,omitempty" json:"accountId,omitempty"`
	Status    *ListOrdersParamsStatus `form:"status,omitempty" json:"status,omitempty"`
	Limit     *int                    `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int                    `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListOrdersParamsStatus defines parameters for ListOrders.
type ListOrdersParamsStatus string

// ListPortfoliosParams defines parameters for ListPortfolios.
type ListPortfoliosParams struct {
	AccountId *int `form:"accountId,omitempty" json:"accountId,omitempty"`
	Limit     *int `form:"limit,omitempty" json:"limit,omitempty"`
	Offset    *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// ListUsersParams defines parameters for ListUsers.
type ListUsersParams struct {
	Limit  *int `form:"limit,omitempty" json:"limit,omitempty"`
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`
}

// CreateAccountJSONRequestBody defines body for CreateAccount for application/json ContentType.
type CreateAccountJSONRequestBody = CreateAccountRequest

// CreateOrderJSONRequestBody defines body for CreateOrder for application/json ContentType.
type CreateOrderJSONRequestBody = CreateOrderRequest

// CreatePortfolioJSONRequestBody defines body for CreatePortfolio for application/json ContentType.
type CreatePortfolioJSONRequestBody = CreatePortfolioRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserRequest

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = UpdateUserRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List accounts
	// (GET /accounts)
	ListAccounts(ctx echo.Context, params ListAccountsParams) error
	// Create an account
	// (POST /accounts)
	CreateAccount(ctx echo.Context) error
	// Get an account
	// (GET /accounts/{accountId})
	GetAccount(ctx echo.Context, accountId int) error
	// List transactions for an account
	// (GET /accounts/{accountId}/transactions)
	ListAccountTransactions(ctx echo.Context, accountId int, params ListAccountTransactionsParams) error
	// List assets
	// (GET /assets)
	ListAssets(ctx echo.Context, params ListAssetsParams) error
	// List orders
	// (GET /orders)
	ListOrders(ctx echo.Context, params ListOrdersParams) error
	// Create an order
	// (POST /orders)
	CreateOrder(ctx echo.Context) error
	// Cancel an order
	// (DELETE /orders/{orderId})
	CancelOrder(ctx echo.Context, orderId int) error
	// Get an order
	// (GET /orders/{orderId})
	GetOrder(ctx echo.Context, orderId int) error
	// List portfolios
	// (GET /portfolios)
	ListPortfolios(ctx echo.Context, params ListPortfoliosParams) error
	// Create a portfolio
	// (POST /portfolios)
	CreatePortfolio(ctx echo.Context) error
	// Get a portfolio
	// (GET /portfolios/{portfolioId})
	GetPortfolio(ctx echo.Context, portfolioId int) error
	// List portfolio holdings
	// (GET /portfolios/{portfolioId}/holdings)
	ListPortfolioHoldings(ctx echo.Context, portfolioId int) error
	// List users
	// (GET /users)
	ListUsers(ctx echo.Context, params ListUsersParams) error
	// Create a user
	// (POST /users)
	CreateUser(ctx echo.Context) error
	// Delete a user
	// (DELETE /users/{userId})
	DeleteUser(ctx echo.Context, userId int) error
	// Get a user
	// (GET /users/{userId})
	GetUser(ctx echo.Context, userId int) error
	// Update a user
	// (PUT /users/{userId})
	UpdateUser(ctx echo.Context, userId int) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListAccounts converts echo context to params.
func (w *ServerInterfaceWrapper) ListAccounts(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListAccountsParams
	// ------------- Optional query parameter "userId" -------------

	err = runtime.BindQueryParameter("form", true, false, "userId", ctx.QueryParams(), &params.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListAccounts(ctx, params)
	return err
}

// CreateAccount converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAccount(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateAccount(ctx)
	return err
}

// GetAccount converts echo context to params.
func (w *ServerInterfaceWrapper) GetAccount(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "accountId" -------------
	var accountId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "accountId", runtime.ParamLocationPath, ctx.Param("accountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAccount(ctx, accountId)
	return err
}

// ListAccountTransactions converts echo context to params.
func (w *ServerInterfaceWrapper) ListAccountTransactions(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "accountId" -------------
	var accountId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "accountId", runtime.ParamLocationPath, ctx.Param("accountId"), &accountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListAccountTransactionsParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListAccountTransactions(ctx, accountId, params)
	return err
}

// ListAssets converts echo context to params.
func (w *ServerInterfaceWrapper) ListAssets(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListAssetsParams
	// ------------- Optional query parameter "assetType" -------------

	err = runtime.BindQueryParameter("form", true, false, "assetType", ctx.QueryParams(), &params.AssetType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter assetType: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListAssets(ctx, params)
	return err
}

// ListOrders converts echo context to params.
func (w *ServerInterfaceWrapper) ListOrders(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListOrdersParams
	// ------------- Optional query parameter "accountId" -------------

	err = runtime.BindQueryParameter("form", true, false, "accountId", ctx.QueryParams(), &params.AccountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	// ------------- Optional query parameter "status" -------------

	err = runtime.BindQueryParameter("form", true, false, "status", ctx.QueryParams(), &params.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter status: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListOrders(ctx, params)
	return err
}

// CreateOrder converts echo context to params.
func (w *ServerInterfaceWrapper) CreateOrder(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateOrder(ctx)
	return err
}

// CancelOrder converts echo context to params.
func (w *ServerInterfaceWrapper) CancelOrder(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "orderId" -------------
	var orderId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "orderId", runtime.ParamLocationPath, ctx.Param("orderId"), &orderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter orderId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CancelOrder(ctx, orderId)
	return err
}

// GetOrder converts echo context to params.
func (w *ServerInterfaceWrapper) GetOrder(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "orderId" -------------
	var orderId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "orderId", runtime.ParamLocationPath, ctx.Param("orderId"), &orderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter orderId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetOrder(ctx, orderId)
	return err
}

// ListPortfolios converts echo context to params.
func (w *ServerInterfaceWrapper) ListPortfolios(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListPortfoliosParams
	// ------------- Optional query parameter "accountId" -------------

	err = runtime.BindQueryParameter("form", true, false, "accountId", ctx.QueryParams(), &params.AccountId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter accountId: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListPortfolios(ctx, params)
	return err
}

// CreatePortfolio converts echo context to params.
func (w *ServerInterfaceWrapper) CreatePortfolio(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreatePortfolio(ctx)
	return err
}

// GetPortfolio converts echo context to params.
func (w *ServerInterfaceWrapper) GetPortfolio(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "portfolioId" -------------
	var portfolioId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "portfolioId", runtime.ParamLocationPath, ctx.Param("portfolioId"), &portfolioId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter portfolioId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPortfolio(ctx, portfolioId)
	return err
}

// ListPortfolioHoldings converts echo context to params.
func (w *ServerInterfaceWrapper) ListPortfolioHoldings(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "portfolioId" -------------
	var portfolioId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "portfolioId", runtime.ParamLocationPath, ctx.Param("portfolioId"), &portfolioId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter portfolioId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListPortfolioHoldings(ctx, portfolioId)
	return err
}

// ListUsers converts echo context to params.
func (w *ServerInterfaceWrapper) ListUsers(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListUsersParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListUsers(ctx, params)
	return err
}

// CreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateUser(ctx)
	return err
}

// DeleteUser converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteUser(ctx, userId)
	return err
}

// GetUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUser(ctx, userId)
	return err
}

// UpdateUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateUser(ctx, userId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/accounts", wrapper.ListAccounts)
	router.POST(baseURL+"/accounts", wrapper.CreateAccount)
	router.GET(baseURL+"/accounts/:accountId", wrapper.GetAccount)
	router.GET(baseURL+"/accounts/:accountId/transactions", wrapper.ListAccountTransactions)
	router.GET(baseURL+"/assets", wrapper.ListAssets)
	router.GET(baseURL+"/orders", wrapper.ListOrders)
	router.POST(baseURL+"/orders", wrapper.CreateOrder)
	router.DELETE(baseURL+"/orders/:orderId", wrapper.CancelOrder)
	router.GET(baseURL+"/orders/:orderId", wrapper.GetOrder)
	router.GET(baseURL+"/portfolios", wrapper.ListPortfolios)
	router.POST(baseURL+"/portfolios", wrapper.CreatePortfolio)
	router.GET(baseURL+"/portfolios/:portfolioId", wrapper.GetPortfolio)
	router.GET(baseURL+"/portfolios/:portfolioId/holdings", wrapper.ListPortfolioHoldings)
	router.GET(baseURL+"/users", wrapper.ListUsers)
	router.POST(baseURL+"/users", wrapper.CreateUser)
	router.DELETE(baseURL+"/users/:userId", wrapper.DeleteUser)
	router.GET(baseURL+"/users/:userId", wrapper.GetUser)
	router.PUT(baseURL+"/users/:userId", wrapper.UpdateUser)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"HBaHfwMAAP//lOQeVZ4pAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}