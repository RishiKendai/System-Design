package url

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lld/url-shortner/counter"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLShortnerService struct {
	mu         sync.RWMutex
	generators []*counter.RangeGenerator
	urls       map[string]*URL
	aliases    map[string]bool
}

func NewURLShortnerService(n int) *URLShortnerService {
	gens := make([]*counter.RangeGenerator, n)

	for i := 0; i < n; i++ {
		gens[i] = counter.NewRangeGenerator(fmt.Sprintf("generator%d", i), 1000)
	}
	return &URLShortnerService{
		urls:       make(map[string]*URL),
		generators: gens,
		aliases:    make(map[string]bool),
	}
}

func (s *URLShortnerService) HealthCheck() error {
	if s == nil || len(s.generators) == 0 {
		return fmt.Errorf("url shortener: not ready")
	}
	return nil
}

func (s *URLShortnerService) GenerateShortCode(longURL string, userID string, customAlias string, expiresAt *time.Time) *URL {
	s.mu.Lock()
	defer s.mu.Unlock()

	var shortCode string
	if customAlias != "" {
		if _, ok := s.aliases[customAlias]; ok {
			return nil
		}
		shortCode = customAlias
		s.aliases[customAlias] = true
	} else {
		gen := s.generators[rand.Intn(len(s.generators))]
		num := gen.Next()
		shortCode = toFixedLength(encodeBase62(num), 7)
		if _, ok := s.aliases[shortCode]; ok {
			return nil
		}
		s.aliases[shortCode] = true
	}

	u := &URL{
		ID:          uuid.New().String(),
		LongURL:     longURL,
		UserID:      userID,
		CustomAlias: customAlias,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		IsActive:    true,
		IsCustom:    customAlias != "",
	}
	if expiresAt != nil {
		u.ExpiresAt = *expiresAt
	}
	s.urls[shortCode] = u
	return u
}

func (s *URLShortnerService) GetURL(shortCode string) string {
	s.mu.RLock()
	u, ok := s.urls[shortCode]
	if !ok || u == nil {
		s.mu.RUnlock()
		return ""
	}
	if u.IsDeleted {
		s.mu.RUnlock()
		return ""
	}
	if !u.ExpiresAt.IsZero() && u.ExpiresAt.Before(time.Now()) {
		s.mu.RUnlock()
		return ""
	}
	if !u.IsActive {
		s.mu.RUnlock()
		return ""
	}
	long := u.LongURL
	s.mu.RUnlock()

	u.recordClick()
	return long
}

func (s *URLShortnerService) DeleteURL(shortCode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.urls[shortCode]
	if !ok {
		return
	}
	u.IsDeleted = true
}

func (s *URLShortnerService) PrintUrl(shortCode string) {
	s.mu.RLock()
	u, ok := s.urls[shortCode]
	s.mu.RUnlock()
	if !ok {
		return
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("URL: ", u.LongURL)
	fmt.Println("Short Code: ", u.ShortCode)
	fmt.Println("--------------------------------")
	fmt.Println("Created At: ", u.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("Expires At: ", u.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Println("Is Active: ", u.IsActive)
	fmt.Println("Is Deleted: ", u.IsDeleted)
	fmt.Println("Is Custom: ", u.IsCustom)
	fmt.Println("Custom Alias: ", u.CustomAlias)
	fmt.Println("Click Count: ", u.Clicks())
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}
