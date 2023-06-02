import requests
import json

data = {
    "date": "TEST",
    "student_name": "Amelia Hamilton",
    "student_number": "777777777",
    "course_code": "SCH4U1",
    "student_section": 4,
    "unit_number": 5
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())