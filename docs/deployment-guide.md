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
