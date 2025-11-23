# Deployment Guide for HCS Lab API

## GitHub Repository
Repository: `https://github.com/zefparis/HCS-LAB`

## Railway Deployment

### Quick Deploy
1. Go to [Railway.app](https://railway.app)
2. Click **"New Project"** → **"Deploy from GitHub repo"**
3. Connect your GitHub account if not already connected
4. Select the repository: **zefparis/HCS-LAB**
5. Railway will automatically detect the Dockerfile and deploy

### Environment Variables
No additional environment variables required. The service will use:
- `PORT`: Automatically set by Railway (usually 8080)

### Post-Deployment
Once deployed, Railway will provide you with a URL like:
```
https://hcs-lab-xxx.railway.app
```

### Test the Deployment
Check health endpoint:
```bash
curl https://your-app.railway.app/health
```

Generate HCS code:
```bash
curl -X POST https://your-app.railway.app/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "dominantElement": "Air",
    "modal": {"cardinal": 0.31, "fixed": 0.23, "mutable": 0.46},
    "cognition": {"fluid": 0.52, "crystallized": 0.13, "verbal": 0.53, "strategic": 0.15, "creative": 0.33},
    "interaction": {"pace": "balanced", "structure": "medium", "tone": "precise"}
  }'
```

## Vercel Dashboard Integration

The API is configured with CORS to accept requests from:
- `https://*.vercel.app`
- `http://localhost:*`
- `https://localhost:*`

### Example Frontend Code
```javascript
const generateHCS = async (profile) => {
  const response = await fetch('https://your-app.railway.app/api/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(profile),
  });
  
  const data = await response.json();
  return {
    codeU3: data.codeU3,
    codeU4: data.codeU4,
    chip: data.chip
  };
};
```

## Local Development

### Prerequisites
- Go 1.22+
- Git

### Clone and Run
```bash
git clone https://github.com/zefparis/HCS-LAB.git
cd HCS-LAB
go mod download
go run ./cmd/hcsapi
```

### Build Docker Image Locally
```bash
docker build -t hcs-lab-api .
docker run -p 8080:8080 hcs-lab-api
```

## Updating the Deployment

### Option 1: Automatic (Recommended)
Railway automatically redeploys when you push to the main branch:
```bash
git add .
git commit -m "Your update message"
git push origin main
```

### Option 2: Manual
1. Go to Railway dashboard
2. Click on your project
3. Click "Redeploy"

## Monitoring

Railway provides:
- Real-time logs
- Metrics (CPU, Memory, Network)
- Deployment history
- Environment variable management

Access these from your Railway project dashboard.

## Security Notes

1. **Salt File**: The `.hcs_salt` file is generated on first run and persists across deployments
2. **No External Dependencies**: The service runs completely offline
3. **Input Validation**: All inputs are validated and sanitized
4. **CORS**: Configured to only accept requests from specified origins

## Support

For issues or questions:
- GitHub Issues: https://github.com/zefparis/HCS-LAB/issues
- Railway Status: https://status.railway.app/
