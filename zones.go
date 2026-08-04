package thunder

import (
	"context"
	"net/http"
	"time"
)

type Zone struct {
	ZoneID      string    `json:"zoneId"`
	DisplayName string    `json:"displayName"`
	NodeCount   int       `json:"nodeCount"`
	ClientCount int       `json:"clientCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateZoneRequest struct {
	DisplayName string `json:"displayName,omitempty"`
}

type CreateZoneResponse struct {
	ZoneID      string `json:"zoneId"`
	OrgID       string `json:"orgId,omitempty"`
	DisplayName string `json:"displayName"`
}

func (c *Client) ListZones(ctx context.Context) ([]Zone, error) {
	var response struct {
		Zones []Zone `json:"zones"`
	}
	path, err := endpointPath("organizationApi.zones.list", nil, nil)
	if err != nil {
		return nil, err
	}
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}
	return response.Zones, nil
}

func (c *Client) CreateZone(ctx context.Context, req CreateZoneRequest) (CreateZoneResponse, error) {
	var response CreateZoneResponse
	path, err := endpointPath("organizationApi.zones.create", nil, nil)
	if err != nil {
		return CreateZoneResponse{}, err
	}
	if err := c.doJSON(ctx, http.MethodPost, path, req, &response); err != nil {
		return CreateZoneResponse{}, err
	}
	return response, nil
}

func (c *Client) DeleteZone(ctx context.Context, zoneID string) error {
	path, err := endpointPath("organizationApi.zones.delete", map[string]string{"zoneId": zoneID}, nil)
	if err != nil {
		return err
	}
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}
