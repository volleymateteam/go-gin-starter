# 📦 Volleymate - AWS Storage Architecture

This document explains the video and scout file storage pipeline used by the Volleymate platform.

---

## 🎥 Video Upload Structure

All match videos are stored in S3 using the following structure:

videos/{season_year}{country}/{competition}{gender}/{match_id}/
├── raw/
│ └── original.mov (.mp4/.mkv)
├── compressed/
│ ├── 1080p/
│ ├── 720p/
│ └── 480p/
└── thumbnails/
└── thumb.jpg

---

## 🧠 Tagging Logic

- **Raw videos** are tagged with:
  Key: storage
  Value: raw

This enables us to target them in lifecycle rules.

---

## ♻️ Lifecycle Rule

Lifecycle Rule: `ArchiveRawVideosAfter30Days`

- **Scope**: only S3 objects with tag `storage=raw`
- **Action**: after 30 days, move to **Glacier Deep Archive**

This keeps costs low while preserving original quality for future use if needed.

---

## 🧾 CloudFront Delivery

- **Compressed videos (1080p/720p/480p)** are delivered via `VIDEO_CLOUDFRONT_DOMAIN`
- **Thumbnails** use the same domain.

---

## 🧪 Next Steps

- [ ] Track processed formats per match (for fallback logic in frontend).
- [ ] Implement usage dashboard for video & storage usage monitoring.
- [ ] Add optional auto-deletion of scout files older than 12 months (future).

---

## S3 Bucket Policy

### Management → Lifecycle Rules

```json
{
  "Rules": [
    {
      "ID": "ArchiveRawVideosAfter30Days",
      "Filter": {
        "Tag": {
          "Key": "storage",
          "Value": "raw"
        }
      },
      "Status": "Enabled",
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "GLACIER_DEEP_ARCHIVE"
        }
      ]
    }
  ]
}
```
