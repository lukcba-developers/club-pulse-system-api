-- =============================================
-- Club Pulse System API
-- Consolidated Initial Schema
-- =============================================

-- Enable Extensions
SET search_path TO public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

-- =============================================
-- Core Identity & User Management
-- =============================================

-- Family Groups (Created before users for reference)
CREATE TABLE IF NOT EXISTS family_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    head_user_id TEXT NOT NULL, -- Reference to users(id) logically, effectively TEXT
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY, 
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    role VARCHAR(50),
    date_of_birth DATE,
    sports_preferences JSONB,
    -- Medical / Health
    medical_cert_status VARCHAR(20) DEFAULT 'PENDING',
    medical_cert_expiry TIMESTAMP WITH TIME ZONE,
    -- Family
    family_group_id UUID REFERENCES family_groups(id) ON DELETE SET NULL,
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- Refresh Tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_id VARCHAR(255),
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);

-- Auth Logs
CREATE TABLE IF NOT EXISTS authentication_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type VARCHAR(50), 
    ip_address VARCHAR(50),
    user_agent TEXT,
    success BOOLEAN,
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Access Logs
CREATE TABLE IF NOT EXISTS access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    facility_id UUID,
    direction VARCHAR(10) NOT NULL, 
    status VARCHAR(50) NOT NULL, 
    reason VARCHAR(255),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_access_logs_user_id ON access_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_access_logs_timestamp ON access_logs (timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_recent ON access_logs (timestamp DESC) WHERE timestamp > CURRENT_TIMESTAMP - INTERVAL '30 days';


-- =============================================
-- Membership & Billing
-- =============================================

-- Membership Tiers
CREATE TABLE IF NOT EXISTS membership_tiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    monthly_fee DECIMAL(10,2) NOT NULL,
    colors VARCHAR(50),
    benefits TEXT[],
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Memberships
CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    membership_tier_id UUID NOT NULL REFERENCES membership_tiers(id),
    status VARCHAR(50) DEFAULT 'PENDING',
    billing_cycle VARCHAR(50) DEFAULT 'MONTHLY',
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    next_billing_date TIMESTAMP WITH TIME ZONE NOT NULL,
    outstanding_balance DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships (user_id);
CREATE INDEX IF NOT EXISTS idx_memberships_status ON memberships (status);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'ARS' NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING' NOT NULL,
    method VARCHAR(50) NOT NULL,
    external_id VARCHAR(255),
    payer_id TEXT NOT NULL,
    reference_id UUID,
    reference_type VARCHAR(50),
    paid_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_payments_payer_id ON payments (payer_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_external_id ON payments (external_id);

-- Wallets
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id),
    balance DECIMAL(10, 2) DEFAULT 0,
    points INT DEFAULT 0,
    transactions JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets (user_id);

-- Subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id),
    membership_id UUID NOT NULL REFERENCES memberships(id),
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'ARS',
    status VARCHAR(20) NOT NULL,
    payment_method_id TEXT,
    next_billing_date DATE,
    last_payment_date TIMESTAMP WITH TIME ZONE,
    fail_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);

-- Scholarships
CREATE TABLE IF NOT EXISTS scholarships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id),
    percentage DECIMAL(5,2) NOT NULL, -- 0.50, 1.00
    reason VARCHAR(255),
    grantor_id TEXT, -- Admin User ID
    valid_until DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_scholarships_user ON scholarships (user_id);

-- =============================================
-- Sports & Training
-- =============================================

-- Disciplines
CREATE TABLE IF NOT EXISTS disciplines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_disciplines_club_id ON disciplines (club_id);

-- Training Groups
CREATE TABLE IF NOT EXISTS training_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    discipline_id UUID NOT NULL REFERENCES disciplines(id),
    category VARCHAR(20) NOT NULL,
    category_year INT,
    coach_id VARCHAR(255),
    schedule VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_training_groups_club_id ON training_groups (club_id);
CREATE INDEX IF NOT EXISTS idx_training_groups_discipline_id ON training_groups (discipline_id);

-- Attendance Lists
CREATE TABLE IF NOT EXISTS attendance_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    group_name VARCHAR(50) NOT NULL,
    coach_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_attendance_lists_date_group ON attendance_lists (date, group_name);

-- Attendance Records
CREATE TABLE IF NOT EXISTS attendance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attendance_list_id UUID NOT NULL REFERENCES attendance_lists(id),
    user_id TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    notes TEXT,
    UNIQUE(attendance_list_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_attendance_records_list_id ON attendance_records (attendance_list_id);

-- =============================================
-- Facilities & Bookings
-- =============================================

-- Facilities
CREATE TABLE IF NOT EXISTS facilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    location VARCHAR(255),
    capacity INT,
    price_per_hour DECIMAL(10, 2),
    amenities JSONB,
    images TEXT[],
    is_active BOOLEAN DEFAULT TRUE,
    embedding vector(256), -- pgvector support
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_facilities_club_id ON facilities (club_id);
-- HNSW Index for Facilities Semantic Search
CREATE INDEX IF NOT EXISTS idx_facilities_embedding 
ON facilities USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Equipment
CREATE TABLE IF NOT EXISTS equipment (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    facility_id UUID REFERENCES facilities(id),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    condition VARCHAR(50),
    status VARCHAR(50),
    is_available BOOLEAN DEFAULT TRUE,
    purchase_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_equipment_facility_id ON equipment (facility_id);

-- Equipment Loans
CREATE TABLE IF NOT EXISTS equipment_loans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    equipment_id UUID NOT NULL REFERENCES equipment(id),
    user_id TEXT NOT NULL REFERENCES users(id),
    loaned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expected_return_at TIMESTAMP WITH TIME ZONE NOT NULL,
    returned_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, RETURNED, OVERDUE, LOST
    condition_on_return TEXT, -- "Good", "Damaged"
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_loans_user ON equipment_loans (user_id);
CREATE INDEX IF NOT EXISTS idx_loans_status ON equipment_loans (status);

-- Bookings
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id),
    facility_id UUID NOT NULL REFERENCES facilities(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'CONFIRMED',
    guest_details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_bookings_club_id ON bookings (club_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_facility_id ON bookings (facility_id);

-- Optimized Booking Indexes
CREATE INDEX IF NOT EXISTS idx_bookings_user_status ON bookings (user_id, status, start_time DESC);
CREATE INDEX IF NOT EXISTS idx_bookings_facility_status_time ON bookings (facility_id, status, start_time, end_time);

-- GIST Index for Overlap Detection
CREATE INDEX IF NOT EXISTS idx_bookings_facility_time 
ON bookings USING GIST (
    facility_id,
    tstzrange(start_time, end_time, '[)') 
);

-- Waitlists
CREATE TABLE IF NOT EXISTS waitlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    resource_id UUID NOT NULL,
    target_date TIMESTAMP WITH TIME ZONE NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_waitlists_resource_date ON waitlists(resource_id, target_date);
CREATE INDEX IF NOT EXISTS idx_waitlists_club_date ON waitlists(club_id, target_date);


-- =============================================
-- Functions & Procedures
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
-- Operational Features (Real-World Pack)
-- =============================================

-- 1. STORE
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INT DEFAULT 0,
    sku VARCHAR(100),
    category VARCHAR(50), -- Merch, Buffet, Equipment
    is_active BOOLEAN DEFAULT TRUE,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_products_club ON products(club_id);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'PAID', -- PAID, PENDING, CANCELLED
    items JSONB NOT NULL, -- [{product_id, qty, unit_price}]
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_club ON orders(club_id);

-- 2. TEAM AVAILABILITY
CREATE TABLE IF NOT EXISTS match_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    training_group_id UUID NOT NULL REFERENCES training_groups(id),
    opponent_name VARCHAR(100),
    location VARCHAR(255), -- "Home", "Away", or specific address
    is_home_game BOOLEAN DEFAULT TRUE,
    meetup_time TIMESTAMP WITH TIME ZONE NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE, -- Optional match start if diff from meetup
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_match_events_group ON match_events(training_group_id);

CREATE TABLE IF NOT EXISTS player_availabilities (
    match_event_id UUID NOT NULL REFERENCES match_events(id),
    user_id TEXT NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL, -- CONFIRMED, DECLINED, MAYBE
    reason TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (match_event_id, user_id)
);

-- 3. SPONSORS
CREATE TABLE IF NOT EXISTS sponsors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    contact_info TEXT, -- JSON or Text
    logo_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_sponsors_club ON sponsors(club_id);

CREATE TABLE IF NOT EXISTS ad_placements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sponsor_id UUID NOT NULL REFERENCES sponsors(id),
    location_type VARCHAR(50), -- WEBSITE_BANNER, PHYSICAL_BANNER, JERSEY
    location_detail VARCHAR(255), -- "Cancha 1 - Fondo"
    contract_start DATE,
    contract_end DATE NOT NULL,
    amount_paid DECIMAL(10,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ad_placements_sponsor ON ad_placements(sponsor_id);

-- 4. INCIDENTS (Y Update Users)
-- Adding columns to users table is handled by migration tool or manually here if safe
-- We use IF NOT EXISTS to be safe in this consolidated file

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS emergency_contact_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS emergency_contact_phone VARCHAR(50),
ADD COLUMN IF NOT EXISTS insurance_provider VARCHAR(100),
ADD COLUMN IF NOT EXISTS insurance_number VARCHAR(100);

CREATE TABLE IF NOT EXISTS incident_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    injured_user_id TEXT REFERENCES users(id), -- Nullable for visitors
    description TEXT NOT NULL,
    witnesses TEXT,
    action_taken TEXT, -- "Ambulance Called", "First Aid"
    reported_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT, -- Staff User ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_incidents_club ON incident_logs(club_id);
