CREATE OR REPLACE FUNCTION increment_redirect_counter(p_link_id INT)
RETURNS redirects
LANGUAGE plpgsql
AS $$
DECLARE
    result redirects;
BEGIN
    INSERT INTO redirects (link_id, redirect_counter)
    VALUES (p_link_id, 1)
    ON CONFLICT (link_id)
    DO UPDATE SET
        redirect_counter = redirects.redirect_counter + 1
    RETURNING * INTO result;

    RETURN result;
END;
$$;