const http = require("http");

const options = {
  hostname: 'localhost',
  port: 8080,
  path: '/api/v1/ad',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
};

// 產生 [min, max] 的隨機整數
function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

const genderSets = [[], ["M"], ["F"]];
const countrySets = [[], ["TW"], ["JP"], ["TW", "JP"]];
const platformSets = [
  [],
  ["android"],
  ["ios"],
  ["web"],
  ["android", "ios"],
  ["ios", "web"],
  ["web", "android"],
  ["android", "ios", "web"],
];

for (let i = 0; i < 1000; i++) {
  const ageStart = randomInt(1, 100);
  const now = new Date();
  const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const todayEnd = new Date(
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
        gender: genderSets[randomInt(0, genderSets.length - 1)],
        country: countrySets[randomInt(0, countrySets.length - 1)],
        platform: platformSets[randomInt(0, platformSets.length - 1)],
      },
    ],
  };

  const req = http.request(options, (res) => {
    let data = '';
    res.on('data', (chunk) => {
      data += chunk;
    });
    res.on('end', () => {
      console.log(data);
    });
  });

  req.on('error', (error) => {
    console.error(error);
  });

  req.write(JSON.stringify(postData));
  req.end();
}
