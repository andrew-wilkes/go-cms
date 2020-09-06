# NGINX Server Block Configuration

Here is an example of how to set up a server block in the **nginx.conf** file or elsewhere that you may do it:

    server {
            server_name gocms.com;
            root /var/www/gocms.com;

            listen 80;

            index index.php index.html;


            location = /favicon.ico {
                log_not_found off;
                access_log off;
            }

            location = /robots.txt {
                allow all;
                log_not_found off;
                access_log off;
            }
            
            location / {
                try_files $uri $uri/ /index.php;
            }

            # static file 404's aren't logged and expires header is set to maximum age
            location ~* \.(jpg|jpeg|gif|css|png|js|ico|html)$ {
                access_log off;
                log_not_found off;
                expires max;
            }
            
            location ~ \.php$ {
                try_files  $uri =404;
                fastcgi_pass   unix:/var/run/php-fpm/php-fpm.sock;
                fastcgi_index  index.php;
                include        fastcgi.conf;
            }
        }

In this example our domain name is: **gocms. com**

Link: https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/

## Server Admin Notes

 Start Nginx Service

    sudo systemctl start nginx

Stop Nginx Service

    sudo systemctl stop nginx

Enable Nginx on boot

    sudo systemctl enable nginx

Get 502 Bad Gateway if php-fpm is not running.

    sudo systemctl start php-fpm
    sudo systemctl stop php-fpm

### PHP-fpm Settings

In */etc/php/php-fpm.d/www.conf* set:

    user = http
    group = www-data
