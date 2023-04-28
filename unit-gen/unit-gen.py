import requests
import json

data = {
    "date": "TEST",
    "student_name": "Diana Hodges",
    "student_number": "110880792",
    "course_code": "TEJ4U1",
    "student_section": 7,
    "unit_number": 8
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())