import * as utils from "/static/utils.js";

const csrfTokenCookie  = "csrf_token";

const doFetch = async (req) => {
  let res, jsonRes;
  try {
    res = await fetch(req);
    jsonRes = await res.json();

  } catch (e) {

    if (e instanceof SyntaxError)
      e = new Error(`status ${res.status}, empty (or invalid) response body`);

    console.error(`api call ${req.method} ${req.url}: unexpected error:`, e);
    throw e;
  }

  if (jsonRes.error) {
    console.error(
      `api call ${req.method} ${req.url}: application error:`,
      res.status,
      jsonRes.error,
    );

    throw jsonRes.error;
  }

  return jsonRes;
}

// may throw
const solvePow = async () => {

  const res = await call('/api/pow/challenge');

  const worker = new Worker('/static/solvePow.js');

  const p = new Promise((resolve, reject) => {
    worker.postMessage({seedHex: res.seed, target: res.target});
    worker.onmessage = resolve;
  });

  const powSol = (await p).data;
  worker.terminate();

  return {seed: res.seed, solution: powSol};
}

const call = async (route, opts = {}) => {
  const {
    method = 'POST',
    body = {},
    requiresPow = false,
  } = opts;

  if (!utils.cookies[csrfTokenCookie]) 
    throw `${csrfTokenCookie} cookie not set, can't make api call`;

  const reqOpts = {
    method,
    headers: {
      "X-CSRF-Token": utils.cookies[csrfTokenCookie],
    },
  };

  if (requiresPow) {
    const {seed, solution} = await solvePow();
    body.powSeed = seed;
    body.powSolution = solution;
  }

  if (Object.keys(body).length > 0) {
    const form = new FormData();
    for (const key in body) form.append(key, body[key]);

    reqOpts.body = form;
  }

  const req = new Request(route, reqOpts);
  return doFetch(req);
}

const ws = async (route, opts = {}) => {
  const {
    requiresPow = false,
    params = {},
  } = opts;

  const docURL = new URL(document.URL);
  const protocol = docURL.protocol == "http:" ? "ws:" : "wss:";

  const fullParams = new URLSearchParams(params);
  const csrfToken = utils.cookies[csrfTokenCookie];

  if (!csrfToken)
    throw `${csrfTokenCookie} cookie not set, can't make api call`;

  fullParams.set("csrfToken", csrfToken);

  if (requiresPow) {
    const {seed, solution} = await solvePow();
    fullParams.set("powSeed", seed);
    fullParams.set("powSolution", solution);
  }

  const rawConn = new WebSocket(`${protocol}//${docURL.host}${route}?${fullParams.toString()}`);

  const conn = {
    next: () => new Promise((resolve, reject) => {
      rawConn.onmessage = (m) => {
        const mj = JSON.parse(m.data);
        resolve(mj);
      };
      rawConn.onerror = reject;
      rawConn.onclose = reject;
    }),

    close: rawConn.close,
  };

  return new Promise((resolve, reject) => {
    rawConn.onopen = () => resolve(conn);
    rawConn.onerror = reject;
  });
}

export {
  call,
  ws
}
