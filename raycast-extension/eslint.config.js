const { defineConfig } = require("eslint/config");
const raycastConfig = require("@raycast/eslint-config");

module.exports = defineConfig([
  ...raycastConfig,
  {
    // tests/ 下是 CommonJS 用例脚本：先用 tsc 把 src/validate-url.ts 编译成 CJS，
    // 再用 node 直接跑断言（不需要 Raycast GUI），因此这里刻意使用 require()。
    files: ["tests/**/*.js"],
    rules: {
      "@typescript-eslint/no-require-imports": "off",
    },
  },
]);
