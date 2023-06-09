package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func generateUnitSubmissionCode(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var unitSubmission *UnitSubmission

	unmarshal_err := json.Unmarshal(body, &unitSubmission)
	if unmarshal_err != nil {
		fmt.Print(unmarshal_err)
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	code := fmt.Sprint(time.Now().Nanosecond())[:6]

	unitSubmission.Code = code
	unitSubmission.Ticks = int(time.Now().Unix())
	unitSubmission.Expiry = 120

	f, err := os.Open("units-in-submission.json")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := io.ReadAll(f)

	file, err := os.Open("./data/units/" + unitSubmission.Student_Number + ".json")
	if err != nil {
		log.Fatal(err)
	}

	fileData, _ := io.ReadAll(file)

	var studentFile *File

	unmarshal_studentfile_err := json.Unmarshal(fileData, &studentFile)
	if unmarshal_studentfile_err != nil {
		fmt.Print(unmarshal_err)
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	for current_iteration, _ := range studentFile.Data.Course_Data {
		if studentFile.Data.Course_Data[current_iteration].UserCourse.Course_Info.Course_Code == unitSubmission.Course_Code {
			unitSubmission.Last_Submitted = studentFile.Data.Course_Data[current_iteration].UserCourse.User_Info.Last_Unit_Date
		}
	}

	// Write current state to slice
	curr := []UnitSubmission{}
	json.Unmarshal(byteValue, &curr)

	for _, current := range curr {
		if current.Code == code {
			code = fmt.Sprint(time.Now().Nanosecond())[:6]
			for _, current := range curr {
				if current.Code == code {
					code = fmt.Sprint(time.Now().Nanosecond())[:6]
				}
			}
		}
	}

	// Append data to the created slice
	curr = append(curr, *unitSubmission)
	JSON, _ := json.MarshalIndent(curr, "", "    ")

	// Write
	_ = os.WriteFile("units-in-submission.json", JSON, 0644)

	json.NewEncoder(response).Encode(code)
}

func confirmUnitSubmission(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var code string

	unmarshal_err := json.Unmarshal(body, &code)
	if unmarshal_err != nil {
		fmt.Print(unmarshal_err)
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	f, err := os.Open("units-in-submission.json")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := io.ReadAll(f)

	// Write current state to slice
	curr := []UnitSubmission{}
	json.Unmarshal(byteValue, &curr)

	for i := 0; i < len(curr); i++ {
		if curr[i].Code == code {
			json.NewEncoder(response).Encode(curr[i])
			return
		}
	}

	result := `{"status":404, "message":"Code does not exist."}`
	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)

	json.NewEncoder(response).Encode(finalResult)
	response.Header().Add("UNIT-SUBMISSION", "FAILED")
}

func acceptUnitSubmission(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var code string

	unmarshal_err := json.Unmarshal(body, &code)
	if unmarshal_err != nil {
		fmt.Print(unmarshal_err)
		http.Error(response, "Bad Request - Wrong Body", http.StatusBadRequest)
		return
	}

	f, err := os.Open("units-in-submission.json")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := io.ReadAll(f)

	// Write current state to slice
	curr := []UnitSubmission{}
	json.Unmarshal(byteValue, &curr)

	for i := 0; i < len(curr); i++ {
		if curr[i].Code == code {
			if curr[i].Ticks+int(curr[i].Expiry) < int(time.Now().Unix()) {
				result := `{"status":400, "message":"Code Expired."}`
				var finalResult map[string]interface{}
				json.Unmarshal([]byte(result), &finalResult)

				curr[i] = curr[len(curr)-1]
				curr = curr[:len(curr)-1]
				JSON, _ := json.MarshalIndent(curr, "", "    ")

				// Write
				_ = os.WriteFile("units-in-submission.json", JSON, 0644)

				json.NewEncoder(response).Encode(finalResult)
				return
			}
			if _, err := os.Stat(path_to_data + curr[i].Student_Number + ".json"); err == nil {
				byteValue, _ := os.ReadFile(path_to_data + curr[i].Student_Number + ".json")

				var file File

				json.Unmarshal(byteValue, &file)

				unit_status := false

				for j := 0; j < len(file.Data.Course_Data); j++ {
					if file.Data.Course_Data[j].UserCourse.Course_Info.Course_Code == curr[i].Course_Code {
						for k := 0; k < len(file.Data.Course_Data[j].UserCourse.User_Info.Units); k++ {
							if file.Data.Course_Data[j].UserCourse.User_Info.Units[k].Unit_Number == curr[i].Unit_Number {
								file.Data.Course_Data[j].UserCourse.User_Info.Units[k].Unit_Completed = true
								file.Data.Course_Data[j].UserCourse.User_Info.Units_Completed_Number += 1
								file.Data.Course_Data[j].UserCourse.User_Info.Units_Uncompleted_Number -= 1
								file.Data.Unit_Data.Units_Completed += 1
								file.Data.Unit_Data.Units_Uncompleted -= 1
								file.Data.Course_Data[j].UserCourse.User_Info.Last_Unit_Date = time.Now().Format(time.RFC822)
								unit_status = true
							}
						}
					}
				}

				if !unit_status {
					curr[i] = curr[len(curr)-1]
					curr = curr[:len(curr)-1]
					JSON, _ := json.MarshalIndent(curr, "", "    ")

					// Write
					_ = os.WriteFile("units-in-submission.json", JSON, 0644)
					result := `{"status":404, "message":"Course/Unit does not exist."}`
					var finalResult map[string]interface{}
					json.Unmarshal([]byte(result), &finalResult)

					json.NewEncoder(response).Encode(finalResult)
					response.Header().Add("UNIT-SUBMISSION", "FAILED")
					return
				}

				name := string(path_to_data + file.Data.Student_Data.Student_Number + ".json")

				file_out, err := json.MarshalIndent(file, "", "    ")
				if err != nil {
					log.Fatal(err)
				}

				for j := 0; j < len(file.Data.Course_Data); j++ {
					if file.Data.Course_Data[j].UserCourse.Course_Info.Course_Code == curr[i].Course_Code {
						_ = os.WriteFile(name, file_out, 0644)
						response.WriteHeader(http.StatusCreated)
						json.NewEncoder(response).Encode(file.Data.Course_Data[j].UserCourse.User_Info.Units[curr[i].Unit_Number])
					}
				}
			} else if os.IsNotExist(err) {
				result := `{"status":404, "message":"User does not exist."}`
				var finalResult map[string]interface{}
				json.Unmarshal([]byte(result), &finalResult)

				json.NewEncoder(response).Encode(finalResult)
				response.Header().Add("UNIT-SUBMISSION", "FAILED")
				return
			}

			curr[i] = curr[len(curr)-1]
			curr = curr[:len(curr)-1]
			JSON, _ := json.MarshalIndent(curr, "", "    ")

			// Write
			_ = os.WriteFile("units-in-submission.json", JSON, 0644)
			return
		}
		fmt.Println("CODE VALIDATED")
	}
	result := `{"status":404, "message":"Code does not exist."}`
	var finalResult map[string]interface{}
	json.Unmarshal([]byte(result), &finalResult)

	json.NewEncoder(response).Encode(finalResult)
	response.Header().Add("UNIT-SUBMISSION", "FAILED")
}
