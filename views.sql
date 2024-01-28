-- Collection of generic views for creating dashboards

CREATE VIEW click_tracking_view AS
SELECT
    d.domain,
    p.page_url,
    JSON_EXTRACT(c.element, '$.tag') as tag,
    JSON_EXTRACT(c.element, '$.href') as href,
    JSON_EXTRACT(c.element, '$.textContent') as content,
    COUNT(*) as count
FROM
    clicks_tb c
        JOIN
    pages_tb p ON c.page_id = p.id
        JOIN
    domains_tb d ON p.domain_id = d.id
GROUP BY
    d.domain, p.page_url, tag, href, content;


CREATE VIEW page_views_view AS
SELECT
    d.domain,
    p.page_url,
    COUNT(*) as count
FROM
    page_views_tb pv
        JOIN
    pages_tb p ON pv.page_id = p.id
        JOIN
    domains_tb d ON p.domain_id = d.id
GROUP BY
    d.domain, p.page_url;

CREATE VIEW utm_tracking_view AS
SELECT
    d.domain,
    p.page_url,
    u.track,
    u.utm_campaign,
    u.utm_medium,
    u.utm_source,
    COUNT(*) as count
FROM
    utm_tb u
        JOIN
    pages_tb p ON p.id = u.page_id
        JOIN
    domains_tb d ON p.domain_id = d.id
GROUP BY
    d.domain, p.page_url, u.track, u.utm_campaign, u.utm_medium, u.utm_source;


