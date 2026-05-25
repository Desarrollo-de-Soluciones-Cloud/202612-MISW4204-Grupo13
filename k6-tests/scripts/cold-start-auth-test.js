import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL, jsonHeaders } from "./config.js";

const data = JSON.parse(open("../data/users.json"));
const users = data.auth_users || [];

export const options = {
  scenarios: {
    cold_start_probe: {
      executor: "constant-vus",
      vus: Number(__ENV.VUS || 5),
      duration: __ENV.DURATION || "2m",
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    "http_req_duration{phase:first_request}": ["p(95)<5000"],
    "http_req_duration{phase:stable}": ["p(95)<2000", "p(99)<4000"],
  },
  summaryTrendStats: ["avg", "min", "med", "max", "p(95)", "p(99)"],
};

export default function () {
  const user = users[(__VU - 1) % users.length];
  const phase = __ITER === 0 ? "first_request" : "stable";

  const response = http.post(
    `${BASE_URL}/api/auth/sign-in`,
    JSON.stringify({
      email: user.email,
      password: user.password,
    }),
    {
      headers: jsonHeaders(),
      tags: { name: "cold_start_sign_in", phase },
    }
  );

  check(response, {
    "cold start sign-in status is 200": (r) => r.status === 200,
    "cold start sign-in returns access token": (r) =>
      r.status === 200 && typeof r.json("access_token") === "string" && r.json("access_token").length > 0,
  });

  sleep(Number(__ENV.SLEEP_SECONDS || 1));
}
