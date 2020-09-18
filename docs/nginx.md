# NGINX Server Block Configuration

Here is an example of how to set up a server block in the **nginx.conf** file or elsewhere that you may do it:
    
	server {
		server_name .gocms.com;

        # routes applied to the App
		location / {
			proxy_set_header Host $host;
			# proxy_set_header X-Real-IP $remote_addr;
			proxy_pass http://localhost:8090;
		}

        # css, js, images etc.
		location ~* \.\w+$ {
			root   /var/www/gocms/gocms.com/static/;
			access_log off;
			expires max;
		}

        # static html
		location ~* \.html$ {
			root   /var/www/gocms/gocms.com/static/;
		}

        location = /robots.txt {
            allow all;
            log_not_found off;
            access_log off;
        }

		location = /favicon.ico {
			log_not_found off;
			access_log off;
		}
	}

In this example our domain name is: **gocms.com**

Link: https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/

## Server Admin Notes

 Start Nginx Service

    sudo systemctl start nginx

Stop Nginx Service

    sudo systemctl stop nginx

Enable Nginx on boot

    sudo systemctl enable nginx
