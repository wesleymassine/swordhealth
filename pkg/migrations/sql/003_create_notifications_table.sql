CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    task_id BIGINT NOT NULL,
    notification_body TEXT NOT NULL,
    notification_status ENUM('sent', 'pending', 'failed') NOT NULL DEFAULT 'pending',
    sent_at DATETIME DEFAULT NULL,
    UNIQUE (task_id),
    CONSTRAINT fk_task_id FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);