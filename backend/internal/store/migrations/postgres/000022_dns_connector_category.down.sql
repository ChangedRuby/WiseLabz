-- 000022_dns_connector_category.down.sql

ALTER TABLE connectors DROP CONSTRAINT connectors_category_check;
ALTER TABLE connectors ADD CONSTRAINT connectors_category_check
    CHECK (category IN ('virtualization','containers_paas','networking'));
