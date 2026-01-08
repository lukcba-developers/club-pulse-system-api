-- Initial Schema for Club Pulse System API
-- Consolidated Migration (001-009)

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clubs Table (for multi-tenancy)
CREATE TABLE IF NOT EXISTS clubs (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    settings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Added from 005_club_public_fields.sql
    slug VARCHAR(255) NOT NULL UNIQUE,
    logo_url TEXT,
    theme_config JSONB
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_clubs_slug ON clubs(slug);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(100) PRIMARY KEY,
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    family_group_id UUID,
    medical_cert_status VARCHAR(50),
    medical_cert_expiry DATE,
    -- Added from 002_user_documents.sql
    is_eligible BOOLEAN DEFAULT FALSE,
    -- GDPR Compliance Fields (from 003_gdpr_compliance.sql)
    terms_accepted_at TIMESTAMPTZ,
    privacy_policy_version VARCHAR(20),
    data_retention_until TIMESTAMPTZ,
    UNIQUE(club_id, email)
);
COMMENT ON COLUMN users.is_eligible IS 'Cached eligibility status based on document validation (updated by background job)';
COMMENT ON COLUMN users.terms_accepted_at IS 'GDPR: Timestamp when user accepted Terms and Conditions';
COMMENT ON COLUMN users.privacy_policy_version IS 'GDPR: Version of Privacy Policy accepted (e.g., 2026-01)';
COMMENT ON COLUMN users.data_retention_until IS 'GDPR: Date until data should be retained before automatic deletion';

-- Family Groups Table (Merged from 010_family_groups.sql)
CREATE TABLE IF NOT EXISTS family_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    head_user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_family_groups_club ON family_groups(club_id);
CREATE INDEX IF NOT EXISTS idx_family_groups_head ON family_groups(head_user_id);

COMMENT ON TABLE family_groups IS 'Groups multiple users (parent + children) for consolidated billing';
COMMENT ON COLUMN family_groups.head_user_id IS 'The primary account holder responsible for billing';


-- User Documents Table (From 002_user_documents.sql)
CREATE TABLE IF NOT EXISTS user_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    file_url TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    expiration_date DATE,
    rejection_notes TEXT,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    validated_at TIMESTAMPTZ,
    validated_by VARCHAR(100) REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_user_documents_user ON user_documents(club_id, user_id);
CREATE INDEX idx_user_documents_status ON user_documents(status);
CREATE INDEX idx_user_documents_expiration ON user_documents(expiration_date) WHERE expiration_date IS NOT NULL;
CREATE INDEX idx_user_documents_type ON user_documents(type);
CREATE INDEX idx_user_documents_user_type ON user_documents(club_id, user_id, type);

COMMENT ON TABLE user_documents IS 'Stores user documents (DNI, medical certificates, insurance, etc.) with expiration tracking';
COMMENT ON COLUMN user_documents.type IS 'Document type: DNI_FRONT, DNI_BACK, EMMAC_MEDICAL, LEAGUE_FORM, INSURANCE';
COMMENT ON COLUMN user_documents.status IS 'Document status: PENDING, VALID, REJECTED, EXPIRED';
COMMENT ON COLUMN user_documents.expiration_date IS 'Date when the document expires (nullable for documents without expiration)';


-- User Stats Table
CREATE TABLE IF NOT EXISTS user_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    matches_played INT DEFAULT 0,
    matches_won INT DEFAULT 0,
    ranking_points INT DEFAULT 0,
    level INT DEFAULT 1,
    experience INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Wallet Table
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(10, 2) DEFAULT 0.0,
    points INT DEFAULT 0,
    transactions JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Sessions Table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255),
    token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Facilities Table
CREATE TABLE IF NOT EXISTS facilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    status VARCHAR(50),
    capacity INT,
    hourly_rate DECIMAL(10, 2),
    -- Added from 009_add_booking_pricing.sql
    guest_fee DECIMAL(10, 2) DEFAULT 0,
    specifications JSONB,
    location JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bookings Table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    facility_id UUID NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    -- Added from 009_add_booking_pricing.sql
    total_price DECIMAL(10, 2) DEFAULT 0,
    guest_details JSONB,
    -- Security Fix (VUL-001)
    payment_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_bookings_payment_expiry ON bookings(payment_expiry) WHERE payment_expiry IS NOT NULL;

-- Membership Tiers Table
CREATE TABLE IF NOT EXISTS membership_tiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    benefits JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Subscriptions (User Memberships) Table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier_id UUID NOT NULL REFERENCES membership_tiers(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payments Table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    payer_id VARCHAR(100) NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    method VARCHAR(50),
    reference_id UUID,
    reference_type VARCHAR(100),
    external_id VARCHAR(255),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Added from 008_add_club_id_to_payments.sql
    club_id TEXT
);
CREATE INDEX idx_payments_club_id ON payments(club_id);

-- Scholarships
CREATE TABLE IF NOT EXISTS scholarships (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    percentage DECIMAL(5, 2) NOT NULL,
    reason TEXT,
    grantor_id VARCHAR(100) REFERENCES users(id),
    valid_until TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_scholarships_user ON scholarships (user_id);

-- Disciplines
CREATE TABLE IF NOT EXISTS disciplines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMptz NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Training Groups
CREATE TABLE IF NOT EXISTS training_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    discipline_id UUID NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    coach_id VARCHAR(100) REFERENCES users(id),
    schedule TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Attendance Lists
CREATE TABLE IF NOT EXISTS attendance_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    training_group_id UUID NOT NULL REFERENCES training_groups(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    coach_id VARCHAR(100) REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(training_group_id, date)
);

-- Attendance Records
CREATE TABLE IF NOT EXISTS attendance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attendance_list_id UUID NOT NULL REFERENCES attendance_lists(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(attendance_list_id, user_id)
);

-- Equipment
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(100) PRIMARY KEY,
    facility_id UUID REFERENCES facilities(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    condition VARCHAR(50),
    status VARCHAR(50),
    is_available BOOLEAN DEFAULT TRUE,
    purchase_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Equipment Loans
CREATE TABLE IF NOT EXISTS equipment_loans (
    id VARCHAR(100) PRIMARY KEY,
    equipment_id VARCHAR(100) NOT NULL REFERENCES equipment(id),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id),
    loaned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expected_return_at TIMESTAMPTZ NOT NULL,
    returned_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL,
    condition_on_return VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Waitlists
CREATE TABLE IF NOT EXISTS waitlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id),
    user_id VARCHAR(100) NOT NULL REFERENCES users(id),
    resource_id UUID NOT NULL, -- e.g., facility_id
    resource_type VARCHAR(50) NOT NULL, -- e.g., "facility"
    target_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_waitlists_resource_date ON waitlists(resource_id, target_date);
CREATE INDEX IF NOT EXISTS idx_waitlists_club_date ON waitlists(club_id, target_date);

-- Store Products
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INT DEFAULT 0,
    sku VARCHAR(100),
    category VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    image_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Store Orders
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id),
    -- Modified from 007_store_guest_orders.sql: user_id is nullable
    user_id VARCHAR(100) REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'PAID',
    items JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    -- Added from 007_store_guest_orders.sql
    guest_name VARCHAR(255),
    guest_email VARCHAR(255)
);

-- Championships Table
CREATE TABLE IF NOT EXISTS championships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    discipline_id UUID NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL, -- e.g., DRAFT, OPEN_FOR_REGISTRATION, IN_PROGRESS, COMPLETED
    format VARCHAR(50) NOT NULL, -- e.g., LEAGUE, KNOCKOUT
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Teams Table
CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    championship_id UUID NOT NULL REFERENCES championships(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    captain_id VARCHAR(100) REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Team Members Table (Junction Table)
CREATE TABLE IF NOT EXISTS team_members (
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
);

-- Matches Table
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    championship_id UUID NOT NULL REFERENCES championships(id) ON DELETE CASCADE,
    home_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    score_home INT,
    score_away INT,
    start_time TIMESTAMPTZ NOT NULL,
    location VARCHAR(255),
    round VARCHAR(100),
    status VARCHAR(50) NOT NULL, -- e.g., SCHEDULED, PLAYED, POSTPONED
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Travel Events (From 003_travel_events.sql)
CREATE TABLE IF NOT EXISTS travel_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL,
    team_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'TRAVEL',
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Detalles del viaje
    destination VARCHAR(255) NOT NULL,
    departure_date TIMESTAMPTZ NOT NULL,
    return_date TIMESTAMPTZ,
    meeting_point VARCHAR(255),
    meeting_time TIMESTAMPTZ NOT NULL,
    
    -- Costos
    estimated_cost DECIMAL(10,2) DEFAULT 0,
    actual_cost DECIMAL(10,2) DEFAULT 0,
    cost_per_person DECIMAL(10,2) DEFAULT 0,
    
    -- Metadata
    max_participants INTEGER,
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);
CREATE INDEX idx_travel_events_club_id ON travel_events(club_id);
CREATE INDEX idx_travel_events_team_id ON travel_events(team_id);
CREATE INDEX idx_travel_events_departure_date ON travel_events(departure_date);
CREATE INDEX idx_travel_events_type ON travel_events(type);

COMMENT ON TABLE travel_events IS 'Eventos de viaje y partidos del equipo';
COMMENT ON COLUMN travel_events.cost_per_person IS 'Costo calculado automáticamente dividiendo costo total entre confirmados';

-- Event RSVPs (From 003_travel_events.sql)
CREATE TABLE IF NOT EXISTS event_rsvps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    
    -- Metadata
    responded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (event_id) REFERENCES travel_events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Un usuario solo puede tener una respuesta por evento
    UNIQUE(event_id, user_id)
);
CREATE INDEX idx_event_rsvps_event_id ON event_rsvps(event_id);
CREATE INDEX idx_event_rsvps_user_id ON event_rsvps(user_id);
CREATE INDEX idx_event_rsvps_status ON event_rsvps(status);
COMMENT ON TABLE event_rsvps IS 'Confirmaciones de asistencia a eventos';
COMMENT ON COLUMN event_rsvps.status IS 'Estados: PENDING, CONFIRMED, DECLINED';

-- Volunteer Assignments (From 004_volunteer_assignments.sql)
CREATE TABLE IF NOT EXISTS volunteer_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(100) NOT NULL,
    match_id UUID NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    role VARCHAR(100) NOT NULL, -- 'BUFFET', 'SECURITY', 'TRANSPORT', etc.
    notes TEXT,
    
    -- Metadata
    assigned_by VARCHAR(100),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Un usuario no puede tener múltiples roles en el mismo partido
    UNIQUE(match_id, user_id)
);
CREATE INDEX idx_volunteer_assignments_club_id ON volunteer_assignments(club_id);
CREATE INDEX idx_volunteer_assignments_match_id ON volunteer_assignments(match_id);
CREATE INDEX idx_volunteer_assignments_user_id ON volunteer_assignments(user_id);
CREATE INDEX idx_volunteer_assignments_role ON volunteer_assignments(role);

COMMENT ON TABLE volunteer_assignments IS 'Asignación de padres como voluntarios en partidos (buffet, seguridad, etc.)';
COMMENT ON COLUMN volunteer_assignments.role IS 'Rol del voluntario: BUFFET, SECURITY, TRANSPORT, FIRST_AID, etc.';

-- News (From 006_club_news.sql)
CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_news_club_id ON news(club_id);
CREATE INDEX idx_news_published ON news(published);

-- ============================================
-- GDPR Compliance Tables (from 003_gdpr_compliance.sql)
-- ============================================

-- Consent Records Table (GDPR Article 7 - Audit trail for all user consents)
CREATE TABLE IF NOT EXISTS consent_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    consent_type VARCHAR(50) NOT NULL, -- 'TERMS', 'PRIVACY', 'MARKETING', 'HEALTH_DATA'
    version VARCHAR(20) NOT NULL,      -- e.g., '2026-01'
    accepted BOOLEAN NOT NULL DEFAULT TRUE,
    accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(45),            -- IPv6 compatible
    user_agent TEXT,
    parent_user_id VARCHAR(255),       -- For minors, the parent who consented
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_consent_records_user_id ON consent_records(user_id);
CREATE INDEX IF NOT EXISTS idx_consent_records_club_id ON consent_records(club_id);
CREATE INDEX IF NOT EXISTS idx_consent_records_type ON consent_records(consent_type);
COMMENT ON TABLE consent_records IS 'GDPR Article 7 - Audit trail for all user consents';

-- Health Data Access Log (GDPR Article 9 compliance)
CREATE TABLE IF NOT EXISTS health_data_access_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id VARCHAR(255) NOT NULL,
    accessed_user_id VARCHAR(255) NOT NULL,    -- User whose health data was accessed
    accessing_user_id VARCHAR(255) NOT NULL,   -- User who accessed the data
    accessing_user_role VARCHAR(50) NOT NULL,
    document_id UUID,
    action VARCHAR(50) NOT NULL,               -- 'VIEW', 'DOWNLOAD', 'VALIDATE', 'DELETE'
    ip_address VARCHAR(45),
    user_agent TEXT,
    accessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_health_access_log_accessed_user ON health_data_access_log(accessed_user_id);
CREATE INDEX IF NOT EXISTS idx_health_access_log_accessing_user ON health_data_access_log(accessing_user_id);
CREATE INDEX IF NOT EXISTS idx_health_access_log_date ON health_data_access_log(accessed_at);
COMMENT ON TABLE health_data_access_log IS 'GDPR Article 9 - Access log for special category data (health)';

-- GDPR Erasure Requests (GDPR Article 17 - Right to erasure request tracking)
CREATE TABLE IF NOT EXISTS gdpr_erasure_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    executed_at TIMESTAMPTZ,
    executed_by VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'COMPLETED', 'FAILED'
    notes TEXT
);
CREATE INDEX IF NOT EXISTS idx_gdpr_erasure_user ON gdpr_erasure_requests(user_id);
COMMENT ON TABLE gdpr_erasure_requests IS 'GDPR Article 17 - Right to erasure request tracking';

-- ============================================
-- Performance Indexes (from 002_performance_indexes.sql)
-- ============================================

-- 1. Booking Availability Indexes
-- Speeds up HasTimeConflict and ListByFacilityAndDate (critical for booking flow)
CREATE INDEX IF NOT EXISTS idx_bookings_availability 
ON bookings (club_id, facility_id, status, start_time, end_time);

-- Speeds up "My Bookings" queries
CREATE INDEX IF NOT EXISTS idx_bookings_user 
ON bookings (user_id, start_time DESC);

-- 2. Championship Module Indexes
-- Fixes slow joins in matches retrieval and standings calculation

-- Indexes for 'tournament_matches' (using name from repository code)
-- Supports GetMatchesByGroup (filtering by group and sorting/joining)
CREATE INDEX IF NOT EXISTS idx_tournament_matches_container 
ON tournament_matches (tournament_id, home_team_id, away_team_id);

CREATE INDEX IF NOT EXISTS idx_tournament_matches_group 
ON tournament_matches (group_id);

-- Indexes for 'teams'
-- Supports joins from matches to teams (finding team names)
CREATE INDEX IF NOT EXISTS idx_teams_championship 
ON teams (championship_id); 

-- Indexes for 'standings'
-- Supports GetStandings (filtering by group and complex sorting)
CREATE INDEX IF NOT EXISTS idx_standings_ranking 
ON standings (group_id, points DESC, goal_difference DESC, goals_for DESC);

-- Indexes for 'team_members'
-- Supports finding players in team (e.g. for user stats updates)
CREATE INDEX IF NOT EXISTS idx_team_members_user 
ON team_members (user_id);