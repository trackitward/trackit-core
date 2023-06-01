import requests
import random
import json

student_number = random.choice([123456789, 000000000, 777777777, 999999999, 987654321])
courses = []

with open(f'../data/units/{student_number}.json', 'r') as file:
    data = json.load(file)

    for course in data['data']["course_data"]:
        course_obj = {
            "course_code": course['user_course']['course_info']['course_code'],
            "course_name": course['user_course']['course_info']['course_name'],
            "course_teacher": course['user_course']['course_info']['course_teacher'],
            "course_total_units": course['user_course']['course_info']['course_total_units']
        }

        courses.append(course_obj)


data = {
    "created_at": 1682798000,
    "student_number": student_number,
    "email": "test@tcdsb.ca",
    "password": "TESTING-PASS",
    "courses": courses
}

headers = {'Content-type': 'application/json', 'Accept': 'text/plain', 'API-PASS': 'PASSTOAPI-TRACKER'}
response = requests.post("http://localhost:31475" + "/post/user/profile/create", data=json.dumps(data), headers=headers)