-- Experiences table (company-level)
CREATE TABLE IF NOT EXISTS experiences (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company TEXT NOT NULL,
  location TEXT NOT NULL,
  total_duration TEXT NOT NULL,
  tags TEXT[] DEFAULT '{}',
  logo TEXT,
  is_current BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Experience jobs table (individual positions within a company)
CREATE TABLE IF NOT EXISTS experience_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  experience_id UUID NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
  position TEXT NOT NULL,
  employment_type TEXT NOT NULL,
  description TEXT NOT NULL,
  competency TEXT[] DEFAULT '{}',
  start_date TEXT NOT NULL,
  end_date TEXT,
  is_current BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_experiences_company ON experiences(company);
CREATE INDEX IF NOT EXISTS idx_experiences_location ON experiences(location);
CREATE INDEX IF NOT EXISTS idx_experience_jobs_experience_id ON experience_jobs(experience_id);
CREATE INDEX IF NOT EXISTS idx_experience_jobs_start_date ON experience_jobs(start_date DESC);
