# Testing

There are several ways of testing available:

- unit testing of packages
- deployment of test version of index.php
- copying of test pages and data

## Unit Testing of Packages
This is the standard way to test functions in a Go project.

Some tests require data to be set up and we must first generate the test data. This is done by running `make g`

This will create `pages.json` in the `files/test/data` folder and many html content files in `files/test/pages`.

## Deployment of test version
It is annoying to test a web app by refreshing a browser with different user input and observing results in the browser. So there is a `test.html` file (in the www folder) which may be used to capture request data for a variety of scenarious, save that data, and use the data in unit tests.

## Copying of Test Pages and Data
By running the `make copyfiles` command, the test page data and content files will be copied to the server after erasing the old data and content files that were on the server before.

Now we can browse the website to check out the navigation and rendering with this dummy test site content.

## Notes
If all of the unit tests pass (although there is not all code coverage by these tests), then the only cause of failure for the website to render should be a problem in the main `gocms.go` file.

In the browser, there is likely development tools available to check Network requests to examine headers and responses. Pressing the F12 key normally activates this functionality.

Running unit tests via `debug test` allows us to set breakpoints and then examine the state of the code.

On a Linux server, we need to ensure that we have the file access permissions set up ok such as the user (you on the computer), web root folder, and web server software are in the same group since running a web app changes the group compared to running scripts from the command line.
