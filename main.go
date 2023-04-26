package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var path_to_data = "../data/"

type Unit struct {
	Course_Code     string      `json:"course_code"`
	Course_Name     string      `json:"course_name"`
	Course_Section  json.Number `json:"course_section"`
	Unit_Number     json.Number `json:"unit_number"`
	Unit_Completed  bool        `json:"unit_completed"`
	Submission_Date string      `json:"submission_date"`
}

type Course struct {
	Course_Code        string `json:"course_code"`
	Course_Name        string `json:"course_name"`
	Course_Teacher     string `json:"course_teacher"`
	Course_Total_Units string `json:"course_total_units"`
}

type User_Course struct {
	Course_Info  Course      `json:"course_info"`
	User_Section json.Number `json:"user_section"`
	User_Info    struct {
		Units_Completed_Number   json.Number `json:"units_completed_number"`
		Units_Uncompleted_Number json.Number `json:"units_uncompleted_number"`
		Units                    []Unit      `json:"units"`
		Last_Unit_Date           string      `json:"last_unit_date"`
	} `json:"user_info"`
}

type File struct {
	Meta struct {
		User_File_Version string `json:"user_file_version"`
		Creation_Date     string `json:"creation_date"`
		Last_Logged_In    string `json:"last_logged_in"`
	} `json:"meta"`
	Data struct {
		Student_Data struct {
			Student_Name      string      `json:"student_name"`
			Student_Number    string      `json:"student_number"`
			Student_Grade     json.Number `json:"student_grade"`
			Student_Ta_Name   string      `json:"student_ta_name"`
			Student_Ta_Number json.Number `json:"student_ta_number"`
		} `json:"student_data"`
		Course_Data []struct {
			UserCourse User_Course `json:"user_course"`
		} `json:"course_data"`
		Unit_Data struct {
			Units_Completed   json.Number `json:"units_completed"`
			Units_Uncompleted json.Number `json:"units_uncompleted"`
			Units_Total       json.Number `json:"units_total"`
		} `json:"unit_data"`
	} `json:"data"`
}

func createUser(response http.ResponseWriter, request *http.Request) {
	if request.Header.Get("API-PASS") != "PASSTOAPI-TRACKER" {
		http.Error(response, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var file *File

	unmarshal_err := json.Unmarshal(body, &file)
	if unmarshal_err != nil {
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	name := string(path_to_data + file.Data.Student_Data.Student_Number + ".json")

	if _, err := os.Stat(name); err == nil {
		http.Error(response, "User Already Exists", http.StatusSeeOther)
	} else if os.IsNotExist(err) {
		os.Create(name)

		file_out, _ := json.MarshalIndent(file, "", "    ")

		_ = ioutil.WriteFile(name, file_out, 0644)
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(file)
		return
	} else {
		log.Fatal(err)
		response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	json.NewEncoder(response).Encode(file)
}

func getTestUser(response http.ResponseWriter, request *http.Request) {
	json_file, err := os.Open("test.json")
	if err != nil {
		log.Fatal(err)
	}
	defer json_file.Close()

	byteValue, _ := ioutil.ReadAll(json_file)

	var file *File

	json.Unmarshal(byteValue, &file)
	fmt.Println(file)
	json.NewEncoder(response).Encode(file)
}

func notFound(response http.ResponseWriter, request *http.Request) {
	result := `{"status": 404, "message": "404 NOT FOUND"}`

	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)

	json.NewEncoder(response).Encode(finalResult)
}

func main() {
	//createTestUser()
	fmt.Println("Starting ByteTrack API...")

	route := mux.NewRouter()
	route.Use(commonMiddleware)

	router := cors.Default().Handler(route)

	route.HandleFunc("/user/post/create", createUser).Methods("POST")
	route.HandleFunc("/user/get/test", getTestUser).Methods(("GET"))
	route.NotFoundHandler = http.HandlerFunc(notFound)

	http.ListenAndServe(":31475", router)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fmt.Print(request.URL.Path)

		var globallimiter = rate.NewLimiter(50, 110)

		if !globallimiter.Allow() {
			ratelimited(response, request)
			return
		}

		response.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(response, request)
	})
}

func ratelimited(response http.ResponseWriter, request *http.Request) {
	result := `{"status":429, "message":"You are requesting too quickly!"}`
	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)
	json.NewEncoder(response).Encode(finalResult)
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusTooManyRequests)
}
