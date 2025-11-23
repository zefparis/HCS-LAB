# HCS Lab API

A standalone Go service for generating Human Coordination System (HCS) codes. This is a private backend API that generates HCS-U3, HCS-U4, and the new HCS-U5 fusion codes that combine Western and Chinese astrological profiles into a unified cognitive signature.

## Overview

The HCS Lab API provides:
- **Deterministic code generation** with persistent salt management
- **HCS-U3 encoding** with structured segments for element, modal balance, cognition, interaction, and CHIP signature
- **HCS-U4 encoding** (base64 stub implementation, ready for future base62)
- **HCS-U5 encoding** NEW: Fusion of Western and Chinese (BaZi) astrology profiles
- **Chinese BaZi computation** Four Pillars of Destiny calculation
- **Fusion profiles** Intelligent merging of Western and Chinese astrological systems
- **RESTful HTTP API** for integration with dashboards and tools
- **CLI tool** for command-line code generation
- **Secure offline operation** with no external network dependencies

## HCS-U5 Format Specification (NEW)

The HCS-U5 code represents a fusion of Western and Chinese astrology:
```
HCS-U5|XX|W:<hex>|C:<hex>|F:<hex>|CHIP:<12-hex>
```

### Segments:
- **XX**: 2-character fusion ID (combined element and balance encoding)
- **W:<hex>**: 16-bit compressed Western profile (4 hex chars)
- **C:<hex>**: 16-bit compressed Chinese BaZi profile (4 hex chars)  
- **F:<hex>**: 16-bit compressed fusion traits (4 hex chars)
- **CHIP**: 12-character hex signature from SHA256 hash

### Example:
```
HCS-U5|A1|W:3c4f|C:8a2d|F:6b91|CHIP:def012345678
```

## HCS-U3 Format Specification

The HCS-U3 code follows this exact format:
```
HCS-U3|E:<ELEM>|MOD:c<CC>f<FF>m<MM>|COG:F<FF>C<CC>V<VV>S<SS>Cr<CR>|INT:PB=<PACE>,SM=<STRUCT>,TN=<TONE>|CHIP:<12-hex>
```

### Segments:
1. **Element (E)**: Single letter mapping
   - Earth → E, Air → A, Water → W, Fire → F

2. **Modal Balance (MOD)**: Percentages (00-99)
   - Cardinal (c), Fixed (f), Mutable (m)

3. **Cognition (COG)**: Percentages (00-99)
   - Fluid (F), Crystallized (C), Verbal (V), Strategic (S), Creative (Cr)

4. **Interaction (INT)**: Preference codes
   - Pace (PB): balanced→B, fast→F, slow→S
   - Structure (SM): low→L, medium→M, high→H
   - Tone (TN): warm→W, neutral→N, sharp→S, precise→P

5. **CHIP**: 12-character hex signature from SHA256 hash

### Example:
```
HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:aae673a93e1f
```

## Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Clone and Build
```bash
git clone https://github.com/corehuman/hcs-lab-api.git
cd hcs-lab-api
go mod download
go build ./cmd/hcsgen
go build ./cmd/hcsapi
```

## Usage

### CLI Tool (hcsgen)

Generate HCS codes from a JSON input file:

```bash
# Basic usage
./hcsgen input.json

# Generate only U3 code
./hcsgen --u3-only input.json

# Generate only U4 code  
./hcsgen --u4-only input.json

# Pretty print JSON output
./hcsgen --pretty --raw-json input.json
```

Output files:
- `input_output.json` - Full HCS output with all fields
- `input_output.hcs` - Just the HCS codes (one per line)

### HTTP API Server

Start the server:
```bash
./hcsapi
# Server starts on port 8080 by default

# Or with custom port
PORT=3000 ./hcsapi
```

#### Endpoints

**Health Check**
```bash
GET /health

Response:
{
  "status": "healthy",
  "version": "1.0.0-hcs-lab",
  "uptime": "2h 15m 30s",
  "secure": true
}
```

**Generate HCS Codes**
```bash
POST /api/generate
Content-Type: application/json

Body:
{
  "dominantElement": "Air",
  "modal": {
    "cardinal": 0.31,
    "fixed": 0.23,
    "mutable": 0.46
  },
  "cognition": {
    "fluid": 0.52,
    "crystallized": 0.13,
    "verbal": 0.53,
    "strategic": 0.15,
    "creative": 0.33
  },
  "interaction": {
    "pace": "balanced",
    "structure": "medium",
    "tone": "precise"
  },
  "birthInfo": {  // Optional: enables HCS-U5 generation
    "year": 1990,
    "month": 6,
    "day": 15,
    "hour": 14,
    "minute": 30,
    "timezone": "UTC"
  }
}

Response:
{
  "input": { ... },
  "codeU3": "HCS-U3|E:A|MOD:c31f23m46|COG:F52C13V53S15Cr33|INT:PB=B,SM=M,TN=P|CHIP:aae673a93e1f",
  "codeU4": "HCS-U4|eyJwcm9maWxlIjp7ImVsZW1lbnQiOiJB...",
  "codeU5": "HCS-U5|A1|W:3c4f|C:8a2d|F:6b91|CHIP:def012345678",  // If birthInfo provided
  "chip": "aae673a93e1f",
  "chineseProfile": {  // If birthInfo provided
    "yearPillar": "Geng-Wu",
    "monthPillar": "Ren-Wu",
    "dayPillar": "Jia-Chen",
    "hourPillar": "Xin-Wei",
    "yinYangBalance": 0.6,
    "elementBalance": {
      "Wood": 0.2,
      "Fire": 0.3,
      "Earth": 0.2,
      "Metal": 0.15,
      "Water": 0.15
    },
    "dayMaster": "Jia",
    "dayMasterStrength": 0.7
  },
  "combinedProfile": {  // If birthInfo provided
    "western": { ... },
    "chinese": { ... },
    "fusion": {
      "elementSignature": { ... },
      "cognitiveFusion": { ... },
      "tempoSignals": { ... },
      "unifiedBalance": 0.65,
      "harmonicResonance": 0.75,
      "fusionId": "A1"
    }
  }
}
```

## Input JSON Format

```json
{
  "dominantElement": "Air",
  "modal": {
    "cardinal": 0.31,
    "fixed": 0.23,
    "mutable": 0.46
  },
  "cognition": {
    "fluid": 0.52,
    "crystallized": 0.13,
    "verbal": 0.53,
    "strategic": 0.15,
    "creative": 0.33
  },
  "interaction": {
    "pace": "balanced",
    "structure": "medium",
    "tone": "precise"
  },
  "birthInfo": {  // Optional: for Chinese BaZi and HCS-U5 generation
    "year": 1990,
    "month": 6,
    "day": 15,
    "hour": 14,
    "minute": 30,
    "timezone": "UTC"
  }
}
```

### Field Constraints:
- **dominantElement**: Must be one of: "Earth", "Air", "Water", "Fire"
- **modal values**: Float between 0.0 and 1.0 (converted to percentages)
- **cognition values**: Float between 0.0 and 1.0 (converted to percentages)
- **interaction.pace**: "balanced", "fast", or "slow"
- **interaction.structure**: "low", "medium", or "high"
- **interaction.tone**: "warm", "neutral", "sharp", or "precise"
- **birthInfo** (optional):
  - **year**: 1900-2100
  - **month**: 1-12
  - **day**: 1-31
  - **hour**: 0-23
  - **minute**: 0-59
  - **timezone**: IANA timezone string (e.g., "UTC", "America/New_York")

## Docker Deployment

### Build Image
```bash
docker build -t hcs-lab-api .
```

### Run Container
```bash
docker run -p 8080:8080 hcs-lab-api
```

### Deploy to Railway
The repository includes a Dockerfile optimized for Railway deployment. Simply connect your GitHub repository to Railway, and it will automatically detect and build using the Dockerfile.

## Testing

Run all tests:
```bash
go test ./tests/... -v
```

Test coverage includes:
- Deterministic code generation
- Salt persistence and management
- HCS-U3 format validation
- HCS-U5 format validation and encoding
- BaZi pillar computation
- Chinese element balance calculation
- Fusion profile generation
- Input validation and error handling
- Percentage rounding and clamping
- CHIP signature generation

## Security Features

- **Persistent Salt**: A 32-byte cryptographic salt is generated on first use and stored in `.hcs_salt`
- **Deterministic Output**: Same input always produces same output (with same salt)
- **Offline Operation**: No external network calls or dependencies
- **Input Validation**: All inputs are validated and clamped to acceptable ranges
- **SHA256 Hashing**: Secure cryptographic hashing for CHIP generation

## Project Structure

```
hcs-lab-api/
├── cmd/
│   ├── hcsgen/          # CLI tool
│   └── hcsapi/          # HTTP API server
├── internal/
│   └── hcs/
│       ├── model.go     # Data structures
│       ├── generator.go # Core generation logic
│       ├── codec_u3.go  # HCS-U3 encoding
│       ├── codec_u4.go  # HCS-U4 encoding (stub)
│       ├── codec_u5.go  # HCS-U5 fusion encoding
│       ├── bazi.go      # Chinese BaZi computation
│       ├── chinese.go   # Chinese profile generation
│       ├── fusion.go    # Western-Chinese fusion logic
│       ├── crypto.go    # SHA256 + CHIP logic
│       └── salt.go      # Salt management
├── tests/               # Test suites
├── examples/            # Sample inputs
├── Dockerfile          # Docker deployment
├── go.mod              # Go dependencies
└── README.md           # This file
```

## Version

Current version: `1.0.0-hcs-lab`

## License

Private repository for CoreHuman/HCS project.
