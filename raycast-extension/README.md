# Short URL Generator（Raycast 扩展）

本仓库的 macOS 客户端。在 Raycast 里输入 URL 生成短链，结果自动复制到剪贴板。

数据流：Raycast → `POST <你的实例>/shorten` → 本仓库的后端服务（见上级目录）。

## 使用

1. 打开 Raycast，输入 `Short URL Generator`
2. 在搜索框粘贴或输入要缩短的 URL
3. 回车 → Toast 提示「短链接已生成」，短链已在剪贴板

输入要求：可解析的 `http` / `https` 地址。没写协议会自动补 `https://`（可直接粘贴 `example.com/path`）。
非法输入（如 `qwq`）会被拦下并提示原因，不会提交。

## 指向你的实例

扩展将所有请求发到偏好项 `apiBaseUrl` 指定的地址。默认值是 `http://localhost:8087/`
（即后端默认端口），所以在本机跑着后端时开箱可用。

用自建实例或远程主机时，在 Raycast 里对该扩展执行 **Configure Extension**，把
`Short URL Service` 改成你的地址即可，例如：

```text
http://192.168.1.10:8087/      内网另一台机器
https://s.example.com/         经反向代理对外提供服务
```

地址会被自动规整（去空格、补结尾斜杠），因此写 `http://host:8087` 或
`https://s.example.com/base` 都能正确拼出 `<base>/shorten`。

## 开发

包管理器用 **pnpm**（仓库内是 `pnpm-lock.yaml`）。

```bash
pnpm install
pnpm run dev      # ray develop：注册到 Raycast，改动热重载
pnpm run build    # ray build
pnpm run lint     # eslint
```

> `pnpm run dev` 只在进程存活期间把扩展注册到 Raycast：关掉终端、注销或重启后扩展就从 Raycast 里消失。
> 需要长期可用时，用 launchd 常驻这个进程（在 macOS 上是常见做法）。

## 测试

`src/validate-url.ts`（输入校验）与 `src/api.ts`（实例地址拼接）都是不依赖 React 的纯函数，
可以脱离 Raycast 单独测：

```bash
# 从仓库根目录
make ext-test
```

共 32 条用例：合法/非法 URL、无协议补全、IP、localhost，以及 base URL 规整与 `/shorten` 拼接。

## 目录

```
src/short.tsx          命令入口与 UI
src/validate-url.ts    输入校验（纯函数）
src/api.ts             实例地址规整与请求地址拼接（纯函数）
tests/                 上述两个纯函数的用例（CommonJS，被 eslint 单独放行 require）
assets/                图标
```

## 与后端的关系

后端在本仓库根目录（Go + Gin）。扩展依赖它的两个约定：

| 约定 | 位置 |
| --- | --- |
| `POST /shorten`，body `{"url":"..."}` | `handlers/shorten.go` |
| 响应 `{"code":200,"data":{"id","short_url"}}` | 同上；`short_url` = `server.domain` + 短码 |

因此后端 `server.domain` 配置改动会影响扩展返回的短链前缀；域名末尾斜杠不能丢。
