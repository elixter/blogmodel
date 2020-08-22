package models

import (
	"time"
)

type Comment struct {
	Id			int		`json: "id" db: "id" gorm: "id"`
	Content			string			`json: "content" db: "content" json: "content"`
	Author			string			`json: "author" db: "author" gorm: "author"`
	Date		time.Time		`json: "date" db: "date" gorm: "date"`
	Updated			time.Time		`json: "updated" db: "updated" gorm: "updated"`
	Pid			int			`json: "pid" db: "pid" gorm: "pid"`
}