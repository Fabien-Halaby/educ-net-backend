--! =============================================
--! SUBJECTS (Matières)
--! =============================================
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL, --! Ex: MATH, PHY, FR
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_subjects_school ON subjects(school_id);

--! =============================================
--! CLASSES
--! =============================================
CREATE TABLE IF NOT EXISTS classes (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, --! Ex: 6ème A, Terminale S2
    level VARCHAR(50) NOT NULL, --! Ex: 6ème, Terminale
    section VARCHAR(50), --! Ex: A, S, L
    capacity INTEGER DEFAULT 40,
    academic_year VARCHAR(20) NOT NULL, -- Ex: 2025-2026
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, name, academic_year)
);

CREATE INDEX idx_classes_school ON classes(school_id);
CREATE INDEX idx_classes_academic_year ON classes(academic_year);

--! =============================================
--! TEACHER_SUBJECTS (Many-to-Many)
--! =============================================
CREATE TABLE IF NOT EXISTS teacher_subjects (
    id SERIAL PRIMARY KEY,
    teacher_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(teacher_id, subject_id)
);

CREATE INDEX idx_teacher_subjects_teacher ON teacher_subjects(teacher_id);
CREATE INDEX idx_teacher_subjects_subject ON teacher_subjects(subject_id);

--! =============================================
--! STUDENT_CLASSES (One student = One class)
--! =============================================
CREATE TABLE IF NOT EXISTS student_classes (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, class_id)
);

CREATE INDEX idx_student_classes_student ON student_classes(student_id);
CREATE INDEX idx_student_classes_class ON student_classes(class_id);

--! =============================================
--! SEED DATA (Matières par défaut)
--! =============================================
COMMENT ON TABLE subjects IS 'Matières enseignées dans chaque école';
COMMENT ON TABLE classes IS 'Classes disponibles dans chaque école';
COMMENT ON TABLE teacher_subjects IS 'Association professeurs ↔ matières';
COMMENT ON TABLE student_classes IS 'Inscription élèves → classe';
