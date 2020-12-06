CREATE TABLE IF NOT EXISTS covid_stats (
    id TEXT PRIMARY KEY,
    name TEXT, 
    slug TEXT,
    confirmed INT,
    deaths INT,
    recovered INT,
    collected_at TIMESTAMP WITH TIME ZONE
)
