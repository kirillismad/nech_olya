package usecases

import "strings"

func normalizeTags(names []string) ([]string, error) {
	result := make([]string, 0, len(names))
	seen := make(map[string]struct{}, len(names))

	for _, raw := range names {
		name := strings.ToLower(strings.TrimSpace(raw))
		if name == "" {
			return nil, ErrValidation
		}
		if _, exists := seen[name]; exists {
			continue
		}

		seen[name] = struct{}{}
		result = append(result, name)
	}

	return result, nil
}
