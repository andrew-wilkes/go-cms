# GO CMS

A website content management system using Go Lang.

An **index.php** page is used to interface with the Go application.

# Introduction

This project contains the Go language source code for the backend of a content management system (CMS) such as might be used for websites and blogs.

It is designed to be simple, but implement the most useful aspects of a CMS. The data is stored in flat files in JSON or HTML format. One user is supported who may log in and out to perform admin features.

The front end will be a Single Page App (SPA) that connects via AJAX requests with JSON encoded data. It will use a REST API.

The compiled app will be copied to the server, and will be able to power multiple websites. The app and the data are kept in a seperate folder to that of the website root.

## Deployed Project Structure

Here is an example of the file structure of a deployed project with one domain:

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
            dashboard/
                front_end_SPA_code

The *dashboard* folder name will be customizable.

# Documentation

* [Installation](docs/Installation.md)
* [How It Works](docs/HowItWorks.md)
* [Testing](docs/Testing.md)
