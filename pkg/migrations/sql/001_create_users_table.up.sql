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
('Super Admin', 'admin@gmail.com', '$2a$10$5UFBURxHbhczo9MB9KkNNeX2OPUcAzFkvQkRAc.rt/YMF8TCU4kAW', 'super_admin'), -- password: secret
('Manager', 'manager@gmail.com', '$2a$10$EeYi/gXMODpLkkLJ5CHUzu7lFU3u6RJJuhtofaacgp/bcyrbryQi2', 'manager'),  -- password: manager
('Technician 1', 'tech1@gmail.com', 'dcb694aa0322f143ed970e275c807bf123bd5db4f73140b94ccc757f42dc8043', 'technician'), -- password: tech
('Technician 2', 'tech2@gmail.com', 'dcb694aa0322f143ed970e275c807bf123bd5db4f73140b94ccc757f42dc8043', 'technician'); -- password: tech
