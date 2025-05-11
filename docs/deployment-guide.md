# deployment-guide.md

## Go Backend Deployment (EC2)

### 1. Clone and Build

```bash
cd ~/volleymate-backend-go
go build -o volleymate-backend main.go
```

### 2. Set Environment File

```bash
cp .env.prod .env
```

### 3. Configure systemd

`/etc/systemd/system/volleymate-backend.service`

```ini
[Unit]
Description=Volleymate Go Backend
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/volleymate-backend-go
ExecStart=/home/ubuntu/volleymate-backend-go/volleymate-backend
EnvironmentFile=/home/ubuntu/volleymate-backend-go/.env
Restart=always
RestartSec=3
Environment=ENV=prod

[Install]
WantedBy=multi-user.target
```

### 4. Reload and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable volleymate-backend.service
sudo systemctl start volleymate-backend.service
```

### 5. Test

```bash
curl http://localhost:8000/health
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
