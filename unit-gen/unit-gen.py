import requests
import json

data = {
    "date": "TEST",
    "student_name": "Donna Luke",
    "student_number": "554974039",
    "course_code": "SBI4U1",
    "student_section": 8,
    "unit_number": 2
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())