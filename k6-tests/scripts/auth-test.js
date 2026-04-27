import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL, buildOptions, jsonHeaders } from "./config.js";

const data = JSON.parse(open("../data/users.json"));
const users = data.auth_users || [];

export const options = buildOptions("auth_sign_in", {
  "http_req_duration{name:auth_sign_in}": ["p(95)<2000", "p(99)<4000"],
});

export default function () {
  const user = users[(__VU - 1) % users.length];

  const response = http.post(
    `${BASE_URL}/api/auth/sign-in`,
    JSON.stringify({
      email: user.email,
      password: user.password,
    }),
    {
      headers: jsonHeaders(),
      tags: { name: "auth_sign_in" },
    }
  );

  check(response, {
    "sign-in status is 200": (r) => r.status === 200,
    "sign-in returns access token": (r) => {
      const body = r.json();
      return body && typeof body.access_token === "string" && body.access_token.length > 0;
    },
  });

  sleep(Number(__ENV.SLEEP_SECONDS || 1));
}
