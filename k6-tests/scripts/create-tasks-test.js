import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL, buildOptions, jsonHeaders } from "./config.js";

const data = JSON.parse(open("../data/tasks.json"));
const creators = data.task_creators || [];

export const options = buildOptions("task_create", {
  "http_req_duration{name:task_create}": ["p(95)<2500", "p(99)<5000"],
});

export function setup() {
  return creators.map((creator) => {
    const signInResponse = http.post(
      `${BASE_URL}/api/auth/sign-in`,
      JSON.stringify({
        email: creator.email,
        password: creator.password,
      }),
      {
        headers: jsonHeaders(),
        tags: { name: "task_creator_sign_in" },
      }
    );

    check(signInResponse, {
      "task creator sign-in status is 200": (r) => r.status === 200,
    });

    return {
      ...creator,
      token: signInResponse.status === 200 ? signInResponse.json("access_token") : null,
    };
  });
}

export default function (sessionData) {
  const creator = sessionData[(__VU - 1) % sessionData.length];
  const uniqueSuffix = `${__VU}-${__ITER}-${Date.now()}`;

  const response = http.post(
    `${BASE_URL}/api/tasks`,
    JSON.stringify({
      assignment_id: creator.assignment_id,
      title: `${creator.title || "Carga"} ${uniqueSuffix}`,
      description: `${creator.description || "Prueba de carga"} ${uniqueSuffix}`,
      status: creator.status || "abierto",
      spent_hours: creator.spent_hours || 1,
      observations: creator.observations || "Generado por k6",
      week_start_date: creator.week_start_date,
    }),
    {
      headers: jsonHeaders(creator.token),
      tags: { name: "task_create" },
    }
  );

  check(response, {
    "task create status is 201": (r) => r.status === 201,
  });

  sleep(Number(__ENV.SLEEP_SECONDS || 1));
}
