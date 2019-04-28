package phases

func (s Service) Apply(ctx SystemContext) {
}
func (k Kubernetes) Apply(ctx SystemContext) {
}
func (runtime ContainerRuntime) Apply(ctx SystemContext) {

}
func (c Container) Apply(ctx SystemContext) {
}

func (u User) Apply(ctx SystemContext) {

}

func (cfg SystemConfig) Apply(ctx SystemContext) {
	for _, container := range cfg.Containers {
		container.Apply(ctx)
	}
	for _, user := range cfg.Users {
		user.Apply(ctx)
	}

	for name, svc := range cfg.Services {
		svc.Name = name
		svc.Apply(ctx)
	}

	if cfg.ContainerRuntime != nil {
		cfg.ContainerRuntime.Apply(ctx)
	}

	if cfg.Kubernetes != nil {
		cfg.Kubernetes.Apply(ctx)
	}

	// cfg.ToScript()
	//	cfg.ToFiles()

}
