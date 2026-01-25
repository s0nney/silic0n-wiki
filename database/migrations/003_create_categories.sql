CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

ALTER TABLE articles ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);
CREATE INDEX IF NOT EXISTS idx_articles_category ON articles(category_id);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);

INSERT INTO categories (slug, name, description) VALUES
('technology', 'Technology', 'Articles about technology, software, and computing'),
('science', 'Science', 'Articles about scientific topics and discoveries'),
('general', 'General', 'General knowledge articles')
ON CONFLICT (slug) DO NOTHING;

UPDATE articles SET category_id = (SELECT id FROM categories WHERE slug = 'technology')
WHERE slug IN ('programming', 'golang', 'postgresql') AND category_id IS NULL;

UPDATE articles SET category_id = (SELECT id FROM categories WHERE slug = 'science')
WHERE slug = 'birds' AND category_id IS NULL;

UPDATE articles SET category_id = (SELECT id FROM categories WHERE slug = 'general')
WHERE slug = 'wiki' AND category_id IS NULL;
