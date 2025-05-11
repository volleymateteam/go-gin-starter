# Volleymate Backend – CloudFront Integration (Internal Docs)

## 🎯 Purpose

To document the setup and usage of AWS CloudFront with our S3 buckets for returning optimized URLs for both **scout files** and **match videos**.

---

## ✅ Current Setup Summary

### CloudFront Distributions

| Use Case     | Domain                                  | CloudFront ID    |
| ------------ | --------------------------------------- | ---------------- |
| Scout Files  | `https://d1b5o37qbj029k.cloudfront.net` | `E2I1LDQ5PKDKHX` |
| Match Videos | `https://d2qk7473q3y7pg.cloudfront.net` | `E1WV6L56OFBZ8E` |

---

## 🔐 S3 Bucket Policy

Bucket: `volleymate-storage`

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowCloudFrontForScoutFiles",
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
      "Sid": "AllowCloudFrontForVideos",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudfront.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::volleymate-storage/videos/*",
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudfront::863518411349:distribution/E1WV6L56OFBZ8E"
        }
      }
    }
  ]
}
```

---

## 🌐 Environment Variables (.env.prod)

```env
SCOUT_CLOUDFRONT_DOMAIN=d1b5o37qbj029k.cloudfront.net
VIDEO_CLOUDFRONT_DOMAIN=d2qk7473q3y7pg.cloudfront.net
```

These are parsed in `config.InitConfig()` and injected into the video and scout file upload services.

---

## 📦 File Upload Result (Example)

PATCH `/api/admin/matches/:id/upload-video`

```json
{
  "video_url": "https://d2qk7473q3y7pg.cloudfront.net/videos/2024-2025_germany/bundesliga_male/{match_id}.mov"
}
```

---

## 🛠️ Related Files

- `pkg/storage/s3.go`
- `services/match_service.go`
- `.env.prod`
- `/etc/systemd/system/volleymate-backend.service`

---

## ✅ Tested Scenarios

- ✅ Uploading `.mov`, `.mp4`, `.mkv` results in correct CloudFront URL.
- ✅ CloudFront properly retrieves files due to bucket policy.
- ✅ Systemd loads `.env.prod` with proper domains.
- ✅ Frontend/mobile receives CloudFront URL.

---

## 🧪 Future Improvements

- Add signed URLs for expiring video access
- Use different cache behavior per file type in CloudFront
- Add logging for CDN hit/miss ratios
