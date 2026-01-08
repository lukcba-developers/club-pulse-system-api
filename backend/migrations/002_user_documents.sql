-- Migration: User Documents Management
-- Description: Adds support for user document management with expiration tracking
-- Date: 2026-01-07

-- Tabla de Documentos de Usuario
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

-- Índices para optimizar consultas
CREATE INDEX idx_user_documents_user ON user_documents(club_id, user_id);
CREATE INDEX idx_user_documents_status ON user_documents(status);
CREATE INDEX idx_user_documents_expiration ON user_documents(expiration_date) WHERE expiration_date IS NOT NULL;
CREATE INDEX idx_user_documents_type ON user_documents(type);

-- Índice compuesto para búsquedas por usuario y tipo
CREATE INDEX idx_user_documents_user_type ON user_documents(club_id, user_id, type);

-- Agregar columna de elegibilidad en users (desnormalizado para performance)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_eligible BOOLEAN DEFAULT FALSE;

-- Comentarios para documentación
COMMENT ON TABLE user_documents IS 'Stores user documents (DNI, medical certificates, insurance, etc.) with expiration tracking';
COMMENT ON COLUMN user_documents.type IS 'Document type: DNI_FRONT, DNI_BACK, EMMAC_MEDICAL, LEAGUE_FORM, INSURANCE';
COMMENT ON COLUMN user_documents.status IS 'Document status: PENDING, VALID, REJECTED, EXPIRED';
COMMENT ON COLUMN user_documents.expiration_date IS 'Date when the document expires (nullable for documents without expiration)';
COMMENT ON COLUMN users.is_eligible IS 'Cached eligibility status based on document validation (updated by background job)';
