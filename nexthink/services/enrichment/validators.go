package enrichment

import "fmt"

// ValidateEnrichmentRequest validates an enrichment request
func ValidateEnrichmentRequest(req *EnrichmentRequest) error {
	if req == nil {
		return fmt.Errorf("enrichment request cannot be nil")
	}

	if req.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	if len(req.Enrichments) == 0 {
		return fmt.Errorf("enrichments is required and must contain at least one enrichment")
	}

	if len(req.Enrichments) > 5000 {
		return fmt.Errorf("enrichments cannot contain more than 5000 items (got %d)", len(req.Enrichments))
	}

	for i, enrichment := range req.Enrichments {
		if err := validateEnrichment(i, &enrichment); err != nil {
			return err
		}
	}

	return nil
}

// validateEnrichment validates a single enrichment
func validateEnrichment(index int, enrichment *Enrichment) error {
	if len(enrichment.Identification) != 1 {
		return fmt.Errorf("enrichments[%d].identification must contain exactly 1 item (got %d)", index, len(enrichment.Identification))
	}

	if err := validateIdentification(index, &enrichment.Identification[0]); err != nil {
		return err
	}

	if len(enrichment.Fields) == 0 {
		return fmt.Errorf("enrichments[%d].fields is required and must contain at least one field", index)
	}

	for j, field := range enrichment.Fields {
		if err := validateField(index, j, &field); err != nil {
			return err
		}
	}

	return nil
}

// validateIdentification validates an identification object
func validateIdentification(enrichmentIndex int, id *Identification) error {
	if id.Name == "" {
		return fmt.Errorf("enrichments[%d].identification[0].name is required", enrichmentIndex)
	}

	if id.Value == "" {
		return fmt.Errorf("enrichments[%d].identification[0].value is required", enrichmentIndex)
	}

	// Validate identification name is one of the allowed values
	validNames := map[string]bool{
		IdentificationDeviceName:  true,
		IdentificationDeviceUID:   true,
		IdentificationUserSID:     true,
		IdentificationUserUID:     true,
		IdentificationUserUPN:     true,
		IdentificationBinaryUID:   true,
		IdentificationPackageUID:  true,
	}

	if !validNames[id.Name] {
		return fmt.Errorf("enrichments[%d].identification[0].name has invalid value: %s", enrichmentIndex, id.Name)
	}

	return nil
}

// validateField validates a field object
func validateField(enrichmentIndex, fieldIndex int, field *Field) error {
	if field.Name == "" {
		return fmt.Errorf("enrichments[%d].fields[%d].name is required", enrichmentIndex, fieldIndex)
	}

	if field.Value == nil {
		return fmt.Errorf("enrichments[%d].fields[%d].value is required", enrichmentIndex, fieldIndex)
	}

	return nil
}
