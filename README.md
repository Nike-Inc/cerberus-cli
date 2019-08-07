# cerberus-cli
A CLI for the Cerberus API.

**Table of Contents**

1. [Installing](#installing)
    1. [Mac](#installing-mac)
    1. [Linux](#installing-linux)
1. [Commands](#commands)
    1. [Help](#commands-help)
    1. [Version](#commands-version)
    1. [Secret](#commands-secret)
    1. [File](#commands-file)
    1. [SDB](#commands-sdb)
    1. [Admin](#commands-admin)
    1. [Logout](#commands-logout)
1. [Authentication](#auth)
1. [Configuration](#configuration)
    1. [Bash Completion](#configuration-autocomplete)

<a name="installing"></a>
## Installing
<a name="installing-mac"></a>
### Mac
We recommend installing `cerberus` via [Homebrew](https://brew.sh).

#### Homebrew
1. [Install Homebrew](https://brew.sh)
1. Add this tap to Homebrew:
	
	    $ brew tap nike-inc/nike
		
1. Install `cerberus`:
	* If you use bash and would like bash completion:
		```
		$ brew install bash-completion
		$ brew install cerberus-cli --with-completion
		```    
			
		Make sure to follow the caveat displayed after installing `bash-completion` by adding this line to your
		`~/.bash_profile`:
		```
		[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"
		```

	* If you don't use bash or would not like bash completion:
		```
		$ brew install cerberus-cli
		```
		
1. Upgrade `cerberus` if needed with:
        
	    $ brew upgrade cerberus-cli
        
    or
    
	    $ brew upgrade cerberus-cli --with-completion

#### Without Homebrew
1. Download a darwin-amd64 [library](https://github.com/Nike-Inc/cerberus-cli/releases/latest).
1. Give the library executable permissions.

    Example:
    
		$ chmod +x ~/Downloads/cerberus-cli-darwin-amd64
		
1. Create a softlink with the downloaded library and a location on your `PATH`, appending the location with `cerberus`.

    Example:
   
		$ ln -s ~/Downloads/cerberus-cli-darwin-amd64 /usr/local/bin/cerberus
		
1. Verify that the `cerberus` command is installed to your path by trying `$ cerberus help`.
1. (Optional) Add a few lines to your `bash_profile` to set any environment variables used by `cerberus`.

    Example:
		
		export CERBERUS_REGION=us-west-2
		export CERBERUS_URL=https://test.cerberus.example.com
		export CERBERUS_EDITOR=code
		
<a name="installing-linux"></a>
### Linux

You can install the lib with our install script (requires `jq` and `curl`)

```sh
curl -s https://raw.githubusercontent.com/Nike-Inc/cerberus-cli/master/install-cerberus-cli-linux.sh | sudo sh
```

Alternatively you can always go to [the latest release page](https://github.com/Nike-Inc/cerberus-cli/releases/latest) and download the linux release and install manually.

#### Docker example

```sh
FROM alpine:latest

RUN apk --no-cache add curl jq
RUN curl -s https://raw.githubusercontent.com/Nike-Inc/cerberus-cli/master/install-cerberus-cli-linux.sh | sh
```

<a name="commands"></a>
## Commands

A list of all commands can be viewed by using the help flags from the root command:

`$ cerberus -h` or `$ cerberus --help`

More details on each individual command can also be displayed by using the same flags:

Example: `$ cerberus secret -h` or `$ cerberus secret --help`

<a name="commands-help"></a>
### Help

Outputs a help screen that displays all possible commands and flags

`$ cerberus help`

<a name="commands-version"></a>
### Version

Outputs the current version of the project

`$ cerberus version`

<a name="commands-secret"></a>
### Secret

Displays all possible commands that can be performed on secrets

`$ cerberus secret`

* #### Read
    Given a complete secure data path, output JSON format of secret to terminal, which can be easily piped into other
    tools like [jq](https://stedolan.github.io/jq/).
    
    Example: `$ cerberus secret read app/mysdb/mysecret` might output
    ```
    {
        "foo": "bar",
        "asdf": "1234"
    }
    ```
    and `$ cerberus secret read app/mysdb/mysecret | jq -r ".foo"` would output `bar`.
        
* #### Write
    Supply a complete secure data path, along with entries in the format of `KEY=VALUE` using the `-e, --entry` flag, 
    to write secrets. The secure data path can already exist in an SDB, or can be a completely new path. If an entry's 
    key already exists in the path, it will be overwritten with this command. A success/failure message will be 
    displayed in the terminal.
    
    Example: `$ cerberus secret write app/mysdb/mysecret -e username=foo -e password=bar`
    
* #### Delete
    Given a complete secure data path, delete the corresponding secret if it exists.
    
    Example: `$ cerberus secret delete app/mysdb/mysecret`
    
* #### Edit
    Given a complete secure data path, temporarily download a secret if it exists, open preferred editor, and upload 
    edits to the same secure data path. Preferred editor can be set using the `CERBERUS_EDITOR` environment variable
    or with `-e, --editor` flags.
    
    When a secret is downloaded, the editor will open a `.yaml`file with the key value pairs of the secret. Make any
    necessary edits using `yaml` formatting. If any errors occur while parsing the edited secret or uploading to
    Cerberus, you will be prompted to open the temporary file again to fix the issue and try uploading again.
    
    Example: `$ cerberus secret edit app/mysdb/mysecret -e atom`
    
    If a path is given that does not already exist in Cerberus, you will be prompted instead to create a new secret
    at that secure data path.

<a name="commands-file"></a>
### File

Displays all possible commands that can be performed on files

`$ cerberus file`

* #### Read
    Outputs content of a file to terminal, provided a complete secure file path.
    
    Example: `$ cerberus file read app/mysdb/myfile.txt`
    
* #### Download
    Downloads a specific file, provided a complete secure file path. Default download directory is the working
    directory, or supply a complete local filepath to download to with `-o, --output` flags.
    
    Example: `$ cerberus file download app/mysdb/myfile.txt --output ~/Downloads/myfile.txt`
        
* #### Edit
    Temporarily download a file, open preferred editor, and upload edits to same secure file path. Preferred editor can
    be set using the `CERBERUS_EDITOR` environment variable, or with `-e, --editor` flags.
    
    Example: `$ cerberus file edit app/mysdb/myfile.txt -e atom` will open the file in Atom and upload edits after
    the file has been saved and closed.

* #### Upload
    Upload a local file to a specified complete secure file path. If the secure file path already exists in Cerberus,
    then the new file will replace the existing one. Otherwise, the local file will be simply uploaded to the secure
    file path. This command takes two required arguments in this order: destination secure file path, and local path to
    source file. A success/failure message will be displayed in the terminal.
    
    Example: `$ cerberus file upload app/mysdb/myfile.txt ~/Desktop/myfile.txt`
    
* #### Delete
    Delete a specific file, provided a complete secure file path. A success/failure message will be displayed in the
    terminal.
    
    Example: `$ cerberus file delete app/mysdb/myfile.txt`	
    
<a name="commands-sdb"></a>
### SDB

Displays all possible commands that can be performed on SDBs

`$ cerberus sdb`

* #### Create
    Create a new SDB. The following flags must be specified: `-n, --name`, `-o, --owner`, and `-c, --category`. 
    Additional flags (`-d, --description`, `-g, --usergroup`, and `-i, --iam`) may also be used.
    
    Use the `-g/-i` flags for each user group permission or IAM Principal permission to add, in the required format of 
    `<NAME>,<ROLE>`. `ROLE` can be `read`, `write`, or `owner`. A success/failure message will be displayed in the 
    terminal.
    
    Example: `$ cerberus sdb create -n mysdb -o Lst.MyTeam -c app`
    
    Example: `$ cerberus sdb create -n mysdb -o Lst.MyTeam -c app -d "SDB for my app" -g Lst.MyTeam,read 
    -i arn:aws:iam::012345678910:role/EXAMPLE.SSO.PowerRole,write`
    
* #### Delete
    Delete an existing SDB. Supply the path of the SDB to delete. A success/failure message will be displayed in the
    terminal.
    
    Example: `$ cerberus sdb delete app/mysdb`

<a name="commands-admin"></a>
### Admin

Displays all possible commands that can be performed as an admin

`$ cerberus admin`

* #### Override SDB Owner
    Override the owner of an existing SDB to a new owner. The following flags must be specified: the name of the sdb
    with `-s, --sdb` and the name of the new owner with `-o, --owner`. Current metadata of the SDB will be displayed,
    as a well as a prompt confirming the change in ownership.
    
    Example: `$ cerberus admin override-owner -s mysdb -o Lst.MyTeam`

<a name="commands-logout"></a>
### Logout

Removes any existing authentication tokens from the [keyring](#auth-notes).

`$ cerberus logout`

<a name="auth"></a>
## Authentication

1. Set the `CERBERUS_REGION` environment variable, or use the `-r, --region` flags.
1. Set the `CERBERUS_URL` environment variable, or use the `-u, --url` flags.

Example: `$ cerberus -r us-west-2 -u https://test.cerberus.example.com`

<a name="auth-notes"></a>
#### Notes:

`cerberus` uses a keyring to store authentication tokens after an authentication attempt is successful. After an 
initial successful authentication, the corresponding token will be used for authentication until the token's validity
expires. This keyring is supported by Linux (dbus), OS X, and Windows. If you want to remove any stored tokens, use the
[logout](#commands-logout) command.

<a name="configuration"></a>
## Configuration
In addition to setting `CERBERUS_URL` and `CERBERUS_REGION` environment variables for authentication,
you can set your preferred editor with the `CERBERUS_EDITOR` environment variable for use with the `file edit` and 
`secret edit` commands. Some good editors to use are `atom`, `subl`, and `code`, provided these shell commands are
installed.

<a name="configuration-autocomplete"></a>
### Bash Completion
If you use a bash shell, adding bash completion to the `cerberus` command can be done by downloading the
`cerberus-completion.sh` script, and adding the following line to your `~/.bash_profile`:

Example: `source ~/Downloads/cerberus-completion.sh`

Note: the `cerberus` command must be [installed](#installing) for bash completion to work.
