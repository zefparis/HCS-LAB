# 🚀 Quick Setup - Connecter Dashboard + APIs

## Étape 1: Variables d'environnement (Vercel)

Dans votre dashboard Vercel, ajoutez ces variables:

```bash
NEXT_PUBLIC_HCS_LAB_API_URL=https://hcs-lab-production.up.railway.app
NEXT_PUBLIC_ASTROLOGY_API_URL=https://your-astrology-api.railway.app
```

## Étape 2: Copier le code d'intégration

Copiez ces fichiers dans votre projet Next.js:

1. **`lib/corehuman-api.js`** (depuis `examples/dashboard-integration.js`)
2. **`components/ProfileGenerator.jsx`** (depuis `examples/ProfileComponent.jsx`)

## Étape 3: Utiliser dans votre app

```jsx
// pages/profile.js ou app/profile/page.js
import ProfileGenerator from '../components/ProfileGenerator';

export default function ProfilePage() {
  return (
    <div>
      <h1>CoreHuman Profile</h1>
      <ProfileGenerator />
    </div>
  );
}
```

## Étape 4: Déployer

```bash
git add .
git commit -m "Add HCS integration"
git push
```

Vercel va automatiquement redéployer avec les nouvelles fonctionnalités.

## 📊 Flux de données

```
User (Dashboard Vercel)
    ↓
    └→ Entre données de naissance
          ↓
          ├→ API Astrology (Swiss Ephemeris)
          │    └→ Retourne chart astrologique
          ↓
          ├→ Transformation → Profil HCS
          ↓
          └→ API HCS Lab
               └→ Retourne HCS-U3 + CHIP
```

## 🧪 Test rapide

```javascript
// Test direct depuis le browser console
fetch('https://hcs-lab-production.up.railway.app/api/generate', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    dominantElement: "Air",
    modal: { cardinal: 0.31, fixed: 0.23, mutable: 0.46 },
    cognition: { fluid: 0.52, crystallized: 0.13, verbal: 0.53, strategic: 0.15, creative: 0.33 },
    interaction: { pace: "balanced", structure: "medium", tone: "precise" }
  })
})
.then(r => r.json())
.then(console.log);
```

## ✅ Checklist

- [ ] Variables d'environnement configurées dans Vercel
- [ ] Code d'intégration copié dans le projet
- [ ] Component React ajouté
- [ ] Test de l'API HCS Lab fonctionnel
- [ ] Déployé sur Vercel

## 🆘 Dépannage

### Erreur CORS
→ Vérifiez que votre domaine Vercel est dans la liste CORS de l'API

### API ne répond pas
→ Testez directement: `curl https://hcs-lab-production.up.railway.app/health`

### Transformation chart → HCS
→ Ajustez la logique dans `transformChartToHCS()` selon vos besoins

## 📱 Support

- **HCS Lab API**: https://hcs-lab-production.up.railway.app
- **GitHub**: https://github.com/zefparis/HCS-LAB
- **Tests**: `.\test-live-api.ps1`
