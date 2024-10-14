package revision

import "oma/contract"

type RevisionConfig struct {
	Type           contract.RevisionRepositoryType
	GitlabPackages GitlabPackagesRevisionRepositoryConfig
	OCI            OCIRevisionRepositoryConfig
}

func (c *RevisionConfig) Validate() error {
	if err := c.Type.Validate(); err != nil {
		return err
	}

	switch c.Type {
	case contract.RevisionTypeGitlabPackages:
		return c.GitlabPackages.Validate()
	case contract.RevisionTypeOCI:
		return c.OCI.Validate()
	}

	return nil
}
