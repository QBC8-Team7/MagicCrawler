-- Insert a new crawl job
-- name: CreateCrawlJob :one
INSERT INTO crawl_jobs (url, source_name, page_type, status)
VALUES (sqlc.arg('url'), sqlc.arg('source_name'), sqlc.arg('page_type'), sqlc.arg('status'))
RETURNING *;


-- name: CheckCrawlJobExists :one
-- url: text
-- statuses: text[]
SELECT EXISTS (
    SELECT 1 FROM crawl_jobs 
    WHERE url = @url 
      AND status = ANY(@statuses::text[])
) AS exists;


-- name: GetFirstMatchingCrawlJob :one
-- url: text
-- statuses: text[]
SELECT * FROM crawl_jobs
WHERE url = sqlc.arg('url')
  AND status = ANY(sqlc.arg('statuses'))
LIMIT 1;

-- name: UpdateCrawlJobStatus :one
UPDATE crawl_jobs 
SET status = sqlc.arg('status') 
WHERE id = sqlc.arg('jobID')
RETURNING id;

-- name: GetFirstCrawlJobByStatus :one
SELECT * FROM crawl_jobs
WHERE status = sqlc.arg('status') 
LIMIT 1;