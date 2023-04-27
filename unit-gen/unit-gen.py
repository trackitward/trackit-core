import requests
import json

data = {
    "date": "TEST",
    "student_name": "RNA",
    "course_code": "RZZ4U1",
    "student_section": 4,
    "unit_number": 1
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.text)