package dto

// RedactedAgentCatalog is an allowlist projection for selection and discovery.
// New launch configuration fields cannot enter this projection implicitly.
func RedactedAgentCatalog(agent AgentDTO) AgentDTO {
	profiles := make([]AgentProfileDTO, 0, len(agent.Profiles))
	for _, p := range agent.Profiles {
		profiles = append(profiles, RedactedProfileCatalog(p))
	}
	return AgentDTO{
		ID: agent.ID, Name: agent.Name, WorkspaceID: agent.WorkspaceID,
		SupportsMCP: agent.SupportsMCP, Profiles: profiles,
		ProfileOrderRevision: agent.ProfileOrderRevision,
		InferenceCapable:     agent.InferenceCapable, CapabilityStatus: agent.CapabilityStatus,
	}
}

// RedactedProfileCatalog is also used for notifications to mixed-authority clients.
func RedactedProfileCatalog(p AgentProfileDTO) AgentProfileDTO {
	return AgentProfileDTO{
		ID: p.ID, AgentID: p.AgentID, Kind: p.Kind, Name: p.Name,
		AgentDisplayName: p.AgentDisplayName, Model: p.Model, Mode: p.Mode,
		Enabled: p.Enabled, CLIPassthrough: p.CLIPassthrough, WorkspaceID: p.WorkspaceID,
	}
}

// RedactedAvailableCatalog excludes discovered host paths and launch commands.
func RedactedAvailableCatalog(agents []AvailableAgentDTO) []AvailableAgentDTO {
	result := make([]AvailableAgentDTO, 0, len(agents))
	for _, agent := range agents {
		result = append(result, AvailableAgentDTO{
			Name: agent.Name, DisplayName: agent.DisplayName, Description: agent.Description,
			SupportsMCP: agent.SupportsMCP, Available: agent.Available,
			Capabilities: agent.Capabilities, UpdatedAt: agent.UpdatedAt,
		})
	}
	return result
}
