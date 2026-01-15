-- 1. Alter Match Scores to Decimal
ALTER TABLE tournament_matches 
ALTER COLUMN home_score TYPE DECIMAL(10,2),
ALTER COLUMN away_score TYPE DECIMAL(10,2);

-- 2. Add Settings to Tournament (JSONB)
ALTER TABLE championships 
ADD COLUMN settings JSONB DEFAULT '{}' NOT NULL;

-- 3. Create Roster Snapshot Table
CREATE TABLE tournament_team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tournament_id UUID NOT NULL REFERENCES championships(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    member_id UUID NOT NULL, -- Logical reference to user/member
    
    player_name VARCHAR(255), -- Snapshot of name at time of storage
    player_number INT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tournament_id, team_id, member_id)
);

CREATE INDEX idx_tournament_team_members_tournament_team ON tournament_team_members(tournament_id, team_id);
