# Dashboard Integration Architecture

## 🏗️ Architecture Overview

```
                 +----------------+
                 | Dashboard UI   |  (Vercel)
                 +-------+--------+
                         |
      -----------------------------------------------
      |                                             |
+-----v-------------------+              +-----------v-------------+
| astrology-platform API |              |   hcs-lab API (new)     |
|  (Swiss Ephemeris)     |              |  HCS-U3/U4 engine       |
|  Railway (existing)    |              |  Railway                |
+------------------------+              +--------------------------+
```

## 📦 Dashboard Configuration

### 1. Environment Variables (Vercel)

Add these to your Vercel dashboard environment variables:

```env
# Existing API
NEXT_PUBLIC_ASTROLOGY_API_URL=https://your-astrology-api.railway.app

# New HCS Lab API
NEXT_PUBLIC_HCS_LAB_API_URL=https://hcs-lab-production.up.railway.app
```

### 2. API Client Service

Create a unified API service in your dashboard:

```javascript
// services/api-client.js

class APIClient {
  constructor() {
    // API endpoints
    this.astrologyAPI = process.env.NEXT_PUBLIC_ASTROLOGY_API_URL;
    this.hcsLabAPI = process.env.NEXT_PUBLIC_HCS_LAB_API_URL;
  }

  // ============= ASTROLOGY API CALLS =============
  async getChart(birthData) {
    const response = await fetch(`${this.astrologyAPI}/api/chart`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(birthData)
    });
    return response.json();
  }

  async getPlanetPositions(date, location) {
    const response = await fetch(`${this.astrologyAPI}/api/planets`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ date, location })
    });
    return response.json();
  }

  // ============= HCS LAB API CALLS =============
  async generateHCS(profile) {
    const response = await fetch(`${this.hcsLabAPI}/api/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(profile)
    });
    
    if (!response.ok) {
      throw new Error(`HCS Generation failed: ${response.status}`);
    }
    
    return response.json();
  }

  async checkHCSHealth() {
    const response = await fetch(`${this.hcsLabAPI}/health`);
    return response.json();
  }

  // ============= COMBINED WORKFLOW =============
  async generateCompleteProfile(birthData) {
    try {
      // Step 1: Get astrological data from astrology-platform
      const chartData = await this.getChart(birthData);
      
      // Step 2: Transform chart data to HCS input profile
      const hcsProfile = this.transformChartToHCSProfile(chartData);
      
      // Step 3: Generate HCS codes
      const hcsResult = await this.generateHCS(hcsProfile);
      
      // Step 4: Combine results
      return {
        astrology: chartData,
        hcs: hcsResult,
        timestamp: new Date().toISOString()
      };
    } catch (error) {
      console.error('Complete profile generation failed:', error);
      throw error;
    }
  }

  // Transform astrology data to HCS input format
  transformChartToHCSProfile(chartData) {
    // Example transformation logic
    return {
      dominantElement: this.calculateDominantElement(chartData),
      modal: this.calculateModalBalance(chartData),
      cognition: this.calculateCognitionProfile(chartData),
      interaction: this.calculateInteractionPreferences(chartData)
    };
  }

  calculateDominantElement(chartData) {
    // Logic to determine dominant element from chart
    const elements = { Fire: 0, Earth: 0, Air: 0, Water: 0 };
    
    // Count planets in each element
    chartData.planets?.forEach(planet => {
      const sign = planet.sign;
      if (['Aries', 'Leo', 'Sagittarius'].includes(sign)) elements.Fire++;
      if (['Taurus', 'Virgo', 'Capricorn'].includes(sign)) elements.Earth++;
      if (['Gemini', 'Libra', 'Aquarius'].includes(sign)) elements.Air++;
      if (['Cancer', 'Scorpio', 'Pisces'].includes(sign)) elements.Water++;
    });
    
    // Return the dominant element
    return Object.keys(elements).reduce((a, b) => 
      elements[a] > elements[b] ? a : b
    );
  }

  calculateModalBalance(chartData) {
    // Calculate cardinal, fixed, mutable balance
    let cardinal = 0, fixed = 0, mutable = 0;
    const total = chartData.planets?.length || 10;
    
    chartData.planets?.forEach(planet => {
      const sign = planet.sign;
      if (['Aries', 'Cancer', 'Libra', 'Capricorn'].includes(sign)) cardinal++;
      if (['Taurus', 'Leo', 'Scorpio', 'Aquarius'].includes(sign)) fixed++;
      if (['Gemini', 'Virgo', 'Sagittarius', 'Pisces'].includes(sign)) mutable++;
    });
    
    return {
      cardinal: cardinal / total,
      fixed: fixed / total,
      mutable: mutable / total
    };
  }

  calculateCognitionProfile(chartData) {
    // Example: Use Mercury, Moon, and aspects for cognitive profile
    const mercury = chartData.planets?.find(p => p.name === 'Mercury');
    const moon = chartData.planets?.find(p => p.name === 'Moon');
    
    return {
      fluid: 0.5,        // Based on Mercury aspects
      crystallized: 0.5, // Based on Saturn aspects
      verbal: 0.5,       // Mercury in Air signs
      strategic: 0.5,    // Mars/Saturn aspects
      creative: 0.5      // Venus/Neptune aspects
    };
  }

  calculateInteractionPreferences(chartData) {
    // Example: Use Ascendant and Moon for interaction style
    const ascendant = chartData.ascendant;
    const moon = chartData.planets?.find(p => p.name === 'Moon');
    
    return {
      pace: "balanced",  // Based on Mars/Mercury
      structure: "medium", // Based on Saturn
      tone: "neutral"    // Based on Venus/Moon
    };
  }
}

export default new APIClient();
```

### 3. React Component Example

```jsx
// components/ProfileGenerator.jsx
import { useState } from 'react';
import apiClient from '../services/api-client';

export default function ProfileGenerator() {
  const [loading, setLoading] = useState(false);
  const [profile, setProfile] = useState(null);
  const [error, setError] = useState(null);

  const generateProfile = async (birthData) => {
    setLoading(true);
    setError(null);
    
    try {
      // This calls both APIs in sequence
      const completeProfile = await apiClient.generateCompleteProfile(birthData);
      
      setProfile({
        // Astrology data
        chart: completeProfile.astrology,
        planets: completeProfile.astrology.planets,
        aspects: completeProfile.astrology.aspects,
        
        // HCS data
        hcsCode: completeProfile.hcs.codeU3,
        hcsChip: completeProfile.hcs.chip,
        
        // Combined insights
        dominantElement: completeProfile.hcs.input.dominantElement,
        modalBalance: completeProfile.hcs.input.modal,
        cognition: completeProfile.hcs.input.cognition,
        interaction: completeProfile.hcs.input.interaction
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Complete Astrological & HCS Profile</h2>
      
      {/* Birth data form */}
      <form onSubmit={(e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        generateProfile({
          name: formData.get('name'),
          date: formData.get('date'),
          time: formData.get('time'),
          location: formData.get('location')
        });
      }}>
        <input name="name" placeholder="Name" required />
        <input name="date" type="date" required />
        <input name="time" type="time" required />
        <input name="location" placeholder="City, Country" required />
        <button type="submit" disabled={loading}>
          {loading ? 'Generating...' : 'Generate Profile'}
        </button>
      </form>

      {/* Results display */}
      {profile && (
        <div className="results">
          <div className="astrology-section">
            <h3>Astrological Data</h3>
            <p>Sun Sign: {profile.chart?.sun_sign}</p>
            <p>Moon Sign: {profile.chart?.moon_sign}</p>
            <p>Ascendant: {profile.chart?.ascendant}</p>
          </div>

          <div className="hcs-section">
            <h3>HCS Profile</h3>
            <p>Dominant Element: {profile.dominantElement}</p>
            <p>HCS Code:</p>
            <code>{profile.hcsCode}</code>
            <p>CHIP: {profile.hcsChip}</p>
          </div>

          <div className="modal-balance">
            <h3>Modal Balance</h3>
            <p>Cardinal: {(profile.modalBalance.cardinal * 100).toFixed(0)}%</p>
            <p>Fixed: {(profile.modalBalance.fixed * 100).toFixed(0)}%</p>
            <p>Mutable: {(profile.modalBalance.mutable * 100).toFixed(0)}%</p>
          </div>

          <div className="cognition">
            <h3>Cognition Profile</h3>
            <p>Fluid: {(profile.cognition.fluid * 100).toFixed(0)}%</p>
            <p>Crystallized: {(profile.cognition.crystallized * 100).toFixed(0)}%</p>
            <p>Verbal: {(profile.cognition.verbal * 100).toFixed(0)}%</p>
            <p>Strategic: {(profile.cognition.strategic * 100).toFixed(0)}%</p>
            <p>Creative: {(profile.cognition.creative * 100).toFixed(0)}%</p>
          </div>
        </div>
      )}

      {error && <div className="error">{error}</div>}
    </div>
  );
}
```

### 4. Next.js API Routes (Optional Proxy)

If you prefer to proxy API calls through your Next.js backend:

```javascript
// pages/api/generate-profile.js
import axios from 'axios';

export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  try {
    const { birthData } = req.body;

    // Call astrology API
    const astrologyResponse = await axios.post(
      `${process.env.ASTROLOGY_API_URL}/api/chart`,
      birthData
    );

    // Transform data for HCS
    const hcsProfile = transformToHCSProfile(astrologyResponse.data);

    // Call HCS Lab API
    const hcsResponse = await axios.post(
      `${process.env.HCS_LAB_API_URL}/api/generate`,
      hcsProfile
    );

    // Return combined result
    res.status(200).json({
      astrology: astrologyResponse.data,
      hcs: hcsResponse.data
    });
  } catch (error) {
    console.error('Profile generation error:', error);
    res.status(500).json({ 
      error: 'Failed to generate profile',
      details: error.message 
    });
  }
}
```

## 🔐 Security Best Practices

1. **Environment Variables**: Never expose API URLs in client code if they contain sensitive endpoints
2. **CORS**: Both APIs should have CORS configured for your Vercel domain
3. **Rate Limiting**: Implement rate limiting on the dashboard side
4. **Error Handling**: Always handle API failures gracefully
5. **Caching**: Cache HCS results since they're deterministic

## 📊 Data Flow

```
User Input (Birth Data)
    ↓
Dashboard UI (Vercel)
    ↓
    ├→ astrology-platform API
    │     ├→ Swiss Ephemeris calculations
    │     └→ Return chart data
    │
    ├→ Transform chart → HCS profile
    │
    └→ hcs-lab API
          ├→ Generate HCS codes
          └→ Return U3/U4 + CHIP

Dashboard combines & displays results
```

## 🚀 Quick Start

1. **Set environment variables** in Vercel dashboard
2. **Copy the API client** to your project
3. **Implement the React component**
4. **Deploy to Vercel**

Your dashboard will now seamlessly integrate both APIs!
