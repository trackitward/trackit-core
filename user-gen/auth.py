import requests
import json

data = {
    "student_number": "264249481",
    "password": "TESTING123"
}

response = requests.post("http://localhost:31475/post/user/profile/auth", data=json.dumps(data))

print(response.text)