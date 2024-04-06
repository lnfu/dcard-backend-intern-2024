import http.client
import json
from random import randint
from datetime import datetime, timedelta
from time import sleep

hostname = "localhost"
port = 8080
path = "/api/v1/ad"
headers = {"Content-Type": "application/json"}


def random_int(min, max):
    return randint(min, max)


gender_sets = [[], ["M"], ["F"]]
country_sets = [[], ["TW"], ["JP"], ["TW", "JP"]]
platform_sets = [
    [],
    ["android"],
    ["ios"],
    ["web"],
    ["android", "ios"],
    ["ios", "web"],
    ["web", "android"],
    ["android", "ios", "web"],
]

connection = http.client.HTTPConnection(hostname, port)

for i in range(1000):
    age_start = random_int(1, 100)
    now = datetime.now()
    today_start = datetime(now.year, now.month, now.day)
    today_end = datetime(now.year, now.month, now.day, 23, 59, 59, 999)
    postData = {
        "title": f"AD {i}",
        "startAt": today_start.isoformat() + 'Z',
        "endAt": today_end.isoformat() + 'Z',
        "conditions": [
            {
                "ageStart": age_start,
                "ageEnd": random_int(age_start, 100),
                "gender": gender_sets[random_int(0, len(gender_sets) - 1)],
                "country": country_sets[random_int(0, len(country_sets) - 1)],
                "platform": platform_sets[random_int(0, len(platform_sets) - 1)],
            }
        ],
    }

    connection.request("POST", path, json.dumps(postData), headers)
    response = connection.getresponse()
    data = response.read().decode()
    print(data)
    sleep(0.1)

connection.close()
