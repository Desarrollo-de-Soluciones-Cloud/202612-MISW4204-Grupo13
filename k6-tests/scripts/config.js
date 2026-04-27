const DEFAULT_HEADERS = {
  "Content-Type": "application/json",
};

export const BASE_URL = __ENV.BASE_URL || "http://localhost";
export const LOAD_PATTERN = __ENV.LOAD_PATTERN || "constant";
export const VUS = Number(__ENV.VUS || 10);
export const DURATION = __ENV.DURATION || "1m";

function constantScenario(name) {
  return {
    [name]: {
      executor: "constant-vus",
      vus: VUS,
      duration: DURATION,
    },
  };
}

function rampScenario(name) {
  return {
    [name]: {
      executor: "ramping-vus",
      startVUs: Number(__ENV.START_VUS || 0),
      stages: [
        {
          duration: __ENV.STAGE_ONE_DURATION || "1m",
          target: Number(__ENV.STAGE_ONE_TARGET || 10),
        },
        {
          duration: __ENV.STAGE_TWO_DURATION || "2m",
          target: Number(__ENV.STAGE_TWO_TARGET || 20),
        },
        {
          duration: __ENV.STAGE_THREE_DURATION || "30s",
          target: Number(__ENV.STAGE_THREE_TARGET || 0),
        },
      ],
    },
  };
}

function spikeScenario(name) {
  return {
    [name]: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        {
          duration: __ENV.SPIKE_UP_DURATION || "10s",
          target: Number(__ENV.SPIKE_TARGET || VUS),
        },
        {
          duration: __ENV.SPIKE_HOLD_DURATION || DURATION,
          target: Number(__ENV.SPIKE_TARGET || VUS),
        },
        {
          duration: __ENV.SPIKE_DOWN_DURATION || "10s",
          target: 0,
        },
      ],
    },
  };
}

export function buildOptions(name, extraThresholds = {}) {
  let scenarios;

  if (LOAD_PATTERN === "ramp") {
    scenarios = rampScenario(name);
  } else if (LOAD_PATTERN === "spike") {
    scenarios = spikeScenario(name);
  } else {
    scenarios = constantScenario(name);
  }

  return {
    scenarios,
    thresholds: {
      http_req_failed: ["rate<0.01"],
      http_req_duration: ["p(95)<2000", "p(99)<4000"],
      ...extraThresholds,
    },
    summaryTrendStats: ["avg", "min", "med", "max", "p(95)", "p(99)"],
  };
}

export function jsonHeaders(token) {
  if (!token) {
    return DEFAULT_HEADERS;
  }

  return {
    ...DEFAULT_HEADERS,
    Authorization: `Bearer ${token}`,
  };
}
