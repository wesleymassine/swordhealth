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
('Super Admin', 'admin@gmail.com', '$2a$10$dwfMRoHw6T42MLGKtyKBW.SFTQ5PSy07wWi6mx2gIwS8gKNNhY/Oi', 'super_admin'), -- password: secret
('Manager', 'manager@gmail.com', '7676aaafb027c825bd9abab78b234070e702752f625b752e55e55b48e607e358', 'manager'),  -- password: manager_pass
('Technician 1', 'tech1@gmail.com', 'dcb694aa0322f143ed970e275c807bf123bd5db4f73140b94ccc757f42dc8043', 'technician'), -- password: tech1_pass
('Technician 2', 'tech2@gmail.com', 'dcb694aa0322f143ed970e275c807bf123bd5db4f73140b94ccc757f42dc8043', 'technician'); -- password: tech2_pass
