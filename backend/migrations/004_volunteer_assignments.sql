-- Migration: 004_volunteer_assignments
-- Description: Agrega tabla para asignación de voluntarios (padres) a fixtures

-- Tabla de asignaciones de voluntarios
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

-- Índices
CREATE INDEX idx_volunteer_assignments_club_id ON volunteer_assignments(club_id);
CREATE INDEX idx_volunteer_assignments_match_id ON volunteer_assignments(match_id);
CREATE INDEX idx_volunteer_assignments_user_id ON volunteer_assignments(user_id);
CREATE INDEX idx_volunteer_assignments_role ON volunteer_assignments(role);

-- Comentarios
COMMENT ON TABLE volunteer_assignments IS 'Asignación de padres como voluntarios en partidos (buffet, seguridad, etc.)';
COMMENT ON COLUMN volunteer_assignments.role IS 'Rol del voluntario: BUFFET, SECURITY, TRANSPORT, FIRST_AID, etc.';
