package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/tagmanager/v2"
)

// Environment variable names for client configuration
const (
	EnvCredentialFile  = "GTM_CREDENTIAL_FILE"
	EnvAccountId       = "GTM_ACCOUNT_ID"
	EnvContainerId     = "GTM_CONTAINER_ID"
	EnvWorkspaceName   = "GTM_WORKSPACE_NAME"
	EnvRetryLimit      = "GTM_RETRY_LIMIT"
	EnvRateLimit       = "GTM_RATE_LIMIT"       // requests per second
	EnvRateBurst       = "GTM_RATE_BURST"       // burst capacity
	EnvThrottleEnabled = "GTM_THROTTLE_ENABLED" // enable/disable throttling
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate float64, burst int) *RateLimiter {
	return &RateLimiter{
		tokens:     float64(burst),
		capacity:   float64(burst),
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can proceed
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()

	// Refill tokens based on elapsed time
	rl.tokens = min(rl.capacity, rl.tokens+elapsed*rl.refillRate)
	rl.lastRefill = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	for !rl.Allow() {
		// Calculate how long to wait for the next token
		rl.mutex.Lock()
		waitTime := time.Duration(1000/rl.refillRate) * time.Millisecond
		rl.mutex.Unlock()

		// Wait at least 10ms, but no more than 1 second
		if waitTime < 10*time.Millisecond {
			waitTime = 10 * time.Millisecond
		} else if waitTime > 1*time.Second {
			waitTime = 1 * time.Second
		}

		time.Sleep(waitTime)
	}
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

type ClientOptions struct {
	CredentialFile  string
	AccountId       string
	ContainerId     string
	RetryLimit      int
	RateLimit       float64 // requests per second
	RateBurst       int     // burst capacity
	ThrottleEnabled bool    // enable/disable throttling
}

// NewClientOptionsFromEnv creates ClientOptions from environment variables
func NewClientOptionsFromEnv() *ClientOptions {
	retryLimit := 10 // Default retry limit
	if os.Getenv(EnvRetryLimit) != "" {
		// Ignoring error handling for simplicity, will use default on error
		if retryLimitVal, err := fmt.Sscanf(os.Getenv(EnvRetryLimit), "%d", &retryLimit); err != nil || retryLimitVal <= 0 {
			retryLimit = 10
		}
	}

	// Default rate limiting: 10 requests per second with burst of 20
	rateLimit := 10.0
	if rateLimitEnv := os.Getenv(EnvRateLimit); rateLimitEnv != "" {
		if val, err := strconv.ParseFloat(rateLimitEnv, 64); err == nil && val > 0 {
			rateLimit = val
		}
	}

	rateBurst := 20
	if rateBurstEnv := os.Getenv(EnvRateBurst); rateBurstEnv != "" {
		if val, err := strconv.Atoi(rateBurstEnv); err == nil && val > 0 {
			rateBurst = val
		}
	}

	throttleEnabled := true // Default to enabled
	if throttleEnv := os.Getenv(EnvThrottleEnabled); throttleEnv != "" {
		if val, err := strconv.ParseBool(throttleEnv); err == nil {
			throttleEnabled = val
		}
	}

	return &ClientOptions{
		CredentialFile:  os.Getenv(EnvCredentialFile),
		AccountId:       os.Getenv(EnvAccountId),
		ContainerId:     os.Getenv(EnvContainerId),
		RetryLimit:      retryLimit,
		RateLimit:       rateLimit,
		RateBurst:       rateBurst,
		ThrottleEnabled: throttleEnabled,
	}
}

type Client struct {
	*tagmanager.Service

	Options     *ClientOptions
	rateLimiter *RateLimiter
}

func NewClient(opts *ClientOptions) (*Client, error) {
	var ctx = context.Background()

	srv, err := tagmanager.NewService(ctx, option.WithCredentialsFile(opts.CredentialFile))
	if err != nil {
		return nil, err
	}

	var rateLimiter *RateLimiter
	if opts.ThrottleEnabled {
		rateLimiter = NewRateLimiter(opts.RateLimit, opts.RateBurst)
	}

	return &Client{
		Service:     srv,
		Options:     opts,
		rateLimiter: rateLimiter,
	}, nil
}

// NewClientFromEnv creates a new client using environment variables
func NewClientFromEnv() (*Client, error) {
	return NewClient(NewClientOptionsFromEnv())
}

func (c *Client) containerPath() string {
	opts := c.Options
	return "accounts/" + opts.AccountId + "/containers/" + opts.ContainerId
}

var ErrNotExist = errors.New("not exist")

func (c *Client) CreateWorkspace(ws *tagmanager.Workspace) (*tagmanager.Workspace, error) {
	return c.getWorkspaceWithRetry(c.Accounts.Containers.Workspaces.Create(c.containerPath(), ws).Do)
}

func (c *Client) ListWorkspaces() ([]*tagmanager.Workspace, error) {
	resp, err := c.getWorkspaceListWithRetry(c.Accounts.Containers.Workspaces.List(c.containerPath()).Do)
	if err != nil {
		return nil, err
	} else {
		return resp.Workspace, nil
	}
}

func (c *Client) Workspace(id string) (*tagmanager.Workspace, error) {
	ws, err := c.getWorkspaceWithRetry(c.Accounts.Containers.Workspaces.Get(c.containerPath() + "/workspaces/" + id).Do)
	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return ws, err
	}
}

func (c *Client) UpdateWorkspaces(id string, ws *tagmanager.Workspace) (*tagmanager.Workspace, error) {
	return c.getWorkspaceWithRetry(c.Accounts.Containers.Workspaces.Update(c.containerPath()+"/workspaces/"+id, ws).Do)
}

func (c *Client) DeleteWorkspace(id string) error {
	return c.executeWithRetry(c.Accounts.Containers.Workspaces.Delete(c.containerPath() + "/workspaces/" + id).Do)
}

func (c *Client) workspacePath(id string) string {
	return c.containerPath() + "/workspaces/" + id
}

// throttle applies rate limiting if enabled
func (c *Client) throttle() {
	if c.rateLimiter != nil {
		c.rateLimiter.Wait()
	}
}

func (c *Client) CreateTag(workspaceId string, tag *tagmanager.Tag) (*tagmanager.Tag, error) {

	return c.getTagWithRetry(c.Accounts.Containers.Workspaces.Tags.Create(c.workspacePath(workspaceId), tag).Do)
}

func (c *Client) ListTags(workspaceId string) ([]*tagmanager.Tag, error) {
	resp, err := c.getTagListWithRetry(c.Accounts.Containers.Workspaces.Tags.List(c.workspacePath(workspaceId)).Do)
	if err != nil {
		return nil, err
	} else {
		return resp.Tag, nil
	}
}

func (c *Client) Tag(workspaceId string, tagId string) (*tagmanager.Tag, error) {
	tag, err := c.getTagWithRetry(c.Accounts.Containers.Workspaces.Tags.Get(c.workspacePath(workspaceId) + "/tags/" + tagId).Do)

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return tag, err
	}
}

func (c *Client) UpdateTag(workspaceId string, tagId string, tag *tagmanager.Tag) (*tagmanager.Tag, error) {
	return c.getTagWithRetry(c.Accounts.Containers.Workspaces.Tags.Update(c.workspacePath(workspaceId)+"/tags/"+tagId, tag).Do)
}

func (c *Client) DeleteTag(workspaceId string, tagId string) error {
	return c.executeWithRetry(c.Accounts.Containers.Workspaces.Tags.Delete(c.workspacePath(workspaceId) + "/tags/" + tagId).Do)
}

func (c *Client) CreateVariable(workspaceId string, variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	return c.getVariableWithRetry(c.Accounts.Containers.Workspaces.Variables.Create(c.workspacePath(workspaceId), variable).Do)
}

func (c *Client) ListVariables(workspaceId string) ([]*tagmanager.Variable, error) {
	resp, err := c.getVariableListWithRetry(c.Accounts.Containers.Workspaces.Variables.List(c.workspacePath(workspaceId)).Do)
	if err != nil {
		return nil, err
	} else {
		return resp.Variable, nil
	}
}

func (c *Client) Variable(workspaceId string, variableId string) (*tagmanager.Variable, error) {
	variable, err := c.getVariableWithRetry(c.Accounts.Containers.Workspaces.Variables.Get(c.workspacePath(workspaceId) + "/variables/" + variableId).Do)

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return variable, err
	}
}

func (c *Client) UpdateVariable(workspaceId string, variableId string, variable *tagmanager.Variable) (*tagmanager.Variable, error) {
	return c.getVariableWithRetry(c.Accounts.Containers.Workspaces.Variables.Update(c.workspacePath(workspaceId)+"/variables/"+variableId, variable).Do)
}

func (c *Client) DeleteVariable(workspaceId string, variableId string) error {
	return c.executeWithRetry(c.Accounts.Containers.Workspaces.Variables.Delete(c.workspacePath(workspaceId) + "/variables/" + variableId).Do)
}

func (c *Client) CreateTrigger(workspaceId string, trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	return c.getTriggerWithRetry(c.Accounts.Containers.Workspaces.Triggers.Create(c.workspacePath(workspaceId), trigger).Do)
}

func (c *Client) ListTriggers(workspaceId string) ([]*tagmanager.Trigger, error) {
	resp, err := c.getTriggerListWithRetry(c.Accounts.Containers.Workspaces.Triggers.List(c.workspacePath(workspaceId)).Do)
	if err != nil {
		return nil, err
	} else {
		return resp.Trigger, nil
	}
}

func (c *Client) Trigger(workspaceId string, triggerId string) (*tagmanager.Trigger, error) {
	trigger, err := c.getTriggerWithRetry(c.Accounts.Containers.Workspaces.Triggers.Get(c.workspacePath(workspaceId) + "/triggers/" + triggerId).Do)

	if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 404 {
		return nil, ErrNotExist
	} else {
		return trigger, err
	}
}

func (c *Client) UpdateTrigger(workspaceId string, triggerId string, trigger *tagmanager.Trigger) (*tagmanager.Trigger, error) {
	return c.getTriggerWithRetry(c.Accounts.Containers.Workspaces.Triggers.Update(c.workspacePath(workspaceId)+"/triggers/"+triggerId, trigger).Do)
}

func (c *Client) DeleteTrigger(workspaceId string, triggerId string) error {
	return c.executeWithRetry(c.Accounts.Containers.Workspaces.Triggers.Delete(c.workspacePath(workspaceId) + "/triggers/" + triggerId).Do)
}

func (c *Client) executeWithRetry(query func(opts ...googleapi.CallOption) error) error {
	retryCount := 0

	for {
		// Apply throttling before making the request
		c.throttle()

		err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := time.Duration(retryCount) * time.Second
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Helper methods for different return types
func (c *Client) getWorkspaceWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.Workspace, error)) (*tagmanager.Workspace, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getWorkspaceListWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.ListWorkspacesResponse, error)) (*tagmanager.ListWorkspacesResponse, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getTagWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.Tag, error)) (*tagmanager.Tag, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getTagListWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.ListTagsResponse, error)) (*tagmanager.ListTagsResponse, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getVariableWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.Variable, error)) (*tagmanager.Variable, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getVariableListWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.ListVariablesResponse, error)) (*tagmanager.ListVariablesResponse, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getTriggerWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.Trigger, error)) (*tagmanager.Trigger, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}

func (c *Client) getTriggerListWithRetry(query func(opts ...googleapi.CallOption) (*tagmanager.ListTriggersResponse, error)) (*tagmanager.ListTriggersResponse, error) {
	retryCount := 0

	for {
		c.throttle()

		resp, err := query()
		if errTyped, ok := err.(*googleapi.Error); ok && errTyped.Code == 429 {
			if retryCount < c.Options.RetryLimit {
				retryCount++
				backoffDuration := 20 * time.Second * time.Duration(retryCount)
				fmt.Printf("Rate limit exceeded. Retrying in %s...\n", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			} else {
				return nil, fmt.Errorf("rate limit exceeded after %d retries", c.Options.RetryLimit)
			}
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}
}
