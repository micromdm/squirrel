<p align="center">
<img src="squirrel.png" alt="squirrel"/><br/>
</p>
Squirrel is a Simple HTTPs server for munki.
See the [legacy](https://github.com/micromdm/squirrel/tree/legacy) branch for the API implementation. I plan on adding it back with Munki v3 support.

# Features

* [X] **Automatic HTTPS** - squirrel provides a built in Let's Encrypt Client.
* [X] Basic Authentication - Basic Auth for munki repo

# Install
Download the [latest release](https://github.com/micromdm/squirrel/releases/latest) from the release page. 

# Quickstart

```
squirrel serve -repo=/path/to/munki_repo -tls-domain=munki.corp.example.com --basic-auth=CHANGEME
```

`-repo` flag must be set to the path of a munki repository.  
`-tls-domain` flag must be set to the domain of your munki repo. This value is used by squirrel to obtain new TLS certificates.  
`-basic-auth` flag must be set to a password which will be used for authentication to the munki repo.  

Once the server starts, you will see a prompt which prints the correct Authorization header that you need to add to your munki configuration profile.

Example:
```
Authorization: Basic c3F1aXJyZWw6Q0hBTkdFTUU=
```

See `squirrel help` for full usage.

# Keep the process running with systemd

For help with systemd see the example/systemd folder.

# Enroll mac hosts:

For help enrolling macOS hosts, check out the example/profile folder.

---
squirrel icon by [Agne Alesiute](https://thenounproject.com/search/?q=squirrel&i=190468) from the [Noun Project](https://thenounproject.com/).
