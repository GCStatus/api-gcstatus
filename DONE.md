# 📝 DONE List

## 🚀 Project Roadmap

### MVP (Minimum Viable Product)

- [x] Set up project repository
  - [x] Initialize a new Git repository
  - [x] Configure README.md, .gitignore, and basic project structure
- [x] Initial project setup
  - [x] Set up Go and Air
  - [x] Set up go lint
  - [x] Set up Dockerfile
  - [x] Set up entrypoint and supervisord
  - [x] Set up releaserc and commitlint
  - [x] Implement hexagonal structure
- [x] Create the project CI
  - [x] Create a CI step to run tests, lint and vet
  - [x] Create a CI step to mark the PR as ready to merge (only on PRs)
  - [x] Create a CI step to run the releaserc (only on main branch push)
  - [x] Create a CI step to build the docker image and push to Docker Hub
  - [x] Create a CI step to run the lint on PR title (to trigger the releaserc and deploy)
  - [x] Create a CI step to deploy the docker image into EC2 server
    - [x] Setup SSH using SSH key
    - [x] Remove old images from EC2 machine
    - [x] Set up the environment variables through CI with docker and github secrets
    - [x] Set up container run
- [x] Deployment
  - [x] Configure Netlify for continuous deployment
  - [x] Set up custom domain and SSL certificate with certbot
  - [x] Set up nginx server
- [x] AWS
  - [x] Integrate EC2
  - [x] Integrate RDS
  - [x] Integrate SES
  - [x] Integrate ElastiCache with Redis OSS
- [x] Auth
  - [x] Create the login method
  - [x] Create the register method
  - [x] Create the password forgot method
  - [x] Create the password reset method
  - [x] Create the password reset method from user profile (with current password validation)
- [x] Documentation
  - [x] Create a detailed readme and how to run API locally
  - [x] Create a contribution guide for open source contributors
- [x] Experience for users
  - [x] User should have the experience quantity
- [x] Coins
  - [x] User should have the coins quantity
  - [x] Create user wallet
  - [x] Migrate the user coins quantity for a has one relation for wallet

### Post MVP
