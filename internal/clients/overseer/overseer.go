package overseer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type OverseerClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func NewOverseerClient(baseURL string, apiKey string) *OverseerClient {
	return &OverseerClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// setHeaders sets the required headers for an HTTP request.
func (oc *OverseerClient) setHeaders(req *http.Request) {
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Api-Key", oc.APIKey)
}

// GetMedia retrieves all media items from the API with pagination.
func (oc *OverseerClient) GetMedia() ([]Media, error) {
	var allMedia []Media
	take := 100
	skip := 0

	for {
		endpoint, err := url.Parse(fmt.Sprintf("%s/media", oc.BaseURL))
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %v", err)
		}

		query := endpoint.Query()
		query.Set("take", strconv.Itoa(take))
		query.Set("skip", strconv.Itoa(skip))
		query.Set("sort", "added")
		endpoint.RawQuery = query.Encode()

		req, err := http.NewRequest("GET", endpoint.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		oc.setHeaders(req)

		resp, err := oc.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("received non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var apiResp GetMediaResponse
		if err = json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		allMedia = append(allMedia, apiResp.Results...)
		skip += take
		if skip >= apiResp.PageInfo.Results {
			break
		}
	}

	return allMedia, nil
}

// GetRequests retrieves all requests from the API with pagination.
func (oc *OverseerClient) GetRequests() ([]Request, error) {
	var allRequests []Request
	take := 100
	skip := 0

	for {
		endpoint, err := url.Parse(fmt.Sprintf("%s/request", oc.BaseURL))
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %v", err)
		}

		query := endpoint.Query()
		query.Set("take", strconv.Itoa(take))
		query.Set("skip", strconv.Itoa(skip))
		endpoint.RawQuery = query.Encode()

		req, err := http.NewRequest("GET", endpoint.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		oc.setHeaders(req)

		resp, err := oc.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("received non-OK HTTP status: %s, body: %s", resp.Status, string(bodyBytes))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		var apiResp getRequestsResponse
		if err = json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		allRequests = append(allRequests, apiResp.Results...)
		skip += take
		if skip >= apiResp.PageInfo.Results {
			break
		}
	}

	return allRequests, nil
}

// DeleteMedia sends a DELETE request to the API to remove a media item by its ID.
func (oc *OverseerClient) DeleteMedia(mediaId int) error {
	endpoint := fmt.Sprintf("%s/media/%d", oc.BaseURL, mediaId)

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	oc.setHeaders(req)

	resp, err := oc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	return nil
}

// UpdateRequest sends a PUT request to the API to update a request item.
func (oc *OverseerClient) UpdateRequest(requestID int, updatedRequest Request) error {
	requestBody, err := json.Marshal(updatedRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request item: %v", err)
	}

	endpoint := fmt.Sprintf("%s/request/%d", oc.BaseURL, requestID)

	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	oc.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := oc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	return nil
}

// DeleteRequest sends a DELETE request to the API to remove a request by its ID.
func (oc *OverseerClient) DeleteRequest(requestID int) error {
	endpoint := fmt.Sprintf("%s/request/%d", oc.BaseURL, requestID)

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	oc.setHeaders(req)

	resp, err := oc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	return nil
}
