--! Messages système pour classes - EducNet Realtime Chat
--! Date: 2026-02-05

BEGIN;

--! Table principale des messages
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    file_url TEXT, -- Pour futures pièces jointes
    is_pinned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

--! Index performances
CREATE INDEX idx_messages_class_id ON messages(class_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_class_created ON messages(class_id, created_at DESC);

--! Trigger updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_messages_updated_at 
    BEFORE UPDATE ON messages 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

--! Vue pour messages avec infos user (WebSocket + API)
CREATE OR REPLACE VIEW messages_view AS
SELECT 
    m.id,
    m.content,
    m.message_type,
    m.file_url,
    m.is_pinned,
    m.created_at,
    m.updated_at,
    u.id AS user_id,
    u.first_name,
    u.last_name,
    (u.first_name || ' ' || u.last_name) AS full_name,
    u.role,
    u.avatar_url,
    c.id AS class_id,
    c.name AS class_name
FROM messages m
JOIN users u ON m.user_id = u.id
JOIN classes c ON m.class_id = c.id;

--! Fonction utilitaire: derniers messages classe
CREATE OR REPLACE FUNCTION get_recent_messages(p_class_id INTEGER, p_limit INTEGER DEFAULT 50)
RETURNS TABLE (
    id BIGINT,
    content TEXT,
    user_id INTEGER,
    full_name TEXT,
    role TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        m.id, 
        m.content, 
        m.user_id, 
        (u.first_name || ' ' || u.last_name) AS full_name,
        u.role, 
        u.avatar_url, 
        m.created_at
    FROM messages m
    JOIN users u ON m.user_id = u.id
    WHERE m.class_id = p_class_id
    ORDER BY m.created_at DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

COMMIT;
