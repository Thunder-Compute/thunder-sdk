package thunder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CreateServerEnrollmentRequest struct {
	ZoneID           string `json:"zoneId"`
	ExpiresInSeconds int64  `json:"expiresInSeconds,omitempty"`
}

type ServerEnrollmentCommandRequest struct {
	EnrollmentToken string
	IP              string
	Zone            string
	PortRange       string
	ServerName      string
}

// ServerConfigurations contains operator-declared controls that affect how
// the control plane uses a server independently of its reported health.
type ServerConfigurations struct {
	Cordoned bool `json:"cordoned"`
}

type Server struct {
	ServerID        string               `json:"hostId"`
	ZoneID          string               `json:"zoneId"`
	DisplayName     string               `json:"displayName"`
	Hostname        string               `json:"hostname"`
	ThunderdVersion *string              `json:"thunderdVersion,omitempty"`
	GPUType         string               `json:"gpuType"`
	GPUCount        uint32               `json:"gpuCount"`
	Status          string               `json:"status"`
	LastSeenAt      *time.Time           `json:"lastSeenAt,omitempty"`
	Configurations  ServerConfigurations `json:"configurations"`
}

type RevokeServerResponse struct {
	ServerID  string    `json:"hostId"`
	RevokedAt time.Time `json:"revokedAt"`
}

func (c *Client) CreateServerEnrollment(ctx context.Context, req CreateServerEnrollmentRequest) (EnrollmentToken, error) {
	body := struct {
		ZoneID           string `json:"zoneId"`
		Role             string `json:"role"`
		ExpiresInSeconds int64  `json:"expiresInSeconds,omitempty"`
	}{ZoneID: req.ZoneID, Role: RoleServer, ExpiresInSeconds: req.ExpiresInSeconds}
	return c.createEnrollment(ctx, body)
}

func (c *Client) EnrollServer(ctx context.Context, req CreateServerEnrollmentRequest) (EnrollmentToken, error) {
	return c.CreateServerEnrollment(ctx, req)
}

func (c *Client) UnenrollServer(ctx context.Context, enrollmentTokenID string) (DeleteEnrollmentServerResponse, error) {
	return c.DeleteEnrollmentServer(ctx, enrollmentTokenID)
}

func (c *Client) ServerEnrollmentCommand(req ServerEnrollmentCommandRequest) string {
	return serverEnrollmentCommand(c.installURL, c.baseURL, req)
}

func serverEnrollmentCommand(installURL, centralURL string, req ServerEnrollmentCommandRequest) string {
	env := []string{
		"THUNDER_INSTALL_MODE=thunderd",
		"THUNDER_CENTRAL_URL=" + shellQuote(centralURL),
		"THUNDER_ENROLLMENT_TOKEN=" + shellQuote(req.EnrollmentToken),
	}
	if strings.TrimSpace(req.IP) != "" {
		env = append(env, "THUNDERD_IP="+shellQuote(req.IP))
	}
	if strings.TrimSpace(req.Zone) != "" {
		env = append(env, "THUNDER_ZONE="+shellQuote(req.Zone))
	}
	if strings.TrimSpace(req.PortRange) != "" {
		env = append(env, "THUNDERD_PORT_RANGE="+shellQuote(req.PortRange))
	}
	if strings.TrimSpace(req.ServerName) != "" {
		env = append(env, "THUNDERD_NODE_NAME="+shellQuote(req.ServerName))
	}
	return fmt.Sprintf("curl -fsSL %s | sudo %s sh", shellQuote(installURL), strings.Join(env, " "))
}

func (c *Client) ListServers(ctx context.Context, zoneID string) ([]Server, error) {
	var response struct {
		Hosts []Server `json:"hosts"`
	}
	path, err := endpointPath("organizationApi.hosts.list", nil, url.Values{"zoneId": []string{zoneID}})
	if err != nil {
		return nil, err
	}
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}
	return response.Hosts, nil
}

func (c *Client) RevokeServer(ctx context.Context, serverID string) (RevokeServerResponse, error) {
	var response RevokeServerResponse
	path, err := endpointPath("organizationApi.hosts.revoke", map[string]string{"hostId": serverID}, nil)
	if err != nil {
		return RevokeServerResponse{}, err
	}
	if err := c.doJSON(ctx, http.MethodPost, path, nil, &response); err != nil {
		return RevokeServerResponse{}, err
	}
	return response, nil
}
