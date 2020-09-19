# Installation

## Server
The server will be running say Apache, IIS, or NGINX.

The server should be configured to send all URIs to the location of the App at the configured port.

Other URI's such as for images, scripts, and other static files should be mapped to their location on port 80 or the https port.

See the example [NGINX Config.](nginx.md) for how we might route specific patterns of requested resources.

Ensure that the firewall rules allow external access from clients.

Finally, https should be configured to encrypt the registration and login details that are sent from client to server.

## Scripts
The **Makefile** contains commands to run scripts to:
- build the code and to place the App in the **/build** folder
- deploy the App and other files (such as templates) to the server
- generate page html files and page data for test purposes and to fill out a dummy website
- copy generated page files and data for test purposes to server

The copy, move, and delete operations are done on the local development computer/server rather than over a network. Folder permissions need to be set up correctly for this to work such as the user being in the same group as the document root for the website.

The `scripts` folder contains scripts that are used for the above tasks. The deploy script has a command line option (-l or --log) to deploy a test version of the installation.

The available **make** commands are:
- make b
- make deploy
- make copyfiles
- make g
- make test

Note that: *b* is for *build* and *g* is for *generate*

The workflow for updating the website code is:

    make b
    make deploy

The `make copyfiles` command will be explained in the Testing section.
