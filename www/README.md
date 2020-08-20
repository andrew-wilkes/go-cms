# `/www`

## Website files

*index.php* calls the *gocms* app with arguments of GET and POST data. The app returns a response of headers and the web page content. On failure, it displays an error code.

### Error Codes
* Code 127 is a refusal to run the *gocms* app because the server user may be in another group compared to the app location
* Code 126 is due to a bug in the command issued to exec