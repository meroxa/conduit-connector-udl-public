package udl

import "github.com/deepmap/oapi-codegen/pkg/securityprovider"

func NewClientWithBasicAuth(url, username, password string) (*Client, error) {
	authProvider, err := generateBasicAuth(username, password)
	if err != nil {
		return nil, err
	}
	return NewClient(url, WithRequestEditorFn(authProvider.Intercept))
}

func generateBasicAuth(username, password string) (*securityprovider.SecurityProviderBasicAuth, error) {
	return securityprovider.NewSecurityProviderBasicAuth(username, password)
}
