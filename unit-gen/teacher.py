import requests
import json

code = input('enter the code: ')

response2 = requests.post('http://localhost:31475/post/unit/submit/validate', data=json.dumps(code))

print(response2.json())