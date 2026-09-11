/**
 * api.ts 的用例测试：
 *   tsc 编译 src/api.ts 后，node tests/api.test.js <编译产物目录>
 */
const path = require("path");

const outDir = process.argv[2];
const { normalizeBaseUrl, buildShortenEndpoint, DEFAULT_API_BASE_URL } = require(path.join(outDir, "api.js"));

const cases = [
  // [输入, 期望的 normalizeBaseUrl 结果, 说明]
  [undefined, DEFAULT_API_BASE_URL, "未设置时回退默认值"],
  ["", DEFAULT_API_BASE_URL, "空串回退默认值"],
  ["   ", DEFAULT_API_BASE_URL, "纯空格回退默认值"],
  ["https://s.example.com", "https://s.example.com/", "无结尾斜杠要补上"],
  ["https://s.example.com/", "https://s.example.com/", "已有结尾斜杠保持不变"],
  ["  https://s.example.com  ", "https://s.example.com/", "去首尾空白并补斜杠"],
  ["http://localhost:8087", "http://localhost:8087/", "本地地址"],
  ["http://192.168.1.10:8087/", "http://192.168.1.10:8087/", "内网 IP 带端口"],
  ["https://s.example.com/base/path", "https://s.example.com/base/path/", "带路径前缀（反代场景）"],
];

let pass = 0;
let fail = 0;

for (const [input, expected, note] of cases) {
  const actual = normalizeBaseUrl(input);
  const ok = actual === expected;
  if (ok) {
    pass += 1;
  } else {
    fail += 1;
  }
  console.log(
    `${ok ? "PASS" : "FAIL"}  normalizeBaseUrl(${JSON.stringify(input)}) -> ${actual}${ok ? "" : `  (期望 ${expected})`}   (${note})`,
  );
}

// buildShortenEndpoint 必须始终以 /shorten 结尾（后端路由就是 /shorten）
const endpointCases = [
  [undefined, `${DEFAULT_API_BASE_URL}shorten`, "默认值 + /shorten"],
  ["https://s.example.com", "https://s.example.com/shorten", "无斜杠补斜杠再拼"],
  ["https://s.example.com/", "https://s.example.com/shorten", "有斜杠不重复"],
  ["https://s.example.com/base", "https://s.example.com/base/shorten", "路径前缀下也要正确拼接"],
];

for (const [input, expected, note] of endpointCases) {
  const actual = buildShortenEndpoint(input);
  const ok = actual === expected;
  if (ok) {
    pass += 1;
  } else {
    fail += 1;
  }
  console.log(
    `${ok ? "PASS" : "FAIL"}  buildShortenEndpoint(${JSON.stringify(input)}) -> ${actual}${ok ? "" : `  (期望 ${expected})`}   (${note})`,
  );
}

console.log("");
console.log(`合计 ${cases.length + endpointCases.length} 例：通过 ${pass}，失败 ${fail}`);
if (fail > 0) {
  process.exit(1);
}
