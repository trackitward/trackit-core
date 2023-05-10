import requests
import json

data = {
    "date": "TEST",
    "student_name": "Amelia Hamilton",
    "student_number": "147655378",
    "course_code": "SCH4U1",
    "student_section": 8,
    "unit_number": 12
}

response = requests.post('http://localhost:31475/post/unit/submit', data=json.dumps(data))

print(response.json())