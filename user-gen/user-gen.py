import requests
import random
import json

base_url = "http://localhost:31475"

for i in range(2):
    first_names = []
    last_names = []

    with open('names.txt', 'r') as f:
        lines = f.readlines()


        for line in lines:
            if line == '\n':
                break
            
            first_names.append(line.rstrip())
            lines.remove(line)
        
        for line in lines:
            last_names.append(line.rstrip())

        student_name = random.choice(first_names) + " " + random.choice(last_names)

    student_number = random.randint(100000000, 999999999)
    student_grade = random.choice([9,10,11,12])
    student_ta_name = random.choice(['Mrs.', 'Mr.']) + random.choice(last_names)
    student_ta_number = random.randint(0,69+1)

    number_of_courses = 8

    courses = []

    total_units_completed = 0
    total_units_uncompleted = 0

    for i in range(number_of_courses):
        with open('courses.txt', 'r') as f:
            lines = f.readlines()

            course = random.choice(lines).split('-')
            course_code = course[0].rstrip()
            course_name = course[1].rstrip()
            course_teacher = random.choice(['Mrs.', 'Mr.']) + random.choice(last_names)
            course_total_units = 18

            user_section = random.randint(1,11)

            units = []

            units_completed = random.randint(0, number_of_courses*course_total_units)
            units_uncompleted = number_of_courses*course_total_units - units_completed

            total_units_completed += units_completed
            total_units_uncompleted += units_uncompleted

            for i in range(course_total_units):
                unit = {
                            "course_code": course_code,
                            "course_name": course_name,
                            "course_section": user_section,
                            "unit_number": i,
                            "unit_completed": random.choice([True, False]),
                            "submission_date": "DATE"
                        }
                
                units.append(unit)
            
            last_submission_date = "DATE"

            course = {
                        "user_course": {
                            "course_info":{
                                "course_code": course_code,
                                "course_name": course_name,
                                "course_teacher": course_teacher,
                                "course_total_units": course_total_units
                            },
                            "user_section": user_section,
                            "user_info": {
                                "units_completed_number": units_completed,
                                "units_uncompleted_number": units_uncompleted,
                                "units": units,
                                "last_submitted_date": "DATE"
                            }, 
                        }
                    }
            courses.append(course)

        

    data = {
        "meta": {
            "user_file_version": "1.0.1",
            "creation_date": "",
            "last_logged_in": ""
        },
        "data":{
            "student_data":{
                "student_name":student_name,
                "student_number": str(student_number),
                "student_grade": student_grade,
                "student_ta_name": student_ta_name,
                "student_ta_number": student_ta_number
            },
            "course_data": courses,
            "unit_data": {
                "units_completed": total_units_completed,
                "units_uncompleted": total_units_uncompleted,
                "units_total": total_units_completed+total_units_uncompleted
            }
        }
    }

    print(data)

    headers = {'Content-type': 'application/json', 'Accept': 'text/plain', 'API-PASS': 'PASSTOAPI-TRACKER'}
    response = requests.post(base_url + "/post/user/create", data=json.dumps(data), headers=headers)

    print(response.text)