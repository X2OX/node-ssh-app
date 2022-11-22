package nsa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	cfToken      = os.Getenv("TOKEN")
	cfIdentifier = os.Getenv("IDENTIFIER")
	domain       = os.Getenv("DOMAIN")
)

func GetTunnel() []Tunnel {
	var tnl TunnelsResp

	if err := makeRequest("GET",
		"/accounts/"+cfIdentifier+"/cfd_tunnel", nil, &tnl,
	); err != nil || !tnl.Success || len(tnl.Result) == 0 {
		return nil
	}
	return tnl.Result
}

func NewTunnel(name string) *Tunnel {
	var (
		tnl TunnelResp
		tcr TunnelConfigResp
	)

	if err := makeRequest("POST",
		"/accounts/"+cfIdentifier+"/cfd_tunnel",
		map[string]string{
			"name":          name,
			"tunnel_secret": base64.RawStdEncoding.EncodeToString([]byte(uuid.New().String())),
			"config_src":    "cloudflare",
		},
		&tnl,
	); err != nil || !tnl.Success {
		return nil
	}

	if err := makeRequest("PUT",
		"/accounts/"+cfIdentifier+"/cfd_tunnel/"+tnl.Result.ID+"/configurations",
		map[string]any{
			"config": ConfigInfo{
				Ingress: []struct {
					Service  string `json:"service"`
					Hostname string `json:"hostname,omitempty"`
				}{
					{Service: "ssh://127.0.0.1", Hostname: fmt.Sprintf("%s.%s", name, domain)},
					{Service: "http_status:404"},
				},
				WarpRouting: struct {
					Enabled bool `json:"enabled"`
				}{false},
			},
		},
		&tcr,
	); err != nil || !tcr.Success {
		return nil
	}

	return &tnl.Result
}

const cfAPIAddr = "https://api.cloudflare.com/client/v4"

func makeRequest(method, url string, body, resp any) error {
	var buf *bytes.Buffer = nil
	if body != nil {
		buf = new(bytes.Buffer)

		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
	}
	var ir io.Reader = nil
	if buf != nil {
		ir = buf
	}

	req, err := http.NewRequest(method, cfAPIAddr+url, ir)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfToken)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if buf == nil {
		buf = new(bytes.Buffer)
	} else {
		buf.Reset()
	}

	_, _ = io.Copy(buf, res.Body)

	switch resp := resp.(type) {
	case *string:
		*resp = buf.String()
	case []byte:
		resp = buf.Bytes()
	default:
		err = json.Unmarshal(buf.Bytes(), resp)
	}

	return err
}

type Response struct {
	Success  bool           `json:"success"`
	Errors   []ResponseInfo `json:"errors"`
	Messages []ResponseInfo `json:"messages"`
}

// ResponseInfo contains a code and message returned by the API as errors or
// informational messages inside the response.
type ResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TunnelResp struct {
	Response
	Result Tunnel `json:"result"`
}

// TunnelConnection represents the connections associated with a tunnel.
type TunnelConnection struct {
	ColoName           string `json:"colo_name"`
	ID                 string `json:"id"`
	IsPendingReconnect bool   `json:"is_pending_reconnect"`
	ClientID           string `json:"client_id"`
	ClientVersion      string `json:"client_version"`
	OpenedAt           string `json:"opened_at"`
	OriginIP           string `json:"origin_ip"`
}

type Tunnel struct {
	ID              string             `json:"id"`
	AccountTag      string             `json:"account_tag"`
	CreatedAt       time.Time          `json:"created_at"`
	Name            string             `json:"name"`
	Connections     []TunnelConnection `json:"connections"`
	ConnsActiveAt   *time.Time         `json:"conns_active_at,omitempty"`
	ConnInactiveAt  *time.Time         `json:"conns_inactive_at,omitempty"`
	TunType         string             `json:"tun_type"`
	Status          string             `json:"status"`
	RemoteConfig    bool               `json:"remote_config"`
	CredentialsFile TunnelCredentials  `json:"credentials_file"`
	Token           string             `json:"token"`
}

type TunnelCredentials struct {
	AccountTag   string `json:"AccountTag"`
	TunnelID     string `json:"TunnelID"`
	TunnelName   string `json:"TunnelName"`
	TunnelSecret string `json:"TunnelSecret"`
}

type TunnelsResp struct {
	Response
	Result []Tunnel `json:"result"`
}

type TunnelConfigResp struct {
	Response
	Result TunnelConfig `json:"result"`
}

type TunnelConfig struct {
	TunnelID  string     `json:"tunnel_id"`
	Version   int        `json:"version"`
	Config    ConfigInfo `json:"config"`
	Source    string     `json:"source"`
	CreatedAt time.Time  `json:"created_at"`
}

type ConfigInfo struct {
	Ingress []struct {
		Service  string `json:"service"`
		Hostname string `json:"hostname,omitempty"`
	} `json:"ingress"`
	WarpRouting struct {
		Enabled bool `json:"enabled"`
	} `json:"warp-routing"`
}
