-- DOWN migration

-- Element
DROP INDEX IF EXISTS idx_element_page_id;
DROP TABLE IF EXISTS Element;

-- Content
DROP INDEX IF EXISTS idx_content_page_id;
DROP TABLE IF EXISTS Content;

-- Page
DROP INDEX IF EXISTS idx_page_url;
DROP INDEX IF EXISTS idx_page_id;
DROP TABLE IF EXISTS Page;