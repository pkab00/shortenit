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