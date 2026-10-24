package dto

// RedactedAgentCatalog is an allowlist projection for selection and discovery.
// New launch configuration fields cannot enter this projection implicitly.
func RedactedAgentCatalog(agent AgentDTO) AgentDTO {
	profiles := make([]AgentProfileDTO, 0, len(agent.Profiles))
	for _, p := range agent.Profiles {
		profiles = append(profiles, AgentProfileDTO{
			ID: p.ID, AgentID: p.AgentID, Kind: p.Kind, Name: p.Name,
			AgentDisplayName: p.AgentDisplayName, Model: p.Model, Mode: p.Mode,
			Enabled: p.Enabled, CLIPassthrough: p.CLIPassthrough,
			WorkspaceID: p.WorkspaceID,
		})
	}
	return AgentDTO{
		ID: agent.ID, Name: agent.Name, WorkspaceID: agent.WorkspaceID,
		SupportsMCP: agent.SupportsMCP, Profiles: profiles,
		ProfileOrderRevision: agent.ProfileOrderRevision,
		InferenceCapable:     agent.InferenceCapable, CapabilityStatus: agent.CapabilityStatus,
	}
}
