package data

const (
	QueryUser = `SELECT id, email, password, role FROM users WHERE email = $1`

	QueryCreateUser = `INSERT INTO users (password, email, role) VALUES ($1, $2, $3)`
)
