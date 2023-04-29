package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func createUser(response http.ResponseWriter, request *http.Request) {
	if request.Header.Get("API-PASS") != "PASSTOAPI-TRACKER" {
		http.Error(response, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var file *File

	unmarshal_err := json.Unmarshal(body, &file)
	if unmarshal_err != nil {
		fmt.Print(unmarshal_err)
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	name := string(path_to_data + file.Data.Student_Data.Student_Number + ".json")

	if _, err := os.Stat(name); err == nil {
		http.Error(response, "User Already Exists", http.StatusSeeOther)
	} else if os.IsNotExist(err) {
		os.Create(name)

		file_out, _ := json.MarshalIndent(file, "", "    ")

		_ = os.WriteFile(name, file_out, 0644)
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

	byteValue, _ := io.ReadAll(json_file)

	var file *File

	json.Unmarshal(byteValue, &file)
	fmt.Println(file)
	json.NewEncoder(response).Encode(file)
}

func getStudentData(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]
	info := params["info"]
	course_code := params["course_code"]

	if _, err := os.Stat(path_to_data + id + ".json"); err == nil {
		json_file, err := os.Open("./data/" + id + ".json")
		if err != nil {
			log.Fatal(err)
		}
		defer json_file.Close()

		byteValue, _ := io.ReadAll(json_file)

		var file File

		json.Unmarshal(byteValue, &file)

		if info == "all_data" {
			json.NewEncoder(response).Encode(file)
		} else if info == "student_data" {
			json.NewEncoder(response).Encode(file.Data.Student_Data)
		} else if info == "all_courses" {
			json.NewEncoder(response).Encode(file.Data.Course_Data)
		} else if info == "unit_data" {
			json.NewEncoder(response).Encode(file.Data.Unit_Data)
		} else if info == "course_code" {
			for i := 0; i < len(file.Data.Course_Data); i++ {
				if file.Data.Course_Data[i].UserCourse.Course_Info.Course_Code == strings.ToUpper(course_code) {
					json.NewEncoder(response).Encode(file.Data.Course_Data[i])
				}
			}
		} else {
			result := `{"status":404, "message":"Wrong parameter included. Query /endpoints for all acceptable endpoints/params."}`
			var finalResult map[string]interface{}
			json.Unmarshal([]byte(result), &finalResult)
			json.NewEncoder(response).Encode(finalResult)
		}

	} else if os.IsNotExist(err) {
		result := `{"status":404, "message":"User does not exist."}`
		var finalResult map[string]interface{}
		json.Unmarshal([]byte(result), &finalResult)

		json.NewEncoder(response).Encode(finalResult)
	}
}
