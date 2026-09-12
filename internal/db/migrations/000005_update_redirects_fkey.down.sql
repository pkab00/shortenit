ALTER TABLE redirects
DROP CONSTRAINT redirects_link_id_fkey;

ALTER TABLE redirects
ADD CONSTRAINT redirects_link_id_fkey 
FOREIGN KEY (link_id) REFERENCES links(link_id) ON DELETE NO ACTION;