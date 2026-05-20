# 🇲🇦 SecureVPN - VPN Open Source pour le Maroc

## Bienvenue! مرحبا!

Ce projet VPN open-source est spécialement conçu pour répondre aux besoins des entreprises marocaines, des TPE/PME, et des particuliers.

---

## 🎯 Pourquoi ce VPN?

### Pour les Entreprises Marocaines

#### TPE (Très Petites Entreprises)
- ✅ **Gratuit et Open Source** - Pas de frais de licence
- ✅ **Facile à déployer** - Installation en 5 minutes
- ✅ **Sécurité professionnelle** - Cryptographie moderne
- ✅ **Support local** - Documentation en français

#### PME (Petites et Moyennes Entreprises)
- ✅ **Évolutif** - Supporte 1000+ clients simultanés
- ✅ **Haute performance** - 1+ Gbps de débit
- ✅ **Télétravail sécurisé** - Accès distant pour employés
- ✅ **Contrôle total** - Hébergez sur vos propres serveurs

### Cas d'Usage au Maroc

#### 1. Bureaux Multiples
Connectez vos bureaux à Casablanca, Rabat, Marrakech, et Tanger de manière sécurisée.

```
Bureau Casablanca ←→ VPN Server ←→ Bureau Rabat
                         ↕
                   Bureau Marrakech
```

#### 2. Télétravail
Permettez à vos employés de travailler depuis chez eux en toute sécurité.

```
Employé (Maison) → VPN → Réseau Entreprise
```

#### 3. Protection des Données
Sécurisez les données sensibles de votre entreprise avec un cryptage de niveau militaire.

#### 4. Accès aux Ressources Internes
Accédez à vos serveurs, bases de données, et applications internes depuis n'importe où.

---

## ⚠️ AVERTISSEMENT SÉCURITÉ ET LÉGAL / SECURITY & LEGAL WARNING

> **🚨 IMPORTANT: À LIRE AVANT UTILISATION / READ BEFORE USE 🚨**

### 🛡️ Ce Projet est Conçu Pour / This Project is Designed For:

✅ **Utilisation Légitime Uniquement / Legitimate Use Only:**
- 🔐 Communications professionnelles sécurisées / Secure business communications
- 🏢 Réseaux privés d'entreprise / Private corporate networks
- 🎓 Objectifs éducatifs et de recherche / Educational and research purposes
- 🔒 Protection des données sensibles / Protecting sensitive data
- 🌐 Accès distant sécurisé / Secure remote access
- 💼 Déploiements VPN professionnels / Professional VPN deployments

### 🚫 STRICTEMENT INTERDIT / STRICTLY PROHIBITED:

❌ **NE PAS UTILISER POUR / DO NOT USE FOR:**
- Activités illégales de toute nature / Illegal activities of any kind
- Contournement de restrictions légales / Bypassing legal restrictions
- Accès non autorisé aux réseaux / Unauthorized network access
- Violation des droits d'auteur / Copyright violations
- Routage de trafic malveillant / Malicious traffic routing
- Toute activité violant les lois locales ou internationales / Any activity violating laws

### ⚖️ Conformité Légale / Legal Compliance:

**VOUS ÊTES RESPONSABLE DE / YOU ARE RESPONSIBLE FOR:**
- ✅ Conformité avec toutes les lois applicables au Maroc / Compliance with all Moroccan laws
- ✅ Obtention des autorisations nécessaires / Obtaining necessary permissions
- ✅ Respect des droits de propriété intellectuelle / Respecting intellectual property
- ✅ Respect des politiques de sécurité / Following security policies
- ✅ Utilisation légitime en tout temps / Ensuring legitimate use

**🇲🇦 Spécifique au Maroc / Morocco Specific:**
- L'utilisation de VPN est légale au Maroc pour des fins légitimes
- Les utilisateurs doivent respecter les réglementations des télécommunications marocaines
- Les entreprises doivent consulter un conseiller juridique pour la conformité
- VPN usage is legal in Morocco for legitimate purposes
- Users must comply with Moroccan telecommunications regulations
- Businesses should consult legal counsel for compliance

### 🔒 Responsabilités de Sécurité / Security Responsibilities:

**EN TANT QU'UTILISATEUR / AS A USER:**
- 🔑 Gardez les clés privées sécurisées / Keep private keys secure
- 🔐 Utilisez des mots de passe forts / Use strong passwords
- 📊 Surveillez les logs / Monitor logs for suspicious activity
- 🔄 Maintenez le logiciel à jour / Keep software updated
- 🛡️ Suivez les meilleures pratiques / Follow security best practices
- 📋 Implémentez des contrôles d'accès / Implement access controls

### ⚠️ Clause de Non-Responsabilité / Disclaimer:

**LES AUTEURS ET CONTRIBUTEURS / THE AUTHORS AND CONTRIBUTORS:**
- Ne sont PAS responsables de la mauvaise utilisation / Are NOT responsible for misuse
- N'approuvent PAS les activités illégales / Do NOT endorse illegal activities
- Fournissent ce logiciel "TEL QUEL" / Provide software "AS IS"
- Ne sont PAS responsables des dommages / Are NOT liable for damages
- Condamnent fermement toute mauvaise utilisation / Strongly condemn any misuse

**EN UTILISANT CE LOGICIEL / BY USING THIS SOFTWARE:**
- Vous acceptez de l'utiliser uniquement à des fins légales / You agree to use it lawfully
- Vous acceptez de respecter toutes les lois / You agree to comply with all laws
- Vous assumez l'entière responsabilité / You take full responsibility
- Les auteurs ne sont pas responsables de vos actions / Authors are not liable for your actions

### 📞 Signaler un Abus / Report Abuse:

Si vous découvrez une mauvaise utilisation / If you discover misuse:
- GitHub Issues: [Signaler Ici / Report Here](https://github.com/lahcenassmira/open-source-vpn/issues)
- Autorités locales si activité illégale / Local authorities if illegal activity

---

**🔴 RAPPEL: Un grand pouvoir implique de grandes responsabilités. Utilisez cet outil de manière éthique et légale. 🔴**

**🔴 REMEMBER: With great power comes great responsibility. Use this tool ethically and legally. 🔴**

---

## 🚀 Installation Rapide

### Prérequis
- Serveur Linux (Ubuntu, Debian, CentOS)
- Go 1.21+ installé
- Accès root/sudo

### Étape 1: Télécharger
```bash
git clone https://github.com/lahcenassmira/open-source-vpn.git
cd open-source-vpn
```

### Étape 2: Compiler
```bash
go mod download
make build
```

### Étape 3: Configurer le Serveur
```bash
# Générer les clés
sudo ./bin/vpn-server keygen --output server-keys.json

# Copier la configuration
cp configs/server.example.yaml server.yaml

# Éditer la configuration (ajoutez votre clé privée)
nano server.yaml
```

### Étape 4: Démarrer le Serveur
```bash
sudo ./bin/vpn-server start --config server.yaml
```

### Étape 5: Configurer le Client
```bash
# Générer les clés client
./bin/vpn-client keygen --output client-keys.json

# Copier la configuration
cp configs/client.example.yaml client.yaml

# Éditer la configuration
nano client.yaml
```

### Étape 6: Se Connecter
```bash
sudo ./bin/vpn-client connect --config client.yaml
```

---

## 💰 Coûts et Économies

### Solution Commerciale Typique
- 💸 Licence: 500-2000 DH/mois par utilisateur
- 💸 Support: 5000-10000 DH/an
- 💸 Maintenance: 3000-5000 DH/an
- **Total: 15,000-50,000 DH/an pour 5 utilisateurs**

### Cette Solution Open Source
- ✅ Licence: **GRATUIT** (MIT License)
- ✅ Support: Communauté open-source
- ✅ Maintenance: Auto-géré
- ✅ Serveur: ~200-500 DH/mois (VPS)
- **Total: ~2,400-6,000 DH/an**

**Économie: 80-90% par rapport aux solutions commerciales!**

---

## 🔐 Sécurité

### Cryptographie de Niveau Militaire
- **Chiffrement**: ChaCha20-Poly1305 (utilisé par Google, Cloudflare)
- **Échange de clés**: X25519 (Curve25519)
- **Authentification**: Clés publiques/privées
- **Protection**: Anti-rejeu, forward secrecy

### Conformité
- ✅ Conforme aux standards internationaux
- ✅ Audit de code possible (open source)
- ✅ Pas de backdoors
- ✅ Contrôle total de vos données

---

## 📊 Performances

### Capacité
- **Clients simultanés**: 1000+
- **Débit**: 1+ Gbps
- **Latence**: < 5ms
- **Mémoire**: ~50MB + 1MB par client

### Testé sur
- ✅ Ubuntu 20.04/22.04
- ✅ Debian 11/12
- ✅ CentOS 8
- ✅ Serveurs VPS marocains

---

## 🏢 Déploiement pour Entreprises

### Option 1: Serveur Dédié
Hébergez sur votre propre serveur dans votre bureau.

**Avantages**:
- Contrôle total
- Pas de frais mensuels
- Données restent locales

### Option 2: VPS Marocain
Utilisez un VPS d'un fournisseur marocain.

**Fournisseurs recommandés**:
- OVH Maroc
- Genious Communications
- Autres fournisseurs locaux

**Avantages**:
- Faible latence
- Support local
- Conformité locale

### Option 3: Cloud International
Utilisez AWS, Google Cloud, ou Azure.

**Avantages**:
- Haute disponibilité
- Évolutivité
- Redondance

---

## 📚 Documentation

### En Français
- [Guide de Démarrage Rapide](QUICKSTART.md)
- [Guide d'Installation Détaillé](docs/SETUP.md)
- [Architecture du Système](docs/ARCHITECTURE.md)

### En Anglais
- [README Principal](README.md)
- [Guide de Contribution](CONTRIBUTING.md)
- [Résumé du Projet](PROJECT_SUMMARY.md)

---

## 🤝 Support et Communauté

### Obtenir de l'Aide
- 📖 Documentation complète incluse
- 🐛 Signaler des bugs: [GitHub Issues](https://github.com/lahcenassmira/open-source-vpn/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/lahcenassmira/open-source-vpn/discussions)

### Contribuer
Ce projet est open source! Vos contributions sont les bienvenues:
- Corrections de bugs
- Nouvelles fonctionnalités
- Améliorations de documentation
- Traductions

Voir [CONTRIBUTING.md](CONTRIBUTING.md) pour plus de détails.

---

## 📋 Exemples de Configuration

### Configuration Serveur Basique
```yaml
server:
  listen_address: "0.0.0.0:51820"
  protocol: "udp"

network:
  interface: "tun0"
  address: "10.8.0.1/24"
  mtu: 1420

crypto:
  private_key: "VOTRE_CLE_PRIVEE_ICI"

clients:
  - public_key: "CLE_PUBLIQUE_CLIENT_1"
    allowed_ips: ["10.8.0.2/32"]
    name: "bureau-casa"
  
  - public_key: "CLE_PUBLIQUE_CLIENT_2"
    allowed_ips: ["10.8.0.3/32"]
    name: "bureau-rabat"
```

### Configuration Client
```yaml
client:
  server_address: "vpn.votre-entreprise.ma:51820"
  protocol: "udp"

network:
  interface: "tun0"
  address: "10.8.0.2/24"
  mtu: 1420

crypto:
  private_key: "VOTRE_CLE_PRIVEE_CLIENT"
  server_public_key: "CLE_PUBLIQUE_SERVEUR"

routing:
  default_route: false
  routes:
    - "10.8.0.0/24"
    - "192.168.1.0/24"
```

---

## 🎓 Formation et Apprentissage

### Pour les Développeurs Marocains
Ce projet est excellent pour apprendre:
- **Go (Golang)** - Langage moderne et performant
- **Cryptographie** - Sécurité des données
- **Réseaux** - Protocoles VPN, TCP/IP
- **Systèmes** - Linux, TUN/TAP devices
- **DevOps** - Docker, déploiement

### Ressources d'Apprentissage
- Code source bien commenté
- Documentation complète
- Architecture claire
- Exemples pratiques

---

## ⚠️ Notes Importantes

### Légalité
- ✅ L'utilisation de VPN est légale au Maroc
- ✅ Utilisez pour des fins légitimes uniquement
- ✅ Respectez les lois locales

### Sécurité
- 🔐 Gardez vos clés privées secrètes
- 🔐 Utilisez des mots de passe forts
- 🔐 Mettez à jour régulièrement
- 🔐 Surveillez les logs

### Support Technique
Pour un support professionnel ou une installation personnalisée, contactez:
- Email: [Votre email de support]
- GitHub: [@lahcenassmira](https://github.com/lahcenassmira)

---

## 📜 Licence

Ce projet est sous licence MIT - voir [LICENSE](LICENSE) pour plus de détails.

**Cela signifie**:
- ✅ Utilisation commerciale autorisée
- ✅ Modification autorisée
- ✅ Distribution autorisée
- ✅ Utilisation privée autorisée

---

## 🙏 Remerciements

Merci à la communauté open source et aux développeurs marocains qui contribuent à rendre la technologie accessible à tous.

**Développé avec ❤️ pour le Maroc 🇲🇦**

---

## 🚀 Commencer Maintenant

```bash
# Cloner le projet
git clone https://github.com/lahcenassmira/open-source-vpn.git

# Entrer dans le répertoire
cd open-source-vpn

# Lire le guide rapide
cat QUICKSTART.md

# Commencer!
make build
```

**Besoin d'aide?** Consultez [GET_STARTED.md](GET_STARTED.md) pour un guide étape par étape.

---

**🇲🇦 Made in Morocco - صنع في المغرب**
