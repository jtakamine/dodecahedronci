package main

type dodecPackage struct {
	Pkg               []byte
	DockerRegistryUrl string
}

var registry = make(map[string]map[string]dodecPackage)

func addPackage(app string, version string, pkg []byte, dockerRegistryUrl string) (err error) {
	if registry[app] == nil {
		registry[app] = make(map[string]dodecPackage)
	}
	registry[app][version] = dodecPackage{Pkg: pkg, DockerRegistryUrl: dockerRegistryUrl}
	return nil
}

func getPackage(app string, version string) (pkg []byte, dockerRegistryUrl string, err error) {
	dPkg := registry[app][version]
	return dPkg.Pkg, dPkg.DockerRegistryUrl, nil
}
