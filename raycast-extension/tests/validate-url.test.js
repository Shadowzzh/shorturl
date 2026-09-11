/**
 * validate-url 的用例测试。
 * 用法：tsc 编译 validate-url.ts 后，node validate-url.test.js <编译产物目录>
 */
const path = require("path");

const outDir = process.argv[2];
const { validateUrl } = require(path.join(outDir, "validate-url.js"));

const cases = [
  // [输入, 是否应通过, 说明]
  ["https://example.com/path?q=1#frag", true, "标准 https"],
  ["http://example.com", true, "标准 http"],
  ["example.com", true, "无协议时补 https"],
  ["example.com/a/b?x=1", true, "无协议的路径+查询"],
  ["sub.domain.co.uk", true, "多级域名"],
  ["https://localhost:3000/x", true, "localhost 带端口"],
  ["http://192.168.1.10:8080/a", true, "内网 IP"],
  ["https://user:pass@example.com/x", true, "带 Basic Auth"],
  ["  https://example.com  ", true, "首尾空格应被裁剪"],
  ["http://qwq", false, "单词主机名（历史脏数据）"],
  ["http:///Volumes/bootfs", false, "空主机（历史脏数据）"],
  ["/Volumes/bootfs", false, "纯本地路径"],
  ["", false, "空输入"],
  ["   ", false, "纯空格"],
  ["ftp://example.com", false, "非 http/https 协议"],
  ["javascript:alert(1)", false, "伪协议"],
  ["http://999.1.1.1", false, "非法 IPv4"],
  ["https://", false, "只有协议头"],
  ["http://", false, "只有协议头 http"],
];

let pass = 0;
let fail = 0;
for (const [input, shouldPass, note] of cases) {
  const result = validateUrl(input);
  const ok = result.ok === shouldPass;
  if (ok) {
    pass += 1;
  } else {
    fail += 1;
  }
  const shown = JSON.stringify(input);
  const detail = result.ok ? `-> ${result.url}` : `-> 拒绝：${result.reason}`;
  console.log(`${ok ? "PASS" : "FAIL"}  ${shown.padEnd(34)} ${String(shouldPass).padEnd(5)} ${detail}   (${note})`);
}

console.log("");
console.log(`合计 ${cases.length} 例：通过 ${pass}，失败 ${fail}`);
if (fail > 0) {
  process.exit(1);
}
