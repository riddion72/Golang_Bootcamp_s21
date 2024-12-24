package repository

import (
	"database/sql"
	"fmt"
	"log"
	"main/internal/models"

	_ "github.com/lib/pq"
)

const (
	col1 string = "article_name"
	col2 string = "article_content"
)

type pgRepo struct {
	db *sql.DB
}

type Repository interface {
	GetArticle(models.GetArticleRequest) (models.GetArticleResponse, error)
	SetArticle(models.SetArticleRequest) error
	GetNumber() (int, error)
}

func NewPGRepository(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) GetArticle(req models.GetArticleRequest) (models.GetArticleResponse, error) {
	var article models.GetArticleResponse
	id := ""

	row, err := r.db.Query(queryGetArticle, req.Limit, req.Ofset)
	if err != nil {
		log.Println("Error getting")
		return article, err
	}

	defer row.Close()

	for i := 0; i < req.Limit && row.Next(); i++ {

		var articleItem models.Article = models.Article{Article_name: "", Article_content: ""}

		err := row.Scan(&id, &articleItem.Article_name, &articleItem.Article_content)
		if err != nil {
			log.Println("Error scanning")
			return article, err
		}
		article.Response = append(article.Response, articleItem)
	}

	return article, nil
}
func (r *pgRepo) GetNumber() (int, error) {
	var number int

	row, err := r.db.Query(queryGetNumber)

	defer row.Close()
	if err != nil {
		log.Println("Error getting number")
		return number, err
	}
	for row.Next() {
		row.Scan(&number)
	}

	return number, nil
}
func (r *pgRepo) SetArticle(req models.SetArticleRequest) error {
	query := fmt.Sprintf(querySetArticle, col1, col2)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Request.Article_name, req.Request.Article_content)
	if err != nil {
		return err
	}
	return nil
}
