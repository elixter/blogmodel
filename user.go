package models

import (
	"time"
	"encoding/json"
	"log"
	"database/sql"
	
	"github.com/gorilla/sessions"
)

type User struct {
	Id			string		`json: "id" db: "id" gorm: "id"`
	Name		string		`json: "name" db: "name" gorm: "name"`
	Admin		int		`json: "admin" db: "admin" gorm: "admin"`		// 1: 관리자, 0: 일반 유저
	Desc		string		`json: "desc" db: "desc" gorm: "desc"`		// user describe
	SessionId	string		`json: "sessionId"`
	ExpiresAt	time.Time	`json: "expiresAt"`							// 세션 만료시간
}

func (u *User) Valid() bool {
	// Check user's session is valid
	log.Println(time.Now())
	log.Println(u.ExpiresAt)
	log.Println(u.ExpiresAt.Sub(time.Now()))
	return u.ExpiresAt.Sub(time.Now()) > 0
}

func (u *User) Check(driver *sql.DB, check ... string) bool {
	// 데이터베이스에 유저 세션이 존재하는지 확인하는 함수
	// 세션아이디 또는 유저아이디로 확인
	// check == "ID" - 유저 아이디로 확인
	// check == "SESSIONS" - 세션 아이디로 확인
	// Default check = "ID"
	var sessCheck int
	
	switch check[0] {
		case "ID":
			err := driver.QueryRow("select count(id) from sessions where uid = ?", u.Id).Scan(&sessCheck)
			if err != nil {
				log.Println(err)
			}
		case "SESSIONS":
			err := driver.QueryRow("select count(id) from sessions where id = ?", u.SessionId).Scan(&sessCheck)
			if err != nil {
				log.Println(err)
			}
		default:
			err := driver.QueryRow("select count(id) from sessions where uid = ?", u.Id).Scan(&sessCheck)
			if err != nil {
				log.Println(err)
			}
	}
	if (sessCheck != 0) {
		return true
	} else {
		return false
	}
}

func (u *User) Refresh() time.Time {
	// Session 30분 연장
	u.ExpiresAt = u.ExpiresAt.Add(time.Minute * 30)
	
	return u.ExpiresAt
}

func GetUser(session *sessions.Session, userKey string) (*User, error) {
	var err error
	u := new(User)
	
	// 세션에서 유저정보 가져오기.
	jUser := session.Values[userKey]
	if jUser != nil {
		err = json.Unmarshal(jUser.([]byte), &u)
		if err != nil {
			log.Println(err)
		}
	}
	
	return u, err
}

func (u *User) GetAdmin(db *sql.DB, session *sessions.Session, userKey string) int {
	if session.Values[userKey] != nil {
		// Database에서 유저세션 정보 확인
		if u.Check(db, "ID") {
			if u.Valid() {
				return 1
			}
		} else {
			return 0
		}
	} else {
		return -1
	}
	// 세션에 유저정보가 없으면 -1 반환
	return -1
}