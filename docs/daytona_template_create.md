## daytona template create

Create a workspace template

```
daytona template create [flags]
```

### Options

```
      --builder BuildChoice          Specify the builder (currently auto/devcontainer/none)
      --custom-image string          Create the workspace with the custom image passed as the flag value; Requires setting --custom-image-user flag as well
      --custom-image-user string     Create the workspace with the custom image user passed as the flag value; Requires setting --custom-image flag as well
      --devcontainer-path string     Automatically assign the devcontainer builder with the path passed as the flag value
      --env stringArray              Specify environment variables (e.g. --env 'KEY1=VALUE1' --env 'KEY2=VALUE2' ...')
      --git-provider-config string   Specify the Git provider configuration ID or alias
      --label stringArray            Specify labels (e.g. --label 'label.key1=VALUE1' --label 'label.key2=VALUE2' ...)
      --manual                       Manually enter the Git repository
      --name string                  Specify the workspace template name
```

### Options inherited from parent commands

```
      --help   help for daytona
```

### SEE ALSO

* [daytona template](daytona_template.md)	 - Manage workspace templates

