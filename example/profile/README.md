To enroll a mac into the `squirrel` server, install an updated version of this profile on the mac.

First, you'll have to make a few changes. Open the profile in your text editor andupdate the values as needed.

You _must_ set the correct value under 
```
<array>
    <string>Authorization: Basic CHANGEME</string>
</array>
```
The correct header will be printed when you run the `squirrel serve` command to start the server.


And SoftwareRepoURL:
```
<key>SoftwareRepoURL</key>
<string>https://munki.corp.micromdm.io/repo</string>
```

Note that `/repo` is a required path for squirrel.


