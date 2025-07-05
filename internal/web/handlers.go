package web

import (
	"context"
	"finance/domain/entities"
	"finance/domain/finance"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Handlers contains all web handlers for the personal finance application
type Handlers struct {
	accountUseCase     *finance.AccountUseCase
	categoryUseCase    *finance.CategoryUseCase
	transactionUseCase *finance.TransactionUseCase
	balanceUseCase     *finance.BalanceUseCase
	templates          *template.Template
}

// NewHandlers creates a new instance of web handlers
func NewHandlers(
	accountUseCase *finance.AccountUseCase,
	categoryUseCase *finance.CategoryUseCase,
	transactionUseCase *finance.TransactionUseCase,
	balanceUseCase *finance.BalanceUseCase,
) *Handlers {
	// Load templates
	templates := template.Must(template.ParseGlob("internal/web/templates/*.html"))

	return &Handlers{
		accountUseCase:     accountUseCase,
		categoryUseCase:    categoryUseCase,
		transactionUseCase: transactionUseCase,
		balanceUseCase:     balanceUseCase,
		templates:          templates,
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

// Dashboard renders the main dashboard page
func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get summary data
	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balances, err := h.balanceUseCase.GetAllBalances(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	ctx := context.Background()

	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	ctx := context.Background()

	account := entities.Account{
		Name:        r.FormValue("name"),
		Type:        entities.AccountType(r.FormValue("type")),
		Description: r.FormValue("description"),
	}

	createdAccount, err := h.accountUseCase.CreateAccount(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	ctx := context.Background()
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

	_, err := h.accountUseCase.UpdateAccount(ctx, account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// DeleteAccount handles account deletion
func (h *Handlers) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	err := h.accountUseCase.DeleteAccount(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated accounts table for HTMX
	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// CategoriesPage renders the categories management page
func (h *Handlers) CategoriesPage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	ctx := context.Background()

	category := entities.Category{
		Name:        r.FormValue("name"),
		Type:        entities.CategoryType(r.FormValue("type")),
		Color:       r.FormValue("color"),
		Description: r.FormValue("description"),
	}

	_, err := h.categoryUseCase.CreateCategory(ctx, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// UpdateCategory handles category updates
func (h *Handlers) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
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

	_, err := h.categoryUseCase.UpdateCategory(ctx, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// DeleteCategory handles category deletion
func (h *Handlers) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err := h.categoryUseCase.DeleteCategory(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated categories table for HTMX
	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// TransactionsPage renders the transactions management page
func (h *Handlers) TransactionsPage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	ctx := context.Background()

	accountID := r.FormValue("account_id")
	if accountID == "" {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	transactionDate, err := time.Parse("2006-01-02", r.FormValue("transaction_date"))
	if err != nil {
		http.Error(w, "Invalid transaction date", http.StatusBadRequest)
		return
	}

	transaction := entities.Transaction{
		AccountID:   accountID,
		CategoryID:  categoryID,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        transactionDate,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	_, err = h.transactionUseCase.CreateTransaction(ctx, transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// UpdateTransaction handles transaction updates
func (h *Handlers) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	accountID := r.FormValue("account_id")
	if accountID == "" {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	transactionDate, err := time.Parse("2006-01-02", r.FormValue("transaction_date"))
	if err != nil {
		http.Error(w, "Invalid transaction date", http.StatusBadRequest)
		return
	}

	transaction := entities.Transaction{
		ID:          id,
		AccountID:   accountID,
		CategoryID:  categoryID,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        transactionDate,
		Status:      entities.TransactionStatus(r.FormValue("status")),
	}

	_, err = h.transactionUseCase.UpdateTransaction(ctx, transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// DeleteTransaction handles transaction deletion
func (h *Handlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err := h.transactionUseCase.DeleteTransaction(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated transactions table for HTMX
	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// HTMX partial handlers

// AccountsTable returns just the accounts table for HTMX updates
func (h *Handlers) AccountsTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	accounts, err := h.accountUseCase.GetAllAccounts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// CategoriesTable returns just the categories table for HTMX updates
func (h *Handlers) CategoriesTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// TransactionsTable returns just the transactions table for HTMX updates
func (h *Handlers) TransactionsTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	transactions, err := h.transactionUseCase.GetAllTransactions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// BalanceSummary returns balance summary for HTMX updates
func (h *Handlers) BalanceSummary(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	balances, err := h.balanceUseCase.GetAllBalances(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Helper functions

// formatCurrency formats a float64 as currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// formatDate formats a time.Time as a date string
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
