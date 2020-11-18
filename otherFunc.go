package main

import (
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func extractHeader(header http.Header)  map[string]string{
	var h = make(map[string]string)
	for key, value := range header {
		if key=="Username" || key=="Password" || key=="Id"{
			h[key] = value[0]
		}else if key=="Token"{
			h[key] = value[0]
			return h
		}
	}
	return h
}

func createID(who bool)  string{
	if who {
		for {
			rand.Seed(time.Now().Unix())
			rnd := rand.Intn(999 - 100)
			id := "69" + strconv.Itoa(rnd)
			err := idExist(id, true)
			if err != nil {
				continue
			} else {
				return id
			}
		}
	}else {
		for {
			rand.Seed(time.Now().Unix())
			rnd := rand.Intn(999 - 100)
			id := "85" + strconv.Itoa(rnd)
			err := idExist(id, false)
			if err != nil {
				continue
			} else {
				return id
			}
		}
	}
}

func extractIdFromToken(token string) (string, error){
	for key, value := range tokens {
		if value.AccessToken == token {
			return key, nil
		}
	}
	return "", errors.New("Bad token ")
}

func checkInfoRequestForCreateSection(info map[string]string) bool{
	if _, ok := info["id"] ; !ok {
		return false
	}else if _, ok := info["begintime"] ; !ok {
		return false
	}else if _, ok := info["endtime"] ; !ok {
		return false
	}else if _, ok := info["day"] ; !ok {
		return false
	}else if _, ok := info["capacity"] ; !ok {
		return false
	}else {
		return true
	}
}

func checkTimeOfSection(t string)  (time.Time,bool) {
	match, _ := regexp.MatchString("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$", t)
	if !match{
		return time.Time{},false
	}else {
		return createTime(strings.Split(t, ":")), true
	}
	return time.Time{},false
}

func createTime(t []string)time.Time{
	hh, _ := strconv.Atoi(t[0])
	mm, _ := strconv.Atoi(t[1])
	return time.Date(2020, 10, 20, hh, mm, 00, 00, time.Local)
}

func checkDay(day string) bool {
	if day=="tuesday" || day=="saturday" || day=="sunday" || day=="monday" || day=="wednesday" || day=="thursday" || day=="friday" {
		return true
	}
	return false
}

func getSectionTime(t time.Time) string {
	str := ""
	if t.Hour()==0 {
		str = "00:"
	}else if t.Hour() < 9 {
		str = "0" + strconv.Itoa(t.Hour()) + ":"
	} else {
		str = strconv.Itoa(t.Hour()) + ":"
	}
	if t.Minute() == 0{
		str = str + "00"
	}else if t.Minute() < 9{
		str = str + "0" + strconv.Itoa(t.Minute())
	}else {
		str = str + strconv.Itoa(t.Minute())
	}
	return str
}