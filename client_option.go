package cas

import (
	"net/http"
	"net/url"
)

// V1ValidateServiceTicketCallbackFunc represents the callback function
// after Service Ticket (ST) verification to CAS server.
//
// Note: the Service Ticket is valid only if isValid == true.
//
type V1ValidateServiceTicketCallbackFunc func(respw http.ResponseWriter, isValid bool, user string)

// V2ValidateServiceTicketCallbackFunc represents the callback function
// after Service Ticket (ST) verification to CAS server.
//
// Note: the Service Ticket is valid only if resp.Success != nil.
//
type V2ValidateServiceTicketCallbackFunc func(respw http.ResponseWriter, serviceResp *ServiceResponse)

// V3ValidateServiceTicketCallbackFunc represents the callback function
// after Service Ticket (ST) verification to CAS server.
//
// Note: the Service Ticket is valid only if resp.Success != nil.
//
// TODO: attributes type is need to define.
type V3ValidateServiceTicketCallbackFunc func(respw http.ResponseWriter, serviceResp *ServiceResponse)

// clientOptions represents the options used to initialize CAS client.
type clientOptions struct {
	// version represents a CAS protocol version to use.
	// it should be a CASVersion* constant.
	version ProtocolVersion

	// serverURL represents the URL address of CAS server, which will be spliced
	// with the URI specified in CAS protocol to generate the final requested URL.
	// For example, the final request address is
	//
	//     https://cas.example.org/cas/login?service=http%3A%2F%2Fwww.example.org%2Fservice
	//
	// Then, the serverURL should be https://cas.example.org/cas.
	serverURL *url.URL

	// httpClient represents the HTTP client used by CAS client to send HTTP request.
	httpClient *http.Client

	// ValidateServiceTicketCallback represents the callback function
	// after Service Ticket (ST) verification to CAS server.
	//
	// Refer to V1ValidateServiceTicketCallbackFunc, etc.
	//
	v1ValidateServiceTicketCallback V1ValidateServiceTicketCallbackFunc
	v2ValidateServiceTicketCallback V2ValidateServiceTicketCallbackFunc
	v3ValidateServiceTicketCallback V3ValidateServiceTicketCallbackFunc
}

// ClientOption is used to set options for initializing CAS client.
type ClientOption interface {
	apply(opts *clientOptions)
}

type httpClientOption struct {
	HTTPClient *http.Client
}

func (o httpClientOption) apply(opts *clientOptions) {
	opts.httpClient = o.HTTPClient
}

// WithHTTPClient set the HTTP client used by CAS client to send HTTP request.
func WithHTTPClient(c *http.Client) ClientOption {
	return httpClientOption{HTTPClient: c}
}

type validateServiceTicketCallbackOption struct {
	ValidateServiceTicketCallback interface{}
}

func (o validateServiceTicketCallbackOption) apply(opts *clientOptions) {
	switch opts.version {
	case CASVersion1:
		if f, ok := o.ValidateServiceTicketCallback.(V1ValidateServiceTicketCallbackFunc); ok {
			opts.v1ValidateServiceTicketCallback = f
			return
		}
		panic("Illegal parameter: type V1ValidateServiceTicketCallbackFunc must be passed")
	case CASVersion2:
		if f, ok := o.ValidateServiceTicketCallback.(V2ValidateServiceTicketCallbackFunc); ok {
			opts.v2ValidateServiceTicketCallback = f
			return
		}
		panic("Illegal parameter: type V2ValidateServiceTicketCallbackFunc must be passed")
	case CASVersion3:
		if f, ok := o.ValidateServiceTicketCallback.(V3ValidateServiceTicketCallbackFunc); ok {
			opts.v3ValidateServiceTicketCallback = f
			return
		}
		panic("Illegal parameter: type V3ValidateServiceTicketCallbackFunc must be passed")
	default:
		// never reach
		panic("Unknown CAS protocol version")
	}
}

// WithValidateServiceTicketCallback accepts a function type that will be called back
// after Service Ticket validation.
// According to the CAS protocol version used by the client, f must be one of
// ValidateServiceTicketCallbackFunc, e.g. V3ValidateServiceTicketCallbackFunc.
// If the illegal parameter f is passed, panic will occur.
func WithValidateServiceTicketCallback(f interface{}) ClientOption {
	return validateServiceTicketCallbackOption{ValidateServiceTicketCallback: f}
}
