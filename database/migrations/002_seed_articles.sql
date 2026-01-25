INSERT INTO articles (slug, title, content) VALUES
('birds', 'Birds', 'Birds are a group of warm-blooded vertebrates constituting the class Aves, characterised by feathers, toothless beaked jaws, the laying of hard-shelled eggs, a high metabolic rate, a four-chambered heart, and a strong yet lightweight skeleton.

Birds live worldwide and range in size from the 5.5 cm (2.2 in) bee hummingbird to the 2.8 m (9 ft 2 in) common ostrich. There are about ten thousand living species, more than half of which are passerine, or perching birds.'),

('programming', 'Programming', 'Computer programming is the process of designing and building an executable computer program to accomplish a specific computing result or to perform a specific task.

Programming involves tasks such as analysis, generating algorithms, profiling algorithms accuracy and resource consumption, and the implementation of algorithms in a chosen programming language.'),

('golang', 'Go Programming Language', 'Go is a statically typed, compiled high-level programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but also has memory safety, garbage collection, structural typing, and CSP-style concurrency.

Go was designed at Google in 2007 to improve programming productivity in an era of multicore, networked machines and large codebases.'),

('postgresql', 'PostgreSQL', 'PostgreSQL, also known as Postgres, is a free and open-source relational database management system emphasizing extensibility and SQL compliance.

PostgreSQL features transactions with ACID properties, automatically updatable views, materialized views, triggers, foreign keys, and stored procedures. It is designed to handle a range of workloads, from single machines to data warehouses or web services with many concurrent users.'),

('wiki', 'Wiki', 'A wiki is a hypertext publication collaboratively edited and managed by its own audience, using a web browser. A typical wiki contains multiple pages for the subjects or scope of the project, and could be either open to the public or limited to use within an organization for maintaining its internal knowledge base.

The term wiki also refers to the collaborative software used to create such a publication.')

ON CONFLICT (slug) DO NOTHING;
