package release_test

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ChartTestSuite struct {
	suite.Suite
}

type rt []repo.Config

func (r *rt) UnmarshalYAML(value *yaml.Node) error {
	var err error
	*r, err = repo.UnmarshalYAML(value)

	return err
}

func (s *ChartTestSuite) SetupSuite() {
	var rs rt
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Len(rs, 1)

	s.Require().NoError(plan.SyncRepositories([]repo.Config(rs)))
}

func (s *ChartTestSuite) TestLocateChartLocal() {
	tmpDir := s.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = filepath.Join(tmpDir, "blabla")

	c, err := rel.GetChart()
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to locate chart")
	s.Require().Nil(c)
}

func (s *ChartTestSuite) TestLoadChartLocal() {
	tmpDir := s.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = tmpDir

	c, err := rel.GetChart()
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to load chart")
	s.Require().Contains(err.Error(), "Chart.yaml file is missing")
	s.Require().Nil(c)
}

//nolint:paralleltest // uses helm repository.yaml flock
func TestChartTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(ChartTestSuite))
}
