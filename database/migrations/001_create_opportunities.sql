CREATE TABLE IF NOT EXISTS opportunities (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('scholarship', 'internship', 'competition')),
    country VARCHAR(100) NOT NULL,
    deadline TIMESTAMP,
    url TEXT NOT NULL,
    source VARCHAR(100) NOT NULL,
    description TEXT,
    eligibility TEXT,
    funding TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(title, source, url)
);

CREATE INDEX idx_opportunities_type ON opportunities(type);
CREATE INDEX idx_opportunities_country ON opportunities(country);
CREATE INDEX idx_opportunities_deadline ON opportunities(deadline);
CREATE INDEX idx_opportunities_title_desc ON opportunities USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));
