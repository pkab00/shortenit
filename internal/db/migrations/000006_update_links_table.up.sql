BEGIN;

ALTER TABLE links RENAME TO links_old;
ALTER SEQUENCE links_link_id_seq RENAME TO links_link_id_seq_old;

CREATE TABLE IF NOT EXISTS links (
	link_id SERIAL PRIMARY KEY,
    link_code VARCHAR(512) UNIQUE NOT NULL,
	url VARCHAR(1024) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now()
);

SELECT setval(pg_get_serial_sequence('links', 'link_id'), coalesce(max(link_id), 1)) FROM links;

INSERT INTO links (link_id, link_code, url, created_at)
SELECT link_id, encode(link_id::text::bytea, 'base64'), url, created_at FROM links_old;

DROP TABLE links_old CASCADE;

/* CAUTION: ALL REDIRECT COUNTERS WILL BE RESET */

COMMIT;