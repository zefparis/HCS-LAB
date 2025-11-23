# HCS Lab API Integration Guide

## 🚀 Live API Endpoint
```
https://hcs-lab-production.up.railway.app
```

## 📡 API Status
- ✅ **Deployed and Running**
- ✅ **Health Check Passing**
- ✅ **HCS Generation Working**
- ✅ **CORS Enabled for Vercel**

## 🔌 Integration Examples

### JavaScript/TypeScript (React/Next.js)

```javascript
// hcs-client.js
const HCS_API_URL = 'https://hcs-lab-production.up.railway.app';

export async function generateHCS(profile) {
  try {
    const response = await fetch(`${HCS_API_URL}/api/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(profile),
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }

    const data = await response.json();
    return {
      success: true,
      codeU3: data.codeU3,
      codeU4: data.codeU4,
      chip: data.chip,
      input: data.input
    };
  } catch (error) {
    return {
      success: false,
      error: error.message
    };
  }
}

// Example usage
const profile = {
  dominantElement: "Air",
  modal: {
    cardinal: 0.31,
    fixed: 0.23,
    mutable: 0.46
  },
  cognition: {
    fluid: 0.52,
    crystallized: 0.13,
    verbal: 0.53,
    strategic: 0.15,
    creative: 0.33
  },
  interaction: {
    pace: "balanced",
    structure: "medium",
    tone: "precise"
  }
};

const result = await generateHCS(profile);
console.log(result.codeU3); // HCS-U3|E:A|MOD:c31f23m46|...
```

### React Component Example

```jsx
import React, { useState } from 'react';
import { generateHCS } from './hcs-client';

function HCSGenerator() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const handleGenerate = async () => {
    setLoading(true);
    setError(null);
    
    const profile = {
      dominantElement: "Earth",
      modal: { cardinal: 0.25, fixed: 0.60, mutable: 0.15 },
      cognition: { 
        fluid: 0.35, 
        crystallized: 0.75, 
        verbal: 0.40, 
        strategic: 0.80, 
        creative: 0.20 
      },
      interaction: { 
        pace: "slow", 
        structure: "high", 
        tone: "warm" 
      }
    };

    const response = await generateHCS(profile);
    
    if (response.success) {
      setResult(response);
    } else {
      setError(response.error);
    }
    
    setLoading(false);
  };

  return (
    <div>
      <button onClick={handleGenerate} disabled={loading}>
        {loading ? 'Generating...' : 'Generate HCS Code'}
      </button>
      
      {result && (
        <div>
          <h3>Generated HCS Code:</h3>
          <code>{result.codeU3}</code>
          <p>CHIP: {result.chip}</p>
        </div>
      )}
      
      {error && (
        <div style={{ color: 'red' }}>
          Error: {error}
        </div>
      )}
    </div>
  );
}

export default HCSGenerator;
```

### Python Example

```python
import requests
import json

HCS_API_URL = "https://hcs-lab-production.up.railway.app"

def generate_hcs(profile):
    """Generate HCS codes from a profile"""
    try:
        response = requests.post(
            f"{HCS_API_URL}/api/generate",
            json=profile,
            headers={"Content-Type": "application/json"}
        )
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None

# Example usage
profile = {
    "dominantElement": "Water",
    "modal": {
        "cardinal": 0.45,
        "fixed": 0.30,
        "mutable": 0.25
    },
    "cognition": {
        "fluid": 0.70,
        "crystallized": 0.30,
        "verbal": 0.65,
        "strategic": 0.40,
        "creative": 0.75
    },
    "interaction": {
        "pace": "fast",
        "structure": "low",
        "tone": "neutral"
    }
}

result = generate_hcs(profile)
if result:
    print(f"HCS-U3: {result['codeU3']}")
    print(f"CHIP: {result['chip']}")
```

### cURL Command Line

```bash
# Generate HCS code
curl -X POST https://hcs-lab-production.up.railway.app/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "dominantElement": "Fire",
    "modal": {"cardinal": 0.50, "fixed": 0.20, "mutable": 0.30},
    "cognition": {"fluid": 0.60, "crystallized": 0.40, "verbal": 0.50, "strategic": 0.55, "creative": 0.65},
    "interaction": {"pace": "fast", "structure": "medium", "tone": "sharp"}
  }'

# Check health
curl https://hcs-lab-production.up.railway.app/health
```

## 📊 Input Validation

### Valid Values

**dominantElement**: 
- `"Earth"`, `"Air"`, `"Water"`, `"Fire"`

**modal/cognition values**: 
- Float between `0.0` and `1.0`

**interaction.pace**: 
- `"balanced"`, `"fast"`, `"slow"`

**interaction.structure**: 
- `"low"`, `"medium"`, `"high"`

**interaction.tone**: 
- `"warm"`, `"neutral"`, `"sharp"`, `"precise"`

## 🔒 Security Notes

1. **HTTPS Only**: The API is served over HTTPS via Railway
2. **CORS Enabled**: Configured for Vercel domains and localhost
3. **Input Validation**: All inputs are validated and sanitized
4. **Deterministic**: Same input always produces same output (per deployment)

## 📈 Response Format

### Success (200)
```json
{
  "input": { /* original input profile */ },
  "codeU3": "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:3117e464594a",
  "codeU4": "HCS-U4|eyJjaGlw...",
  "chip": "3117e464594a"
}
```

### Error (400)
```json
{
  "error": "Validation error",
  "message": "invalid dominant element: InvalidElement",
  "code": 400
}
```

## 🎯 Best Practices

1. **Cache Results**: HCS codes are deterministic, so you can cache them
2. **Batch Processing**: Send multiple requests in parallel if needed
3. **Error Handling**: Always handle 400 (validation) and 500 (server) errors
4. **Rate Limiting**: Be reasonable with request frequency (no hard limit currently)

## 🔗 Useful Links

- **API Base URL**: https://hcs-lab-production.up.railway.app
- **GitHub Repository**: https://github.com/zefparis/HCS-LAB
- **Railway Dashboard**: https://railway.app (for monitoring)

## 📞 Support

For issues or questions:
- GitHub Issues: https://github.com/zefparis/HCS-LAB/issues
- API Status: Check `/health` endpoint
