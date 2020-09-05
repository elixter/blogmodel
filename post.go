package models

import (
	"time"
	"regexp"	
	"database/sql"
	"log"
	"strings"
	
	"github.com/grokify/html-strip-tags-go"
	"github.com/PuerkitoBio/goquery"
)

const (
	src				=		2
)

type Post struct {
	Id			int			`json: "id" db: "id" gorm: "id"`
	Author		string		`json: "author" db: "author" gorm: "author"`
	UDesc		string		`json: "udesc" db: "udesc" gorm: "udesc"`
	Title		string		`json: "title" db: "title" gorm: "title"`
	Thumbnail	string		`json: "thumbnail" db: "thumbnail" gorm: "thumbnail"`
	Content		string		`json: "content" db: "content" gorm: "content"`
	Summary		string		`json: "summary" db: "summary" gorm: "summary"`
	Date	time.Time	`json: "date" db: "date" gorm: "date"`
	Updated	time.Time	`json: "updated" db: "update" gorm: "updated"`
	Category string `json: category" db: "category" gorm: "category"`
	HashTags []string
}

type Page struct {
	CurrentPage int
	Length int
}

// Helper Functions
func GetCategories(db *sql.DB) ([]string, error) {
	var categories []string

	// Database에서 카테고리 가져오기
	rows, err := db.Query("select category from categories;")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		if err != nil {
			log.Println(err)
		}
		categories = append(categories, tmp)
	}

	return categories, err
}

// Get Total count of posts from database
// if something wrong while querying return -1 and error.
// else return count of posts and nil.
func GetPostCount(db *sql.DB) (int, error) {
	var postCount int
	err := db.QueryRow("select count(*) as postCount from posts;").Scan(&postCount)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	
	return postCount, nil
}

func GetPostCountByHashTag(db *sql.DB, hashTag string) (int, error) {
	var postCount int
	err := db.QueryRow("select count(*) as postCount from posts where hashtag like ?", hashTag).Scan(&postCount)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	
	return postCount, nil
}

func GetPostCountByCategory(db *sql.DB, category string) (int, error) {
	var postCount int
	err := db.QueryRow("select count(*) as postCount from posts where category = ?", category).Scan(&postCount)
	if err != nil {
		log.Println(err)
		
		return -1, err
	}
	
	return postCount, nil
}


// Get posts in current page from database
func GetCurrentPagePosts (db *sql.DB, currentPage int, rowsPerPage int) ([]Post, error) {
	var posts []Post
	var hashs sql.NullString
	var thumbnail sql.NullString
	
	// 현재 페이지에 해당하는 게시글만 쿼리
	Rows, err := db.Query("select * from posts where id > ? and id  <= ? order by id desc;", (currentPage - 1) * rowsPerPage, currentPage * rowsPerPage)
	if err != nil {
		log.Fatal(err)
		
		return nil, err
	}
	defer Rows.Close()
	
	for Rows.Next() {
		p := Post{}
		err := Rows.Scan(&p.Id, &p.Author, &p.UDesc, &p.Title, &thumbnail, &p.Content, &p.Summary, &p.Date, &p.Updated, &p.Category, &hashs)
		if err != nil {
			log.Println(err)
			
			return nil, err
		}

		p.convertHashTag(hashs.String)
		p.Thumbnail = thumbnail.String
		// posts에 새로 받아온 post append
		posts = append(posts, p)
	}
	
	return posts, nil
}

func GetCurrentPagePostsByHashTag(db *sql.DB, currentPage int, rowsPerPage int, hashTag string) ([]Post, error) {
	var posts []Post
	
	Rows, err := db.Query(`select id, author, udesc, title, content, summary, date, updated, category, hashtag 
									from (
										select @num := @num + 1 as num,
											p.id, p.author, p.udesc, p.title, p.content, p.summary, p.date, p.updated, p.category, p.hashtag
											from (select @num:=0) a, posts p where hashtag like ?) post
									where post.num > ? and post.num <= ? order by id desc;`, hashTag, (currentPage - 1) * rowsPerPage, currentPage * rowsPerPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer Rows.Close()
	
	for Rows.Next() {
		p := Post{}
		var hashs sql.NullString		// 하나로 합쳐져있는 해쉬태그 문자열
		err := Rows.Scan(&p.Id, &p.Author, &p.UDesc, &p.Title, &p.Content, &p.Summary, &p.Date, &p.Updated, &p.Category, &hashs)
		if err != nil {
			log.Println(err)
		}
		p.convertHashTag(hashs.String)

		// posts에 새로 받아온 post append
		posts = append(posts, p)
	}
	
	return posts, nil
}

func GetCurrentPagePostsByCategory(db *sql.DB, currentPage int, rowsPerPage int, category string) ([]Post, error){
	var posts []Post
	
	Rows, err := db.Query(`select id, author, udesc, title, content, summary, date, updated, category, hashtag 
									from (
										select @num := @num + 1 as num,
										p.id, p.author, p.udesc, p.title, p.content, p.summary, p.date, p.updated, p.category, p.hashtag
										from (select @num:=0) a, posts p where category = ?) post
									where post.num > ? and post.num <= ? order by id desc;`, category, (currentPage - 1) * rowsPerPage, currentPage * rowsPerPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer Rows.Close()
	
	for Rows.Next() {
		p := Post{}
		var hashs sql.NullString		// 하나로 합쳐져있는 해쉬태그 문자열
		err := Rows.Scan(&p.Id, &p.Author, &p.UDesc, &p.Title, &p.Content, &p.Summary, &p.Date, &p.Updated, &p.Category, &hashs)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		p.convertHashTag(hashs.String)

		// posts에 새로 받아온 post append
		posts = append(posts, p)
	}
	
	return posts, nil
}

// Member Functions
func (p *Post) convertHashTag(hashTags string) []string {
	// "#해쉬태그 #테스트" 와 같은 문자열을
	// [해쉬태그, 테스트] 로 바꿔주는 함수
	regExp := regexp.MustCompile("#")

	regedHash := regExp.ReplaceAllLiteralString(hashTags, "")		// 정규화된 해쉬태그들을 공백으로 토큰화
	p.HashTags = strings.Split(regedHash, " ")
	
	return p.HashTags
}

func (p *Post) CreateSummary(summaryLength int) string {
	// Content에서 특정길이 문자열에서 html태그 제거한것.
	var sumText string
	
	if len(p.Content) >= summaryLength {
		if (strings.Contains(p.Content, "&nbsp;")) {
			sumText = strings.Split(p.Content, "&nbsp;")[0]
		} else {
			sumText = p.Content[:summaryLength]
		}
	}  else {
		if (strings.Contains(p.Content, "&nbsp;")) {
			sumText = strings.Split(p.Content, "&nbsp;")[0]
		} else {
			sumText = p.Content
		}
	}
	p.Summary = strip.StripTags(sumText)
	
	return p.Summary
}

func (p *Post) CreateThumbnail() (string, error) {
	var thumbnail string
	contentDoc, err := goquery.NewDocumentFromReader(strings.NewReader(p.Content))
	if err != nil {
		log.Println(err)
		
		return "nil", err
	}
	firstImage := contentDoc.Find("img").First()
	if len(firstImage.Nodes) != 0 {
		thumbnail = firstImage.Nodes[0].Attr[src].Val
	} else {
		thumbnail = ""
	}
	
	p.Thumbnail = thumbnail

	return p.Thumbnail, nil
}

// Database에 게시글 저장하는 함수
func (p *Post) NewPost(db *sql.DB, hashTags string) error {
	_, err := db.Exec(`insert into posts values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, p.Id, p.Author, p.UDesc, p.Title, p.Thumbnail, p.Content, p.Summary, p.Date, p.Updated, p.Category, hashTags)
	if err != nil {
		log.Println(err)
		
		return err
	}
	
	// 해쉬태그부분 나눠서 리팩토링중
	var pid int
	err = db.QueryRow(`select id from posts where title = ? AND content = ?`, p.Title, p.Content).Scan(&pid)
	if err != nil {
		log.Println(err)
	}
	
	hashTagArr := p.convertHashTag(hashTags)
	
	for _, val := range(hashTagArr) {
		go func() {
			_, err = db.Exec(`INSERT INTO hashtag(tag) SELECT * FROM (SELECT ?) AS tmp WHERE NOT EXISTS(SELECT * FROM hashtag WHERE tag = ?) LIMIT 1;`, val, val)
			if err != nil {
				log.Println(err)
			}


			_, err = db.Exec(`INSERT INTO post_hashtag(hid, pid) SELECT id, ? FROM hashtag where tag = ?;`, pid, val)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	// -------------------------
	
	log.Printf("Post \"%s\" is posted on %s\n", p.Title, p.Date.Format("2006-01-02 15:04:05"))
	
	return nil
}

// Update Post
func (p *Post) UpdatePost(db *sql.DB, newHashTags string, pid int) error {
	_, err := db.Exec("update posts set title = ?,thumbnail = ?, content = ?, category = ?, hashtag = ?, updated = ? where id = ?", p.Title, p.Thumbnail, p.Content, p.Category, newHashTags, p.Updated, pid)
	if err != nil {
		log.Println(err)
		
		return err
	}
	
	return nil
}

// Database에서 하나의 게시글을 가져오는 함수
func (p *Post) GetPostFromDB(db *sql.DB, pid int) error {
	var hashs sql.NullString
	var thumbnail sql.NullString
	
	err := db.QueryRow("SELECT * FROM posts WHERE id = ?", pid).Scan(&p.Id, &p.Author, &p.UDesc, &p.Title, &thumbnail, &p.Content, &p.Summary, &p.Date, &p.Updated, &p.Category, &hashs)
	if err != nil {
		log.Println(err)
		
		return err
	}
	p.convertHashTag(hashs.String)
	p.Thumbnail = thumbnail.String
	
	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	_, err := db.Exec("delete from posts where id = ?", p.Id)
	if err != nil {
		log.Println(err)
	}

	// auto_increment initialize and sort.
	_, err = db.Exec("ALTER TABLE posts AUTO_INCREMENT=1;")
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec("SET @COUNT = 0;")
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec("UPDATE posts SET id = @COUNT := @COUNT+1;")
	if err != nil {
		log.Println(err)
	}
	
	return nil
}