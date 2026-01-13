package sso

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/crewjam/saml"
	"github.com/hanyouqing/openauth/internal/models"
)

func ParseCertificate(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	return cert, nil
}

func ParsePrivateKey(keyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		rsaKey, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		return rsaKey, nil
	}
	return key, nil
}

func BuildSAMLResponse(samlConfig *models.SAMLConfig, user *models.User, requestID string) (*saml.Response, error) {
	_, err := ParseCertificate(samlConfig.Certificate)
	if err != nil {
		return nil, err
	}

	_, err = ParsePrivateKey(samlConfig.PrivateKey)
	if err != nil {
		return nil, err
	}

	now := saml.TimeNow()
	response := &saml.Response{
		Destination:  samlConfig.SSOURL,
		ID:           fmt.Sprintf("id-%d", time.Now().UnixNano()),
		InResponseTo: requestID,
		IssueInstant: now,
		Version:      "2.0",
		Issuer: &saml.Issuer{
			Value: samlConfig.EntityID,
		},
		Status: saml.Status{
			StatusCode: saml.StatusCode{
				Value: saml.StatusSuccess,
			},
		},
	}

	// Build assertion
	assertion := &saml.Assertion{
		ID:           fmt.Sprintf("assertion-%d", time.Now().UnixNano()),
		IssueInstant: now,
		Version:      "2.0",
		Issuer: saml.Issuer{
			Value: samlConfig.EntityID,
		},
		Subject: &saml.Subject{
			NameID: &saml.NameID{
				Format: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
				Value:  user.Email,
			},
			SubjectConfirmations: []saml.SubjectConfirmation{
				{
					Method: "urn:oasis:names:tc:SAML:2.0:cm:bearer",
					SubjectConfirmationData: &saml.SubjectConfirmationData{
						NotOnOrAfter: now.Add(5 * time.Minute),
						Recipient:    samlConfig.SSOURL,
					},
				},
			},
		},
		Conditions: &saml.Conditions{
			NotBefore:    now.Add(-5 * time.Minute),
			NotOnOrAfter: now.Add(5 * time.Minute),
			AudienceRestrictions: []saml.AudienceRestriction{
				{
					Audience: saml.Audience{
						Value: samlConfig.EntityID,
					},
				},
			},
		},
		AuthnStatements: []saml.AuthnStatement{
			{
				AuthnInstant: now,
				SessionIndex: fmt.Sprintf("session-%d", time.Now().UnixNano()),
				SubjectLocality: &saml.SubjectLocality{
					Address: "127.0.0.1",
				},
				AuthnContext: saml.AuthnContext{
					AuthnContextClassRef: &saml.AuthnContextClassRef{
						Value: "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport",
					},
				},
			},
		},
	}

	// Add attributes based on attribute map
	if samlConfig.AttributeMap != nil {
		attributes := []saml.Attribute{}
		for samlAttr, userAttr := range samlConfig.AttributeMap {
			var value string
			switch userAttr {
			case "email":
				value = user.Email
			case "username":
				value = user.Username
			case "name":
				value = user.Username
			default:
				continue
			}
			attributes = append(attributes, saml.Attribute{
				Name:       samlAttr,
				NameFormat: "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
				Values: []saml.AttributeValue{
					{Value: value},
				},
			})
		}
		assertion.AttributeStatements = []saml.AttributeStatement{
			{
				Attributes: attributes,
			},
		}
	}

	response.Assertion = assertion

	return response, nil
}

func BuildSAMLMetadata(entityID, ssoURL, sloURL string, cert *x509.Certificate) (*saml.EntityDescriptor, error) {
	// Build SAML metadata
	metadata := &saml.EntityDescriptor{
		EntityID: entityID,
		IDPSSODescriptors: []saml.IDPSSODescriptor{
			{
				SingleSignOnServices: []saml.Endpoint{
					{
						Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
						Location: ssoURL,
					},
					{
						Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
						Location: ssoURL,
					},
				},
			},
		},
	}

	if sloURL != "" {
		metadata.IDPSSODescriptors[0].SingleLogoutServices = []saml.Endpoint{
			{
				Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
				Location: sloURL,
			},
		}
	}

	// Add certificate to key descriptor
	if cert != nil {
		metadata.IDPSSODescriptors[0].KeyDescriptors = []saml.KeyDescriptor{
			{
				Use: "signing",
				KeyInfo: saml.KeyInfo{
					X509Data: saml.X509Data{
						X509Certificates: []saml.X509Certificate{
							{
								Data: base64.StdEncoding.EncodeToString(cert.Raw),
							},
						},
					},
				},
			},
		}
	}

	return metadata, nil
}
