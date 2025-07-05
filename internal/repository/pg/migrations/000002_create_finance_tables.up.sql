BEGIN TRANSACTION;

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "type" TEXT NOT NULL CHECK (type IN ('checking', 'savings', 'credit', 'investment', 'cash')),
    "description" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "type" TEXT NOT NULL CHECK (type IN ('income', 'expense')),
    "description" TEXT,
    "color" TEXT DEFAULT '#6B7280',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    "id" UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    "account_id" UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    "category_id" UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    "amount" DECIMAL(12,2) NOT NULL,
    "description" TEXT NOT NULL,
    "date" DATE NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'cleared' CHECK (status IN ('pending', 'cleared', 'cancelled')),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create balances table for caching account balances
CREATE TABLE IF NOT EXISTS balances (
    "account_id" UUID NOT NULL PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
    "current_balance" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "pending_balance" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "available_balance" DECIMAL(12,2) NOT NULL DEFAULT 0,
    "last_calculated" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- Create unique constraint for category names within type
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name_type ON categories(name, type);

-- Insert default categories
INSERT INTO categories (name, type, description, color) VALUES
    ('Salary', 'income', 'Regular salary income', '#10B981'),
    ('Business', 'income', 'Business income', '#059669'),
    ('Investment', 'income', 'Investment returns', '#34D399'),
    ('Other Income', 'income', 'Miscellaneous income', '#6EE7B7'),
    
    ('Groceries', 'expense', 'Food and groceries', '#EF4444'),
    ('Dining Out', 'expense', 'Restaurants and takeout', '#F87171'),
    ('Transportation', 'expense', 'Gas, public transport, car maintenance', '#F59E0B'),
    ('Housing', 'expense', 'Rent, mortgage, utilities', '#8B5CF6'),
    ('Healthcare', 'expense', 'Medical expenses', '#EC4899'),
    ('Entertainment', 'expense', 'Movies, games, hobbies', '#06B6D4'),
    ('Shopping', 'expense', 'Clothes, electronics, general shopping', '#84CC16'),
    ('Credit Card', 'expense', 'Credit card payments', '#6B7280'),
    ('Other Expense', 'expense', 'Miscellaneous expenses', '#9CA3AF')
ON CONFLICT (name, type) DO NOTHING;

-- Function to update balances after transaction changes
CREATE OR REPLACE FUNCTION update_account_balance(account_uuid UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO balances (account_id, current_balance, pending_balance, available_balance, last_calculated)
    SELECT 
        account_uuid,
        COALESCE(SUM(CASE WHEN status = 'cleared' THEN amount ELSE 0 END), 0) as current_balance,
        COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_balance,
        COALESCE(SUM(CASE WHEN status IN ('cleared', 'pending') THEN amount ELSE 0 END), 0) as available_balance,
        NOW()
    FROM transactions 
    WHERE account_id = account_uuid
    ON CONFLICT (account_id) 
    DO UPDATE SET 
        current_balance = EXCLUDED.current_balance,
        pending_balance = EXCLUDED.pending_balance,
        available_balance = EXCLUDED.available_balance,
        last_calculated = EXCLUDED.last_calculated;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update balances when transactions change
CREATE OR REPLACE FUNCTION transaction_balance_trigger()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM update_account_balance(OLD.account_id);
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        PERFORM update_account_balance(NEW.account_id);
        IF OLD.account_id != NEW.account_id THEN
            PERFORM update_account_balance(OLD.account_id);
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        PERFORM update_account_balance(NEW.account_id);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER transaction_balance_update
    AFTER INSERT OR UPDATE OR DELETE ON transactions
    FOR EACH ROW EXECUTE FUNCTION transaction_balance_trigger();

COMMIT; 