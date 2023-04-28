import requests
import json

data = {
    "date": "TEST",
    "student_name": "Diana Hodges",
    "student_number": "110880792",
    "course_code": "HRE4M1",
    "student_section": 1,
    "unit_number": 3
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())

code = input('enter the code: ')

response2 = requests.post('http://localhost:31475/post/unit/submit/validate', data=json.dumps(code))

print(response2.json())