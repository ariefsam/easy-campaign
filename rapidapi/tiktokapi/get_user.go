package tiktokapi

import (
	"campaign/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (t *tiktokApi) GetUser(ctx context.Context, username string) (user User, err error) {
	host := os.Getenv("TIKTOK_API_HOST")
	key := os.Getenv("RAPIDAPI_KEY")
	url := fmt.Sprintf("https://%s/api/user/info?uniqueId=%s", host, username)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", key)
	req.Header.Add("x-rapidapi-host", host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if err := json.Unmarshal(body, &user); err != nil {
		logger.Error(err)
		return user, err
	}
	return
}
