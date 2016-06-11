Squirrel is a secure and easy to use webserver for [munki](https://github.com/munki/munki).
Squirrel is built on top of the [caddy](https://caddyserver.com/) webserver and adds munki specific features through plugins.

Below is a list of features. Some are immediately usable, while others are in various stages of completion.

# Features

* [√] *Automatic HTTPS* - squirrel provides a built in Let's Encrypt Client(through caddy). You can also provide your own certs.
* [√] *Built in [SCEP](https://tools.ietf.org/html/draft-nourse-scep-23) server* - The `scepclient` can request client certificates in a munki preflight script.
* [√] *HTTP/2* - Automatically supported by the server and NSURLSession on OS X.
* [√] *git/git-fat/lfs sync* - syncing a repo on a time interval. provided by the caddy [addon](https://caddyserver.com/docs/git)
* [In Progress] *API* - A REST API for managing a munki repo remotely. Mostly complete. Porting over from `https://github.com/groob/ape`
* [In Progress] *apiimport* - A custom `munkiimport` tool which allows importing packages using the API instead of mounting the repo.
* [In Progress] *Web UI* - A web interface for managing the munki repo. 
* [In Progress] *dynamic catalogs* - currently possible to run `makecatalogs` after a git pull, but the server will also support this feature natively.
* [In Progress] *autopromotion/sharding* - part of having dynamic catalogs. The server will allow configuration of promotion between catalogs and [sharding](http://grahamgilbert.com/blog/2015/11/23/releasing-changes-with-sharding/) support.
* [In Progress] *monitoring* - structured logging and prometheus metrics. 
* [Future] DEP/MDM integration - as [micromdm](https://github.com/micromdm/micromdm) is developed, integrations will be added where they make sense. For example - ability to create manifests or validate SCEP requests based on DEP membership.
* [Future] rsync - another way to sync a repo at an interval for those who don't use git.
* [Future] [The Update Framework](https://theupdateframework.github.io/) - investigating TUF/[notary](https://github.com/docker/notary) as a way to validate catalogs and manifests.
