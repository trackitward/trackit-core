package main

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func hashPassword(password string) []byte {
	h := crypto.SHA256.New()
	h.Write([]byte(password))
	return h.Sum(nil)
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

		user.CreatedAt = time.Now().UTC().Second()
		hashed_password := hex.EncodeToString(hashPassword(user.Password))
		user.Password = hashed_password
		fmt.Print(hashed_password)

		file_out, _ := json.MarshalIndent(user, "", "    ")

		_ = ioutil.WriteFile(name, file_out, 0644)
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
