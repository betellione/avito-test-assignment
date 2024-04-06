CREATE TABLE IF NOT EXISTS banners (
     banner_id SERIAL PRIMARY KEY,
     feature_id INTEGER NOT NULL,
     title VARCHAR(255) NOT NULL,
     text VARCHAR(255) NOT NULL,
     url VARCHAR(255) NOT NULL,
     is_active BOOLEAN NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS banner_tags (
    banner_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (banner_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES banners(banner_id),
    FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
);

CREATE TABLE IF NOT EXISTS features (
    feature_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL
);

CREATE INDEX idx_banner_feature ON banners (feature_id);
CREATE INDEX idx_banner_tag ON banner_tags (banner_id, tag_id);
