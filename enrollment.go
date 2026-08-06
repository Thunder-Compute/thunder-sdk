package thunder

import (
	"context"
	"net/http"
	"time"
)

type EnrollmentToken struct {
	EnrollmentTokenID string     `json:"enrollmentTokenId"`
	EnrollmentToken   string     `json:"enrollmentToken"`
	OrgID             string     `json:"orgId"`
	ZoneID            string     `json:"zoneId,omitempty"`
	Role              string     `json:"role"`
	GPUType           string     `json:"gpuType,omitempty"`
	GPUCount          uint32     `json:"gpuCount,omitempty"`
	ExpiresAt         *time.Time `json:"expiresAt,omitempty"`
}

type DeleteEnrollmentServerResponse struct {
	EnrollmentTokenID string    `json:"enrollmentTokenId"`
	Role              string    `json:"role"`
	ClientID          string    `json:"clientId,omitempty"`
	ServerID          string    `json:"hostId,omitempty"`
	ServerDeleted     bool      `json:"nodeDeleted"`
	DeletedAt         time.Time `json:"deletedAt"`
}

func (c *Client) createEnrollment(ctx context.Context, body any) (EnrollmentToken, error) {
	var response EnrollmentToken
	path, err := endpointPath("organizationApi.enrollmentTokens.create", nil, nil)
	if err != nil {
		return EnrollmentToken{}, err
	}
	if err := c.doJSON(ctx, http.MethodPost, path, body, &response); err != nil {
		return EnrollmentToken{}, err
	}
	return response, nil
}

func (c *Client) DeleteEnrollmentServer(ctx context.Context, enrollmentTokenID string) (DeleteEnrollmentServerResponse, error) {
	var response DeleteEnrollmentServerResponse
	path, err := endpointPath("organizationApi.enrollmentTokens.deleteServer", map[string]string{"enrollmentTokenId": enrollmentTokenID}, nil)
	if err != nil {
		return DeleteEnrollmentServerResponse{}, err
	}
	if err := c.doJSON(ctx, http.MethodDelete, path, nil, &response); err != nil {
		return DeleteEnrollmentServerResponse{}, err
	}
	return response, nil
}
