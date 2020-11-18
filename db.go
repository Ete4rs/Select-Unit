package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var  (
	db = &gorm.DB{}
)

func createDBConnection() {
	content, err := ioutil.ReadFile("userpass.txt")
	if err!=nil{
		log.Fatal("Error in reading user & password file :", err.Error())
	}
	info := strings.Split(string(content), " ")
	dsn := "user=" + info[0] + " password=" + info[1] + " dbname=fartak	 port=5432 TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err!=nil{
		log.Fatal("Error in opening database", err.Error())
	}
	db.AutoMigrate(&Master{}, &Student{}, &Lesson{}, &Section{}, &SelectUnit{})
}

func signUPStudent(id, name, lastName, password string) error{
	return db.Create(&Student{ID: id, Name: name, LastName: lastName, Password: password}).Error
}

func signUPMaster(id, name, lastName, password string) error{
	return db.Create(&Master{ID: id, Name: name, LastName: lastName, Password: password}).Error
}

func addLessonDB(id, name string, unit int)  error{
	return db.Create(&Lesson{ID: id, Name: name, Unit: unit}).Error
}

func isExist(id, password string, who bool) error{
	if who {
		return db.First(&Master{}, "id = ?", id, "password = ?", password).Error
	}else {
		return db.First(&Student{}, "id = ?", id, "password = ?", password).Error
	}
}

func idExist(id string, who bool)  error{
	if who {
		return db.First(&Master{}, "id = ?", id).Error
	}else {
		return db.First(&Student{}, "id = ?", id).Error
	}
}

func getAllLessons()  []Lesson{
	var lesson []Lesson
	db.Table("lessons").Select("id, name, unit").Scan(&lesson)
	return lesson
}

func createSection(masterId, lessonId , day string, beginTime, endTime time.Time, capacity int)  error{
	return db.Create(&Section{MasterID: masterId, LessonID: lessonId,
		BeginTime: beginTime, EndTime: endTime, Day: day, Capacity: capacity}).Error
}

func checkMasterSection(masterId, day string, beginTime, endTime time.Time) error{
	return db.First(&Section{}, "master_id = ?", masterId, "day = ?", day, "begin_time = ?", beginTime, "end_time = ?", endTime).Error
}

func getAllSection()  []AllSection{
	var sectionsMap []AllSection
	var sections []Section
	db.Table("sections").Select("id, lesson_id, master_id, begin_time, end_time, day, capacity").Scan(&sections)
	for _, section := range sections {
		var (
			master Master
			lesson Lesson
		)
		db.First(&Master{}, "id = ?", section.MasterID).Scan(&master)
		db.First(&Lesson{}, "id = ?", section.LessonID).Scan(&lesson)
		sectionsMap = append(sectionsMap, AllSection{SectionId: strconv.Itoa(int(section.ID)) , LessonName: lesson.Name, Unit: lesson.Unit,
			MasterName: master.Name + " " + master.LastName, BeginTime: getSectionTime(section.BeginTime),
				EndTime: getSectionTime(section.EndTime), Day: section.Day, Capacity: section.Capacity})
	}
	return sectionsMap
}

func checkLessonExistence(id string)  bool{
	err := db.First(&Lesson{}, "id = ?", id)
	if err==nil{
		return false
	}
	return true
}

func selectUnit(sectionID int64, studentID string) error {
	var section Section
	err := db.First(&Section{}, "id = ?", sectionID).Scan(&section)
	if err.Error != nil {
		return errors.New("section not found")
	}
	if section.Capacity > 0 {
		err := db.Create(&SelectUnit{SectionID: sectionID, StudentID: studentID}).Error
		if err != nil {
			return errors.New("can not take this section ")
		}
		db.Model(&Section{}).Where("id = ?", sectionID).Update("capacity", section.Capacity-1)
		return nil
	}else {
		return errors.New("this section is full ")
	}
}

func getSectionOfStudent(studentID string) ([]MySecions,error) {
	var mySection []MySecions
	var myUnits []SelectUnit
	err := db.Model(&SelectUnit{}).Where("student_id = ?", studentID).Scan(&myUnits)
	if err.Error != nil {
		return nil, err.Error
	}
	fmt.Println(myUnits)
	for _, unit := range myUnits {
		var (
			sec Section
			mas Master
			less Lesson
		)
		db.Model(&Section{}).First("id = ?", unit.SectionID).Scan(&sec)
		db.Model(&Master{}).First("master_id = ?", sec.MasterID).Scan(&mas)
		db.Model(&Lesson{}).First("lesson_id = ?", sec.LessonID).Scan(&less)
		mySection = append(mySection, MySecions{MasterName: mas.Name + " " + mas.LastName, BeginTime: getSectionTime(sec.BeginTime),
			EndTime: getSectionTime(sec.EndTime), Day: sec.Day, Unit: less.Unit, LessonName: less.Name})
	}
	return mySection, nil
}