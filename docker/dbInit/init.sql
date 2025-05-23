CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    password_hash TEXT NOT NULL,
    stripe_customer_id TEXT UNIQUE,
    stripe_subscription_id TEXT UNIQUE,
    subscription_status VARCHAR(50) DEFAULT 'free',  -- free, basic, pro, agency
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE competitors (
    id SERIAL PRIMARY KEY,
    competitor_name VARCHAR(50),
    url TEXT UNIQUE NOT NULL,
    last_scraped_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tracked_competitors (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    competitor_id INT REFERENCES competitors(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS price_changes (
    id SERIAL PRIMARY KEY,
    competitor_id INT REFERENCES competitors(id) ON DELETE CASCADE,
    detected_change TEXT NOT NULL,  -- Example: "Pro Plan: $49 → $59"
    ai_summary TEXT,  -- OpenAI-generated summary (NULL initially)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    competitor_url TEXT NOT NULL,
    plan_name TEXT NOT NULL,
    price TEXT NOT NULL,
    billing_cycle TEXT NOT NULL,
    extracted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
