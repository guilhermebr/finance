package web

import (
	"bytes"
	"encoding/json"
	"finance/domain/entities"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/guilhermebr/gox/monetary"
)

// Response DTOs that match the API contracts
type AccountResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Asset       string               `json:"asset"`
	Description string               `json:"description"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type CategoryResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Type        entities.CategoryType `json:"type"`
	Description string                `json:"description"`
	Color       string                `json:"color"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

type TransactionResponse struct {
	ID          string                     `json:"id"`
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      string                     `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
	CreatedAt   string                     `json:"created_at"`
	UpdatedAt   string                     `json:"updated_at"`
	Account     *AccountResponse           `json:"account,omitempty"`
	Category    *CategoryResponse          `json:"category,omitempty"`
}

type BalanceResponse struct {
	AccountID        string           `json:"account_id"`
	CurrentBalance   string           `json:"current_balance"`
	PendingBalance   string           `json:"pending_balance"`
	AvailableBalance string           `json:"available_balance"`
	LastCalculated   string           `json:"last_calculated"`
	Account          *AccountResponse `json:"account,omitempty"`
}

type BalanceSummaryResponse struct {
	TotalAssets      string `json:"total_assets"`
	TotalLiabilities string `json:"total_liabilities"`
	NetWorth         string `json:"net_worth"`
	LastCalculated   string `json:"last_calculated"`
}

// Handlers contains all web handlers for the personal finance application
type Handlers struct {
	apiBaseURL string
	httpClient *http.Client
	templates  *template.Template
}

// NewHandlers creates a new instance of web handlers
func NewHandlers(apiBaseURL string) *Handlers {
	// Load templates individually to avoid naming conflicts
	templates := template.New("")

	// Parse each template file individually
	templateFiles := map[string]string{
		"dashboard.html":          "internal/web/templates/dashboard.html",
		"accounts.html":           "internal/web/templates/accounts.html",
		"categories.html":         "internal/web/templates/categories.html",
		"transactions.html":       "internal/web/templates/transactions.html",
		"accounts-table.html":     "internal/web/templates/accounts-table.html",
		"categories-table.html":   "internal/web/templates/categories-table.html",
		"transactions-table.html": "internal/web/templates/transactions-table.html",
		"balance-summary.html":    "internal/web/templates/balance-summary.html",
	}

	for name, file := range templateFiles {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse template %s: %v", file, err))
		}
		// Associate each template with its intended name
		templates, _ = templates.AddParseTree(name, tmpl.Tree)
	}

	return &Handlers{
		apiBaseURL: apiBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		templates: templates,
	}
}

// Router returns the HTTP router for the web application
func (h *Handlers) Router() http.Handler {
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static/"))))

	// Web routes
	r.HandleFunc("/", h.Dashboard).Methods("GET")
	r.HandleFunc("/accounts", h.AccountsPage).Methods("GET")
	r.HandleFunc("/accounts/create", h.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", h.UpdateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{id}", h.DeleteAccount).Methods("DELETE")

	r.HandleFunc("/categories", h.CategoriesPage).Methods("GET")
	r.HandleFunc("/categories/create", h.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id}", h.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", h.DeleteCategory).Methods("DELETE")

	r.HandleFunc("/transactions", h.TransactionsPage).Methods("GET")
	r.HandleFunc("/transactions/create", h.CreateTransaction).Methods("POST")
	r.HandleFunc("/transactions/{id}", h.UpdateTransaction).Methods("PUT")
	r.HandleFunc("/transactions/{id}", h.DeleteTransaction).Methods("DELETE")

	// HTMX partial routes
	r.HandleFunc("/htmx/accounts", h.AccountsTable).Methods("GET")
	r.HandleFunc("/htmx/categories", h.CategoriesTable).Methods("GET")
	r.HandleFunc("/htmx/transactions", h.TransactionsTable).Methods("GET")
	r.HandleFunc("/htmx/balance-summary", h.BalanceSummary).Methods("GET")

	return r
}

// Helper method to make GET requests to the API
func (h *Handlers) apiGet(endpoint string, result interface{}) error {
	url := h.apiBaseURL + endpoint
	resp, err := h.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// Helper method to make POST requests to the API
func (h *Handlers) apiPost(endpoint string, payload interface{}, result interface{}) error {
	url := h.apiBaseURL + endpoint

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := h.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// Helper method to make PUT requests to the API
func (h *Handlers) apiPut(endpoint string, payload interface{}, result interface{}) error {
	url := h.apiBaseURL + endpoint

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// Helper method to make DELETE requests to the API
func (h *Handlers) apiDelete(endpoint string) error {
	url := h.apiBaseURL + endpoint

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Dashboard renders the main dashboard page
func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	var accounts []AccountResponse
	var categories []CategoryResponse
	var transactions []TransactionResponse
	var balances []BalanceResponse

	// Get data from API
	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		// Don't fail if balances can't be loaded, just use empty slice
		balances = []BalanceResponse{}
	}

	data := struct {
		Accounts     []AccountResponse
		Categories   []CategoryResponse
		Transactions []TransactionResponse
		Balances     []BalanceResponse
		Title        string
		CurrentPage  string
	}{
		Accounts:     accounts,
		Categories:   categories,
		Transactions: transactions,
		Balances:     balances,
		Title:        "Personal Finance Dashboard",
		CurrentPage:  "dashboard",
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AccountsPage renders the accounts management page
func (h *Handlers) AccountsPage(w http.ResponseWriter, r *http.Request) {
	var accounts []AccountResponse

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Accounts    []AccountResponse
		Title       string
		CurrentPage string
	}{
		Accounts:    accounts,
		Title:       "Manage Accounts",
		CurrentPage: "accounts",
	}

	if err := h.templates.ExecuteTemplate(w, "accounts.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateAccount handles account creation
func (h *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the asset
	assetName := r.FormValue("asset")
	if assetName == "" {
		assetName = "BRL" // Default to BRL if no asset is provided
	}

	asset, ok := monetary.FindAssetByName(assetName)
	if !ok {
		http.Error(w, "Invalid asset", http.StatusBadRequest)
		return
	}

	// Create request payload that matches API expectations
	requestPayload := struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Asset       string `json:"asset"`
		Description string `json:"description"`
	}{
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Asset:       asset.Asset,
		Description: r.FormValue("description"),
	}

	var createdAccount AccountResponse
	if err := h.apiPost("/api/v1/accounts", requestPayload, &createdAccount); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	var accounts []AccountResponse
	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	var balances []BalanceResponse
	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		// Don't fail if balances can't be loaded, just use empty slice
		balances = []BalanceResponse{}
	}

	data := struct {
		Accounts []AccountResponse
		Balances []BalanceResponse
	}{
		Accounts: accounts,
		Balances: balances,
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("account-created-%s", createdAccount.ID))
}

// UpdateAccount handles account updates
func (h *Handlers) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Parse and validate the asset
	assetName := r.FormValue("asset")
	if assetName == "" {
		assetName = "BRL" // Default to BRL if no asset is provided
	}

	asset, ok := monetary.FindAssetByName(assetName)
	if !ok {
		http.Error(w, "Invalid asset", http.StatusBadRequest)
		return
	}

	// Create request payload that matches API expectations
	requestPayload := struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Asset       string `json:"asset"`
		Description string `json:"description"`
	}{
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Asset:       asset.Asset,
		Description: r.FormValue("description"),
	}

	var updatedAccount AccountResponse
	if err := h.apiPut("/api/v1/accounts/"+id, requestPayload, &updatedAccount); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update account: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	var accounts []AccountResponse
	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	var balances []BalanceResponse
	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		// Don't fail if balances can't be loaded, just use empty slice
		balances = []BalanceResponse{}
	}

	data := struct {
		Accounts []AccountResponse
		Balances []BalanceResponse
	}{
		Accounts: accounts,
		Balances: balances,
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("account-updated-%s", updatedAccount.ID))
}

// DeleteAccount handles account deletion
func (h *Handlers) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	if err := h.apiDelete("/api/v1/accounts/" + id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete account: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	var accounts []entities.Account
	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Accounts []entities.Account
		Balances []entities.Balance
	}{
		Accounts: accounts,
		Balances: []entities.Balance{}, // Empty for now due to API issue
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("account-deleted-%s", id))
}

// CategoriesPage renders the categories management page
func (h *Handlers) CategoriesPage(w http.ResponseWriter, r *http.Request) {
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories  []CategoryResponse
		Title       string
		CurrentPage string
	}{
		Categories:  categories,
		Title:       "Manage Categories",
		CurrentPage: "categories",
	}

	if err := h.templates.ExecuteTemplate(w, "categories.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCategory handles category creation
func (h *Handlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Create request payload that matches API expectations
	requestPayload := struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}{
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Color:       r.FormValue("color"),
		Description: r.FormValue("description"),
	}

	var createdCategory CategoryResponse
	if err := h.apiPost("/api/v1/categories", requestPayload, &createdCategory); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create category: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	var categories []CategoryResponse
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []CategoryResponse
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("category-created-%s", createdCategory.ID))
}

// UpdateCategory handles category updates
func (h *Handlers) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Create request payload that matches API expectations
	requestPayload := struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Color       string `json:"color"`
		Description string `json:"description"`
	}{
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Color:       r.FormValue("color"),
		Description: r.FormValue("description"),
	}

	var updatedCategory CategoryResponse
	if err := h.apiPut("/api/v1/categories/"+id, requestPayload, &updatedCategory); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update category: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	var categories []CategoryResponse
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []CategoryResponse
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("category-updated-%s", updatedCategory.ID))
}

// DeleteCategory handles category deletion
func (h *Handlers) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.apiDelete("/api/v1/categories/" + id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete category: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	var categories []CategoryResponse
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []CategoryResponse
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("category-deleted-%s", id))
}

// TransactionsPage renders the transactions management page
func (h *Handlers) TransactionsPage(w http.ResponseWriter, r *http.Request) {
	var transactions []TransactionResponse
	var accounts []AccountResponse
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []TransactionResponse
		Accounts     []AccountResponse
		Categories   []CategoryResponse
		Title        string
		CurrentPage  string
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
		Title:        "Manage Transactions",
		CurrentPage:  "transactions",
	}

	if err := h.templates.ExecuteTemplate(w, "transactions.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateTransaction handles transaction creation
func (h *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	amountStr := r.FormValue("amount")
	// Validate amount format by trying to parse it as float
	if _, err := strconv.ParseFloat(amountStr, 64); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Validate date format but send as string to match API expectations
	dateStr := r.FormValue("transaction_date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Create request payload that matches API expectations
	requestPayload := struct {
		AccountID   string                     `json:"account_id"`
		CategoryID  string                     `json:"category_id"`
		Amount      string                     `json:"amount"`
		Description string                     `json:"description"`
		Date        string                     `json:"date"`
		Status      entities.TransactionStatus `json:"status"`
	}{
		AccountID:   r.FormValue("account_id"),
		CategoryID:  r.FormValue("category_id"),
		Amount:      amountStr,
		Description: r.FormValue("description"),
		Date:        dateStr,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	var createdTransaction TransactionResponse
	if err := h.apiPost("/api/v1/transactions", requestPayload, &createdTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	var transactions []TransactionResponse
	var accounts []AccountResponse
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []TransactionResponse
		Accounts     []AccountResponse
		Categories   []CategoryResponse
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("transaction-created-%s", createdTransaction.ID))
}

// UpdateTransaction handles transaction updates
func (h *Handlers) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	// Validate amount format by trying to parse it as float
	if _, err := strconv.ParseFloat(amountStr, 64); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Validate date format but send as string to match API expectations
	dateStr := r.FormValue("transaction_date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	// Create request payload that matches API expectations
	requestPayload := struct {
		AccountID   string                     `json:"account_id"`
		CategoryID  string                     `json:"category_id"`
		Amount      string                     `json:"amount"`
		Description string                     `json:"description"`
		Date        string                     `json:"date"`
		Status      entities.TransactionStatus `json:"status"`
	}{
		AccountID:   r.FormValue("account_id"),
		CategoryID:  r.FormValue("category_id"),
		Amount:      amountStr,
		Description: r.FormValue("description"),
		Date:        dateStr,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	var updatedTransaction TransactionResponse
	if err := h.apiPut("/api/v1/transactions/"+id, requestPayload, &updatedTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	var transactions []TransactionResponse
	var accounts []AccountResponse
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []TransactionResponse
		Accounts     []AccountResponse
		Categories   []CategoryResponse
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("transaction-updated-%s", updatedTransaction.ID))
}

// DeleteTransaction handles transaction deletion
func (h *Handlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.apiDelete("/api/v1/transactions/" + id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	var transactions []TransactionResponse
	var accounts []AccountResponse
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []TransactionResponse
		Accounts     []AccountResponse
		Categories   []CategoryResponse
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("transaction-deleted-%s", id))
}

// AccountsTable renders the accounts table partial for HTMX
func (h *Handlers) AccountsTable(w http.ResponseWriter, r *http.Request) {
	var accounts []AccountResponse
	var balances []BalanceResponse

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		// Don't fail if balances can't be loaded, just use empty slice
		balances = []BalanceResponse{}
	}

	data := struct {
		Accounts []AccountResponse
		Balances []BalanceResponse
	}{
		Accounts: accounts,
		Balances: balances,
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CategoriesTable renders the categories table partial for HTMX
func (h *Handlers) CategoriesTable(w http.ResponseWriter, r *http.Request) {
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []CategoryResponse
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TransactionsTable renders the transactions table partial for HTMX
func (h *Handlers) TransactionsTable(w http.ResponseWriter, r *http.Request) {
	var transactions []TransactionResponse
	var accounts []AccountResponse
	var categories []CategoryResponse

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []TransactionResponse
		Accounts     []AccountResponse
		Categories   []CategoryResponse
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BalanceSummary renders the balance summary partial for HTMX
func (h *Handlers) BalanceSummary(w http.ResponseWriter, r *http.Request) {
	var balances []BalanceResponse

	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		// Don't fail if balances can't be loaded, just use empty slice
		balances = []BalanceResponse{}
	}

	data := struct {
		Balances []BalanceResponse
	}{
		Balances: balances,
	}

	if err := h.templates.ExecuteTemplate(w, "balance-summary.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
