package thunder

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

//go:embed endpoints.json
var endpointContract []byte

var (
	endpointContractOnce sync.Once
	endpointPaths        map[string]string
	endpointContractErr  error
)

func endpointPath(name string, params map[string]string, query url.Values) (string, error) {
	path, err := endpointTemplate(name)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for i := 0; i < len(path); {
		if path[i] != '{' {
			out.WriteByte(path[i])
			i++
			continue
		}
		end := strings.IndexByte(path[i+1:], '}')
		if end < 0 {
			return "", fmt.Errorf("endpoint %s has unterminated path parameter", name)
		}
		param := path[i+1 : i+1+end]
		value, ok := params[param]
		if !ok || strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("endpoint %s requires path parameter %s", name, param)
		}
		out.WriteString(url.PathEscape(value))
		i += end + 2
	}
	if len(query) > 0 {
		out.WriteByte('?')
		out.WriteString(query.Encode())
	}
	return out.String(), nil
}

func endpointTemplate(name string) (string, error) {
	endpointContractOnce.Do(func() {
		var raw map[string]any
		if err := json.Unmarshal(endpointContract, &raw); err != nil {
			endpointContractErr = err
			return
		}
		endpointPaths = make(map[string]string)
		flattenEndpointPaths("", raw)
	})
	if endpointContractErr != nil {
		return "", fmt.Errorf("load endpoint contract: %w", endpointContractErr)
	}
	path, ok := endpointPaths[name]
	if !ok {
		return "", fmt.Errorf("endpoint %s is not defined", name)
	}
	return path, nil
}

func flattenEndpointPaths(prefix string, values map[string]any) {
	for key, value := range values {
		name := key
		if prefix != "" {
			name = prefix + "." + key
		}
		switch typed := value.(type) {
		case string:
			endpointPaths[name] = typed
		case map[string]any:
			flattenEndpointPaths(name, typed)
		}
	}
}
