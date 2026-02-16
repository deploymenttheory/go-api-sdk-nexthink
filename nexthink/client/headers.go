package client

import "resty.dev/v3"

// applyHeaders applies headers to a request with proper precedence:
// 1. Global headers are applied first
// 2. Per-request headers override global headers with the same key
func (c *Client) applyHeaders(req *resty.Request, requestHeaders map[string]string) {
	// Apply global headers first
	for k, v := range c.globalHeaders {
		if v != "" {
			req.SetHeader(k, v)
		}
	}

	// Apply per-request headers (overrides global headers)
	for k, v := range requestHeaders {
		if v != "" {
			req.SetHeader(k, v)
		}
	}
}
