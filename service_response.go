package cas

import (
	"encoding/xml"
	"time"
)

type ServiceResponse struct {
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas serviceResponse"`

	Failure *AuthenticationFailure
	Success *AuthenticationSuccess
}

type AuthenticationFailure struct {
	XMLName xml.Name `xml:"authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",innerxml"`
}

type AuthenticationSuccess struct {
	XMLName             xml.Name        `xml:"authenticationSuccess"`
	User                string          `xml:"user"`
	ProxyGrantingTicket string          `xml:"proxyGrantingTicket,omitempty"`
	Proxies             *Proxies        `xml:"proxies"`
	Attributes          *Attributes     `xml:"attributes"`
	ExtraAttributes     []*AnyAttribute `xml:",any"`
}

type Proxies struct {
	XMLName xml.Name `xml:"proxies"`
	Proxies []string `xml:"proxy"`
}

func (p *Proxies) AddProxy(proxy string) {
	p.Proxies = append(p.Proxies, proxy)
}

type Attributes struct {
	XMLName                                xml.Name  `xml:"attributes"`
	AuthenticationDate                     time.Time `xml:"authenticationDate"`
	LongTermAuthenticationRequestTokenUsed bool      `xml:"longTermAuthenticationRequestTokenUsed"`
	IsFromNewLogin                         bool      `xml:"isFromNewLogin"`
	MemberOf                               []string  `xml:"memberOf"`
	UserAttributes                         *UserAttributes
	ExtraAttributes                        []*AnyAttribute `xml:",any"`
}

type UserAttributes struct {
	XMLName       xml.Name          `xml:"userAttributes"`
	Attributes    []*NamedAttribute `xml:"attribute"`
	AnyAttributes []*AnyAttribute   `xml:",any"`
}

type NamedAttribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:",innerxml"`
}

type AnyAttribute struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
