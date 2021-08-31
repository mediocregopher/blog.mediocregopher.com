import * as utils from "/assets/utils.js";

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

  const worker = new Worker('/assets/solvePow.js');

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

  if (!utils.cookies["csrf_token"]) 
    throw "csrf_token cookie not set, can't make api call";

  const reqOpts = {
    method,
    headers: {
      "X-CSRF-Token": utils.cookies["csrf_token"],
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

export {
  call,
}
