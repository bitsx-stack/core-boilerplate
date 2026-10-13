CREATE TABLE IF NOT EXISTS streets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    continent_code VARCHAR(10) NOT NULL,
    country_code VARCHAR(30) NOT NULL,
    region_code VARCHAR(30) NOT NULL,
    district_code VARCHAR(30) NOT NULL,
    wards_code VARCHAR(30) NOT NULL,
    village_code VARCHAR(30) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    deleted_by BIGINT
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;