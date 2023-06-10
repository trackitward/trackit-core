import requests
import json

data = {
    "student_number": "777777777",
    "password": "q1!"
}


headers = {'API-PASS': 'PASSTOAPI-TRACKER'}
response = requests.post("http://localhost:31475/post/user/profile/create", data=json.dumps(data), headers=headers)

print(response.text)