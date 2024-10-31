-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION get_users_with_tag_counts(org_id UUID)
RETURNS TABLE (
    user_id UUID,
    tags_created_count INTEGER,
    email VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        u.id AS user_id,
        tag_counts.tags_created_count,
        u.email
    FROM
        users u
    JOIN (
        SELECT
            created_by AS user_id,
            COUNT(id) AS tags_created_count
        FROM
            tags
        GROUP BY
            created_by
    ) AS tag_counts ON u.id = tag_counts.user_id
    WHERE
        u.organization_id = org_id::UUID
    ORDER BY
        tag_counts.tags_created_count DESC;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS get_users_with_tag_counts();
-- +goose StatementEnd
