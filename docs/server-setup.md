# server-setup.md

## EC2 Setup

### 1. Instance Type

- Ubuntu 24.04 LTS
- t3.medium

### 2. Security Groups

- Allow 22 (SSH), 80 (HTTP), 443 (HTTPS)
- Parser microservice runs on `127.0.0.1:8001` only (internal only)

### 3. nginx Reverse Proxy

`/etc/nginx/sites-available/volleymate_api`

```nginx
  GNU nano 7.2                                          /etc/nginx/sites-available/volleymate_api
server {
    listen 80;
    server_name api.volleymate.app;
    client_max_body_size 5000M;

    # Add timeout directives
    proxy_read_timeout 7200;
    proxy_connect_timeout 7200;
    proxy_send_timeout 7200;
    client_body_timeout 7200;
    keepalive_timeout 7200;

    # Redirect all HTTP to HTTPS
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name api.volleymate.app;
    client_max_body_size 5000M;

    # Add timeout directives
    proxy_read_timeout 7200;
    proxy_connect_timeout 7200;
    proxy_send_timeout 7200;
    client_body_timeout 7200;
    keepalive_timeout 7200;


    ssl_certificate /etc/letsencrypt/live/api.volleymate.app/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.volleymate.app/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    location /api/v1/ {
        proxy_pass http://127.0.0.1:8000/api/v1/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://127.0.0.1:8000/health;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. Setup SSL (Let’s Encrypt)

```bash
sudo certbot --nginx -d api.volleymate.app
```

### 5. Log Locations

- Backend logs: `/home/ubuntu/volleymate-backend-go/backend.log`
- Parser logs: `journalctl -u scout-parser.service -f`
