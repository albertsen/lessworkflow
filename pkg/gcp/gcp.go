package gcp

import "os"

var Project string

func init() {
	Project = os.Getenv("LW_GCP_PROJECT_ID")
	if Project == "" {
		Project = "sap-se-commerce-arch"
	}
}
