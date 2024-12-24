package repository

const (
	queryGetArticle = `SELECT * FROM articles LIMIT $1 OFFSET $2`
	queryGetNumber  = `SELECT COUNT(*) FROM articles`
	querySetArticle = `INSERT INTO articles (%s, %s) VALUES ($1, $2)
	ON CONFLICT (article_name) DO UPDATE SET article_content = $2`
)

// ON CONFLICT (article_name) DO UPDATE SET article_content = $2
