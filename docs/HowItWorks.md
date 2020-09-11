# How It Works
In this section we will present an overview of how the App works and why.

## End Points
The website is accessed via URLs and the part after the domain name is the URI such as `/` or `/blog` which are used as end points into the App.

In most cases, they will be requests for content such as a web page. In other cases, they will be sent from the Admin page which will be a seperate Javascript App.

The admin end points will be according to a very simple REST API as follows:

    /api/user
    /api/page
    /api/pages

The request method (GET or POST) determines the kind of action requested such as get data or save data.

## Request Handling
When a request is received by the server, the **index.php** script first processes this request.

The server global variable $_SERVER contains most of the request data and server vars except for posted data objects, so we add that data as another key value pair in the associated array (map).

    $_SERVER['RAW_DATA'] = file_get_contents("php://input");

Next, we JSON encode this data and execute the Go App with this data as a command line parameter. If all goes well, the App prints to standard output a response or to the error output if there is a panic (exception).

We add a redirection of error output to standard output on the command line so that we are able to see the error message in the web page.

With a correct response we will get zero or more headers to respond with to the client, and the HTML content to output.

## App Packages

### main
The main entry point receives the data from the **index.php** calling script.

The first command line argument contains the path to the App and we use this to set the path to our data files. This is necessary because the location of these files is different in our development environment where we run tests using dummy data.

The second command line argument contains the JSON data string. We decode this to extract the data that we are interested in with the `ProcessInput` function.

The URI is further processed with the `ParseURI` function to extract GET vars that may be present e.g. `/contact?a=2&b=3`

The data is put into a `reqest.Info` structure.

Then we pass the data to the **router** package for processing into return values of headers and content.

### request
This package simply contains an Info data structure which is depended on by other packages.

    type Info struct {
        Domain    string
        Route     string
        SubRoutes []string
        Method    string
        Scheme    string
        GetArgs   map[string]string
        PostData  map[string]json.RawMessage
    }



### router
The router package is concerned with routes, which are paths to pages. Some pages are special and utilize sub-routes that provide data to the page such as the **archive** page. The sub routes here are for year and month.

Example archive page route: `/archive/2020/03`

The `ExtractSubRoutes` function extracts the sub routes into the Info struct and returns the page route which is used to match with a page file by name. If the page is not found then an error 404 header is returned. Otherwise, we now have a page.Info struct.

Now we get the template for the page. The name of a template file is obtained from the page Info.

Next we replace tokens in the template with content.

Finally, we return the header and content values.

Some other packages were used here: files, page, and content.

### files
This package contains the `Root` path variable to the data files and a `GetTemplate` function to return page templates.

Also, the files folder where we place this package file is a convenient place to store test data files. These files are for example JSON input data, and generated web pages.

### page
This package is responsible for saving and loading page data and content. It also has functions for returning a page by ID or Route, and to provide lists of pages. Also, it has the page.Info structure definition.

### content
This package is responsible for assembling content and replacing tokens in templates with the content.

### settings
This package is responsible for CRUD operations on application settings.

### user
This package is responsible for REST interactions related to the user login between the Admin app and the main App.

The user states are:
- logged out
- loggin in
- unknown

### api
This package is responsible for REST interactions between the Admin app and the main App end points.

