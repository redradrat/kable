package repositories

import "encoding/base64"

var DemoRegistry = RepoRegistry{
	Repositories: DemoRepositories,
	Auths:        DemoAuths,
}

var DemoHttpsUrl = "https://github.com/redradrat/demo-concepts.git"
var DemoSshUrl = "git@github.com:redradrat/demo-concepts.git"

var DemoHttpsRepository = Repository{
	GitRepository: GitRepository{
		URL:    trimUrl(DemoHttpsUrl),
		GitRef: masterGitRef,
	},
	Name: "demo-https",
}

var DemoSshRepository = Repository{
	GitRepository: GitRepository{
		URL:    trimUrl(DemoSshUrl),
		GitRef: masterGitRef,
	},
	Name: "demo-ssh",
}

var DemoRepositories = Repositories{
	DemoHttpsRepository.Name: DemoHttpsRepository,
	DemoSshRepository.Name:   DemoSshRepository,
}

var DemoAuths = Auths{
	trimUrl(DemoHttpsUrl): Auth{
		Basic: base64.StdEncoding.EncodeToString([]byte(Authstring)),
	},
}

var Authstring = `{
	username: "testuser",
	password: "testpass",
}`
