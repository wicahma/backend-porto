-- Auth domain: admin_users table
CREATE TABLE IF NOT EXISTS admin_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  github_id TEXT UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
