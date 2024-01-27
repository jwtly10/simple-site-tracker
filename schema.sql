CREATE SCHEMA IF NOT EXISTS tracker_db;

CREATE TABLE IF NOT EXISTS domains_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    siteKey VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pages_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain_id INT,
    page_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (domain_id) REFERENCES domains_tb(id)
);

CREATE TABLE IF NOT EXISTS ip_addresses_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ip_address VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS utm_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    page_id INT,
    ip_address_id INT,
    utm_source VARCHAR(255),
    utm_medium VARCHAR(255),
    utm_campaign VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (page_id) REFERENCES pages_tb(id),
    FOREIGN KEY (ip_address_id) REFERENCES ip_addresses_tb(id)
);

CREATE TABLE IF NOT EXISTS clicks_tb (
    id INT AUTO_INCREMENT PRIMARY KEY,
    element VARCHAR(255),
    page_id INT,
    ip_address_id INT,
    clicked_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (page_id) REFERENCES pages_tb(id),
    FOREIGN KEY (ip_address_id) REFERENCES ip_addresses_tb(id)
);
