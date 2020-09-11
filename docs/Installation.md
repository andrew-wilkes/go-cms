# Installation

## Configuration
Some configuration options need to be set in a **config.toml** file. An example file is provided.


    [website]
    web_root = "/var/www/gocms.com"
    domain = "gocms.com"

    [server]
    server_script_path = "/var/www/apps/gocms"

The `web_root` is the document root for the website.

The `domain` value is used to form URLs, so should be the same as the website domain name. It is also used in the file path to data, to seperate the data of different target domains.

The `server_script_path` is that of the application, and is embedded in the **index.php** file in order to run the App.

## Server
The server will be running say Apache, IIS, or NGINX. Also, PHP should be installed since we have based this project off a PHP entry point to simply handle HTTP requests and interface with the server.

The server should be configured to send all URIs to index.php unless a file exists to match the URI.

Example URI passed to **index.php**: `/blog/post-1`

Finally, https should be configured to encrypt the registration and login details that are sent from client to server.

## Scripts
The **Makefile** contains commands to run scripts to:
- build the code and to place the App in the **/build** folder
- deploy the App and other files (such as templates) to the server
- copy generated page files and data for test purposes to server

The copy, move, and delete operations are done on the local development computer/server rather than over a network. Folder permissions need to be set up correctly for this to work such as the user being in the same group as the document root for the website.

The `scripts` folder contains php scripts that are used for the above tasks. The deploy script has a command line option (-l or --log) to deploy a test version of the installation.

The available **make** commands are:
- make b
- make deploy
- make copyfiles

Note that: *b* is for *build*

The workflow for updating the website code is:

    make b
    make deploy

The `make copyfiles` command will be explained in the Testing section.
