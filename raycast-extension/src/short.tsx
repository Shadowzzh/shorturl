import { useState } from "react";
import { Action, ActionPanel, List, showToast, Toast, Clipboard, getPreferenceValues } from "@raycast/api";
import { validateUrl } from "./validate-url";
import { buildShortenEndpoint } from "./api";

interface Preferences {
  apiBaseUrl?: string;
}

interface ShortenResponse {
  code: number;
  data: {
    id: string;
    short_url: string;
  };
  msg: string;
}

export default function Command() {
  const [searchText, setSearchText] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submittedUrl, setSubmittedUrl] = useState<string>("");
  const [shortUrl, setShortUrl] = useState<string>("");

  const handleSubmit = async () => {
    const validated = validateUrl(searchText);
    if (!validated.ok) {
      await showToast({
        style: Toast.Style.Failure,
        title: "无效的 URL",
        message: validated.reason,
      });
      return;
    }

    setIsSubmitting(true);

    try {
      // 接口地址来自扩展偏好 apiBaseUrl；未设置时回退到默认实例
      const preferences = getPreferenceValues<Preferences>();
      const endpoint = buildShortenEndpoint(preferences.apiBaseUrl);

      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: validated.url }),
      });

      const res = (await response.json()) as ShortenResponse;

      if (res.code === 200) {
        setSubmittedUrl(validated.url);
        setShortUrl(res.data.short_url);
        await showToast({
          style: Toast.Style.Success,
          title: "短链接已生成",
          message: `${res.data.short_url} (已复制到剪贴板)`,
        });
        await Clipboard.copy(res.data.short_url);
      } else {
        throw new Error(res.msg || "生成短链接失败");
      }
    } catch (error) {
      await showToast({
        style: Toast.Style.Failure,
        title: "错误",
        message: error instanceof Error ? error.message : "网络请求失败",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <List
      searchText={searchText}
      onSearchTextChange={setSearchText}
      searchBarPlaceholder="输入要缩短的URL，然后按回车"
      actions={
        <ActionPanel>
          <Action title={isSubmitting ? "生成中…" : "生成短链接"} onAction={handleSubmit} icon="🔗" />
        </ActionPanel>
      }
    >
      <List.EmptyView
        icon="🔗"
        title="输入 URL 生成短链接"
        description={
          searchText ? `输入的URL: ${searchText}\n按回车键或点击按钮生成短链接` : "在搜索框中输入要缩短的URL"
        }
        actions={
          searchText.trim() ? (
            <ActionPanel>
              <Action title={isSubmitting ? "生成中…" : "生成短链接"} onAction={handleSubmit} icon="🔗" />
            </ActionPanel>
          ) : undefined
        }
      />

      {shortUrl && (
        <List.Item
          title={submittedUrl}
          subtitle={shortUrl}
          icon="✅"
          actions={
            <ActionPanel>
              <Action title={isSubmitting ? "生成中…" : "生成短链接"} onAction={handleSubmit} icon="🔗" />
              <Action.CopyToClipboard content={shortUrl} />
              <Action
                title="清除结果"
                onAction={() => {
                  setShortUrl("");
                  setSubmittedUrl("");
                }}
                icon="🗑"
              />
            </ActionPanel>
          }
        />
      )}
    </List>
  );
}
