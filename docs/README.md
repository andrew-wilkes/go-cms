# Introduction

This project contains the Go language source code for the backend of a content management system (CMS) such as might be used for websites and blogs.

It is designed to be simple, but implement the most useful aspects of a CMS. The data is stored in flat files in JSON or HTML format. One user is supported who may log in and out to perform admin features.

The front end will be a Single Page App (SPA) that connects via AJAX requests with JSON encoded data.

The compiled app will be copied to the server, and will be able to power multiple websites. The app and the data are kept in a seperate folder to that of the website root.

The website will use an **index.php** file to interface to the app.

## Deployed Project Structure

Here is an example of the file structure of a deployed procect with one domain:

    app_folder/app
        domain.tld/
            data/
                pages.json
                status.json
            templates/
                post.html
                page.html
                header.html
                footer.html
                .
                .
            pages/
                1.html
                2.html
                .
                .
    web_folder/
        domain.tld/
            index.php
            js/
                script.js
                .
                .
            css/
                main.css
                .
                .
            images/
                logo.png
                .
                .
            robots.txt

## Installation

Some configuration options need to be set in a **config.toml** file. An example file is provided.

The **Makefile** contains commands to run scripts to:
- build the code to place the app in the **/build** folder
- deploy the app and other files to the server
- copy generated page files for test purposes to the server

The scripts folder contains php scripts that are used for the above tasks. The deploy script has a command line option (-l or --log) to deploy a test version of the installation.

## Testing

The test version runs javascript code to make various kinds of web requests for testing and logging of the data. This data may be copied from the generated log file and used in unit tests with the go code.

Some unit tests scripts generate dummy page data which may be copied to the server for integration testing.

## Index page

The server needs to be configured to direct all requests to non-existant files or the home page to the **index.php** script when the website domain is accessed in a browser.

The **index.php** captures all the server vars containing the requested page URI, method (POST/GET), and data. It then runs the Go app to process this data. The Go app then responds back to the PHP with response headers and page content.

During the deployment phase, the path to the Go app is baked into the PHP script.

*More to follow ...*
