BEGIN TRANSACTION;

-- Drop triggers first
DROP TRIGGER IF EXISTS transaction_balance_update ON transactions;

-- Drop functions
DROP FUNCTION IF EXISTS transaction_balance_trigger();
DROP FUNCTION IF EXISTS update_account_balance(UUID);

-- Drop indexes
DROP INDEX IF EXISTS idx_categories_name_type;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_category_id;
DROP INDEX IF EXISTS idx_transactions_account_id;

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;

COMMIT; 