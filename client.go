package cas

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client represents a CAS client.
type Client struct {
	// options holds CAS client options
	options clientOptions
}

// Make sure that Client implements the http.Handler interface.
var _ http.Handler = (*Client)(nil)

// NewClient initialize a CAS client.
func NewClient(version ProtocolVersion, serverURL string, opts ...ClientOption) (*Client, error) {
	if version == CASVersionUndefined || version > CASVersion3 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidProtocolVersion, version)
	}
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidServerURL, err.Error())
	}

	options := clientOptions{
		version:    version,
		serverURL:  parsedURL,
		httpClient: nil,
	}

	// Apply custom options.
	for _, o := range opts {
		o.apply(&options)
	}

	// Check the validity of the options, and use the default options for the unspecified options.
	if options.httpClient == nil {
		options.httpClient = &http.Client{}
	}

	return &Client{
		options: options,
	}, nil
}

// Version returns the CAS protocol version used by the client.
func (c *Client) Version() ProtocolVersion {
	return c.options.version
}

func pretreatRequest(req *http.Request) {
	// For reverse proxy, i.e. nginx
	if schema := req.Header.Get("X-Forwarded-Proto"); schema != "" {
		req.URL.Scheme = schema
	} else if req.TLS != nil {
		req.URL.Scheme = "https"
	} else {
		req.URL.Scheme = "http"
	}

	// For reverse proxy, i.e. nginx
	if host := req.Header.Get("X-Forwarded-Host"); host != "" {
		req.URL.Host = host
	} else {
		req.URL.Host = req.Host
	}
}

// RedirectToServerOption specifies the option to pass to the CAS server
// when login to CAS Server.
type RedirectToServerOption struct {
	// Service specifies that the user should be redirected here after logging in
	Service string

	// Renew if this parameter is set, single sign-on will be bypassed. In this case,
	// CAS will require the client to present credentials regardless of the existence
	// of a single sign-on session with CAS. This parameter is not compatible with the
	// gateway parameter.
	Renew bool

	// Gateway if this parameter is set, CAS will not ask the client for credentials.
	// If the client has a pre-existing single sign-on session with CAS, or if a single
	// sign-on session can be established through non-interactive means, CAS MAY redirect
	// the client to the URL specified by the service parameter, appending a valid service
	// ticket. If the client does not have a single sign-on session with CAS, and a non-interactive
	// authentication cannot be established, CAS MUST redirect the client to the URL specified
	// by the service parameter with no “ticket” parameter appended to the URL.
	Gateway bool

	// Method [CAS 3.0] The method to be used when sending responses. While native
	// HTTP redirects (GET) may be utilized as the default method, applications that require a
	// POST response can use this parameter to indicate the method type. A HEADER method may
	// also be specified to indicate the CAS final response such as service and ticketshould
	// be returned in form of HTTP response headers. It is up to the CAS server implementation
	// to determine whether or not POST or HEADER responses are supported.
	Method string
}

// RedirectToServer uses HTTP 302 response to redirect the user to CAS server.
func (c *Client) RedirectToServer(opts *RedirectToServerOption) http.HandlerFunc {
	if opts.Service == "" {
		panic("Service option is required")
	}
	if opts.Service[0] != '/' && !strings.HasPrefix(opts.Service, "http") {
		panic("Invalid Service value: " + opts.Service)
	}
	if opts.Renew && opts.Gateway {
		panic("Renew and Gateway option cannot be set at the same time")
	}
	switch opts.Method {
	case "", "GET", "POST", "HEAD":
	default:
		panic("Unsupported Method option: " + opts.Method)
	}

	return func(resp http.ResponseWriter, req *http.Request) {
		pretreatRequest(req)

		query := url.Values{}
		if opts.Service[0] == '/' {
			query.Add("service", req.URL.Scheme+"://"+req.URL.Host+opts.Service)
		} else {
			query.Add("service", opts.Service)
		}

		if opts.Renew {
			query.Add("renew", "true")
		}
		if opts.Gateway {
			query.Add("gateway", "true")
		}
		if opts.Method != "" {
			query.Add("method", opts.Method)
		}

		serverURL := *c.options.serverURL
		serverURL.Path = CASLoginURI
		serverURL.RawQuery = query.Encode()

		// Write response
		resp.Header().Add("Location", serverURL.String())
		resp.WriteHeader(http.StatusFound)
	}
}

// ValidateServiceTicketOption specifies the option to pass to the CAS server
// when Service Ticket (ST) is validated.
type ValidateServiceTicketOption struct {
	// Renew if this parameter is true, ticket validation will only succeed
	// if the service ticket was issued from the presentation of the user’s
	// primary credentials. It will fail if the ticket was issued from a
	// single sign-on session.
	Renew bool

	// PGTUrl the URL of the proxy callback.
	PGTUrl string

	// Format if this parameter is set, ticket validation response MUST be produced
	// based on the parameter value. Supported values are XML and JSON. If this parameter
	// is not set, the default XML format will be used.
	Format string
}

// ValidateServiceTicket verifies the validity of Service Ticket (ST) to CAS server.
//
func (c *Client) ValidateServiceTicket(opts *ValidateServiceTicketOption) http.HandlerFunc {
	// check ValidateServiceTicketOption
	switch opts.Format {
	case "", "XML", "JSON":
	default:
		panic("Unsupported format: " + opts.Format)
	}

	return func(respw http.ResponseWriter, req *http.Request) {
		// Verify the request
		ticket := req.URL.Query().Get("ticket")
		if len(ticket) < 16 || !strings.HasPrefix(ticket, "ST-") {
			// invalid ticket
			respw.WriteHeader(http.StatusBadRequest)
			return
		}

		pretreatRequest(req)

		sendURL := *c.options.serverURL
		sendQuery := url.Values{}
		sendQuery.Add("service", req.URL.String())
		sendQuery.Add("ticket", ticket)
		if opts.Renew {
			sendQuery.Add("renew", "yes")
		}

		switch c.options.version {
		case CASVersion1:
			sendURL.Path = CASValidateURI
		case CASVersion2:
			sendURL.Path = CASVersion2ServiceValidateURI
			if opts.PGTUrl != "" {
				sendQuery.Add("pgtUrl", opts.PGTUrl)
			}
			if opts.Format != "" {
				sendQuery.Add("format", opts.Format)
			}
		case CASVersion3:
			sendURL.Path = CASVersion3ServiceValidateURI
			if opts.PGTUrl != "" {
				sendQuery.Add("pgtUrl", opts.PGTUrl)
			}
			if opts.Format != "" {
				sendQuery.Add("format", opts.Format)
			}
		default:
			// never reach
			panic("Unknown CAS protocol version")
		}
		sendURL.RawQuery = sendQuery.Encode()

		sendRequest := &http.Request{
			Method: "GET",
			URL:    &sendURL,
		}
		validateResp, err := c.options.httpClient.Do(sendRequest)
		if err != nil {
			// CAS Server error
			respw.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		defer validateResp.Body.Close()

		// Parse CAS Server's response
		switch c.options.version {
		case CASVersion1:
			result, err := ioutil.ReadAll(validateResp.Body)
			if err != nil {
				// CAS Server error
				respw.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			// Due to compiler optimization, there is no additional memory allocation
			if len(result) >= 4 && string(result[0:4]) == "yes\n" {
				r := bytes.IndexByte(result[4:], 0x0A)
				if r < 1 {
					// CAS server does not implement the protocol correctly
					respw.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				// Authentication Success
				c.options.v1ValidateServiceTicketCallback(respw, true, string(result[4:4+r]))
				return
			}
			// Authentication Failed
			c.options.v1ValidateServiceTicketCallback(respw, false, "")
			return

		case CASVersion2, CASVersion3:
			sericeResp := ServiceResponse{}
			switch opts.Format {
			case "XML":
				err := xml.NewDecoder(validateResp.Body).Decode(&sericeResp)
				if err != nil {
					// CAS Server error
					respw.WriteHeader(http.StatusServiceUnavailable)
					return
				}
			case "JSON":
				err := json.NewDecoder(validateResp.Body).Decode(&sericeResp)
				if err != nil {
					// CAS Server error
					respw.WriteHeader(http.StatusServiceUnavailable)
					return
				}
			default:
				// never reach
				respw.WriteHeader(http.StatusInternalServerError)
				return
			}
			c.options.v2ValidateServiceTicketCallback(respw, &sericeResp)
			return
		default:
			// never reach
			respw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (c *Client) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

}
