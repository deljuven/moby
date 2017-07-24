package daemon

import (
	"testing"
	//"github.com/docker/docker/pkg/registrar"
	//"path/filepath"
	//"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/registry"
	//"fmt"
	//"github.com/docker/docker/pkg/stringid"
	"context"
)

const (
	ROOT_DIR = "/var/lib/docker"
)

func TestImageInspect(t *testing.T) {
	//id := fmt.Sprintf("d%s", stringid.TruncateID(stringid.GenerateRandomID()))
	//dir := filepath.Join(ROOT_DIR, id)
	//daemonFolder, _ := filepath.Abs(dir)
	//imageRoot := filepath.Join(ROOT_DIR, "image", "aufs")
	//distributionMetadataStore, err := metadata.NewFSMetadataStore(filepath.Join(imageRoot, "distribution"))
	daemon := &Daemon{}
	serviceOptions := registry.ServiceOptions{}
	registryService := registry.NewService(serviceOptions)
	daemon.RegistryService = registryService
	daemon.ImageInspect(context.Background(), "alpine")
}

func TestQueryLayersByImage(t *testing.T) {
	//id := fmt.Sprintf("d%s", stringid.TruncateID(stringid.GenerateRandomID()))
	//dir := filepath.Join(ROOT_DIR, id)
	//daemonFolder, _ := filepath.Abs(dir)
	//imageRoot := filepath.Join(ROOT_DIR, "image", "aufs")
	//distributionMetadataStore, err := metadata.NewFSMetadataStore(filepath.Join(imageRoot, "distribution"))
	daemon := &Daemon{}
	serviceOptions := registry.ServiceOptions{
		InsecureRegistries: []string{},
		Mirrors:            []string{},
	}
	registryService := registry.NewService(serviceOptions)
	daemon.RegistryService = registryService
	layers, _ := daemon.QueryLayersByImage(context.Background(), "alpine", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlDTHpDQ0FkU2dBd0lCQWdJQkFEQUtCZ2dxaGtqT1BRUURBakJHTVVRd1FnWURWUVFERXp0Uk5Gb3pPa2RYTjBrNldGUlFSRHBJVFRSUk9rOVVWRmc2TmtGRlF6cFNUVE5ET2tGU01rTTZUMFkzTnpwQ1ZrVkJPa2xHUlVrNlExazFTekFlRncweE56QTFNREl5TWpBME5UZGFGdzB4T0RBMU1ESXlNakEwTlRkYU1FWXhSREJDQmdOVkJBTVRPMDFPTms0NlJraFVWenBKV0VWSE9rOUpOMUU2UVRWWFJqcFpSRVUwT2pkRE4wNDZSMWRKVVRvMVZ6STNPa2hPTlVvNlZVNURRVG95U0UxQ01Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRU5KRklhQ1hHNWYxSk9BZnZSaTJDU081K1Q5RVpKd2doai9SUXgzNW9Uc3Q4RnhXY0dRc3ZOMG5sdW5DVVdIbENxN2I4NFJRTXV0WUVIUnY4MVhweTU2T0JzakNCcnpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RHdZRFZSMGxCQWd3QmdZRVZSMGxBREJFQmdOVkhRNEVQUVE3VFU0MlRqcEdTRlJYT2tsWVJVYzZUMGszVVRwQk5WZEdPbGxFUlRRNk4wTTNUanBIVjBsUk9qVlhNamM2U0U0MVNqcFZUa05CT2pKSVRVSXdSZ1lEVlIwakJEOHdQWUE3VVRSYU16cEhWemRKT2xoVVVFUTZTRTAwVVRwUFZGUllPalpCUlVNNlVrMHpRenBCVWpKRE9rOUdOemM2UWxaRlFUcEpSa1ZKT2tOWk5Vc3dDZ1lJS29aSXpqMEVBd0lEU1FBd1JnSWhBSTJVUlpMQVRTM3R4bjNpNTY0SXVQSFEwQU1Mb1g5cTZCMmdnN01KSHJuTkFpRUE0Q3lzbmtENHhjQm42amdobVdnQzczQjdGVkszenFnOTV4ZjNRK2xGVHlrPSJdfQ.eyJhY2Nlc3MiOlt7InR5cGUiOiJyZXBvc2l0b3J5IiwibmFtZSI6ImxpYnJhcnkvYWxwaW5lIiwiYWN0aW9ucyI6WyJwdWxsIl19XSwiYXVkIjoicmVnaXN0cnkuZG9ja2VyLmlvIiwiZXhwIjoxNTAwODgxOTk2LCJpYXQiOjE1MDA4ODE2OTYsImlzcyI6ImF1dGguZG9ja2VyLmlvIiwianRpIjoiRjVMcU9UM0JTZm1pQk1xUXVMS1QiLCJuYmYiOjE1MDA4ODEzOTYsInN1YiI6IiJ9.El1EDN0Xhk4Zf-oM9ij6mP7VK1T0LvnCY6bXYyAvlt7NR6mMeqqlhV2x2z9dO-c5-SgpIwCIEUhfcTkMtxs17A")
	t.Logf("%v", layers)
}
