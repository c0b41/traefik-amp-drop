package ampdrop

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Config the plugin configuration.
type Config struct {
	Querys []string `json:"querys,omitempty"`
	Status int      `json:"status,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Querys: make([]string, 0),
		Status: 301, // Default status code for http redirect.
	}
}

// Amp Drop plugin.
type AmpDrop struct {
	next   http.Handler
	querys []string
	status int
	name   string
}

// New created a new Amp Drop plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Querys) == 0 {
		return nil, fmt.Errorf("querys cannot be empty")
	}

	return &AmpDrop{
		querys: config.Querys,
		status: config.Status,
		next:   next,
		name:   name,
	}, nil
}

func (c *AmpDrop) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	check := false
	q, _ := url.ParseQuery(req.URL.RawQuery) // Get Query params

	for _, v := range c.querys {
		if q.Has(v) {
			q.Del(v) // remove matched params
			check = true
		}
	}

	if check { // only redirect if params match.
		req.URL.RawQuery = q.Encode()
		http.Redirect(rw, req, req.URL.String(), c.status) // redirect new path and params
		return
	}

	c.next.ServeHTTP(rw, req)

}
