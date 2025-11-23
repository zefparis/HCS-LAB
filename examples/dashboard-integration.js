// Quick integration example for Vercel Dashboard
// Copy this to your Next.js/React project

const API_CONFIG = {
  astrology: process.env.NEXT_PUBLIC_ASTROLOGY_API_URL,
  hcsLab: process.env.NEXT_PUBLIC_HCS_LAB_API_URL || 'https://hcs-lab-production.up.railway.app'
};

// Main integration class
class CoreHumanAPI {
  async getAstrologyChart(birthData) {
    const response = await fetch(`${API_CONFIG.astrology}/api/chart`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(birthData)
    });
    return response.json();
  }

  async generateHCSCode(profile) {
    const response = await fetch(`${API_CONFIG.hcsLab}/api/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(profile)
    });
    return response.json();
  }

  // Complete workflow
  async generateFullProfile(birthData) {
    // Step 1: Get astrology data
    const chart = await this.getAstrologyChart(birthData);
    
    // Step 2: Transform to HCS format
    const hcsProfile = {
      dominantElement: this.getDominantElement(chart),
      modal: this.getModalBalance(chart),
      cognition: {
        fluid: 0.5,
        crystallized: 0.5,
        verbal: 0.5,
        strategic: 0.5,
        creative: 0.5
      },
      interaction: {
        pace: "balanced",
        structure: "medium",
        tone: "neutral"
      }
    };
    
    // Step 3: Generate HCS code
    const hcsResult = await this.generateHCSCode(hcsProfile);
    
    return {
      astrology: chart,
      hcsCode: hcsResult.codeU3,
      chip: hcsResult.chip,
      profile: hcsResult.input
    };
  }

  getDominantElement(chart) {
    // Count planets in each element
    const elements = { Fire: 0, Earth: 0, Air: 0, Water: 0 };
    const elementSigns = {
      Fire: ['Aries', 'Leo', 'Sagittarius'],
      Earth: ['Taurus', 'Virgo', 'Capricorn'],
      Air: ['Gemini', 'Libra', 'Aquarius'],
      Water: ['Cancer', 'Scorpio', 'Pisces']
    };
    
    chart.planets?.forEach(planet => {
      for (const [element, signs] of Object.entries(elementSigns)) {
        if (signs.includes(planet.sign)) {
          elements[element]++;
        }
      }
    });
    
    // Return dominant
    return Object.keys(elements).reduce((a, b) => 
      elements[a] > elements[b] ? a : b
    );
  }

  getModalBalance(chart) {
    const modalities = { cardinal: 0, fixed: 0, mutable: 0 };
    const modalitySigns = {
      cardinal: ['Aries', 'Cancer', 'Libra', 'Capricorn'],
      fixed: ['Taurus', 'Leo', 'Scorpio', 'Aquarius'],
      mutable: ['Gemini', 'Virgo', 'Sagittarius', 'Pisces']
    };
    
    const total = chart.planets?.length || 10;
    
    chart.planets?.forEach(planet => {
      for (const [modality, signs] of Object.entries(modalitySigns)) {
        if (signs.includes(planet.sign)) {
          modalities[modality]++;
        }
      }
    });
    
    return {
      cardinal: modalities.cardinal / total,
      fixed: modalities.fixed / total,
      mutable: modalities.mutable / total
    };
  }
}

// Export for use in React components
export default new CoreHumanAPI();
