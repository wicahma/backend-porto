-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE articles ADD COLUMN IF NOT EXISTS embedding vector(768);

CREATE TABLE IF NOT EXISTS topic_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    original_url TEXT,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    embedding vector(768),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_topic_queue_status ON topic_queue(status);
CREATE INDEX IF NOT EXISTS idx_topic_queue_embedding ON topic_queue USING hnsw (embedding vector_cosine_ops);
CREATE INDEX IF NOT EXISTS idx_articles_embedding ON articles USING hnsw (embedding vector_cosine_ops);

-- +goose Down
DROP TABLE IF EXISTS topic_queue;
ALTER TABLE articles DROP COLUMN IF EXISTS embedding;
DROP EXTENSION IF EXISTS vector;
