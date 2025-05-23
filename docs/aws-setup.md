# AWS Infrastructure Setup for Volleymate

This document outlines the finalized AWS architecture and configuration decisions made for the Volleymate platform. Each decision aims to optimize cost, performance, and maintainability.

⸻

✅ Completed Steps

1. S3 Lifecycle Rule
   • Rule: Move raw videos to Glacier Deep Archive after 30 days.
   • Why: Saves up to 80% storage cost.
   • Path: videos/.../raw/
   • Status: Enabled via Lifecycle rule with prefix videos/.

2. EC2 Instance Configuration
   • Instance type: t3.small
   • Region: eu-north-1 (Stockholm — cheapest in EU)
   • Storage: 30 GB gp3
   • Purpose: Runs Go Backend + FastAPI Parser
   • Status: ✅ Running with working HTTPS and services

3. SSL with Let’s Encrypt
   • Domain: <https://api.volleymate.app>
   • Tool: Certbot
   • Location: /etc/letsencrypt/live/api.volleymate.app/
   • NGINX: Proxy + SSL Config Enabled
   • Status: ✅ Confirmed live

4. CloudFront Media Delivery
   • Domain: <https://d2qk7473q3y7pg.cloudfront.net>
   • Used for: Videos + Thumbnails
   • Status: Already integrated into upload pipeline

5. AWS Budgets + Alerts
   • Budget: $50/month
   • Alerts:
   • 50% Usage (Warn)
   • 80% Usage (Alert)
   • Status: ✅ Created with monthly recurring rule

⸻

📦 Cost Estimate (Monthly)

Item Estimated Cost (100 matches/mo)
S3 Storage ~$20 (Standard + Glacier after 30d)
CloudFront ~$5–10 for ~50–100 GB usage
EC2 (t3.small) ~$12
SQS <$1
Route 53 ~$1
Total $40–55 / month

⸻

🔜 Next Steps (Future Improvements)
• Add Cache-Control headers for media
• Integrate CloudWatch custom metrics (SQS queue length, parser performance)
• Enable CDN caching rules for thumbnails
• Automate lifecycle policies via Terraform (optional)
• Set up CloudWatch alerts for S3 usage or EC2 CPU/memory

⸻

🤖 Notes
• .dvw scout files are NOT archived — always accessible.
• Multiple compressed video formats are kept: 1080p, 720p, 480p
• Default returned to frontend: video_url (720p), plus video_urls map for UI flexibility

⸻

Last updated: 2025-05-23
