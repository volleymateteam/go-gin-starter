# Volleymate AWS Infrastructure Setup

This document tracks the current AWS architecture and services used for the Volleymate backend (`go-volleymate`). It includes video processing infrastructure, cost-saving decisions, and future plans.

---

## ✅ Services in Use

| AWS Service | Purpose                                                                         |
| ----------- | ------------------------------------------------------------------------------- |
| S3          | Stores raw videos, compressed versions, thumbnails, scout files (JSON and .dvw) |
| SQS         | Queues video processing jobs (asynchronous)                                     |
| CloudFront  | Serves videos and thumbnails securely via CDN                                   |
| EC2         | Runs the Go backend and FastAPI microservice                                    |
| Route 53    | Custom domain routing (`*.volleymate.app`)                                      |
| IAM         | Fine-grained permissions for `volleymate-s3-user`                               |
| CloudWatch  | Logs and monitoring for background processing                                   |

---

## 📦 S3 Folder Structure

videos/{season_year}{country}/{competition}{gender}/{match_id}/
├── raw/ # Original video (uploaded)
├── compressed/ # Default compressed video (720p)
├── compressed/1080p/
├── compressed/720p/
├── compressed/480p/
└── thumbnails/ # JPG thumbnail generated from raw video

Example full path:

videos/2024-2025_germany/bundesliga_male/2f9b3d1c…/compressed/720p/abc123.mp4

---

## 🎬 Video Processing Pipeline (Go Backend)

1. User uploads a raw match video via `PATCH /admin/matches/:id/upload-video`.
2. Raw file is stored in `videos/.../raw/`.
3. A `VideoProcessingJob` is enqueued in SQS.
4. Background worker (Go app) compresses the video to 1080p, 720p, and 480p using `ffmpeg`.
5. Thumbnail is generated and uploaded to S3.
6. Compressed URLs and thumbnail URL are saved to the `Match` record in PostgreSQL.

---

## 💸 Cost Control & Storage Plan

| Folder        | Storage Plan                       | Transition Rule        |
| ------------- | ---------------------------------- | ---------------------- |
| `raw/`        | S3 Standard → Glacier Deep Archive | After 30 days          |
| `compressed/` | S3 Standard (keep)                 | No transition needed   |
| `thumbnails/` | S3 Standard (lightweight)          | No transition needed   |
| `scout/`      | S3 Standard                        | Future: maybe compress |

---

## 💰 Cost Estimate (Monthly)

| Item           | Estimated Cost (100 matches/mo)     |
| -------------- | ----------------------------------- |
| S3 Storage     | ~$20 (Standard + Glacier after 30d) |
| CloudFront     | ~$5–10 for ~50–100 GB usage         |
| EC2 (t3.small) | ~$10–12                             |
| SQS            | <$1                                 |
| Route 53       | ~$1                                 |
| Total          | **$40–55 / month**                  |

---

## ✅ What’s Already Done

- S3 bucket and CloudFront setup
- Secure video URL structure
- Thumbnail generation and storage
- Multi-resolution encoding
- `video_urls` and `thumbnail_url` added to match API
- Purged test messages from SQS
- `/test-enqueue` removed
