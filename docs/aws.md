# 📊 Volleymate Production Cost Breakdown (Estimates)

We’ll split this into 4 parts:

1️⃣ Storage (S3) — Match Videos & Scout Files

🔹 What you’re storing:

| Type               | Est. Size per File | Monthly Uploads | Total Monthly |
| ------------------ | ------------------ | --------------- | ------------- |
| Match Videos (MP4) | 500 MB – 2 GB      | 50              | \~100 GB      |
| Scout Files (.dvw) | 100 KB – 1 MB      | 500             | \~500 MB      |
| JSON Outputs       | 500 KB – 5 MB      | 500             | \~2 GB        |

🧮 Monthly S3 Estimate:
Storage: 100–200 GB

Cost: ~$2.30/month per 100GB in S3 Standard → $4.60/month

Put/List/Get Requests: negligible (under $1)

CloudFront egress bandwidth: depends on traffic

✅ Plan for 500 GB to start → ~$12/month max (S3 + CloudFront combined)

2️⃣ Compute (EC2) — Backend + Microservices
🧠 Your setup:
Service Type Instance Monthly
Go Backend (API) t3.medium 2 vCPU, 4GB ~$30–35
Python Parser (FastAPI) t3.small 2GB ~$17
Future workers/batch jobs Reserve 1 - -

🧠 You may consolidate Go + Python to same t3.medium at first for cost savings

✅ Estimate: $30–50/month depending on usage

3️⃣ Network & Traffic (CloudFront, S3 Bandwidth)
AWS gives 1 GB free per month for egress

CloudFront is cheaper than direct S3:

~$0.085 per GB (Europe region)

500 GB egress = ~$42.50

If users stream video or download data heavily, plan for ~$40–50/month in bandwidth

✅ Estimate: $50 for moderate traffic

4️⃣ Other AWS Services
Service Usage Cost
Route 53 DNS (e.g., api.volleymate.app) ~$0.50/month
ACM (SSL) Free with CloudFront $0
SES (optional) Email invites (first 62k free) $0–$5
CloudWatch Logs, monitoring ~$2–3

✅ Estimate: ~$5–10/month

📦 Mobile Store Costs
Store One-time / Annual Notes
Apple $99/year Required for App Store
Google $25 one-time Required for Play Store
App assets ~$0 Icons, splash, etc. you design
Privacy hosting ~$2/month On S3 or Vercel/Netlify

Netlify

💰 Total Monthly Estimate (Starting Point)
Category Monthly
S3 + CloudFront ~$12–15
EC2 + Python ~$30–50
Network egress ~$40–50
Extras (Route53, logs, etc.) ~$5
Total ~$90–120/month

📈 Scaling Considerations
Metric Upgrade When...
EC2 CPU/RAM Requests feel slow / memory errors
S3 Storage You hit >500GB
CloudFront bandwidth >1TB/month
Parser load Scout file queue increases
