-- Rename password_hash to password (if it exists)
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'password_hash') THEN
        ALTER TABLE users RENAME COLUMN password_hash TO password;
    END IF;
END $$;

-- Drop email column (if it exists)
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email') THEN
        ALTER TABLE users DROP COLUMN email;
    END IF;
END $$;

-- Rename phone_number to phone (if it exists)
DO $$ 
BEGIN 
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'phone_number') THEN
        ALTER TABLE users RENAME COLUMN phone_number TO phone;
    END IF;
END $$;

-- Handle name column transformation
DO $$ 
BEGIN 
    -- Add name column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'name') THEN
        ALTER TABLE users ADD COLUMN name VARCHAR(100);
        -- Combine first_name and last_name into name
        UPDATE users SET name = CONCAT(COALESCE(first_name, ''), ' ', COALESCE(last_name, ''));
        ALTER TABLE users ALTER COLUMN name SET NOT NULL;
    END IF;

    -- Drop first_name and last_name if they exist
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'first_name') THEN
        ALTER TABLE users DROP COLUMN first_name;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'last_name') THEN
        ALTER TABLE users DROP COLUMN last_name;
    END IF;
END $$;
