package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/byebyebruce/wadu/internal/server"
	"github.com/byebyebruce/wadu/model"
)

func PostRawBook(ctx context.Context, host string, b *model.RawBook) error {
	api := host + server.APIPathCreateBook
	return postJSON(ctx, api, b)
}

func postJSON(ctx context.Context, api string, v interface{}) error {
	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", api, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http code: %d\nbody:%s", resp.StatusCode, string(b))
	}
	return nil
}
