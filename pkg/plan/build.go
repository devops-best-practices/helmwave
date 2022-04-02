package plan

import (
	log "github.com/sirupsen/logrus"
)

// Build plan with yml and tags/matchALL options.
func (p *Plan) Build(yml string, tags []string, matchAll bool, templater string) error {
	p.templater = templater

	// Create Body
	body, err := NewBody(p.fsys, yml)
	if err != nil {
		return err
	}
	p.body = body

	// Build Releases
	log.Info("Building releases...")
	p.body.Releases = buildReleases(tags, p.body.Releases, matchAll)
	if len(p.body.Releases) == 0 {
		return nil
	}

	// Build graphs
	log.Info("Building graphs...")
	p.graphMD = buildGraphMD(p.body.Releases)
	log.Infof("Depends On:\n%s", buildGraphASCII(p.body.Releases))

	// Build Values
	log.Info("Building values...")
	err = p.buildValues()
	if err != nil {
		return err
	}

	// Build Repositories
	log.Info("Building repositories...")
	_, err = p.buildRepositories()
	if err != nil {
		return err
	}

	// Sync Repositories
	err = SyncRepositories(p.body.Repositories)
	if err != nil {
		return err
	}

	// Build Registries
	log.Info("Building registries...")
	_, err = p.buildRegistries()
	if err != nil {
		return err
	}

	// Sync Registries
	err = p.syncRegistries()
	if err != nil {
		return err
	}

	// Build Manifest
	log.Info("Building manifests...")
	err = p.buildManifest()
	if err != nil {
		return err
	}

	return nil
}
