package meta

var version = "dev"

const serviceName = "xyzservice"

// Version tag
func Version() string {
	return version
}

// ServiceName returns the service name
func ServiceName() string {
	return serviceName
}
