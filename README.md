# 🎮 Welcome to GCStatus - The Ultimate Gamer Hub 🎮

### About GCStatus

**_Welcome to GCStatus, where gaming isn’t just a hobby—it's a way of life! Whether you're a casual button-masher or a hardcore pro, GCStatus is your ultimate portal to everything gaming. Find the latest news, reviews, and tips & tricks to level up your game. And don't worry, we've got plenty of memes to keep things lit 🔥._**

## Features

- 🕹️ Latest Gaming News - Stay updated with the hottest game releases and updates;
- 🎮 Game Reviews - Honest, in-depth reviews by real gamers, for real gamers;
- 🛠️ Tips & Tricks - Pro-level strategies to make sure you never rage-quit again;
- 💬 Community Forum - Discuss, debate, and maybe even find your next co-op buddy;
- 🚀 A lot more!

#### Because we know you'd rather be gaming than reading instructions, here's how to get GCStatus running on your local machine in no time!

## Prerequisites

**Before you start, make sure you have these bad boys installed:**

- Docker - [Get it here](https://docs.docker.com/engine/install/);
- Git - Because version control is life. [Get it here](https://git-scm.com/downloads)

## Installation

**Clone the repo:**

```bash
git clone https://github.com/felipebrsk/gcstatus-api.git
```

**Enter the project folder:**

```bash
cd gcstatus
```

**Install dependencies:**

```bash
docker buildx build --platform linux/amd64 -p 8080:8080 -t gcstatus-api .
```

**Run the container up command:**

```bash
docker run -d --name gcstatus-api -p 8080:8080 gcstatus-api
```

**Open your browser:**

_Navigate to http://localhost:8080 and voila! Welcome to GCStatus API._

If you want to enter the container to see or change something, just run:

```bash
docker exec -it gcstatus-api zsh
```

## Technologies Used

- ⚛️ **GoLang**: The backbone of your API, providing robust performance and concurrency handling;
- 🗃️ **MySQL**: The reliable database management system storing your application's data;
- 🔐 **GoLang JWT**: The library that powers your secure authentication with JSON Web Tokens;
- 🧩 **Redis**: The high-speed caching and queuing solution, enhancing performance and scalability;
- 🔄 **GORM**: The powerful ORM for Go, simplifying database interactions with elegant query handling;
- 🚀 **Gin**: The fast and flexible web framework that serves as the foundation of your API server;
- 🌐 **Gin CORS**: The middleware for managing CORS, ensuring smooth and secure cross-origin requests;
- 🔒 **X Crypto**: The library for encryption and decryption, safeguarding sensitive data with robust cryptographic techniques.

## 📊 Status

GCStatus API is currently under development. That means it's mostly stable, but we might still break things occasionally. If you find any bugs, you can either:

- Blame your internet connection 🐢
- Create an issue on GitHub 🐛

## 📚 Documentation & Tutorials

Check out our extensive documentation to get the most out of GCStatus:

- [API Reference](https://google.com) - For when you need to get technical.

## 🎉 Contributing

**Want to contribute? Awesome! Here's how you can get started:**

- Fork the repo;
- Create a new branch (e.g., feature/your-feature);
- Commit your changes (git commit -m 'feat: Add some feature');
- Push to the branch (git push origin -u HEAD);
- Open a Pull Request;
- We welcome PRs, especially those that fix bugs, improve performance, or add new features!

## 🛠 Maintenance

GCStatus API is maintained by a small but passionate team of gamers. We’re always looking for feedback, so feel free to reach out!

## 📄 License

This project is licensed under the MIT License.

### 🕹️ Game On and Have Fun!

Remember, life’s more fun when you’re gaming. So, get out there and show the world your skills. And if you’re not sure how to do something, just press all the buttons and hope for the best — it usually works! 😎
