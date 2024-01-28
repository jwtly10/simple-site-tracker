-- Collection of generic views for creating dashboards

-- Create a view for clicks per domain/page
CREATE VIEW click_tracking_view AS
SELECT
    d.domain,
    p.page_url,
    JSON_EXTRACT(c.element, '$.tag') as tag,
    JSON_EXTRACT(c.element, '$.href') as href,
    JSON_EXTRACT(c.element, '$.textContent') as content,
    c.created_at as timestamp
FROM
    clicks_tb c
        JOIN
    pages_tb p ON c.page_id = p.id
        JOIN
    domains_tb d ON p.domain_id = d.id;

-- Create a view for page views per domain/page
CREATE VIEW page_views_view AS
SELECT
    d.domain,
    p.page_url,
    pv.created_at as timestamp
FROM
    page_views_tb pv
        JOIN
    pages_tb p ON pv.page_id = p.id
        JOIN
    domains_tb d ON p.domain_id = d.id;

-- Create a view for the number of utms and which pages they led to
CREATE VIEW utm_tracking_view AS
SELECT
    d.domain,
    p.page_url,
    u.track,
    u.utm_campaign,
    u.utm_medium,
    u.utm_source,
    u.created_at as timestamp
FROM
    utm_tb u
        JOIN
    pages_tb p ON p.id = u.page_id
        JOIN
    domains_tb d ON p.domain_id = d.id;

