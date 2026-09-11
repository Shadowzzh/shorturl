/**
 * URL 校验与规整。
 *
 * 背景：扩展原先只判断"输入非空"，导致 `http:///Volumes/bootfs`、`http://qwq`
 * 这类无意义内容也会被提交并落库。这里做最小充分校验：
 *  - 语法必须能被 URL 解析（http/https 且有主机名）
 *  - 只接受 http / https
 *  - 主机名必须是像域名的形式（含点）、localhost、或 IP，挡掉 `qwq` 这类单词
 *
 * 刻意不改变用户输入的内容：只补协议头，不重写路径与查询串。
 */

export type UrlValidationResult = { ok: true; url: string } | { ok: false; reason: string };

const SCHEME_RE = /^[a-zA-Z][a-zA-Z0-9+.-]*:\/\//;
const IPV4_RE = /^\d{1,3}(\.\d{1,3}){3}$/;

export function validateUrl(input: string): UrlValidationResult {
  const raw = input.trim();
  if (!raw) {
    return { ok: false, reason: "请输入 URL" };
  }

  // 没写协议时补 https://（例如粘贴 example.com/path）
  const candidate = SCHEME_RE.test(raw) ? raw : `https://${raw}`;

  let parsed: URL;
  try {
    parsed = new URL(candidate);
  } catch {
    return { ok: false, reason: `不是合法的 URL：${raw}` };
  }

  if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
    return { ok: false, reason: `只支持 http/https，收到 ${parsed.protocol.replace(":", "")}` };
  }

  const host = parsed.hostname;
  const isIpv6 = host.includes(":");
  const isIpv4 = IPV4_RE.test(host);
  if (!host.includes(".") && host !== "localhost" && !isIpv4 && !isIpv6) {
    return { ok: false, reason: `域名不完整：${host}（例：example.com）` };
  }
  if (isIpv4 && host.split(".").some((part) => Number(part) > 255)) {
    return { ok: false, reason: `IP 地址无效：${host}` };
  }

  return { ok: true, url: candidate };
}
