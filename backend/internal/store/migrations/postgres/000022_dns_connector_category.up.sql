-- 000022_dns_connector_category.up.sql — allow the "dns" connector category
-- (DNS Resolver, Pi-hole) for cross-service linking (#235)

ALTER TABLE connectors DROP CONSTRAINT connectors_category_check;
ALTER TABLE connectors ADD CONSTRAINT connectors_category_check
    CHECK (category IN ('virtualization','containers_paas','networking','dns'));
