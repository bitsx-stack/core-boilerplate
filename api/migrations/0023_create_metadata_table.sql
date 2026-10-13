CREATE TABLE IF NOT EXISTS metadata (
    id BIGSERIAL PRIMARY KEY,
    
    entity_type VARCHAR(50) NOT NULL,   
    entity_id BIGINT NOT NULL,          
    
    meta_key VARCHAR(100) NOT NULL,     
    meta_value TEXT,                    
    
    value_type VARCHAR(20),             
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    deleted_at TIMESTAMP NULL,
);