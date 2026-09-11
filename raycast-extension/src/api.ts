/**
 * 与后端 API 交互的纯函数（不依赖 @raycast/api，因此可脱离 Raycast GUI 单独测试）。
 *
 * 后端接口约定：
 *   POST <base>/shorten   body {"url":"..."}
 *   GET  <base>/<短码>     302 跳转
 *
 * base 来自 Raycast 扩展偏好 `apiBaseUrl`，默认指向本机默认端口（后端默认 8087），
 * 因此在本机跑服务的用户开箱可用；用自建实例（或远程主机）时在扩展设置里改成自己的地址。
 */

export const DEFAULT_API_BASE_URL = "http://localhost:8087/";

/**
 * 规整 base URL：
 *  - 去首尾空白
 *  - 空值回退到 DEFAULT_API_BASE_URL
 *  - 补齐结尾斜杠（后端返回的 short_url 也是 domain + code 直接拼接，同样要求结尾斜杠）
 */
export function normalizeBaseUrl(input?: string): string {
  const raw = (input ?? "").trim();
  const value = raw.length > 0 ? raw : DEFAULT_API_BASE_URL;
  return value.endsWith("/") ? value : `${value}/`;
}

/** 由 base URL 得到创建短链的完整接口地址 */
export function buildShortenEndpoint(input?: string): string {
  return `${normalizeBaseUrl(input)}shorten`;
}
