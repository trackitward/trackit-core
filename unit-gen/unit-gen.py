import requests
import json

data = {
    "date": "TEST",
    "student_name": "Ruth Martin",
    "student_number": "973419933",
    "course_code": "ENG4U1",
    "student_section": 2,
    "unit_number": 6
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())

code = input('enter the code: ')

response2 = requests.post('http://localhost:31475/post/unit/submit/validate', data=json.dumps(code))

print(response.json())