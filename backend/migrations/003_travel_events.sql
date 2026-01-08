-- Migration: 003_travel_events
-- Description: Crea tablas para gestión de eventos de viaje y confirmaciones (RSVPs)

-- Tabla de eventos de viaje
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

-- Índices para travel_events
CREATE INDEX idx_travel_events_club_id ON travel_events(club_id);
CREATE INDEX idx_travel_events_team_id ON travel_events(team_id);
CREATE INDEX idx_travel_events_departure_date ON travel_events(departure_date);
CREATE INDEX idx_travel_events_type ON travel_events(type);

-- Tabla de confirmaciones de asistencia (RSVPs)
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

-- Índices para event_rsvps
CREATE INDEX idx_event_rsvps_event_id ON event_rsvps(event_id);
CREATE INDEX idx_event_rsvps_user_id ON event_rsvps(user_id);
CREATE INDEX idx_event_rsvps_status ON event_rsvps(status);

-- Comentarios
COMMENT ON TABLE travel_events IS 'Eventos de viaje y partidos del equipo';
COMMENT ON TABLE event_rsvps IS 'Confirmaciones de asistencia a eventos';
COMMENT ON COLUMN travel_events.cost_per_person IS 'Costo calculado automáticamente dividiendo costo total entre confirmados';
COMMENT ON COLUMN event_rsvps.status IS 'Estados: PENDING, CONFIRMED, DECLINED';
