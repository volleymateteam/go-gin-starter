# deployment-guide.md

## Go Backend Deployment (EC2)

### 1. Clone and Build

```bash
cd ~/volleymate-backend-go
go build -o volleymate-backend main.go
```

### 2. Configure systemd

`/etc/systemd/system/volleymate-backend.service`

There are two methods to configure environment variables:

#### Method 1: Using EnvironmentFile (Recommended)

```ini
[Unit]
Description=Volleymate Go Backend
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/volleymate-backend-go
ExecStart=/home/ubuntu/volleymate-backend-go/volleymate-backend
EnvironmentFile=/home/ubuntu/volleymate-backend-go/.env.prod
Restart=always
RestartSec=3
Environment=ENV=prod

[Install]
WantedBy=multi-user.target
```

#### Method 2: Direct Environment Variables (Alternative)

```ini
[Unit]
Description=Volleymate Go Backend
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/volleymate-backend-go
Environment=SCOUT_CLOUDFRONT_DOMAIN=d1b5o37qbj029k.cloudfront.net
Environment=VIDEO_CLOUDFRONT_DOMAIN=d2qk7473q3y7pg.cloudfront.net
Environment=ASSET_CLOUDFRONT_DOMAIN=d1b5o37qbj029k.cloudfront.net
Environment=ENV=prod
ExecStart=/home/ubuntu/volleymate-backend-go/volleymate-backend
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

### 3. Reload and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable volleymate-backend.service
sudo systemctl start volleymate-backend.service
```

### 4. Test

```bash
curl http://localhost:8000/health
```

### 5. Troubleshooting

#### Environment Variables

1. Verify environment variables are loaded:

   ```bash
   sudo systemctl show -p Environment volleymate-backend.service
   ```

2. If variables aren't loading from EnvironmentFile:
   - Check file permissions: `sudo chmod 644 .env.prod`
   - Or use Method 2 with direct Environment variables in service file

#### Common Issues

- "no Host in request URL": Check CloudFront domain environment variables
- Service not starting: Check logs with `sudo journalctl -u volleymate-backend.service`
- Environment variables not loading: Verify .env.prod file exists and has correct permissions

---

## Deployment Workflow (Updates)

### 1. Update from Git

```bash
cd ~/volleymate-backend-go
git checkout main
git pull origin main
```

### 2. Set Environment and Build

```bash
export ENV=prod
source .env.prod
go build -o volleymate-backend main.go
```

### 3. Restart Service

```bash
sudo systemctl restart volleymate-backend.service
```

---

## Python Microservice (Parser)

### 1. Virtualenv

```bash
cd ~/scout_parser_microservice
python3 -m venv env
source env/bin/activate
pip install -r requirements.txt
```

### 2. Service File

`/etc/systemd/system/scout-parser.service`

```ini
[Unit]
Description=Volleymate Scout Parser Microservice
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/scout_parser_microservice
ExecStart=/home/ubuntu/scout_parser_microservice/env/bin/uvicorn main:app --host 127.0.0.1 --port 8001
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

### 3. Enable and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable scout-parser.service
sudo systemctl start scout-parser.service
```

### 4. Confirm

```bash
curl http://127.0.0.1:8001/health
```

---

## 🎥 S3 Video Access Configuration (CORS & Bucket Policy)

This section documents how to configure AWS S3 and CORS to allow secure video access from the Volleymate mobile app.

### 1. Bucket Policy for CloudFront and App Access

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowCloudFrontServicePrincipal",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudfront.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::volleymate-storage/*",
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudfront::863518411349:distribution/E2I1LDQ5PKDKHX"
        }
      }
    },
    {
      "Sid": "AllowAppAccess",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::volleymate-storage/*",
      "Condition": {
        "StringLike": {
          "aws:Referer": [
            "http://localhost:8081",
            "capacitor://localhost",
            "http://localhost",
            "http://localhost:3000"
          ]
        }
      }
    }
  ]
}
```

### 2. CORS Configuration

```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "HEAD"],
    "AllowedOrigins": [
      "http://localhost:8081",
      "capacitor://localhost",
      "http://localhost",
      "http://localhost:3000"
    ],
    "ExposeHeaders": [],
    "MaxAgeSeconds": 3000
  }
]
```
