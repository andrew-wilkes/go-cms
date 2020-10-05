# How It Works
In this section we will present an overview of how the App works and why.

## End Points
The website is accessed via URLs and the part after the domain name is the URI such as `/` or `/blog` which are used as end points into the App.

In most cases, they will be requests for content such as a web page. In other cases, they will be sent from the Admin page which will be a seperate Javascript App.

The admin end points will be according to a very simple REST API as follows:

    /api
    /api/user/logon
    /api/user/logoff
    /api/user/register
    /api/page/load
    /api/page/save
    /api/pages/load
    /api/pages/save 

The `/api` end point responds with `Msg: "register"` or `Msg: "logon"` To determine what the Dashboard should initially display to the user.

## App Packages

### main
We pass the data to the **router** package for processing into return values of headers and content.

Then we assemle the response headers, status code, and content to output as a server on port 8090.

The decision was made to not load static files from the server and output them from here. So we configure the server to react to any URI with an extension (.???) as a static file. This App ony handles routes without a file extension.

### response
This package simply contains an Info data structure for the output response data from functions feeding back to the final output response of the App.

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
This package is responsible for saving and loading page data and content. It also has functions for returning a page by ID or Route, and to provide lists of pages. Also, it has the page.Info structure definition. This includes a Format specifier which is used as the file extension. Valid values are: html, md, txt, css, js, and log.

### content
This package is responsible for assembling content and replacing tokens in templates with the content.

Template tokens are comprised of upper-case text surrounded by hashes e.g. #TITLE# which would be replaced by the page title text.

Other tokens include a data value such as a tag or digit count value. Examples are: #SIDE_MENU# and #PAGES_2# which specify the actual menu or the depth of the indented list of pages in a tree of parent-child relationships.

It's best to have a look at the `content.go` source code to see all of the available tokens and the code that produces the associated content. It includes recursive functions and slice magic.

`md` format causes the content to be transformed from Markdown (GIT flavour) to HTML.

### settings
This package is responsible for storing user credential verification data and the session id and expiry. And maybe Dashboard info such as the project name.

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
