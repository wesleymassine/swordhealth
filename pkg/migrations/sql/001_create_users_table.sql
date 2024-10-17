CREATE DATABASE IF NOT EXISTS task
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;


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

CREATE INDEX idx_users_role ON users(role);

-- Insert sample roles into the users table
INSERT INTO users (name, email, password_hash, role) VALUES 
('Super Admin', 'admin@gmail.com', '$2a$10$xK/OeQdzzlOvpFjIej/vAe3QFKOo6a6EKIhUzfn4mr1iVFMifkbuy', 'super_admin'), -- pass: admin
('Manager', 'manager@gmail.com', '$2a$10$VdZObwc80B/d4O2VBA9ofetnjg.y6zdju95.tHa3Fz3FSKZ6x9EK2', 'manager'), -- pass: manager
('Technician 1', 'tech1@gmail.com', '$2a$10$RmvreKRznDEVsZ9proTZXOiYORy9R1iAeq4N9SKZ.JgJl/gGcQz6G', 'technician'), -- pass: tech1
('Technician 2', 'tech2@gmail.com', '$2a$10$uYf4eUgrt90402CpO9.XYeBqNSWBpsBuIscBlfIeLYH/lPdao2kye', 'technician'); -- pass: tech2