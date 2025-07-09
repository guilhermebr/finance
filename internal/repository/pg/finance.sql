-- =============================================================================
-- ACCOUNTS
-- =============================================================================

-- name: CreateAccount :one
INSERT INTO accounts (name, type, description, asset)
VALUES ($1, $2, $3, $4)
RETURNING id, name, type, description, asset, created_at, updated_at;

-- name: GetAccountByID :one
SELECT id, name, type, description, asset, created_at, updated_at
FROM accounts
WHERE id = $1;

-- name: GetAllAccounts :many
SELECT id, name, type, description, asset, created_at, updated_at
FROM accounts
ORDER BY name;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $2, type = $3, description = $4, asset = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, name, type, description, asset, created_at, updated_at;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- =============================================================================
-- CATEGORIES
-- =============================================================================

-- name: CreateCategory :one
INSERT INTO categories (name, type, description, color)
VALUES ($1, $2, $3, $4)
RETURNING id, name, type, description, color, created_at, updated_at;

-- name: GetCategoryByID :one
SELECT id, name, type, description, color, created_at, updated_at
FROM categories
WHERE id = $1;

-- name: GetAllCategories :many
SELECT id, name, type, description, color, created_at, updated_at
FROM categories
ORDER BY type, name;

-- name: GetCategoriesByType :many
SELECT id, name, type, description, color, created_at, updated_at
FROM categories
WHERE type = $1
ORDER BY name;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, type = $3, description = $4, color = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, name, type, description, color, created_at, updated_at;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- =============================================================================
-- TRANSACTIONS
-- =============================================================================

-- name: CreateTransaction :one
INSERT INTO transactions (account_id, category_id, amount, description, date, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, account_id, category_id, amount, description, date, status, created_at, updated_at;

-- name: GetTransactionByID :one
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
WHERE id = $1;

-- name: GetAllTransactions :many
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByAccount :many
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
WHERE account_id = $1
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByCategory :many
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
WHERE category_id = $1
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByDateRange :many
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
WHERE date >= $1 AND date <= $2
ORDER BY date DESC, created_at DESC;

-- name: GetTransactionsByAccountAndDateRange :many
SELECT id, account_id, category_id, amount, description, date, status, created_at, updated_at
FROM transactions
WHERE account_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC, created_at DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET account_id = $2, category_id = $3, amount = $4, description = $5, date = $6, status = $7, updated_at = NOW()
WHERE id = $1
RETURNING id, account_id, category_id, amount, description, date, status, created_at, updated_at;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, account_id, category_id, amount, description, date, status, created_at, updated_at;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = $1;

-- =============================================================================
-- BALANCES
-- =============================================================================

-- name: GetBalanceByAccountID :one
SELECT account_id, current_balance, pending_balance, available_balance, last_calculated
FROM balances
WHERE account_id = $1;

-- name: GetAllBalances :many
SELECT account_id, current_balance, pending_balance, available_balance, last_calculated
FROM balances
ORDER BY account_id;

-- name: RefreshAccountBalance :exec
SELECT update_account_balance($1);

-- name: GetBalanceSummary :one
SELECT 
    COALESCE(SUM(CASE WHEN a.type IN ('checking', 'savings', 'investment', 'cash') THEN b.current_balance ELSE 0 END), 0) as total_assets,
    COALESCE(SUM(CASE WHEN a.type = 'credit' THEN ABS(b.current_balance) ELSE 0 END), 0) as total_liabilities,
    COALESCE(SUM(CASE WHEN a.type IN ('checking', 'savings', 'investment', 'cash') THEN b.current_balance ELSE -ABS(b.current_balance) END), 0) as net_worth,
    NOW() as last_calculated
FROM balances b
JOIN accounts a ON b.account_id = a.id;

-- =============================================================================
-- JOINED QUERIES FOR DETAILED VIEWS
-- =============================================================================

-- name: GetTransactionWithDetails :one
SELECT 
    t.id, t.account_id, t.category_id, t.amount, t.description, t.date, t.status, t.created_at, t.updated_at,
    a.name as account_name, a.type as account_type, a.asset as account_asset,
    c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE t.id = $1;

-- name: GetTransactionsWithDetails :many
SELECT 
    t.id, t.account_id, t.category_id, t.amount, t.description, t.date, t.status, t.created_at, t.updated_at,
    a.name as account_name, a.type as account_type, a.asset as account_asset,
    c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
ORDER BY t.date DESC, t.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAccountWithBalance :one
SELECT 
    a.id, a.name, a.type, a.description, a.asset, a.created_at, a.updated_at,
    COALESCE(b.current_balance, 0) as current_balance,
    COALESCE(b.pending_balance, 0) as pending_balance,
    COALESCE(b.available_balance, 0) as available_balance
FROM accounts a
LEFT JOIN balances b ON a.id = b.account_id
WHERE a.id = $1;

-- name: GetAccountsWithBalances :many
SELECT 
    a.id, a.name, a.type, a.description, a.asset, a.created_at, a.updated_at,
    COALESCE(b.current_balance, 0) as current_balance,
    COALESCE(b.pending_balance, 0) as pending_balance,
    COALESCE(b.available_balance, 0) as available_balance
FROM accounts a
LEFT JOIN balances b ON a.id = b.account_id
ORDER BY a.name; 