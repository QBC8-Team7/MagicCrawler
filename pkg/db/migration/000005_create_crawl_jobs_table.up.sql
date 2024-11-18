-- Table for crawl jobs
CREATE TABLE if NOT EXISTS crawl_jobs (
    id BIGSERIAL PRIMARY KEY,
    url VARCHAR(500) NOT NULL,
    source_name VARCHAR(100) NOT NULL,
    page_type VARCHAR(10) CHECK (page_type IN ('archive', 'single')) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('waiting', 'picked', 'done', 'failed')) DEFAULT 'waiting' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);