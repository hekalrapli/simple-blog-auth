package database

import (
	"database/sql"
	"log"
	"simple-blog-auth/models"
	_ "modernc.org/sqlite"
)

func Insert(db *sql.DB, username string, password string) error {
	sql := `INSERT INTO users (username, password)
			VALUES (?, ?);`
	_, err := db.Exec(sql, username, password)
	if err != nil {
		return err
	}
	return nil
}

func Select(db *sql.DB, username string) (int64, string, error) {
	var userId int64
	var password string
	err := db.QueryRow(`SELECT id, password from users where username = ?`, username).Scan(&userId,&password)
	if err != nil {
		return 0,"",err
	}
	return userId,password,nil
}
func GetUsername(db *sql.DB, userId int64) (string, error) {
	var username string
	err := db.QueryRow(`SELECT username from users where id = ?`, userId).Scan(&username)
	if err != nil {
		return "",err
	}
	return username,nil
}

func InsertArticle(db *sql.DB, title, content string, authorID int64) error {
	sql := `INSERT INTO articles (title,content,author_id) VALUES (?,?,?);`
	_,err := db.Exec(sql,title,content,authorID)
	if err != nil {
		return err
	}
	return nil
}

func GetArticles(db *sql.DB, userId int64) ([]models.Articles, error) {
	var Articles []models.Articles
	rows, err := db.Query(`SELECT a.id, a.title, a.content, a.author_id, u.username, strftime('%Y-%m-%d %H:%M', a.created_at, '+7 hours'), strftime('%Y-%m-%d %H:%M', a.updated_at, '+7 hours') from ARTICLES as a INNER JOIN users as u ON a.author_id = u.id where a.author_id = ? ORDER BY a.updated_at desc`,userId)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.Articles
		if err := rows.Scan(&a.Id, &a.Title, &a.Content, &a.AuthorID, &a.AuthorName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil,err
		}
		Articles = append(Articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil,err
	}
	return Articles, nil
}

func GetAllArticles(db *sql.DB) ([]models.Articles, error) {
	var Articles []models.Articles
	rows, err := db.Query(`SELECT a.id, a.title, a.content, a.author_id, u.username, strftime('%Y-%m-%d %H:%M', a.created_at, '+7 hours'), strftime('%Y-%m-%d %H:%M', a.updated_at, '+7 hours') from ARTICLES as a INNER JOIN users as u ON a.author_id = u.id  ORDER BY a.updated_at desc`)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.Articles
		if err := rows.Scan(&a.Id, &a.Title, &a.Content, &a.AuthorID, &a.AuthorName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil,err
		}
		Articles = append(Articles, a)
	}
	if err := rows.Err(); err != nil {
		return nil,err
	}
	return Articles, nil
}


func UpdateArticles(db *sql.DB, articleId int64,title string, content string) error {

	sql := `UPDATE articles SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`
	_,err := db.Exec(sql,title,content,articleId)
	if err != nil {
		log.Printf("Gagal mengupdate data: %v",err)
		return err
	}
	return nil
}

func DeleteArticle(db *sql.DB, articleId int64) error {
	sql := `DELETE from articles where id = ?;`
	_,err := db.Exec(sql,articleId)
	if err != nil {
		log.Printf("Gagal menghapus data: %v",err)
		return err
	}
	return nil
}



func GetArticle(db *sql.DB, articleId int64) (*models.Articles, error) {
	var article models.Articles
	err := db.QueryRow(`SELECT a.id,a.title,a.content, a.author_id,u.username, strftime('%Y-%m-%d %H:%M', a.created_at, '+7 hours'), strftime('%Y-%m-%d %H:%M', a.updated_at, '+7 hours') from articles a inner join users u on a.author_id = u.id where a.id = ?`, articleId).Scan(&article.Id,&article.Title, &article.Content, &article.AuthorID,&article.AuthorName, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		log.Printf("Gagal mengambil data: %v",err)
		return nil,err
	}
	
	return &article, nil
}


func InitializeSchema(db *sql.DB) error {
	const schema = `
	 CREATE TABLE IF NOT EXISTS users(
	 	id INTEGER PRIMARY KEY ,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	 );
	 CREATE TABLE IF NOT EXISTS articles(
	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id)
	 ); 

	`
	_, err := db.Exec(schema)
	return err
}

func InitDB() *sql.DB {
	db,err := sql.Open("sqlite", "database/app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v",err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to open database: %v",err)
	}
	
	if err := InitializeSchema(db);  err != nil {
		log.Fatalf("Failed to initialize schema: %v",err)
	}

	log.Println("Succesfully connected to the database")
	return db
}