// tsproxy 用 tsnet 加入 tailnet,并在本机提供 SOCKS5 代理,
// 让不在 tailnet 内的机器无需系统级 Tailscale 即可访问 tailnet 主机。
package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	socks5 "github.com/armon/go-socks5"
	"tailscale.com/tsnet"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// authKey 优先读 TS_AUTHKEY_FILE 指向的文件,其次读 TS_AUTHKEY。
// 状态目录里已有登录状态时可以都不设。
func authKey() string {
	if f := os.Getenv("TS_AUTHKEY_FILE"); f != "" {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read TS_AUTHKEY_FILE: %v", err)
		}
		return strings.TrimSpace(string(b))
	}
	return os.Getenv("TS_AUTHKEY")
}

func main() {
	home, _ := os.UserHomeDir()
	srv := &tsnet.Server{
		Hostname: env("TS_HOSTNAME", "tsproxy"),
		AuthKey:  authKey(),
		Dir:      env("TS_STATE_DIR", home+"/.local/share/tsproxy"),
	}
	defer srv.Close()

	s, err := socks5.New(&socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return srv.Dial(ctx, network, addr)
		},
		Resolver: passthrough{},
	})
	if err != nil {
		log.Fatal(err)
	}
	listen := env("TS_LISTEN", "127.0.0.1:1080")
	log.Printf("socks5 listening on %s", listen)
	log.Fatal(s.ListenAndServe("tcp", listen))
}

// passthrough 不在本地解析域名,交给 tsnet(MagicDNS)解析。
type passthrough struct{}

func (passthrough) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	return ctx, nil, nil
}
