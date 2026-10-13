 CREATE TABLE IF NOT EXISTS people (
    id BIGSERIAL PRIMARY KEY,

    type VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,

    tin VARCHAR(50) NULL,
    vrn VARCHAR(50) NULL,
    nin VARCHAR(50) NULL,

    email VARCHAR(150) NULL,
    location VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(50) NOT NULL,

    status BOOLEAN DEFAULT true,

    logo_url TEXT NULL,
    document_header_url TEXT NULL,

    po_box VARCHAR(50) NULL,
    tags TEXT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by BIGINT NOT NULL,

    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL
);