package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/bitbucket"
)

// Bitbucket

type bitbucketProvider struct {
	*oauth2.Config
}

func NewBitbucketProvider(key, secret string) Provider {
	return &bitbucketProvider{
		&oauth2.Config{
			ClientID:     key,
			ClientSecret: secret,
			Endpoint:     bitbucket.Endpoint,
		},
	}
}

func (g bitbucketProvider) GetOAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.Exchange(ctx, code)
}

func (g bitbucketProvider) GetUserEmail(ctx context.Context, tok *oauth2.Token) (string, error) {
	client := g.Client(ctx, tok)
	userRes, err := client.Get("https://api.bitbucket.org/1.0/user")
	if err != nil {
		return "", err
	}
	defer userRes.Body.Close()

	username := ""
	u := map[string]interface{}{}
	if err := json.NewDecoder(userRes.Body).Decode(&u); err != nil {
		return "", err
	}

	x, ok := u["user"]
	u, ok2 := x.(map[string]interface{})
	if !ok || !ok2 {
		return "", errors.New("Invalid response when requesting email address from bitbucket")

	}

	if name, ok := u["username"]; ok {
		username, ok = name.(string)
		if !ok {
			return "", errors.New("Invalid response when requesting email address from bitbucket")
		}
	}

	return getUserEmail(ctx, tok, fmt.Sprintf("https://api.bitbucket.org/1.0/users/%s/emails", username), g.Config)
}
