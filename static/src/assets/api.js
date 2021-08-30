
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

  const res = await call('GET', '/api/pow/challenge');

  const worker = new Worker('/assets/solvePow.js');

  const p = new Promise((resolve, reject) => {
    worker.postMessage({seedHex: res.seed, target: res.target});
    worker.onmessage = resolve;
  });

  const powSol = (await p).data;
  worker.terminate();

  return {seed: res.seed, solution: powSol};
}

const call = async (method, route, opts = {}) => {
  const { body = {}, requiresPow = false } = opts;

  const reqOpts = { method };

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
