<?php
$_SERVER['RAW_DATA'] = file_get_contents("php://input");
$_SERVER['POST_DATA'] = (object)$_POST;
$_SERVER['GET_DATA'] = (object)$_GET;

file_put_contents(sprintf("log/args%s.json", $_SERVER['REQUEST_TIME_FLOAT']), json_encode($_SERVER));

$cmd = sprintf("server_script_path/gocms %s", json_encode($_SERVER));
exec($cmd, $output, $code);
if ($code == 0) {
  $response = json_decode($output[0]);
  foreach ($response->Headers as $header)
    header($header);
  echo $response->Content;
}
else
  echo "Error code: " . $code;
?>
<!DOCTYPE html>
<html>
  <head>
    <title>Go CMS Test Page</title>
    <meta charset="utf-8">
    <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
    <script src="https://code.jquery.com/jquery-1.12.4.min.js"></script>
  </head>
  <body>
    <script>
      // Do a GET request for a user with a given ID
      axios.get('/user?ID=12345')
        .then(function (response) {
          // handle success
          console.log(response);
        })
        .catch(function (error) {
          // handle error
          console.log(error);
        })
        .then(function () {
          // always executed
        });

      // Do a POST request with application/json content
      axios.post('/user', {
        firstName: 'Fred',
        lastName: 'Flintstone'
      })
      .then(function (response) {
        console.log(response);
      })
      .catch(function (error) {
        console.log(error);
      });

      // Do a POST of regular form data
      $.post("/jquery-user", { name: "John", time: "2pm" });

      // Do a multi-part form data submission such as used for file uploads
      let data = new FormData();
      data.append('binary', Uint8Array.from('test', c => c.charCodeAt(0)));
      axios.post('/form-user', data)
        .then(function (response) {
          console.log(response);
        })
        .catch(function (error) {
          console.log(error);
        });
    </script>
  </body>
</html>