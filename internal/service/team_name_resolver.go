package service

import (
	"context"
	"strings"
	"sync"

	"github.com/jendrix/worldcup-stats-service/internal/domain"
)

const defaultTeamNameLanguage = "es"

// TeamTranslationProvider defines the dependency needed to load team names.
type TeamTranslationProvider interface {
	ListTeamTranslations(ctx context.Context) ([]domain.TeamTranslation, error)
}

// TeamNameResolver resolves localized team names by code.
type TeamNameResolver interface {
	Resolve(ctx context.Context, code string, language string) (string, error)
}

type cachedTeamNameResolver struct {
	provider TeamTranslationProvider
	mu       sync.RWMutex
	loaded   bool
	names    map[string]map[string]string
}

// NewCachedTeamNameResolver creates a resolver backed by team translations.
func NewCachedTeamNameResolver(provider TeamTranslationProvider) TeamNameResolver {
	return &cachedTeamNameResolver{provider: provider}
}

func (r *cachedTeamNameResolver) Resolve(ctx context.Context, code string, language string) (string, error) {
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))
	if normalizedCode == "" {
		return "", nil
	}

	if err := r.ensureLoaded(ctx); err != nil {
		return "", err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if namesByCode := r.names[normalizeTeamNameLanguage(language)]; namesByCode != nil {
		if name := namesByCode[normalizedCode]; name != "" {
			return name, nil
		}
	}
	if namesByCode := r.names[defaultTeamNameLanguage]; namesByCode != nil {
		if name := namesByCode[normalizedCode]; name != "" {
			return name, nil
		}
	}
	return normalizedCode, nil
}

func (r *cachedTeamNameResolver) ensureLoaded(ctx context.Context) error {
	r.mu.RLock()
	if r.loaded {
		r.mu.RUnlock()
		return nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.loaded {
		return nil
	}

	translations, err := r.provider.ListTeamTranslations(ctx)
	if err != nil {
		return err
	}

	names := make(map[string]map[string]string)
	for _, translation := range translations {
		language := normalizeTeamNameLanguage(translation.Language)
		code := strings.ToUpper(strings.TrimSpace(translation.TeamCode))
		if language == "" || code == "" {
			continue
		}
		if _, ok := names[language]; !ok {
			names[language] = make(map[string]string)
		}
		names[language][code] = translation.Name
	}
	r.names = names
	r.loaded = true

	return nil
}

func normalizeTeamNameLanguage(language string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		return defaultTeamNameLanguage
	}
	return language
}
