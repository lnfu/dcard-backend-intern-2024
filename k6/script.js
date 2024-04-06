import http from "k6/http";
import { SharedArray } from "k6/data";
import { sleep, check } from "k6";

const host = "http://localhost:8080";

export let options = {
  scenarios: {
    contacts: {
      executor: "constant-arrival-rate",
      rate: 24000,
      duration: "1m30s",
      preAllocatedVUs: 100,
			maxVUs: 15000,
    },
  },
};

// 產生 [min, max] 的隨機整數
function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

const limits = [10, 20, 30];
const genders = ["M", "F"];
const countries = ["TW", "JP"];
const platforms = ["android", "ios", "web"];

export default function () {
  const gender = genders[randomInt(0, genders.length - 1)]
  const country = countries[randomInt(0, countries.length - 1)]
  const platform = platforms[randomInt(0, platforms.length - 1)]
  const offset = 0
  const limit = limits[randomInt(0, limits.length - 1)]
  http.get(`${host}/api/v1/ad?gender=${gender}&country=${country}&platform=${platform}&offset=${offset}&limit=${limit}`);
  // console.log(res)
  sleep(1);
}
