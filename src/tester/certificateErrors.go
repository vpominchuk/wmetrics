package tester

import "fmt"

type CertificateFileError struct {
	FileName string
	Err      error
}

type CertificateFileFormatError struct {
	FileName string
	Err      error
}

func (r *CertificateFileError) Error() string {
	return fmt.Sprintf("Failed to read client certificate file: %s. Error: %v", r.FileName, r.Err)
}

func (r *CertificateFileFormatError) Error() string {
	return fmt.Sprintf("Unable to load client certificate and key pair, file: %s. Error: %v", r.FileName, r.Err)
}
