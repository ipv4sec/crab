package crabe

type Template struct {
	Init struct {
		Parameter string
		Template  string
	}
	Workloads map[string]struct {
		Parameter string
		Template  string
	}
}

type Param struct {
	RootDomain string
	InstanceId string

	Authorization []struct {
		Service   string
		Namespace string
		Resources []struct {
			Uri    string
			Action []string
		}
	}
	ServiceEntry []struct {
		Name     string
		Host     string
		Address  string
		Port     int
		Protocol string
	}
	Configs []struct {
		Path    string
		SubPath string
		Data    []struct {
			Name  string
			Value string
		}
	}
}
