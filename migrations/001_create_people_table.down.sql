-- Drop trigger
DROP TRIGGER IF EXISTS update_people_updated_at ON people;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop index
DROP INDEX IF EXISTS idx_people_name;

-- Drop table
DROP TABLE IF EXISTS people;

-- Made with Bob
