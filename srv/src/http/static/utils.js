const cookies = {};
const cookieKVs = document.cookie
  .split(';')
  .map(cookie => cookie.trim().split('=', 2));

for (const i in cookieKVs) {
  cookies[cookieKVs[i][0]] = cookieKVs[i][1];
}

export {
  cookies,
}
