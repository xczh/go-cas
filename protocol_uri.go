package cas

const (
	// CASLoginURI represents credential requestor / acceptor
	CASLoginURI = "/login"

	// CASLogoutURI represents destroy CAS session (logout)
	CASLogoutURI = "/logout"

	// CASValidateURI represents service ticket validation
	CASValidateURI = "/validate"

	// CASVersion2ServiceValidateURI represents service ticket validation [CAS 2.0]
	CASVersion2ServiceValidateURI = "/serviceValidate"

	// CASVersion2ProxyValidateURI represents service/proxy ticket validation [CAS 2.0]
	CASVersion2ProxyValidateURI = "/proxyValidate"

	// CASVersion2ProxyURI represents proxy ticket service [CAS 2.0]
	CASVersion2ProxyURI = "/proxy"

	// CASVersion3ServiceValidateURI represents service ticket validation [CAS 3.0]
	CASVersion3ServiceValidateURI = "/p3/serviceValidate"

	// CASVersion3ProxyValidateURI represents service/proxy ticket validation [CAS 3.0]
	CASVersion3ProxyValidateURI = "/p3/proxyValidate"
)
