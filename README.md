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

# Repo provider
Squirrel has pluggable storage for munki repo, specified with the `-provider` flag. 
The default provider is set to `filesystem`, which will expect a local path to munki repo.

Other supported providers are `s3` and `gcs`(google cloud storage).

See the provider specific sections below for how to connect to S3 or GCS.

## S3 provider

First, export the necessary credentials of a IAM user and region as environment variables. (You can also use the ~/.aws/credentials config file as described [here](https://github.com/aws/aws-sdk-go#aws-shared-config-file-awsconfig))

```
export AWS_ACCESS_KEY_ID=AKID1234567890
export AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
export AWS_REGION=us-east-1
```

Now serve squirrel. Use the AWS bucket name for the `-repo` flag.
```
squirrel serve -repo=awsbucketname -tls-domain=munki.corp.example.com --basic-auth=CHANGEME -provider=s3
```

To try the config locally on port 8080, you can run

```
squirrel serve -basic-auth="CHANGEME" -repo=awsbucketname -tls=false -provider=s3
```
Which will make your munki repo available at `http://localhost:8080/repo/`.
Go to `http://localhost:8080/repo/catalogs/all` to get a list of available credentials.

## Google Cloud Storage Provider 

To use squirrel with GCS, you'll need a GCP service account file. 
```
squirrel serve \
    -repo=gcsbucketname \
    -tls-domain=munki.corp.example.com \
    -basic-auth=CHANGEME \
    -provider=gcs \
    -gcs-credentials /Users/groob/Downloads/groob-gcs-credentials.json 
```


To try the config locally on port 8080, you can run

```
squirrel serve \
    -repo=gcsbucketname \
    -tls=false \
    -basic-auth=CHANGEME \
    -provider=gcs \
    -gcs-credentials /Users/groob/Downloads/groob-gcs-credentials.json 
```

Which will make your munki repo available at `http://localhost:8080/repo/`.
Go to `http://localhost:8080/repo/catalogs/all` to get a list of available credentials.

--
squirrel icon by [Agne Alesiute](https://thenounproject.com/search/?q=squirrel&i=190468) from the [Noun Project](https://thenounproject.com/).
