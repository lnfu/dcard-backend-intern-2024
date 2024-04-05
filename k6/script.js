import http from "k6/http";
import { SharedArray } from "k6/data";
import { sleep, check } from "k6";

const host = "http://localhost:8080";

export let options = {
  setupTimeout: "500s",
  scenarios: {
    contacts: {
      executor: "constant-arrival-rate",
      rate: 10000, // 200 RPS, since timeUnit is the default 1s
      duration: "1m",
      preAllocatedVUs: 50,
      maxVUs: 10000,
    },
  },
};

// 產生 [min, max] 的隨機整數
function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

const genders = [[], ["M"], ["F"]];
const countries = [[], ["TW"], ["JP"], ["TW", "JP"]];
const platforms = [
  [],
  ["android"],
  ["ios"],
  ["web"],
  ["android", "ios"],
  ["ios", "web"],
  ["web", "android"],
  ["android", "ios", "web"],
];

export function setup() {
  for (let i = 0; i < 1000; i++) {
    const ageStart = randomInt(1, 100);
    const now = new Date();
    const todayStart = new Date( // 今天的 00:00:00
      now.getFullYear(),
      now.getMonth(),
      now.getDate()
    );
    const todayEnd = new Date( // 今天的 23:59:59
      now.getFullYear(),
      now.getMonth(),
      now.getDate(),
      23,
      59,
      59,
      999
    );

    let postData = {
      title: `AD ${i}`,
      startAt: todayStart.toISOString(),
      endAt: todayEnd.toISOString(),
      conditions: [
        {
          ageStart: ageStart,
          ageEnd: randomInt(ageStart, 100),
          gender: genders[randomInt(0, genders.length - 1)],
          country: countries[randomInt(0, countries.length - 1)],
          platform: platforms[randomInt(0, platforms.length - 1)],
        },
      ],
    };
    let res = http.post(`${host}/api/v1/ad`, JSON.stringify(postData), {
      headers: {
        "Content-Type": "application/json",
      },
    });
    console.log(res.body);
    check(res, {
      "status is 200": (r) => r.status === 200,
    });
    sleep(0.1);
  }
}

export default function () {
  let res = http.get(`${host}/api/v1/ad`);
  sleep(1);
}
