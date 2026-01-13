# Getting Started with Superfly

Welcome! This guide helps you choose the right path based on your needs.

---

## 🎯 What Do You Want to Do?

### I want to deploy an app on a remote server → [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)
**Best for**: Production deployments, first-time users

You'll learn:
- ✅ How to set up a Debian 13 server from scratch
- ✅ Install all dependencies automatically
- ✅ Deploy your first app with HTTPS
- ✅ Manage apps via API
- ✅ Troubleshoot common issues
- ✅ Security hardening

**Time**: 15-20 minutes

---

### I just want the commands → [TLDR.md](TLDR.md)
**Best for**: Experienced users who know what they're doing

Minimal explanation, just commands:
```bash
./dev-setup.sh
make init && go mod tidy && make migrate && make sqlc-generate && make build
sudo systemctl start superfly-api
curl -X POST localhost:8080/api/apps -d '{...}'
```

**Time**: 10 minutes

---

### I want to understand the architecture → [VISUAL_GUIDE.md](VISUAL_GUIDE.md)
**Best for**: Visual learners, architects, team leads

You'll see:
- 📊 System diagrams
- 🔄 Data flow charts
- 🗺️ Network topology
- 📈 Deployment timeline
- 🎨 Command visuals

**Time**: 5 minutes reading

---

### I'm developing Superfly locally → [QUICKSTART.md](QUICKSTART.md)
**Best for**: Contributors, local testing

You'll get:
- ✅ Fast local setup
- ✅ Live reload workflow
- ✅ Testing with curl
- ✅ Making code changes

**Time**: 10 minutes

---

### I need API documentation → [API.md](API.md)
**Best for**: Integration, automation, CI/CD

Complete reference:
- 📚 All endpoints
- 📋 Request/response formats
- ⚠️ Error codes
- 💡 Usage examples
- 🔐 Authentication (future)

---

### I want real-world examples → [EXAMPLES.md](EXAMPLES.md)
**Best for**: Learning by example

Covers:
- 🌐 Static websites
- 🔌 APIs (Node.js, Python, Go)
- 🗄️ Databases (PostgreSQL, Redis)
- 🏗️ Microservices
- 📊 Full-stack apps

---

### I want to contribute code → [DEVELOPMENT.md](DEVELOPMENT.md)
**Best for**: Contributors, feature developers

Learn about:
- 🏗️ Project structure
- 🔧 Development workflow
- 🧪 Testing
- 📦 Database migrations
- 🎨 Code generation with sqlc

---

### I want to understand how it works → [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)
**Best for**: Curious minds, maintainers

Deep dive into:
- 📁 File organization
- 🔄 Data flow
- 🎯 Design decisions
- 🛠️ Tech stack
- 🔮 Future features

---

## 🚦 Quick Decision Tree

```
Start Here
    │
    ├─ Do you have a server ready?
    │   │
    │   ├─ Yes → Want detailed explanations?
    │   │   │
    │   │   ├─ Yes → [DEPLOYMENT_GUIDE.md]
    │   │   └─ No  → [TLDR.md]
    │   │
    │   └─ No → Setting up server or using local?
    │       │
    │       ├─ Server → [DEPLOYMENT_GUIDE.md]
    │       └─ Local → [QUICKSTART.md]
    │
    ├─ Want to understand first?
    │   │
    │   ├─ Visual learner → [VISUAL_GUIDE.md]
    │   └─ Text learner  → [PROJECT_STRUCTURE.md]
    │
    ├─ Need API docs?
    │   │
    │   └─ [API.md]
    │
    └─ Want examples?
        │
        └─ [EXAMPLES.md]
```

---

## 📚 All Available Guides

| Guide | Purpose | Time | Audience |
|-------|---------|------|----------|
| [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | Complete production setup | 20 min | Everyone |
| [TLDR.md](TLDR.md) | Just the commands | 10 min | Experienced users |
| [VISUAL_GUIDE.md](VISUAL_GUIDE.md) | Diagrams & flowcharts | 5 min | Visual learners |
| [QUICKSTART.md](QUICKSTART.md) | Local development | 10 min | Developers |
| [API.md](API.md) | API reference | Reference | API users |
| [EXAMPLES.md](EXAMPLES.md) | Real-world usage | Reference | New users |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Contributing guide | Reference | Contributors |
| [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) | Architecture deep-dive | 15 min | Architects |
| [README.md](README.md) | Project overview | 5 min | Everyone |

---

## 🎬 Recommended Learning Path

### For Complete Beginners
1. Start with [README.md](README.md) - Get the big picture
2. Read [VISUAL_GUIDE.md](VISUAL_GUIDE.md) - Understand the flow
3. Follow [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) - Deploy your first app
4. Check [EXAMPLES.md](EXAMPLES.md) - Deploy different types of apps
5. Explore [API.md](API.md) - Learn all the features

### For Experienced Developers
1. Skim [README.md](README.md) - Quick overview
2. Run through [TLDR.md](TLDR.md) - Get it running fast
3. Reference [API.md](API.md) - Integrate into your workflow
4. Check [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Understand internals

### For Contributors
1. Read [README.md](README.md) - Project goals
2. Follow [QUICKSTART.md](QUICKSTART.md) - Set up locally
3. Study [DEVELOPMENT.md](DEVELOPMENT.md) - Development workflow
4. Deep dive [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Codebase structure
5. Make changes and test!

---

## 💬 Common Questions

### Q: I'm not a DevOps person. Can I use this?
**A**: Yes! Start with [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md). It explains everything step-by-step.

### Q: I just want to see if it works. What's the fastest way?
**A**: Use [TLDR.md](TLDR.md). Copy-paste commands, you'll have an app running in 10 minutes.

### Q: I learn best from diagrams. Where should I start?
**A**: Check out [VISUAL_GUIDE.md](VISUAL_GUIDE.md) first, then [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md).

### Q: How do I deploy a specific type of app (Node.js, Python, etc)?
**A**: See [EXAMPLES.md](EXAMPLES.md). We have examples for all common scenarios.

### Q: I want to integrate this into my CI/CD. Where's the API docs?
**A**: [API.md](API.md) has complete documentation with curl examples.

### Q: Can I contribute? How do I set up locally?
**A**: Yes! Follow [QUICKSTART.md](QUICKSTART.md) then read [DEVELOPMENT.md](DEVELOPMENT.md).

### Q: What happens under the hood?
**A**: Read [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for a complete technical overview.

### Q: Something's not working. Help!
**A**: Check the Troubleshooting sections in [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md).

---

## 🆘 Need Help?

1. **Check the docs** - Most questions are answered in the guides above
2. **Search issues** - Someone might have had the same problem
3. **Ask in discussions** - Community support
4. **Open an issue** - Bug reports and feature requests

---

## 🚀 Ready to Start?

Pick your guide from above and let's deploy some apps! 

**Most popular starting point**: [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) 

**Fastest route**: [TLDR.md](TLDR.md)

**Best overview**: [VISUAL_GUIDE.md](VISUAL_GUIDE.md)

---

Good luck! 🎉
