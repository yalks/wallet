package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// RequestContext holds the request context information
type RequestContext struct {
	Source    string            `json:"source"`
	IP        string            `json:"ip"`
	UserAgent string            `json:"user_agent"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// RequestContextKey is the key for storing request context in context.Context
type RequestContextKey struct{}

// ExtractRequestContext extracts request context information from the context
func ExtractRequestContext(ctx context.Context) *RequestContext {
	// Check if we already have a RequestContext in the context
	if reqCtx, ok := ctx.Value(RequestContextKey{}).(*RequestContext); ok {
		return reqCtx
	}

	// Try to extract from ghttp.Request if available
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		return &RequestContext{
			Source:    extractSource(r),
			IP:        "", // IP must be explicitly passed, not extracted from headers
			UserAgent: r.Header.Get("User-Agent"),
			Metadata:  extractMetadata(r),
		}
	}

	// Try to extract from g.Map in context
	if ctxData := g.RequestFromCtx(ctx); ctxData != nil {
		if mapData := ctxData.Get("request_context"); mapData != nil {
			if gMap := mapData.Map(); gMap != nil {
				return extractFromGMap(gMap)
			}
		}
	}

	// Return default context
	return &RequestContext{
		Source:    "unknown",
		IP:        "",
		UserAgent: "",
		Metadata:  make(map[string]string),
	}
}

// WithRequestContext adds request context to context.Context
func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, RequestContextKey{}, reqCtx)
}

// extractSource determines the request source
func extractSource(r *ghttp.Request) string {
	// Check custom header first
	if source := r.Header.Get("X-Request-Source"); source != "" {
		return source
	}

	// Check API key header to determine if it's API
	if r.Header.Get("X-API-Key") != "" {
		return "api"
	}

	// Check referer
	referer := r.Header.Get("Referer")
	if strings.Contains(referer, "telegram") {
		return "telegram"
	}
	if strings.Contains(referer, "admin") {
		return "admin"
	}

	// Check user agent
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	if strings.Contains(userAgent, "telegram") {
		return "telegram"
	}

	// Default to web
	return "web"
}


// extractMetadata extracts additional metadata from the request
func extractMetadata(r *ghttp.Request) map[string]string {
	metadata := make(map[string]string)

	// Add request ID if available
	if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
		metadata["request_id"] = reqID
	}

	// Add session ID if available
	if sessionID := r.Header.Get("X-Session-ID"); sessionID != "" {
		metadata["session_id"] = sessionID
	}

	// Add device type if available
	if deviceType := r.Header.Get("X-Device-Type"); deviceType != "" {
		metadata["device_type"] = deviceType
	}

	// Add platform if available
	if platform := r.Header.Get("X-Platform"); platform != "" {
		metadata["platform"] = platform
	}

	return metadata
}

// extractFromGMap extracts RequestContext from a g.Map
func extractFromGMap(data g.Map) *RequestContext {
	if data == nil {
		return &RequestContext{
			Source:    "unknown",
			IP:        "",
			UserAgent: "",
			Metadata:  make(map[string]string),
		}
	}

	reqCtx := &RequestContext{
		Source:    gconv.String(data["source"]),
		IP:        gconv.String(data["ip"]),
		UserAgent: gconv.String(data["user_agent"]),
		Metadata:  make(map[string]string),
	}

	// If values are empty, use defaults
	if reqCtx.Source == "" {
		reqCtx.Source = "unknown"
	}

	// Extract metadata if present
	if metadataRaw := data["metadata"]; metadataRaw != nil {
		switch metadata := metadataRaw.(type) {
		case map[string]string:
			reqCtx.Metadata = metadata
		case map[string]interface{}:
			for k, v := range metadata {
				reqCtx.Metadata[k] = gconv.String(v)
			}
		}
	}

	return reqCtx
}


// ToJSON converts RequestContext to JSON string
func (rc *RequestContext) ToJSON() string {
	data, _ := json.Marshal(rc)
	return string(data)
}

// MetadataToJSON converts metadata map to JSON string
func (rc *RequestContext) MetadataToJSON() string {
	if len(rc.Metadata) == 0 {
		return "{}"
	}
	data, _ := json.Marshal(rc.Metadata)
	return string(data)
}

// CreateRequestContext creates a RequestContext with the given parameters
func CreateRequestContext(source, ip, userAgent string, metadata map[string]string) *RequestContext {
	if metadata == nil {
		metadata = make(map[string]string)
	}
	return &RequestContext{
		Source:    source,
		IP:        ip,
		UserAgent: userAgent,
		Metadata:  metadata,
	}
}