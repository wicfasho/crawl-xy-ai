-- UP migration

-- Page Table
CREATE TABLE Page (
    PageID UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    URL VARCHAR(255) NOT NULL,
    Title VARCHAR(255),
    Meta_Description TEXT,
    Meta_Keywords TEXT,
    Last_Modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
);

-- Index in Page Table
CREATE UNIQUE INDEX idx_page_url ON Page (URL);
CREATE UNIQUE INDEX idx_page_id ON Page (PageID);

-- Content Table
CREATE TABLE Content (
    ContentID UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    PageID UUID,
    Content_Type VARCHAR(20) NOT NULL,
    Content_Value TEXT,
    FOREIGN KEY (PageID) REFERENCES Page(PageID) ON DELETE CASCADE
);

-- Index in Content Table
CREATE UNIQUE INDEX idx_content_page_id ON Content (PageID);

-- Element Table
CREATE TABLE Element (
    ElementID UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    PageID UUID,
    Element_Type VARCHAR(20) NOT NULL,
    Position_X INT,
    Position_Y INT,
    Content_Value TEXT,
    FOREIGN KEY (PageID) REFERENCES Page(PageID) ON DELETE CASCADE
);

-- Index in Element Table
CREATE UNIQUE INDEX idx_element_page_id ON Element (PageID);