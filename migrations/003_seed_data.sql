
--!Insert subjects
INSERT INTO subjects (school_id, name, code) VALUES
(5, 'Mathématiques', 'MATH'),
(5, 'Physique', 'PHY'),
(5, 'Chimie', 'CHI'),
(5, 'Français', 'FR'),
(5, 'Anglais', 'ENG'),
(5, 'Histoire-Géo', 'HISTGEO');

--! Insert classes
INSERT INTO classes (school_id, name, level, section, academic_year) VALUES
(5, 'Terminale C', 'Terminale', 'C', '2025-2026'),
(5, 'Terminale A', 'Terminale', 'A', '2025-2026'),
(5, 'Terminale L', 'Terminale', 'L', '2025-2026'),
(5, 'Terminale OSE', 'Terminale', 'OSE', '2025-2026'),
(5, 'Terminale D', 'Terminale', 'D', '2025-2026');
