# tsproxy

用 [tsnet](https://pkg.go.dev/tailscale.com/tsnet) 加入 tailnet,并在本机起一个 SOCKS5 代理。
不需要系统级 Tailscale、不需要 sudo、不改路由;只有经过代理的流量进入 tailnet。

## 构建

    go build -o tsproxy .

## 配置(环境变量)

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `TS_AUTHKEY_FILE` | - | Auth Key 所在文件(推荐,权限 600) |
| `TS_AUTHKEY` | - | 直接给 Auth Key;状态目录已登录过则可不设 |
| `TS_HOSTNAME` | `tsproxy` | 在 tailnet 中显示的设备名,多台机器要各自不同 |
| `TS_STATE_DIR` | `~/.local/share/tsproxy` | 状态目录,固定它才不会每次新增设备 |
| `TS_LISTEN` | `127.0.0.1:1080` | 本地监听地址,端口需 >1024 |

## 使用

    printf '%s' 'tskey-auth-...' > ~/.config/ts-proxy.key && chmod 600 ~/.config/ts-proxy.key
    TS_HOSTNAME=mac-pro TS_AUTHKEY_FILE=~/.config/ts-proxy.key ./tsproxy

用完整地址并走代理访问 tailnet 主机:

    curl -I --proxy socks5h://127.0.0.1:1080 https://<host>.<tailnet>.ts.net
    env HTTPS_PROXY=socks5h://127.0.0.1:1080 rex -S https://<host>.<tailnet>.ts.net --timeout 60s whoami

要点:

- 代理地址用 `socks5h://`,让域名由 tsnet 解析。
- Go 程序只认 `HTTPS_PROXY` / `HTTP_PROXY`,不认 `ALL_PROXY`。
- 不要全局导出 `HTTPS_PROXY`,否则普通 HTTPS 流量也会被送进 tailnet 而失败;只对需要的命令设置。
- **首次连接慢**(可达十几秒):要先登录并建立到对端的路径。rex 默认超时 10 秒,会失败,加 `--timeout 60s`,或先用 curl 预热一次。之后延迟是几十毫秒。
- 每台机器单独运行一份,且用不同的 `TS_HOSTNAME` 和 `TS_STATE_DIR`。
- Auth Key 等同凭证:别入库(`.gitignore` 已忽略 `*.key`),用完可在控制台撤销;已注册的设备不受影响。
