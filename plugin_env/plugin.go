package plugin_env

import (
	"context"
	"dhswt.de/drone-github-extensions/shared"
	"github.com/Masterminds/semver"
	"github.com/drone/drone-go/plugin/environ"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// New returns a new secret plugin.
func New(config *shared.AppConfig) environ.Plugin {
	return &plugin{
		config: config,
	}
}

type plugin struct {
	config *shared.AppConfig
}

func (p *plugin) List(ctx context.Context, req *environ.Request) ([]*environ.Variable, error) {
	logrus.Infof("[env] request for build=%d %s/%s commit=%s", req.Build.ID, req.Repo.Namespace, req.Repo.Name, req.Build.After)

	logrus.Debugf("environment plugin request received: build=%+v repo=%+v", req.Build, req.Repo)

	envVars := []*environ.Variable{}

	if p.config.EmulateCIPrefixedVariables {
		ciVariables := []*environ.Variable{
			// mirror various gitlab CI_ variables
			{Name: "CI_PROJECT_NAME", Data: req.Repo.Name, Mask: false},
			{Name: "CI_PROJECT_NAMESPACE", Data: req.Repo.Namespace, Mask: false},

			//{Name: "CI_REGISTRY", Data: p.giteaDockerRegistryHost, Mask: false},
			//{Name: "CI_REGISTRY_IMAGE", Data: p.giteaDockerRegistryHost + "/" + req.Repo.Namespace + "/" + req.Repo.Name, Mask: false},
			//{Name: "CI_REGISTRY_USER", Data: token.Name, Mask: false},
			//{Name: "CI_REGISTRY_PASSWORD", Data: token.Token, Mask: true},
		}
		envVars = append(envVars, ciVariables...)
	}

	if p.config.EnvAddTagSemver && strings.HasPrefix(req.Build.Ref, "refs/tags/") {
		tag := strings.TrimPrefix(req.Build.Ref, "refs/tags/")
		v, err := semver.NewVersion(tag)
		if err != nil {
			logrus.Debugf("failed to ref as semver: %s", req.Build.Ref)
		}

		semverVars := []*environ.Variable{
			// mirror various gitlab CI_ variables
			{Name: "SEMVER_MAJOR", Data: strconv.FormatInt(v.Major(), 10), Mask: false},
			{Name: "SEMVER_MINOR", Data: strconv.FormatInt(v.Minor(), 10), Mask: false},
			{Name: "SEMVER_PATCH", Data: strconv.FormatInt(v.Patch(), 10), Mask: false},
			{Name: "SEMVER_PRERELEASE", Data: v.Prerelease(), Mask: false},
			{Name: "SEMVER_METADATA", Data: v.Metadata(), Mask: false},
			{
				Name: "SEMVER_MAJOR_MINOR",
				Data: strconv.FormatInt(v.Major(), 10) +
					"." + strconv.FormatInt(v.Minor(), 10),
				Mask: false,
			},
			{
				Name: "SEMVER_MAJOR_MINOR_PATCH",
				Data: strconv.FormatInt(v.Major(), 10) +
					"." + strconv.FormatInt(v.Minor(), 10) +
					"." + strconv.FormatInt(v.Patch(), 10),
				Mask: false,
			},
		}
		envVars = append(envVars, semverVars...)
	}

	return envVars, nil
}
