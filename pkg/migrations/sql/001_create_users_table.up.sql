CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('super_admin', 'manager', 'technician') NOT NULL
);

-- Create index for the role column to optimize role-based queries
CREATE INDEX idx_users_role ON users(role);

-- Insert sample users into the users table
INSERT INTO users (name, email, password_hash, role) VALUES 
('Super Admin', 'admin@gmail.com', '$2a$10$2424gsLTbMzht15Avjpnd.AN1eBelOF11ATmHYTbYfLCTCuAZIvk6', 'super_admin'), -- password: secret
('Manager', 'manager@gmail.com', '$2a$10$UAaBpSxOq7OuEAyHYL3xOu.styvAJoJ/DDcbu2y7/7NC1a27/iL2u', 'manager'),  -- password: manager
('Technician 1', 'tech1@gmail.com', '$2a$10$u.yzb9Z37xnJkud4UJMvXeECBpOy5Vd.3SRFHVe4Np8KwHnhXOABm', 'technician'), -- password: tech
('Technician 2', 'tech2@gmail.com', '$2a$10$u.yzb9Z37xnJkud4UJMvXeECBpOy5Vd.3SRFHVe4Np8KwHnhXOABm', 'technician'); -- password: tech
