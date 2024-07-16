-- name: InsertPage :one
INSERT INTO Page (URL, Title, Meta_Description, Meta_Keywords, Last_Modified)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
ON CONFLICT (URL)
DO UPDATE SET
    Title = EXCLUDED.Title,
    Meta_Description = EXCLUDED.Meta_Description,
    Meta_Keywords = EXCLUDED.Meta_Keywords,
    Last_Modified = CURRENT_TIMESTAMP
RETURNING PageID;
