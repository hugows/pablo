import random
import requests
import time


def post(json):
    r = requests.post('http://localhost:5000/chart', json=json)
    print(r.status_code)
    try:
        print(r.json())
    except:
        pass


def main():
    count = 0.5
    while True:
        data = {'pressure': count, "temp": random.randint(50, 100)}
        count += 1.0
        try:
            post(data)
        except:
            print "Could not POST"
        time.sleep(1)


if __name__ == "__main__":
    main()
