CREATE TABLE IF NOT EXISTS redirects (
    id SERIAL PRIMARY KEY,
    link_id INT UNIQUE REFERENCES links(link_id),
    redirect_counter BIGINT NOT NULL
);

CREATE OR REPLACE FUNCTION increment_redirect_counter(p_link_id INT)
RETURNS redirects
LANGUAGE plpgsql
AS $$
DECLARE
    result redirects;
BEGIN
    UPDATE redirects
    SET redirect_counter = redirect_counter + 1
    WHERE link_id = p_link_id
    RETURNING * INTO result;

    RETURN result;
END;
$$;