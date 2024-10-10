package revision

import "oma/contract"

type RevisionConfig struct {
	Type                contract.RevisionRepositoryType
	GitlabPackages      GitlabPackagesRevisionRepositoryConfig
	PolicyProxyPackages PolicyProxyRevisionRepositoryConfig
}

func (c *RevisionConfig) Validate() error {
	if err := c.Type.Validate(); err != nil {
		return err
	}

	switch c.Type {
	case contract.GitlabPackages:
		return c.GitlabPackages.Validate()
	case contract.PolicyProxy:
		return c.PolicyProxyPackages.Validate()

	}

	return nil
}
