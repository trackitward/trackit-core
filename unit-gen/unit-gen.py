import requests
import json

data = {
    "date": "TEST",
    "student_name": "Lauren Young",
    "student_number": "111021531",
    "course_code": "SCH4U1",
    "student_section": 10,
    "unit_number": 3
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())