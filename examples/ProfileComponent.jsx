// React component for Vercel Dashboard
// Place this in your Next.js components folder

import React, { useState } from 'react';

const ProfileGenerator = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  // API URLs from environment variables
  const ASTROLOGY_API = process.env.NEXT_PUBLIC_ASTROLOGY_API_URL;
  const HCS_API = process.env.NEXT_PUBLIC_HCS_LAB_API_URL || 
                   'https://hcs-lab-production.up.railway.app';

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const formData = new FormData(e.target);
    const birthData = {
      name: formData.get('name'),
      date: formData.get('date'),
      time: formData.get('time'),
      location: formData.get('location'),
      timezone: formData.get('timezone')
    };

    try {
      // Step 1: Get astrological chart (if astrology API is available)
      let chartData = null;
      if (ASTROLOGY_API) {
        const chartRes = await fetch(`${ASTROLOGY_API}/api/chart`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(birthData)
        });
        chartData = await chartRes.json();
      }

      // Step 2: Prepare HCS profile (with or without chart data)
      const hcsProfile = chartData ? 
        transformChartToHCS(chartData) : 
        getDefaultHCSProfile();

      // Step 3: Generate HCS code
      const hcsRes = await fetch(`${HCS_API}/api/generate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(hcsProfile)
      });

      if (!hcsRes.ok) {
        throw new Error(`HCS API error: ${hcsRes.status}`);
      }

      const hcsData = await hcsRes.json();

      // Combine results
      setResult({
        birthData,
        chart: chartData,
        hcs: hcsData
      });

    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Transform astrological chart to HCS profile
  const transformChartToHCS = (chart) => {
    // This is a simplified example - adjust based on your logic
    const elements = countElements(chart.planets || []);
    const modalities = countModalities(chart.planets || []);
    
    return {
      dominantElement: getDominantElement(elements),
      modal: {
        cardinal: modalities.cardinal / 10,
        fixed: modalities.fixed / 10,
        mutable: modalities.mutable / 10
      },
      cognition: {
        fluid: 0.5,        // TODO: Calculate from Mercury
        crystallized: 0.5, // TODO: Calculate from Saturn
        verbal: 0.5,       // TODO: Calculate from Mercury/3rd house
        strategic: 0.5,    // TODO: Calculate from Mars/10th house
        creative: 0.5      // TODO: Calculate from Venus/Neptune
      },
      interaction: {
        pace: "balanced",   // TODO: Derive from Mars/Mercury
        structure: "medium", // TODO: Derive from Saturn
        tone: "neutral"     // TODO: Derive from Venus/Moon
      }
    };
  };

  // Default profile when no astrology data
  const getDefaultHCSProfile = () => ({
    dominantElement: "Air",
    modal: { cardinal: 0.33, fixed: 0.33, mutable: 0.34 },
    cognition: { 
      fluid: 0.5, crystallized: 0.5, verbal: 0.5, 
      strategic: 0.5, creative: 0.5 
    },
    interaction: { pace: "balanced", structure: "medium", tone: "neutral" }
  });

  const countElements = (planets) => {
    const counts = { Fire: 0, Earth: 0, Air: 0, Water: 0 };
    const signs = {
      Fire: ['Aries', 'Leo', 'Sagittarius'],
      Earth: ['Taurus', 'Virgo', 'Capricorn'],
      Air: ['Gemini', 'Libra', 'Aquarius'],
      Water: ['Cancer', 'Scorpio', 'Pisces']
    };
    
    planets.forEach(planet => {
      for (const [element, signList] of Object.entries(signs)) {
        if (signList.includes(planet.sign)) counts[element]++;
      }
    });
    
    return counts;
  };

  const countModalities = (planets) => {
    const counts = { cardinal: 0, fixed: 0, mutable: 0 };
    const signs = {
      cardinal: ['Aries', 'Cancer', 'Libra', 'Capricorn'],
      fixed: ['Taurus', 'Leo', 'Scorpio', 'Aquarius'],
      mutable: ['Gemini', 'Virgo', 'Sagittarius', 'Pisces']
    };
    
    planets.forEach(planet => {
      for (const [modality, signList] of Object.entries(signs)) {
        if (signList.includes(planet.sign)) counts[modality]++;
      }
    });
    
    return counts;
  };

  const getDominantElement = (counts) => {
    return Object.keys(counts).reduce((a, b) => 
      counts[a] > counts[b] ? a : b
    );
  };

  return (
    <div className="profile-generator">
      <h2>CoreHuman Profile Generator</h2>
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input name="name" type="text" required />
        </div>
        
        <div className="form-group">
          <label>Birth Date:</label>
          <input name="date" type="date" required />
        </div>
        
        <div className="form-group">
          <label>Birth Time:</label>
          <input name="time" type="time" required />
        </div>
        
        <div className="form-group">
          <label>Birth Location:</label>
          <input name="location" type="text" placeholder="City, Country" required />
        </div>
        
        <div className="form-group">
          <label>Timezone:</label>
          <select name="timezone" required>
            <option value="UTC">UTC</option>
            <option value="EST">EST</option>
            <option value="PST">PST</option>
            <option value="CET">CET</option>
          </select>
        </div>
        
        <button type="submit" disabled={loading}>
          {loading ? 'Generating...' : 'Generate Profile'}
        </button>
      </form>

      {error && (
        <div className="error">
          Error: {error}
        </div>
      )}

      {result && (
        <div className="results">
          <h3>Generated Profile</h3>
          
          {/* HCS Results */}
          <div className="hcs-section">
            <h4>HCS Code</h4>
            <code className="hcs-code">{result.hcs.codeU3}</code>
            <p><strong>CHIP:</strong> {result.hcs.chip}</p>
            <p><strong>Element:</strong> {result.hcs.input.dominantElement}</p>
          </div>

          {/* Modal Balance */}
          <div className="modal-section">
            <h4>Modal Balance</h4>
            <ul>
              <li>Cardinal: {(result.hcs.input.modal.cardinal * 100).toFixed(0)}%</li>
              <li>Fixed: {(result.hcs.input.modal.fixed * 100).toFixed(0)}%</li>
              <li>Mutable: {(result.hcs.input.modal.mutable * 100).toFixed(0)}%</li>
            </ul>
          </div>

          {/* Cognition Profile */}
          <div className="cognition-section">
            <h4>Cognition Profile</h4>
            <ul>
              <li>Fluid: {(result.hcs.input.cognition.fluid * 100).toFixed(0)}%</li>
              <li>Crystallized: {(result.hcs.input.cognition.crystallized * 100).toFixed(0)}%</li>
              <li>Verbal: {(result.hcs.input.cognition.verbal * 100).toFixed(0)}%</li>
              <li>Strategic: {(result.hcs.input.cognition.strategic * 100).toFixed(0)}%</li>
              <li>Creative: {(result.hcs.input.cognition.creative * 100).toFixed(0)}%</li>
            </ul>
          </div>

          {/* Astrological Data (if available) */}
          {result.chart && (
            <div className="astrology-section">
              <h4>Astrological Chart</h4>
              <p>Sun: {result.chart.sun_sign}</p>
              <p>Moon: {result.chart.moon_sign}</p>
              <p>Rising: {result.chart.ascendant}</p>
            </div>
          )}
        </div>
      )}

      <style jsx>{`
        .profile-generator {
          max-width: 600px;
          margin: 0 auto;
          padding: 20px;
        }

        .form-group {
          margin-bottom: 15px;
        }

        label {
          display: block;
          margin-bottom: 5px;
          font-weight: bold;
        }

        input, select {
          width: 100%;
          padding: 8px;
          border: 1px solid #ddd;
          border-radius: 4px;
        }

        button {
          background: #0070f3;
          color: white;
          padding: 10px 20px;
          border: none;
          border-radius: 4px;
          cursor: pointer;
          font-size: 16px;
        }

        button:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        .error {
          background: #fee;
          color: #c00;
          padding: 10px;
          border-radius: 4px;
          margin-top: 10px;
        }

        .results {
          margin-top: 20px;
          padding: 20px;
          background: #f5f5f5;
          border-radius: 8px;
        }

        .hcs-code {
          display: block;
          padding: 10px;
          background: #000;
          color: #0f0;
          font-family: monospace;
          margin: 10px 0;
          word-break: break-all;
        }

        h4 {
          color: #333;
          border-bottom: 2px solid #0070f3;
          padding-bottom: 5px;
        }
      `}</style>
    </div>
  );
};

export default ProfileGenerator;
