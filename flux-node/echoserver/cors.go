/*
The MIT License (MIT)

Copyright (c) 2017 LabStack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package echoserver

import (
	"github.com/bytepowered/flux/flux-node"
	"net/http"
	"strconv"
	"strings"
)

type CorsConfig struct {
	Skipper          flux.WebSkipper
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func NewCORSInterceptor() flux.WebInterceptor {
	return NewCORSMiddlewareWith(CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	})
}

func NewCORSMiddlewareWith(config CorsConfig) flux.WebInterceptor {
	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)
	return func(next flux.WebHandler) flux.WebHandler {
		return func(webex *flux.WebExchange) error {
			if config.Skipper != nil && config.Skipper(webex) {
				return next(webex)
			}

			origin := webex.HeaderVar(flux.HeaderOrigin)
			allowOrigin := ""

			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
				if matchSubdomain(origin, o) {
					allowOrigin = origin
					break
				}
			}
			rwHeader := webex.ResponseWriter().Header()
			// Simple request
			if webex.Method() != http.MethodOptions {
				rwHeader.Add(flux.HeaderVary, flux.HeaderOrigin)
				rwHeader.Set(flux.HeaderAccessControlAllowOrigin, allowOrigin)
				if config.AllowCredentials {
					rwHeader.Set(flux.HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					rwHeader.Set(flux.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				return next(webex)
			}

			// Preflight request
			rwHeader.Add(flux.HeaderVary, flux.HeaderOrigin)
			rwHeader.Add(flux.HeaderVary, flux.HeaderAccessControlRequestMethod)
			rwHeader.Add(flux.HeaderVary, flux.HeaderAccessControlRequestHeaders)
			rwHeader.Set(flux.HeaderAccessControlAllowOrigin, allowOrigin)
			rwHeader.Set(flux.HeaderAccessControlAllowMethods, allowMethods)
			if config.AllowCredentials {
				rwHeader.Set(flux.HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				rwHeader.Set(flux.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := webex.HeaderVar(flux.HeaderAccessControlRequestHeaders)
				if h != "" {
					rwHeader.Set(flux.HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge > 0 {
				rwHeader.Set(flux.HeaderAccessControlMaxAge, maxAge)
			}
			return &flux.ServeError{
				StatusCode: http.StatusNoContent,
				Message:    "NO_CONTENT",
			}
		}
	}
}

func matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func matchSubdomain(domain, pattern string) bool {
	if !matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}
