# How It Works
In this section we will present an overview of how the App works and why.

## End Points
The website is accessed via URLs and the part after the domain name is the URI such as `/` or `/blog` which are used as end points into the App.

In most cases, they will be requests for content such as a web page. In other cases, they will be sent from the Admin page which will be a seperate Javascript App.

The admin end points will be according to a very simple REST API as follows:

    /api/user/logon
    /api/user/logoff
    /api/user/register
    /api/page/load
    /api/page/save
    /api/pages/load
    /api/pages/save 

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
This package simply contains an Info data structure for the input request data to functions which is depended on by other packages.

    type Info struct {
        Domain    string
        Route     string
        SubRoutes []string
        Method    string
        Scheme    string
        IPAddr    string
        GetArgs   map[string]string
        PostData  map[string]json.RawMessage
    }

### response
This package simply contains an Info data structure for the output response data from functions feeding pack to the final output response of the App.

    type Info struct {
        ID   string
        Data string
        Msg  string
    }

The ID is empty for failed tasks or when the session has expired, and equal to the current session ID otherwise.

The Data is for any expected response data such as for page data.

The Msg is for message text that may be displayed to the user in the client App or a simple "ok" to indicate success completing the action.

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

Template tokens are comprised of upper-case text surrounded by hashes e.g. #TITLE# which would be replaced by the page title text.

Other tokens include a data value such as a tag or digit count value. Examples are: #SIDE_MENU# and #PAGES_2# which specify the actual menu or the depth of the indented list of pages in a tree of parent-child relationships.

It's best to have a look at the `content.go` source code to see all of the available tokens and the code that produces the associated content. It includes recursive functions and slice magic.

### settings
This package is responsible for CRUD operations on application settings. At this moment, the actual settings that will be available are undefined.

### user
This package is responsible for user login and status features.

The user's Email value will be the thing that relates to stored previous status data to check if they are known or not.

They will log in with password and email identity or register their details (userName, email, password) to initialize the admin user account.

The `userName` will be applied to the #AUTHOR# token in templates.

The user data will be saved to a `data/settings.json` file.

A session ID is provided to the client on a successful login, and this must be included (`?id=value`) with each subsequent request sent to the server, otherwise the user is assumed to be unknown or logged out. The session ID has an expiration time out of 8 hours.

User actions are:
- LogOn (Compares the supplied credentials, and starts a new session)
- LogOff (Deletes the session ID)
- Register (Also has the effect of updating the user details if they are loggin in)

### api
This package is responsible for REST interactions between the Dashboard, Content Editor, and the main App end points.

In order to process actions, it will compare the session ID to that supplied as an ID from the client. If not authorized, it will prompt for log-in credentials.
