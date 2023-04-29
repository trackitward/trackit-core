package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) []byte {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return hashed_password
}

func createUserProfile(response http.ResponseWriter, request *http.Request) {
	if request.Header.Get("API-PASS") != "PASSTOAPI-TRACKER" {
		http.Error(response, "Unauthorized", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var user *UserProfile

	unmarshal_err := json.Unmarshal(body, &user)
	if unmarshal_err != nil {
		log.Fatal(err)
	}
	name := string(path_to_profiles + string(user.StudentNumber) + ".json")
	fmt.Print(name)

	if _, err := os.Stat(name); err == nil {
		http.Error(response, "User Already Exists", http.StatusSeeOther)
	} else if os.IsNotExist(err) {
		os.Create(name)

		user.CreatedAt = int(time.Now().Unix())
		hashed_password := string(hashPassword(user.Password))
		user.Password = hashed_password
		fmt.Print(hashed_password)

		file_out, _ := json.MarshalIndent(user, "", "    ")

		_ = os.WriteFile(name, file_out, 0644)
		response.WriteHeader(http.StatusCreated)
		json.NewEncoder(response).Encode(user)
		return
	} else {
		log.Fatal(err)
		response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	json.NewEncoder(response).Encode(user.CreatedAt)
}

func authorizeUser(response http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var user *UserProfile

	unmarshal_err := json.Unmarshal(body, &user)
	if unmarshal_err != nil {
		log.Fatal(err)
	}

	name := string(path_to_profiles + string(user.StudentNumber) + ".json")

	if _, err := os.Stat(name); err == nil {

		json_file, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer json_file.Close()

		byteValue, _ := io.ReadAll(json_file)

		var user_profile *UserProfile

		json.Unmarshal(byteValue, &user_profile)

		fmt.Print(user.Password, "\n", user_profile.Password)

		err = bcrypt.CompareHashAndPassword([]byte(user_profile.Password), []byte(user.Password))
		if err != nil {
			log.Println(err)
			http.Error(response, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(response).Encode(user_profile)
	} else {
		http.Error(response, "User Does Not Exist", http.StatusNotFound)
	}
}
