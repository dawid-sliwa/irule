package data

const (
	QueryUser = `SELECT id, email, password, role, organization_id FROM users WHERE email = $1`

	QueryCreateUser = `INSERT INTO users (password, email, role) VALUES ($1, $2, $3)`

	QueryCreateUserOrg = `INSERT INTO users (password, email, role, organization_id) VALUES ($1, $2, $3, $4)`

	QueryDocumentationsByOrg = `
	SELECT 
		d.id,
		d.name,
		d.content,
		(
			SELECT COUNT(*)
			FROM tags t
			WHERE t.documentation_id = d.id
		) AS tag_count
	FROM documentations d
	WHERE d.organization_id = $1`

	QueryDocumentationByID = `
	SELECT 
		d.id,
		d.name,
		d.content,
		(
			SELECT COUNT(*)
			FROM tags t
			WHERE t.documentation_id = d.id
		) AS tag_count
	FROM documentations d
	WHERE d.id = $1 AND d.organization_id = $2`

	InsertDocumentation = `INSERT INTO documentations (name, content, organization_id) VALUES ($1, $2, $3) RETURNING id, name, content`

	UpdateDocumentation = `UPDATE documentations SET name = $1, content = $2 WHERE id = $3 RETURNING id, name, content`

	DeleteDocumentation = `DELETE FROM documentations WHERE id = $1`

	QueryTagsByDoc = `SELECT id, name FROM tags WHERE documentation_id = $1`

	QueryTagByID = `SELECT id, name FROM tags WHERE id = $1 AND documentation_id = $2`

	InsertTag = `INSERT INTO tags (name, documentation_id, created_by) VALUES ($1, $2, $3
	) RETURNING id, name`

	UpdateTag = `UPDATE tags SET name = $1 WHERE id = $2 RETURNING id, name`

	DeleteTag = `DELETE FROM tags WHERE id = $1`

	GetUserSttats = `
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
	u.organization_id = $1
ORDER BY
	tag_counts.tags_created_count DESC;
`
)
