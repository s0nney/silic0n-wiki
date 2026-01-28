CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(slug, category_id)
);

CREATE INDEX IF NOT EXISTS idx_tags_slug ON tags(slug);
CREATE INDEX IF NOT EXISTS idx_tags_category ON tags(category_id);

CREATE TABLE IF NOT EXISTS article_tags (
    article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_article_tags_article ON article_tags(article_id);
CREATE INDEX IF NOT EXISTS idx_article_tags_tag ON article_tags(tag_id);

-- Seed example tags for Technology category
INSERT INTO tags (slug, name, category_id) VALUES
('linux', 'Linux', (SELECT id FROM categories WHERE slug = 'technology')),
('databases', 'Databases', (SELECT id FROM categories WHERE slug = 'technology')),
('web-development', 'Web Development', (SELECT id FROM categories WHERE slug = 'technology')),
('programming-languages', 'Programming Languages', (SELECT id FROM categories WHERE slug = 'technology'))
ON CONFLICT (slug, category_id) DO NOTHING;

-- Seed example tags for Science category
INSERT INTO tags (slug, name, category_id) VALUES
('mammals', 'Mammals', (SELECT id FROM categories WHERE slug = 'science')),
('birds', 'Birds', (SELECT id FROM categories WHERE slug = 'science')),
('physics', 'Physics', (SELECT id FROM categories WHERE slug = 'science')),
('biology', 'Biology', (SELECT id FROM categories WHERE slug = 'science'))
ON CONFLICT (slug, category_id) DO NOTHING;

-- Seed example tags for General category
INSERT INTO tags (slug, name, category_id) VALUES
('reference', 'Reference', (SELECT id FROM categories WHERE slug = 'general')),
('history', 'History', (SELECT id FROM categories WHERE slug = 'general'))
ON CONFLICT (slug, category_id) DO NOTHING;

-- Associate existing articles with tags
INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'golang' AND t.slug = 'programming-languages'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'postgresql' AND t.slug = 'databases'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'programming' AND t.slug = 'programming-languages'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'birds' AND t.slug = 'birds'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'birds' AND t.slug = 'biology'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id FROM articles a, tags t
WHERE a.slug = 'wiki' AND t.slug = 'reference'
ON CONFLICT DO NOTHING;
