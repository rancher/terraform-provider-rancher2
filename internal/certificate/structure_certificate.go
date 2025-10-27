package rancher2

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

const (
	certificateRSAPrivateKey = "RSA PRIVATE KEY"
	certificatePrivateKey    = "PRIVATE KEY"
	certificatePublicCert    = "CERTIFICATE"
)

// Flatteners

func flattenProjectCertificate(d *schema.ResourceData, in *projectClient.Certificate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("certs", Base64Encode(in.Certs))
	d.Set("project_id", in.ProjectID)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

func flattenNamespacedCertificate(d *schema.ResourceData, in *projectClient.NamespacedCertificate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("certs", Base64Encode(in.Certs))
	d.Set("project_id", in.ProjectID)
	d.Set("namespace_id", in.NamespaceId)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

func flattenCertificate(d *schema.ResourceData, in interface{}) error {
	namespaceID := d.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return flattenNamespacedCertificate(d, in.(*projectClient.NamespacedCertificate))
	}

	return flattenProjectCertificate(d, in.(*projectClient.Certificate))

}

// Expanders

func expandCertificate(in *schema.ResourceData) (interface{}, error) {
	namespaced := false
	if in == nil {
		return nil, nil
	}
	temp := map[string]interface{}{}
	if v := in.Id(); len(v) > 0 {
		temp["id"] = v
	}

	if v, ok := in.Get("certs").(string); ok && len(v) > 0 {
		certs, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: certs is not base64 encoded: %s", v)
		}
		temp[projectClient.NamespacedCertificateFieldCerts] = certs
	}

	if v, ok := in.Get("key").(string); ok && len(v) > 0 {
		key, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: key is not base64 encoded: %s", v)
		}
		temp[projectClient.NamespacedCertificateFieldKey] = key
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	temp[projectClient.NamespacedCertificateFieldProjectID] = projectID

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		temp[projectClient.NamespacedCertificateFieldDescription] = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		temp[projectClient.NamespacedCertificateFieldName] = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		temp[projectClient.NamespacedCertificateFieldAnnotations] = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		temp[projectClient.NamespacedCertificateFieldLabels] = toMapString(v)
	}

	if v, ok := in.Get("namespace_id").(string); ok && len(v) > 0 {
		temp[projectClient.NamespacedCertificateFieldNamespaceId] = v
		namespaced = true
	}

	err := expandCertificateData(temp)
	if err != nil {
		return nil, fmt.Errorf("expanding certificate: %s", err)
	}
	certYAML, err := interfaceToYAML(temp)
	if err != nil {
		return nil, fmt.Errorf("expanding certificate: yaml marshalling %s", err)
	}

	if namespaced {
		obj := &projectClient.NamespacedCertificate{}
		err = yamlToInterface(certYAML, obj)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: yaml unmarshalling %s", err)
		}
		return obj, nil
	}

	obj := &projectClient.Certificate{}
	err = yamlToInterface(certYAML, obj)
	if err != nil {
		return nil, fmt.Errorf("expanding certificate: yaml unmarshalling %s", err)
	}
	return obj, nil
}

func expandCertificateData(input map[string]interface{}) error {
	inputCert, ok := input[projectClient.NamespacedCertificateFieldCerts].(string)
	if !ok || len(inputCert) == 0 {
		return fmt.Errorf("expanding certificate: Certs is nil")
	}
	inputKey, ok := input[projectClient.NamespacedCertificateFieldKey].(string)
	if !ok || len(inputKey) == 0 {
		return fmt.Errorf("expanding certificate: Key is nil")
	}
	certDER, _ := pem.Decode([]byte(inputCert))
	if certDER == nil || certDER.Type != certificatePublicCert {
		return fmt.Errorf("expanding certificate: failed to decode PEM certificates block")
	}
	certificate, err := x509.ParseCertificate(certDER.Bytes)
	if err != nil {
		return fmt.Errorf("expanding certificate: parsing certificates: %s", err)
	}

	keyDER, _ := pem.Decode([]byte(inputKey))
	if keyDER == nil || (keyDER.Type != certificateRSAPrivateKey && keyDER.Type != certificatePrivateKey) {
		return fmt.Errorf("expanding certificate: failed to decode PEM private key block")
	}
	key := &rsa.PrivateKey{}
	if keyDER.Type == certificatePrivateKey {
		k, err := x509.ParsePKCS8PrivateKey(keyDER.Bytes)
		if err != nil {
			return fmt.Errorf("expanding certificate: parsing private key: %s", err)
		}
		v, ok := k.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("expanding certificate: parsed private key is nil")
		}
		key = v
	} else {
		key, err = x509.ParsePKCS1PrivateKey(keyDER.Bytes)
		if err != nil {
			return fmt.Errorf("expanding certificate: parsing rsa private key: %s", err)
		}
	}
	input[projectClient.NamespacedCertificateFieldAlgorithm] = certificate.PublicKeyAlgorithm.String()
	input[projectClient.NamespacedCertificateFieldCN] = certificate.Subject.CommonName
	sum := ""
	for i, data := range sha1.Sum(certificate.Raw) {
		if i > 0 {
			sum = fmt.Sprintf("%s:%02X", sum, data)
			continue
		}
		sum = fmt.Sprintf("%02X", data)
	}
	input[projectClient.NamespacedCertificateFieldCertFingerprint] = string(sum[:])
	input[projectClient.NamespacedCertificateFieldExpiresAt] = certificate.NotAfter.String()
	input[projectClient.NamespacedCertificateFieldIssuedAt] = certificate.NotBefore.String()
	input[projectClient.NamespacedCertificateFieldIssuer] = certificate.Issuer.CommonName
	input[projectClient.NamespacedCertificateFieldKeySize] = strconv.Itoa(key.Size())
	input[projectClient.NamespacedCertificateFieldSerialNumber] = certificate.SerialNumber.String()
	var altNames []string
	if certificate.DNSNames != nil && len(certificate.DNSNames) > 0 {
		altNames = certificate.DNSNames
	}
	for i := range certificate.IPAddresses {
		altNames = append(altNames, certificate.IPAddresses[i].String())
	}
	if len(altNames) > 0 {
		input[projectClient.NamespacedCertificateFieldSubjectAlternativeNames] = altNames
	}
	input[projectClient.NamespacedCertificateFieldVersion] = strconv.Itoa(certificate.Version)

	return nil
}
