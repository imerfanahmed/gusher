CREATE TABLE webhooks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    app_id VARCHAR(255) NOT NULL,
    event VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    api_token VARCHAR(255),
    UNIQUE KEY unique_webhook (app_id, event, url),
    FOREIGN KEY (app_id) REFERENCES apps(id)
);