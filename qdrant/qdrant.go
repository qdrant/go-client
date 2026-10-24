package qdrant

import (
	"context"
)

// Check liveness of the service.
func (c *Client) HealthCheck(ctx context.Context) (*HealthCheckReply, error) {
	resp, err := c.GetQdrantClient().HealthCheck(ctx, &HealthCheckRequest{})
	if err != nil {
		return nil, newQdrantErr(err, "HealthCheck")
	}
	return resp, nil
}
