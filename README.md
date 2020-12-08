blogmodel
=========

* download
<pre>
<code>
go get github.com/elixter/blogmodel
</code>
</pre>

* 현재 post, user모델만 구현

- ## Post.go
   <pre><code>
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
   </code></pre>
   - Id : Post_id
   - Author : 게시글 작성자
   - UDesc : 작성자 소개(?)
   - Title : 게시글 제목
   - Thumbnail : 게시글 목록에 보여지는 썸네일 (없을경우 nil)
   - Content : 게시글 내용
   - Summary : 게시글 요약내용 게시글 목록에서 보여짐.
   - Date : 게시글 작성시간
   - Updated : 게시글 수정시간
   - Category : 게시글 카테고리
   - HashTags : 게시글 해쉬태그
   
 - #### func GetCurrentPagePosts(db *sql.DB, currentPage int, rowsPerPage int) ([]Post, error)
    ##### parameters
    - db : DB driver
    - currentPage : 현재 페이지번호
    - rowsPerPage : 게시글목록 페이지당 보여지는 게시글 수
    ##### returns
    - []Post : 현재 페이지에서 보여지는 게시글 객체 슬라이스 return
    - error : 성공적일 경우 nil, db쿼리에서 문제있을경우 db error return
    - GetCurrentPagePostsBy... 함수들도 같은원리로 작동 다만 어떤 카테고리 또는 어떤 해쉬태그가 들어간 게시글만 가져옴
    
 - #### func (p *Post) convertHashTag(hashTags string) []string
    ##### parameters
    - hashTags : DB에는 해쉬태그들이 공백으로 구분되어 하나의 string으로 저장되어있음. 그 string을 받는 매개변수
    ##### returns
    - string : 받아온 문자열에서 #기호들을 모두 제거 후 공백을 기준으로 토큰화하여 Post 구조체의 HashTags 변수에 저장후 p.HashTag return
    
 - #### func (p *Post) createSummary(summaryLength int) string
    ##### parameters
    - summaryLength : 게시글의 첫문자부터 summaryLength만큼의 문자를 가져오기위한 설정
    ##### returns
    - string : db로 부터 받아온 게시글에서 HTML태그들을 제거한 후 summaryLength만큼의 문자를 추출해서 p.Summary에 저장 후 p.Summary return
    
 - #### func (p *Post) CreateThumbnail() (string, error)
    ##### parameters
    - None
    ##### returns
    - string : 게시글에서 가장 처음나오는 이미지태그를 p.Thumbnail에 저장후 p.Thumbnail return
    - error : goquery로 게시글 내용을 가져오는과정에서 문제가 있을경우 error return, 정상일 경우 nil
