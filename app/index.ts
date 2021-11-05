import express from "express";

// rest of the code remains same
const app = express();
const PORT = 8000;
app.set('view engine', 'pug')
app.set('views', './views')

serveStatics(app)
servePages(app)

app.listen(PORT, () => {
    console.log(`⚡️[server]: Server is running at http://localhost:${PORT}`);
});

function serveStatics(app: express.Express) {
    app.use("/css", express.static('public/css'))
    app.use("/img", express.static('public/img'))
    app.use("/js", express.static('public/js'))
}

function servePages(app: express.Express) {
// Serve Index
    app.get('/', function (req, res) {
        res.render('index')
    })

    app.get('/repos', function (req, res) {
        let repos = [{
            name: "Elkoss Combine",
            private: true,
            url: "https://github.com/elkcom/concepts",
            ref: "refs/heads/master"
        }, {
            name: "Aldrin Labs",
            url: "https://github.com/aldrinlabs/infrastructure",
            ref: "refs/heads/master"
        }, {
            name: "Serrice Council",
            url: "https://github.com/serrice/concepts",
            ref: "refs/heads/master"
        }, {name: "Hahne Kedar", url: "https://github.com/hkmanufacturing/infra", ref: "refs/heads/master"}];

        let privrepos = 0
        let pubrepos = 0
        for (let repo of repos) {
            if (repo.private) {
                privrepos += 1
            } else {
                pubrepos += 1
            }
        }

        res.render('repos', {repos: repos, privrepos: privrepos, pubrepos: pubrepos})
    })

    app.get('/concepts', function (req, res) {
        let concepts = [{
            id: "storage_postgresql@elkcom",
            name: "PostgreSQL",
            type: "jsonnet",
            version: "1.1.0-beta4",
            maintainer: "Michele Tarantino"
        }, {id: "storage_mysql@elkcom", name: "MySQL", type: "jsonnet", version: "1.0.0", maintainer: "Trostan Mírsson"}, {
            id: "storage_redis@elkcom",
            name: "Redis",
            type: "jsonnet",
            version: "1.3.0",
            maintainer: "Mateo Valdueza"
        }, {id: "storage_memcached@aldrinlabs", name: "Memcached", type: "helm", version: "2.3.1", maintainer: "Unknown"}]

        let stableconcepts = 2
        let betaconcepts = 1
        let alphaconcepts = 0

        res.render('concepts', {
            concepts: concepts,
            stableconcepts: stableconcepts,
            betaconcepts,
            alphaconcepts: alphaconcepts
        })
    })

    app.get('/concepts/:conceptid', function (req, res) {
        let cname = req.params.conceptid

        let concept = {
            displayName: "PostgreSQL",
            maintainer: {
                name: "Ralph Kühnert",
                email: "kuehnert.ralph@gmail.com"
            },
            version: "1.1.0-beta4",
            type: "helm",
            inputs: {
                mandatory: {
                    dbname: {
                        type: "string"
                    },
                    dbuser: {
                        type: "string"
                    },
                    dbpass: {
                        type: "string"
                    },
                },
                optional: {
                    dbscheme: {
                        type: "string"
                    }
                }
            }
        }
        let conceptrepo = {
            name: "Elkoss Combine",
            id: "elkcom"
        }

        res.render('concept-details', {
            name: cname,
            concept: concept,
            repo: conceptrepo
        })
    })

    app.get('/kubeapps', function (req, res) {
        res.render('kubeapps')
    })

    app.get('/stats', function (req, res) {
        res.render('stats')
    })
}