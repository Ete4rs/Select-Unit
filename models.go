package main

import (
	"gorm.io/gorm"
	"time"
)

type Student struct{
	gorm.Model
	ID string `gorm:"primaryKey; unique; default=false; not null;"`
	Name string `gorm:"not null; size:256"`
	LastName string `gorm:"not null; size:256"`
	Password string `gorm:"not null; size:256"`
}

type Master struct{
	gorm.Model
	ID 		 string `gorm:"primaryKey; unique; default=false; not null; size:7"`
	Name     string `gorm:"not null; size:18"`
	LastName string `gorm:"not null; size:28"`
	Password string `gorm:"not null; size:256"`
}

type Lesson struct {
	gorm.Model
	ID string `gorm:"primaryKey; unique; default=false; not null; size:5"`
	Name string `gorm:"default=false; not null; size:16"`
	Unit int `gorm:"default=false; not null; size:3"`
}

type Section struct {
	gorm.Model
	ID        int64 `gorm:"unique; primaryKey; autoIncrement; not null"`
	LessonID  string
	Lesson    Lesson `gorm:"references:ID ;foreignKey:LessonID"`
	MasterID  string
	Master    Master `gorm:"references:ID ;foreignKey:MasterID"`
	BeginTime time.Time
	EndTime   time.Time
	Day       string
	Capacity  int	`gorm:"not null; default=false; size=5"`
}

type SelectUnit struct {
	gorm.Model
	SectionID int64
	Section   Section `gorm:"references:ID ;foreignKey:SectionID"`
	StudentID string
	Student   Student `gorm:"references:ID ;foreignKey:StudentID"`

}

type TokenDetails struct {
	ID uint64
	AccessToken  string
	person string
	AtExpires    time.Time
}

type AllSection struct {
	SectionId  string `json:"sectionid"`
	LessonName string `json:"lessonname"`
	Unit       int    `json:"unit"`
	MasterName string `json:"mastername"`
	BeginTime  string `json:"begintime"`
	EndTime    string `json:"endtime"`
	Day        string `json:"day"`
	Capacity   int    `json:"capacity"`
}

type PersonJson struct {
	Name		string	`json:"name"`
	LastName	string	`json:"lastname"`
	password 	string	`json"password"`
}

type MySecions struct {
	MasterName string	`json:"mastername"`
	BeginTime string	`json:"begintime"`
	EndTime string		`json:"endtime"`
	Day string			`json:"day"`
	Unit int			`json:"unit"`
	LessonName string	`json:"lessonname"`
}