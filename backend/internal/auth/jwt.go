// Package auth provides JWT token issuance, validation, and HTTP middleware.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenAudienceAccess    = "access"
	tokenAudienceRefresh   = "refresh"
	tokenAudienceElevation = "elevation"
)

// Claims represents the JWT claims for access and refresh tokens.
// InstanceAdmin reflects only the flat, non-connector-scoped role (user
// management, API keys, granting connector permissions); per-connector
// access is looked up per-request from the store, not carried in the token,
// since it can change per connector at any time (see auth.RequireConnectorRole).
type Claims struct {
	jwt.RegisteredClaims
	UserID        string `json:"uid"`
	InstanceAdmin bool   `json:"admin"`
}

// APIKeyClaims represents the identity and lifecycle fields needed to
// authenticate an opaque API key. Like Claims.InstanceAdmin, InstanceAdmin
// here is the flat instance-wide role only.
type APIKeyClaims struct {
	KeyID         string
	UserID        string
	InstanceAdmin bool
	ExpiresAt     string
	RevokedAt     string
	LastUsedAt    string
}

// ElevationClaims represents a short-lived step-up token for destructive actions.
type ElevationClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"uid"`
	Action string `json:"action"` // e.g. "connector.delete"
}

// TokenPair is the response for a successful login or refresh.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int    `json:"expiresIn"` // seconds until access token expires
}

// ElevationToken is the response for a successful step-up re-authentication.
type ElevationToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Service handles JWT creation and validation.
type Service struct {
	secret       []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
	elevationTTL time.Duration
}

// NewService creates a new Service.
func NewService(secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		secret:       []byte(secret),
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
		elevationTTL: 60 * time.Second, // hardcoded: elevation tokens live 60s
	}
}

// IssuePair creates a new access + refresh token pair.
func (s *Service) IssuePair(userID string, instanceAdmin bool) (*TokenPair, error) {
	now := time.Now()

	access, err := s.issue(Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{tokenAudienceAccess},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			ID:        newTokenID(),
		},
		UserID:        userID,
		InstanceAdmin: instanceAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("issue access token: %w", err)
	}

	refresh, err := s.issue(Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{tokenAudienceRefresh},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			ID:        newTokenID(),
		},
		UserID:        userID,
		InstanceAdmin: instanceAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("issue refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(s.accessTTL.Seconds()),
	}, nil
}

// ValidateAccess validates an access token and returns its claims.
func (s *Service) ValidateAccess(tokenString string) (*Claims, error) {
	return s.validate(tokenString, tokenAudienceAccess)
}

// ValidateRefresh validates a refresh token and returns its claims.
func (s *Service) ValidateRefresh(tokenString string) (*Claims, error) {
	return s.validate(tokenString, tokenAudienceRefresh)
}

// IssueElevation creates a short-lived elevation token scoped to a single action.
func (s *Service) IssueElevation(userID, action string) (*ElevationToken, error) {
	now := time.Now()
	expiresAt := now.Add(s.elevationTTL)

	token, err := s.issueElevation(ElevationClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{tokenAudienceElevation},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        newTokenID(),
		},
		UserID: userID,
		Action: action,
	})
	if err != nil {
		return nil, fmt.Errorf("issue elevation token: %w", err)
	}

	return &ElevationToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateElevation validates an elevation token and checks it matches the required action.
func (s *Service) ValidateElevation(tokenString, action, userID string) (*ElevationClaims, error) {
	claims, err := s.validateElevation(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Action != action {
		return nil, fmt.Errorf("elevation token is for action %q, not %q", claims.Action, action)
	}
	if claims.UserID != userID {
		return nil, fmt.Errorf("elevation token belongs to a different user")
	}
	return claims, nil
}

func (s *Service) issue(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) validate(tokenString, audience string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	c, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if !hasAudience(c.Audience, audience) {
		return nil, fmt.Errorf("token is not for %s", audience)
	}
	return c, nil
}

func (s *Service) issueElevation(claims ElevationClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) validateElevation(tokenString string) (*ElevationClaims, error) {
	claims := &ElevationClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse elevation token: %w", err)
	}
	c, ok := token.Claims.(*ElevationClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid elevation token")
	}
	if !hasAudience(c.Audience, tokenAudienceElevation) {
		return nil, fmt.Errorf("token is not an elevation token")
	}
	return c, nil
}

func hasAudience(audiences jwt.ClaimStrings, audience string) bool {
	for _, value := range audiences {
		if value == audience {
			return true
		}
	}
	return false
}

func newTokenID() string {
	return uuid.NewString()
}
