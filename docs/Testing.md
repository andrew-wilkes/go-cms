

## Installation

Some configuration options need to be set in a **config.toml** file. An example file is provided.

The **Makefile** contains commands to run scripts to:
- build the code to place the app in the **/build** folder
- deploy the app and other files to the server
- copy generated page files for test purposes to the server

The scripts folder contains php scripts that are used for the above tasks. The deploy script has a command line option (-l or --log) to deploy a test version of the installation.

The available **make** commands are:
- make b
- make deploy
- make copy files

## Testing

The test version runs javascript code to make various kinds of web requests for testing and logging of the data. This data may be copied from the generated log file and used in unit tests with the go code.

Some unit tests scripts generate dummy page data which may be copied to the server for integration testing.

## Index page

The server needs to be configured to direct all requests to non-existant files or the home page to the **index.php** script when the website domain is accessed in a browser.

The **index.php** captures all the server vars containing the requested page URI, method (POST/GET), and data. It then runs the Go app to process this data. The Go app then responds back to the PHP with response headers and page content.

During the deployment phase, the path to the Go app is baked into the PHP script.

*More to follow ...*
