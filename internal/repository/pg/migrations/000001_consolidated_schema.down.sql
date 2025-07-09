BEGIN TRANSACTION;

-- =============================================================================
-- DROP TRIGGERS AND FUNCTIONS
-- =============================================================================
DROP TRIGGER IF EXISTS transaction_balance_update ON transactions;
DROP FUNCTION IF EXISTS transaction_balance_trigger();
DROP FUNCTION IF EXISTS update_account_balance(UUID);

-- =============================================================================
-- DROP INDEXES
-- =============================================================================
DROP INDEX IF EXISTS idx_categories_name_type;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_category_id;
DROP INDEX IF EXISTS idx_transactions_account_id;

-- =============================================================================
-- DROP TABLES (in reverse dependency order)
-- =============================================================================
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS examples;

COMMIT; 