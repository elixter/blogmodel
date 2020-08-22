package models

import (
	"time"
)

type Session struct {
	Id			string			`json: "id" db: "id" gorm: "id"`
	Uid			string			`json: "uid" db: "uid" gorm: "uid"`
	CreatedAt		time.Time		`json: "createdAt" db: "createdAt" gorm: "createdAt"`
	ExpiresAt		time.Time		`json: "expiresAt" db: "expiresAt" gorm: "expiresAt"`
}