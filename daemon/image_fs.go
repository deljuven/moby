package daemon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/registry"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
	"strings"
)

// ImageInspect returns image info with the specified image
func (daemon *Daemon) ImageInspect(ctx context.Context, image string) (*types.ImageInspect, error) {
	return daemon.LookupImage(image)
}

func (daemon *Daemon) getLayerDigests(ctx context.Context, encodedAuth string) ([]string, error) {
	// retrieve auth config from encoded auth
	v2MetadataService := metadata.NewV2MetadataService(daemon.distributionMetadataStore)
	return v2MetadataService.List()
}

func (daemon *Daemon) queryManifestByImage(ctx context.Context, image string, authConfig *types.AuthConfig) ([]string, error) {
	digests := make([]string, 0)

	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, err
	}
	image = reference.FamiliarString(reference.TagNameOnly(ref))
	image = strings.TrimSuffix(image, ":")
	ref, err = reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, err
	}
	repoInfo, err := daemon.RegistryService.ResolveRepository(ref)
	if err := distribution.ValidateRepoName(repoInfo.Name); err != nil {
		return nil, err
	}
	endpoints, err := daemon.RegistryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	var (
		lastErr     error
		confirmedV2 bool
		manifest    dist.Manifest
		tagOrDigest string // Used for logging/progress only
	)
	for _, endpoint := range endpoints {
		if confirmedV2 && endpoint.Version == registry.APIVersion1 {
			log.G(ctx).Debugf("Skipping v1 endpoint %s because v2 registry was detected", endpoint.URL)
			continue
		}
		repo, confirmedV2, err := distribution.NewV2Repository(ctx, repoInfo, endpoint, make(map[string][]string), authConfig, "pull")
		if err != nil {
			log.G(ctx).Warnf("Error getting v2 registry: %v", err)
			continue
		}
		if !confirmedV2 {
			continue
		}
		manSvc, err := repo.Manifests(ctx)
		if err != nil {
			continue
		}
		if tagged, isTagged := ref.(reference.NamedTagged); isTagged {
			manifest, err = manSvc.Get(ctx, "", dist.WithTag(tagged.Tag()))
			if err != nil {
				continue
			}
			tagOrDigest = tagged.Tag()
		} else if digested, isDigested := ref.(reference.Canonical); isDigested {
			manifest, err = manSvc.Get(ctx, digested.Digest())
			if err != nil {
				continue
			}
			tagOrDigest = digested.Digest().String()
		} else {
			lastErr = fmt.Errorf("internal error: reference has neither a tag nor a digest: %s", reference.FamiliarString(ref))
			log.G(ctx).Errorf("internal error: %v", lastErr)
			continue
		}

		if manifest == nil {
			lastErr = fmt.Errorf("image manifest does not exist for tag or digest %q", tagOrDigest)
			log.G(ctx).Errorf("internal error: %v", lastErr)
			continue
		}
		if m, ok := manifest.(*schema2.DeserializedManifest); ok {
			var allowedMediatype bool
			for _, t := range distribution.ImageTypes {
				if m.Manifest.Config.MediaType == t {
					allowedMediatype = true
					break
				}
			}
			if !allowedMediatype {
				configClass := "unknown"
				lastErr = fmt.Errorf("Encountered remote %q(%s) when fetching", m.Manifest.Config.MediaType, configClass)
				log.G(ctx).Errorf("internal error: %v", lastErr)
				continue
			}
			layers := m.Layers
			for _, layer := range layers {
				digests = append(digests, layer.Digest.String())
			}
			lastErr = nil
			break
		} else {
			log.G(ctx).Debugf("unsupported manifest type, only support schema 2 right now")
			continue
		}
	}

	if len(digests) == 0 {
		digests = nil
	}
	return digests, lastErr
}

func (daemon *Daemon) queryLayerDigestsByImages(ctx context.Context, images []string, encodedAuth string) (imgDigests map[string][]string, err error) {
	authConfig := &types.AuthConfig{}
	if encodedAuth != "" {
		if err := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(encodedAuth))).Decode(authConfig); err != nil {
			log.G(ctx).Warnf("invalid authconfig: %v", err)
		}
	}
	for _, image := range images {
		digests, err := daemon.queryManifestByImage(ctx, image, authConfig)
		if err != nil {
			log.G(ctx).Errorf("error occurs during querying manifest for image %v : %v", image, err)
			continue
		}
		imgDigests[image] = digests
	}
	return
}

// ImageList return all images on the underlying node
func (daemon *Daemon) ImageList(ctx context.Context) ([]types.ImageSummary, error) {
	images, err := daemon.Images(filters.NewArgs(), true, false)
	if err != nil {
		return nil, err
	}
	imgs := make([]types.ImageSummary, len(images))
	for i, img := range images {
		imgs[i] = *img
	}

	return imgs, nil
}

// GetLayers return all layers digests on the node
func (daemon *Daemon) GetLayers(ctx context.Context, encodedAuth string) ([]string, error) {
	return daemon.getLayerDigests(ctx, encodedAuth)
}

// QueryLayersByImage return layer digests of specified image on the underlying node
func (daemon *Daemon) QueryLayersByImage(ctx context.Context, image, encodedAuth string) ([]string, error) {
	authConfig := &types.AuthConfig{}
	if encodedAuth != "" {
		if err := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(encodedAuth))).Decode(authConfig); err != nil {
			log.G(ctx).Warnf("invalid authconfig: %v", err)
		}
	}
	return daemon.queryManifestByImage(ctx, image, authConfig)
}
