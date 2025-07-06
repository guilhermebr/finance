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
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Handlers contains all web handlers for the personal finance application
type Handlers struct {
	apiBaseURL string
	httpClient *http.Client
	templates  *template.Template
}

// NewHandlers creates a new instance of web handlers
func NewHandlers(apiBaseURL string) *Handlers {
	// Load templates
	templates := template.Must(template.ParseGlob("internal/web/templates/*.html"))

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
	var accounts []entities.Account
	var categories []entities.Category
	var transactions []entities.Transaction
	var balances []entities.Balance

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
		http.Error(w, fmt.Sprintf("Failed to get balances: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Accounts     []entities.Account
		Categories   []entities.Category
		Transactions []entities.Transaction
		Balances     []entities.Balance
		Title        string
	}{
		Accounts:     accounts,
		Categories:   categories,
		Transactions: transactions,
		Balances:     balances,
		Title:        "Personal Finance Dashboard",
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AccountsPage renders the accounts management page
func (h *Handlers) AccountsPage(w http.ResponseWriter, r *http.Request) {
	var accounts []entities.Account

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Accounts []entities.Account
		Title    string
	}{
		Accounts: accounts,
		Title:    "Manage Accounts",
	}

	if err := h.templates.ExecuteTemplate(w, "accounts.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateAccount handles account creation
func (h *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := entities.Account{
		Name:        r.FormValue("name"),
		Type:        entities.AccountType(r.FormValue("type")),
		Description: r.FormValue("description"),
	}

	var createdAccount entities.Account
	if err := h.apiPost("/api/v1/accounts", account, &createdAccount); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusBadRequest)
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
	}{
		Accounts: accounts,
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

	account := entities.Account{
		ID:          id,
		Name:        r.FormValue("name"),
		Type:        entities.AccountType(r.FormValue("type")),
		Description: r.FormValue("description"),
	}

	var updatedAccount entities.Account
	if err := h.apiPut("/api/v1/accounts/"+id, account, &updatedAccount); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update account: %v", err), http.StatusBadRequest)
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
	}{
		Accounts: accounts,
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
	}{
		Accounts: accounts,
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("account-deleted-%s", id))
}

// CategoriesPage renders the categories management page
func (h *Handlers) CategoriesPage(w http.ResponseWriter, r *http.Request) {
	var categories []entities.Category

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []entities.Category
		Title      string
	}{
		Categories: categories,
		Title:      "Manage Categories",
	}

	if err := h.templates.ExecuteTemplate(w, "categories.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCategory handles category creation
func (h *Handlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := entities.Category{
		Name:        r.FormValue("name"),
		Type:        entities.CategoryType(r.FormValue("type")),
		Color:       r.FormValue("color"),
		Description: r.FormValue("description"),
	}

	var createdCategory entities.Category
	if err := h.apiPost("/api/v1/categories", category, &createdCategory); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create category: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	var categories []entities.Category
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []entities.Category
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

	category := entities.Category{
		ID:          id,
		Name:        r.FormValue("name"),
		Type:        entities.CategoryType(r.FormValue("type")),
		Color:       r.FormValue("color"),
		Description: r.FormValue("description"),
	}

	var updatedCategory entities.Category
	if err := h.apiPut("/api/v1/categories/"+id, category, &updatedCategory); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update category: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	var categories []entities.Category
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []entities.Category
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
	var categories []entities.Category
	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []entities.Category
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
	var transactions []entities.Transaction
	var accounts []entities.Account
	var categories []entities.Category

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
		Transactions []entities.Transaction
		Accounts     []entities.Account
		Categories   []entities.Category
		Title        string
	}{
		Transactions: transactions,
		Accounts:     accounts,
		Categories:   categories,
		Title:        "Manage Transactions",
	}

	if err := h.templates.ExecuteTemplate(w, "transactions.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateTransaction handles transaction creation
func (h *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	transaction := entities.Transaction{
		AccountID:   r.FormValue("account_id"),
		CategoryID:  r.FormValue("category_id"),
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        date,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	var createdTransaction entities.Transaction
	if err := h.apiPost("/api/v1/transactions", transaction, &createdTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	var transactions []entities.Transaction
	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []entities.Transaction
	}{
		Transactions: transactions,
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

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	transaction := entities.Transaction{
		ID:          id,
		AccountID:   r.FormValue("account_id"),
		CategoryID:  r.FormValue("category_id"),
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        date,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	var updatedTransaction entities.Transaction
	if err := h.apiPut("/api/v1/transactions/"+id, transaction, &updatedTransaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	var transactions []entities.Transaction
	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []entities.Transaction
	}{
		Transactions: transactions,
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
	var transactions []entities.Transaction
	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []entities.Transaction
	}{
		Transactions: transactions,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf("transaction-deleted-%s", id))
}

// AccountsTable renders the accounts table partial for HTMX
func (h *Handlers) AccountsTable(w http.ResponseWriter, r *http.Request) {
	var accounts []entities.Account

	if err := h.apiGet("/api/v1/accounts", &accounts); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Accounts []entities.Account
	}{
		Accounts: accounts,
	}

	if err := h.templates.ExecuteTemplate(w, "accounts-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CategoriesTable renders the categories table partial for HTMX
func (h *Handlers) CategoriesTable(w http.ResponseWriter, r *http.Request) {
	var categories []entities.Category

	if err := h.apiGet("/api/v1/categories", &categories); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get categories: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []entities.Category
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
	var transactions []entities.Transaction

	if err := h.apiGet("/api/v1/transactions", &transactions); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transactions: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Transactions []entities.Transaction
	}{
		Transactions: transactions,
	}

	if err := h.templates.ExecuteTemplate(w, "transactions-table.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BalanceSummary renders the balance summary partial for HTMX
func (h *Handlers) BalanceSummary(w http.ResponseWriter, r *http.Request) {
	var balances []entities.Balance

	if err := h.apiGet("/api/v1/balances", &balances); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get balances: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Balances []entities.Balance
	}{
		Balances: balances,
	}

	if err := h.templates.ExecuteTemplate(w, "balance-summary.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Template helper functions
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
