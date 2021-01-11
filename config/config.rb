daemon off;
# Heroku dynos have at least 4 cores.
worker_processes <%= ENV['NGINX_WORKERS'] || 4 %>;

events {
	use epoll;
	accept_mutex on;
	worker_connections <%= ENV['NGINX_WORKER_CONNECTIONS'] || 1024 %>;
}

http {
	gzip on;
	gzip_comp_level 2;
	gzip_min_length 512;

	server_tokens off;
	proxy_cache_path /tmp/cache keys_zone=cache:100m levels=1:2 inactive=600s max_size=500m;
	log_format l2met 'measure#nginx.service=$request_time request_id=$http_x_request_id';
	access_log <%= ENV['NGINX_ACCESS_LOG_PATH'] || 'logs/nginx/access.log' %> l2met;
	error_log <%= ENV['NGINX_ERROR_LOG_PATH'] || 'logs/nginx/error.log' %>;

	include mime.types;
	default_type application/octet-stream;
	sendfile off;

	# Must read the body in 5 seconds.
	client_body_timeout 5;

	upstream app_server {
		server unix:/tmp/nginx.socket fail_timeout=0;
	}

	server {
		listen <%= ENV["PORT"] %>;
		server_name _;
		keepalive_timeout 5;
		proxy_redirect off;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;

		location / {
			proxy_pass http://app_server;
		}

        # v2
		location /v2/market/info {
		    expires 1h;
		    proxy_cache cache;
            proxy_cache_valid 15m;
          	proxy_pass http://app_server;
        }

        location /v2/market/charts {
            expires 1m;
            proxy_cache cache;
            proxy_cache_valid 1m;
            proxy_pass http://app_server;
        }

        location /v2/market/ticker {
            expires 1m;
            proxy_cache cache;
            proxy_cache_valid 1m;
            proxy_pass http://app_server;
        }

        location /v2/market/tickers {
            expires 1m;
            proxy_cache cache;
            proxy_cache_valid 1m;
            proxy_pass http://app_server;
        }
	}
}
