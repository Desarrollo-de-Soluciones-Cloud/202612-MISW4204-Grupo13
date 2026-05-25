import http from "k6/http";
import { check, group, sleep } from "k6";
import { BASE_URL, buildOptions, jsonHeaders } from "./config.js";

const data = JSON.parse(open("../data/users.json"));
const professorView = data.professor_view || {};

export const options = buildOptions("professor_view", {
  "http_req_duration{name:professor_view_workspaces}": ["p(95)<2500"],
  "http_req_duration{name:professor_view_weeks}": ["p(95)<2500"],
  "http_req_duration{name:professor_view_tasks}": ["p(95)<3000"],
  "http_req_duration{name:professor_view_assignments}": ["p(95)<3000"],
});

export function setup() {
  const signInResponse = http.post(
    `${BASE_URL}/api/auth/sign-in`,
    JSON.stringify({
      email: professorView.email,
      password: professorView.password,
    }),
    {
      headers: jsonHeaders(),
      tags: { name: "professor_view_sign_in" },
    }
  );

  check(signInResponse, {
    "professor sign-in status is 200": (r) => r.status === 200,
  });

  return {
    token: signInResponse.status === 200 ? signInResponse.json("access_token") : null,
    periodId: professorView.period_id,
    assignmentUserIds: professorView.assignment_user_ids || [],
  };
}

export default function (session) {
  group("professor_view", function () {
    const commonParams = {
      headers: jsonHeaders(session.token),
    };

    const workspacesResponse = http.get(
      `${BASE_URL}/api/workspaces?period_id=${session.periodId}`,
      {
        ...commonParams,
        tags: { name: "professor_view_workspaces" },
      }
    );

    check(workspacesResponse, {
      "workspaces status is 200": (r) => r.status === 200,
    });

    const weeksResponse = http.get(
      `${BASE_URL}/api/weeks/periods/${session.periodId}`,
      {
        ...commonParams,
        tags: { name: "professor_view_weeks" },
      }
    );

    check(weeksResponse, {
      "weeks status is 200": (r) => r.status === 200,
    });

    const tasksResponse = http.get(`${BASE_URL}/api/tasks`, {
      ...commonParams,
      tags: { name: "professor_view_tasks" },
    });

    check(tasksResponse, {
      "tasks status is 200": (r) => r.status === 200,
    });

    for (const userId of session.assignmentUserIds) {
      const assignmentsResponse = http.get(
        `${BASE_URL}/api/assignments?user_id=${userId}`,
        {
          ...commonParams,
          tags: { name: "professor_view_assignments" },
        }
      );

      check(assignmentsResponse, {
        "assignments status is 200": (r) => r.status === 200,
      });
    }
  });

  sleep(Number(__ENV.SLEEP_SECONDS || 1));
}
