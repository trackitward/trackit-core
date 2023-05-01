package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var path_to_data = "./data/units/"
var path_to_profiles = "./data/profiles/"

type Unit struct {
	Course_Code     string `json:"course_code"`
	Course_Name     string `json:"course_name"`
	Course_Section  int    `json:"course_section"`
	Unit_Number     int    `json:"unit_number"`
	Unit_Completed  bool   `json:"unit_completed"`
	Submission_Date string `json:"submission_date"`
}

type Course struct {
	Course_Code        string `json:"course_code"`
	Course_Name        string `json:"course_name"`
	Course_Teacher     string `json:"course_teacher"`
	Course_Total_Units int    `json:"course_total_units"`
}

type User_Course struct {
	Course_Info  Course `json:"course_info"`
	User_Section int    `json:"user_section"`
	User_Info    struct {
		Units_Completed_Number   int    `json:"units_completed_number"`
		Units_Uncompleted_Number int    `json:"units_uncompleted_number"`
		Units                    []Unit `json:"units"`
		Last_Unit_Date           string `json:"last_unit_date"`
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
			Student_Name      string `json:"student_name"`
			Student_Number    string `json:"student_number"`
			Student_Grade     int    `json:"student_grade"`
			Student_Ta_Name   string `json:"student_ta_name"`
			Student_Ta_Number int    `json:"student_ta_number"`
		} `json:"student_data"`
		Course_Data []struct {
			UserCourse User_Course `json:"user_course"`
		} `json:"course_data"`
		Unit_Data struct {
			Units_Completed   int `json:"units_completed"`
			Units_Uncompleted int `json:"units_uncompleted"`
			Units_Total       int `json:"units_total"`
		} `json:"unit_data"`
	} `json:"data"`
}

type UnitSubmission struct {
	Code            string `json:"code,omitempty"`
	Date            string `json:"date"`
	Ticks           int    `json:"ticks,omitempty"`
	Student_Number  string `json:"student_number"`
	Student_Name    string `json:"student_name"`
	Course_Code     string `json:"course_code"`
	Student_Section int    `json:"student_section"`
	Unit_Number     int    `json:"unit_number"`
	Expiry          int    `json:"expiry,omitempty"`
}

type UserProfile struct {
	CreatedAt     int      `json:"created_at,omitempty"`
	StudentNumber string   `json:"student_number"`
	Email         string   `json:"email,omitempty"`
	Password      string   `json:"password"`
	Courses       []Course `json:"courses,omitempty"`
}

type UnitConfirmation struct {
	Date           string `json:"date"`
	Student_Name   string `json:"student_name"`
	Student_Grade  int    `json:"student_grade"`
	Course_Code    string `json:"course_code"`
	Course_Name    string `json:"course_name"`
	Unit_Number    int    `json:"unit_number"`
	Last_Submitted string `json:"last_submitted"`
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

	//GET ROUTES
	route.HandleFunc("/get/user/test", getTestUser).Methods("GET")
	route.HandleFunc("/get/user/{id}/{info}", getStudentData).Methods("GET")
	route.HandleFunc("/get/user/{id}/{info}/{course_code}", getStudentData).Methods("GET")

	route.HandleFunc("/post/user/create", createUser).Methods("POST")
	route.HandleFunc("/post/user/profile/create", createUserProfile).Methods("POST")
	route.HandleFunc("/post/user/profile/auth", authorizeUser).Methods("POST")

	route.HandleFunc("/post/unit/submit", generateUnitSubmissionCode).Methods("POST")
	route.HandleFunc("/post/unit/submit/validate", acceptUnitSubmission).Methods("POST")
	route.NotFoundHandler = http.HandlerFunc(notFound)

	http.ListenAndServe(":31475", router)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)

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
