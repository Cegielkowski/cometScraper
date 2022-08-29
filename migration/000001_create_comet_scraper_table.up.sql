CREATE TABLE IF NOT EXISTS comet_scraper (
    uuid VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    applicant JSONB,
    time_taken INTERVAL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP                             
);
