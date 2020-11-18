package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/*
	############################################################
								Student
	############################################################
*/
func handleStudentSignUP(response http.ResponseWriter, request *http.Request)  {
	var person PersonJson
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}else {
		err = json.Unmarshal(requestBody, &person)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
		}else {
			id := createID(false)
			err = signUPStudent(id, person.Name, person.LastName, person.password)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
			}else {
				uid, _ := strconv.Atoi(id)
				td, err := CreateToken(uint64(uid), "student")
				if err!=nil{
					http.Error(response, err.Error(), http.StatusBadRequest)
				}else {
					tokens[id] = td
					response.WriteHeader(http.StatusCreated)
					response.Write([]byte(td.AccessToken))
				}
			}
		}
	}
}

func handleStudentLogin(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "student")
		if err!=nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			response.WriteHeader(http.StatusOK)
		}
	} else {
		id, ok := header["Id"]
		password , okP  := header["Password"]
		if  ok && okP{
			err := isExist(id, password, false)
			if err!=nil{
				http.Error(response, "no such id and password ", http.StatusUnauthorized)
			}else {
				userID, err := strconv.Atoi(id)
				if err!=nil {
					http.Error(response, "Bad request parameter ", http.StatusBadRequest)
				}
				td, err := CreateToken(uint64(userID), "student")
				if err!=nil{
					http.Error(response, err.Error(), http.StatusBadRequest)
				}else {
					tokens[id] = td
					response.WriteHeader(http.StatusCreated)
					response.Write([]byte(td.AccessToken))
				}
			}

		}else {
			http.Error(response, "Error in parameter of header of request ", http.StatusBadRequest)
		}
	}
}

func handleStudentGetAllSection(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "student")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			sections:= getAllSection()
			if sections==nil {
				http.Error(response, "no section found", http.StatusInternalServerError)
			}else {
				jsonBody, err := json.Marshal(sections)
				if err != nil {
					http.Error(response, err.Error(), http.StatusInternalServerError)
				}else {
					response.WriteHeader(http.StatusOK)
					response.Write(jsonBody)
				}
			}
		}
	}else {
		http.Error(response, "no token found ", http.StatusBadRequest)
	}
}

func handleStudentSelectUnit(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "student")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		} else {
			type Unit struct {
				SectionID int `json:"id"`
			}
			var unit Unit
			requestBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
			}else {
				err = json.Unmarshal(requestBody, &unit)
				if err != nil {
					http.Error(response, err.Error(), http.StatusBadRequest)
				}else {
					err = selectUnit(int64(unit.SectionID), strconv.Itoa(unit.SectionID))
					if err!=nil{
						http.Error(response, err.Error(), http.StatusBadRequest)
					}else {
						response.WriteHeader(http.StatusOK)
					}

				}
			}
		}
	}
}

func handleStudentGetUnits(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "student")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		} else {
			id, err := extractIdFromToken(token)
			if err!=nil{
				http.Error(response, err.Error(), http.StatusUnauthorized)
			}else {
				sections, err := getSectionOfStudent(id)
				if err !=nil{
					http.Error(response, err.Error(), http.StatusBadRequest)
				}else {
					jsonBody, err := json.Marshal(sections)
					if err != nil {
						http.Error(response, err.Error(), http.StatusInternalServerError)
					}else {
						response.WriteHeader(http.StatusOK)
						response.Write(jsonBody)
					}
				}
			}
		}
	}
}
/*
	############################################################
								Master
	############################################################
*/
func handleMasterSignUP(response http.ResponseWriter, request *http.Request)  {
	var person PersonJson
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
	}else {
		err = json.Unmarshal(requestBody, &person)
		if err != nil {
			http.Error(response, err.Error(), http.StatusBadRequest)
		}else {
			id := createID(true)
			err = signUPMaster(id, person.Name, person.LastName, person.password)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
			}else {
				uid, _ := strconv.Atoi(id)
				td, err := CreateToken(uint64(uid), "master")
				if err!=nil{
					http.Error(response, err.Error(), http.StatusBadRequest)
				}else {
					tokens[id] = td
					response.WriteHeader(http.StatusCreated)
					response.Write([]byte(td.AccessToken))
				}
			}
		}
	}
}

func handleMasterLogin(response http.ResponseWriter, request *http.Request){
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "master")
		if err!=nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			response.WriteHeader(http.StatusOK)
		}
	} else {
		id, ok := header["Id"]
		password , okP  := header["Password"]
		if  ok && okP{
			err := isExist(id, password, true)
			if err!=nil{
				http.Error(response, "no such id and password ", http.StatusUnauthorized)
			}else {
				userID, err := strconv.Atoi(id)
				if err!=nil {
					http.Error(response, "Bad request parameter ", http.StatusBadRequest)
				}
				td, err := CreateToken(uint64(userID), "master")
				if err!=nil{
					http.Error(response, err.Error(), http.StatusBadRequest)
				}
				tokens[id] = td
				response.WriteHeader(http.StatusCreated)
				response.Write([]byte(td.AccessToken))
			}
		}else {
			http.Error(response, "Error in parameter of header of request ", http.StatusBadRequest)
		}
	}
}

func handleMasterGetAllLessons(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "master")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			lessons := getAllLessons()
			jsonBody, err :=json.Marshal(lessons)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
			}else {
				response.WriteHeader(http.StatusOK)
				response.Write(jsonBody)
			}
		}
	}else {
		http.Error(response, "no token found ", http.StatusBadRequest)
	}
}

func handleMasterCreateSection(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "master")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			requestBody := map[string]string{}
			jsonbody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				response.WriteHeader(400)
			}else {
				json.Unmarshal(jsonbody, &requestBody)
				err := checkInfoRequestForCreateSection(requestBody)
				if !err {
					http.Error(response, errors.New("Bad request parameter ").Error(), http.StatusBadRequest)
				}else {
					id, err := extractIdFromToken(header["Token"])
					if err != nil {
						http.Error(response, "Bad token ", http.StatusUnauthorized)
					}
					less := checkLessonExistence(requestBody["id"])
					beginTime, beginBool := checkTimeOfSection(requestBody["begintime"])
					endTime, endBool := checkTimeOfSection(requestBody["endtime"])
					dayBool := checkDay(requestBody["day"])
					capacity, err := strconv.Atoi(requestBody["capacity"])
					if !less {
						http.Error(response, errors.New("Bad lesson ID ").Error(), http.StatusBadRequest)
					} else if err != nil || !dayBool || !beginBool || !endBool {
						http.Error(response, errors.New("Bad parameter ").Error(), http.StatusBadRequest)
					}else {
						err := checkMasterSection(id, requestBody["day"], beginTime, endTime)
						if err == nil {
							http.Error(response, "Bad time for the master ", http.StatusBadRequest)
						}else {
							err = createSection(id, requestBody["id"], requestBody["day"], beginTime, endTime, capacity)
							if err != nil {
								http.Error(response, err.Error(), http.StatusBadRequest)
							}else {
								response.WriteHeader(http.StatusCreated)
							}
						}
					}
				}
			}
		}
	}else {
		http.Error(response, "no token found ", http.StatusBadRequest)
	}
}

/*
	############################################################
								Admin
	############################################################
 */
func handleAdminLogin(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"] ; ok {
		err := checkToken(token, "admin")
		if err!=nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			response.WriteHeader(http.StatusOK)
		}
	}else {
		_, userNameBool := header["Username"]
		password, passwordBool := header["Password"]
		if userNameBool && passwordBool {
			err := checkAdmin(password)
			if err!=nil {
				http.Error(response, err.Error(), http.StatusUnauthorized)
			}
			rand.Seed(time.Now().Unix())
			admin = rand.Intn(99999 - 10000)
			td, err := CreateToken(uint64(admin), "admin")
			if err!=nil{
				http.Error(response, err.Error(), http.StatusBadRequest)
			}
			tokens[string(admin)] = td
			response.WriteHeader(http.StatusCreated)
			response.Write([]byte(td.AccessToken))
		}else {
			http.Error(response, "Bad header request ", http.StatusBadRequest)
		}
	}
}

func handleAdminAddLesson(response http.ResponseWriter, request *http.Request)  {
	header := extractHeader(request.Header)
	if token, ok := header["Token"]; ok {
		err := checkToken(token, "admin")
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
		}else {
			id := request.FormValue("id")
			name := request.FormValue("name")
			unit, err := strconv.Atoi(request.FormValue("unit"))
			if err != nil {
				http.Error(response, "Bad unit number ", http.StatusBadRequest)
			}
			err = addLessonDB(id, name, unit)
			if err != nil {
				http.Error(response, "Bad unit number ", http.StatusBadRequest)
			}
			response.WriteHeader(http.StatusCreated)
		}
	}else {
		response.WriteHeader(http.StatusBadRequest)
	}
}




func main() {
	createDBConnection()

	http.HandleFunc("/admin/login", handleAdminLogin)
	http.HandleFunc("/admin/addlesson", handleAdminAddLesson)

	http.HandleFunc("/master/signup", handleMasterSignUP)
	http.HandleFunc("/master/login", handleMasterLogin)
	http.HandleFunc("/master/alllessons", handleMasterGetAllLessons)
	http.HandleFunc("/master/createsection", handleMasterCreateSection)

	http.HandleFunc("/student/signup", handleStudentSignUP)
	http.HandleFunc("/student/login", handleStudentLogin)
	http.HandleFunc("/student/sections", handleStudentGetAllSection)
	http.HandleFunc("/student/selectunit", handleStudentSelectUnit)
	http.HandleFunc("/student/units", handleStudentGetUnits)

	log.Fatal(http.ListenAndServe(":8090", nil))
}