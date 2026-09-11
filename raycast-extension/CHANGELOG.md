# shorturl Changelog

## [2026-09-11] 首个版本

- 在 Raycast 里输入地址即创建短链，结果自动复制到剪贴板
- 提交前校验输入：仅接受可解析的 `http`/`https` 地址；无协议时自动补 `https://`
- 实例地址由偏好项 `apiBaseUrl` 配置（默认指向本机 `http://localhost:8087/`）
