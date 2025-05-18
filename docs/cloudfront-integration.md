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
| Assets       | `https://d3assetdomain.cloudfront.net`  | `ExxxxxxxxxxxX`  |

---

## 🔍 Verifying CloudFront Setup

### 1. Environment Variables

```bash
# Check in running service
sudo systemctl show -p Environment volleymate-backend.service

# Should show all CloudFront domains:
- SCOUT_CLOUDFRONT_DOMAIN
- VIDEO_CLOUDFRONT_DOMAIN
- ASSET_CLOUDFRONT_DOMAIN
```

### 2. URL Construction

- Scout files: `https://{SCOUT_CLOUDFRONT_DOMAIN}/scout-files/{file_id}.json`
- Videos: `https://{VIDEO_CLOUDFRONT_DOMAIN}/videos/{path}`
- Assets: `https://{ASSET_CLOUDFRONT_DOMAIN}/{asset_type}/{file}`

### 3. Testing URLs

1. Scout Files:

   ```bash
   curl -I https://d1b5o37qbj029k.cloudfront.net/scout-files/example.json
   ```

2. Videos:

   ```bash
   curl -I https://d2qk7473q3y7pg.cloudfront.net/videos/example.mp4
   ```

---

## S3 Bucket Policy

Bucket: `volleymate-storage`

Current bucket policy with detailed access control for different asset types:

```json
{
  "Version": "2008-10-17",
  "Statement": [
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
    },
    {
      "Sid": "AllowCloudFrontForAssets",
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudfront.amazonaws.com"
      },
      "Action": "s3:GetObject",
      "Resource": [
        "arn:aws:s3:::volleymate-storage/scout-files/*",
        "arn:aws:s3:::volleymate-storage/avatars/*",
        "arn:aws:s3:::volleymate-storage/logos/*",
        "arn:aws:s3:::volleymate-storage/logos/seasons/*",
        "arn:aws:s3:::volleymate-storage/logos/seasons/defaults/*"
      ],
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudfront::863518411349:distribution/E2I1LDQ5PKDKHX"
        }
      }
    }
  ]
}
```

### Policy Breakdown

1. **Video Distribution** (`E1WV6L56OFBZ8E`):

   - Handles all match videos
   - Path pattern: `/videos/*`
   - Uses dedicated CloudFront distribution for video content

2. **Asset Distribution** (`E2I1LDQ5PKDKHX`):
   - Handles multiple asset types:
     - Scout files: `/scout-files/*`
     - User avatars: `/avatars/*`
     - Team/Club logos: `/logos/*`
     - Season logos: `/logos/seasons/*`
     - Default season logos: `/logos/seasons/defaults/*`
   - Uses shared CloudFront distribution for all static assets

### Security Notes

- Each CloudFront distribution has its own IAM conditions
- Access is strictly controlled by path patterns
- No direct S3 access is allowed; all requests must go through CloudFront

---

## 🌐 Environment Variables (.env.prod)

```env
SCOUT_CLOUDFRONT_DOMAIN=d1b5o37qbj029k.cloudfront.net
VIDEO_CLOUDFRONT_DOMAIN=d2qk7473q3y7pg.cloudfront.net
ASSET_CLOUDFRONT_DOMAIN=d3assetdomain.cloudfront.net
```

These are parsed in `config.InitConfig()` and injected into the video and scout file upload services.

### Asset Domain Usage

The `ASSET_CLOUDFRONT_DOMAIN` is used specifically for:

- Team/Club/Player logos
- Season logos
- User avatars
- Other static assets

Example implementation in `services/season_service.go`:

```go
LogoURL: fmt.Sprintf("https://%s/logos/seasons/%s", config.AssetCloudFrontDomain, season.Logo)
```

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
