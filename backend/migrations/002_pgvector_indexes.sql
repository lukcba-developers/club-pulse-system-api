-- =============================================
-- Club Pulse System API
-- Migration: pgvector & Performance Indexes
-- =============================================

-- Enable pgvector extension for semantic search
CREATE EXTENSION IF NOT EXISTS vector;

-- Add embedding column to facilities table
-- 256-dimensional vector for semantic search
ALTER TABLE facilities ADD COLUMN IF NOT EXISTS embedding vector(256);

-- HNSW index for fast vector similarity search (pgvector)
-- This provides sub-millisecond search performance
CREATE INDEX IF NOT EXISTS idx_facilities_embedding 
ON facilities USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- =============================================
-- Booking Performance Optimization
-- =============================================

-- GIST index for time range overlap detection
-- This allows PostgreSQL to efficiently detect booking conflicts
-- Uses exclusion constraint pattern for bullet-proof conflict detection
CREATE INDEX IF NOT EXISTS idx_bookings_facility_time 
ON bookings USING GIST (
    facility_id,
    tstzrange(start_time, end_time, '[)') 
);

-- Composite index for common booking queries
CREATE INDEX IF NOT EXISTS idx_bookings_user_status 
ON bookings (user_id, status, start_time DESC);

-- Composite index for facility availability queries
CREATE INDEX IF NOT EXISTS idx_bookings_facility_status_time 
ON bookings (facility_id, status, start_time, end_time);

-- =============================================
-- Audit Log Optimization  
-- =============================================

-- Index for audit log queries by user
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_action 
ON audit_logs (user_id, action, timestamp DESC);

-- Partial index for recent audit logs (last 30 days optimization)
CREATE INDEX IF NOT EXISTS idx_audit_logs_recent 
ON audit_logs (timestamp DESC) 
WHERE timestamp > CURRENT_TIMESTAMP - INTERVAL '30 days';

-- =============================================
-- Function: Check Booking Overlap (Optimized)
-- =============================================

CREATE OR REPLACE FUNCTION check_booking_overlap(
    p_facility_id UUID,
    p_start_time TIMESTAMPTZ,
    p_end_time TIMESTAMPTZ,
    p_exclude_booking_id UUID DEFAULT NULL
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 
        FROM bookings 
        WHERE facility_id = p_facility_id
          AND status IN ('confirmed', 'pending')
          AND (p_exclude_booking_id IS NULL OR id != p_exclude_booking_id)
          AND tstzrange(start_time, end_time, '[)') && tstzrange(p_start_time, p_end_time, '[)')
    );
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================
-- Function: Semantic Search Facilities
-- =============================================

CREATE OR REPLACE FUNCTION search_facilities_by_embedding(
    query_embedding vector(256),
    result_limit INT DEFAULT 10
) RETURNS TABLE (
    id UUID,
    name VARCHAR,
    type VARCHAR,
    status VARCHAR,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        f.id::UUID,
        f.name::VARCHAR,
        f.type::VARCHAR,
        f.status::VARCHAR,
        1 - (f.embedding <=> query_embedding) as similarity
    FROM facilities f
    WHERE f.embedding IS NOT NULL
      AND f.status = 'active'
    ORDER BY f.embedding <=> query_embedding
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================
-- Comments
-- =============================================

COMMENT ON FUNCTION check_booking_overlap IS 
'Checks if a time range overlaps with existing bookings. Uses GIST index for O(log n) performance.';

COMMENT ON FUNCTION search_facilities_by_embedding IS 
'Semantic search using cosine similarity. Uses HNSW index for approximate nearest neighbor search.';
