# 📝 TODO List

## 🚀 Project Roadmap

### MVP (Minimum Viable Product)

- [ ] User authentication
  - [ ] Implement OAuth for Google and Facebook
- [ ] Review the entire error throwing of platform to user-friendly
- [ ] Email verification
  - [ ] Block user access if email is not verified
  - [ ] If user change email, set verified_at to null and block access again
- [ ] Abstract all use cases and rules to services
- [ ] Integrate with Kibana for logs
- [ ] Create jobs to run some sevices async
  - [ ] Email send
  - [ ] Create secondary records, such as transactions or notifications
- [ ] Review all exceptions thrown
- [ ] Missions and titles
  - [ ] Review the action keys for missions and titles
  - [ ] Planning all missions and titles type
  - [ ] Planning all daily/weekly/monthly missions
  - [ ] Create association for games or DLCs
    - [ ] Create torrent websites (such as firgitl, skidrow etc)
    - [ ] Create the king of protections (such as Denuvo, Steam, GOG etc)
    - [ ] Create game developers (such as Game Science)
    - [ ] Create game publishers (such as Game Science)
    - [ ] Create requirement types
      - [ ] A requirement type should have a potential column, that should be enum, with minimum, maximum or recommended
      - [ ] A requirement type should have a type column, that should be enum, with windows, mac or linux
    - [ ] Create requirements and associate with requirement types
- [ ] Create game outline system
  - [ ] Add game DLCs
  - [ ] Add game critics
  - [ ] Add game torrents
  - [ ] Add game crack
  - [ ] Add game reviews
  - [ ] Add game galleries
  - [ ] Add game publishers
  - [ ] Add game developers
  - [ ] Add game requirements
  - [ ] Add game messages (for torrents section)
  - [ ] Add game support

### Post-MVP

- [ ] Testing & QA
  - [ ] Write unit tests for all packages
  - [ ] Set up end-to-end testing with Cypress or Playwright
- [ ] Documentation
  - [ ] Document API endpoints
  - [ ] Create an issue template
- [ ] Filters for users
  - [ ] Create a helper to clean bad words
    - [ ] Implement on user nickname
    - [ ] Implement on user name
    - [ ] Implement on user email
- [ ] Filters for posts, comments or messages
  - [ ] Posts
  - [ ] Comments
  - [ ] Messages (?)
- [ ] Refresh token
- [ ] MFA - 2FA

### Future Ideas

- [ ] Integration with external APIs
  - [ ] Pull in live data from gaming APIs like Twitch or Steam
  - [ ] Display trending streams or game stats
- [ ] Add some games on the website
- [ ] Quiz & polls: about the most exciting game release on future or something else
- [ ] Integrate with social medias
- [ ] Coin system
  - [ ] Create a system to reward user to do something on platform, like comment in a game
    - [ ] Block user to earn coins for commenting the same game (or reduce the amount for each task doing)
  - [ ] Add possibility to buy coins
  - [ ] Some quizzez can reward with coins
- [ ] Title system
  - [ ] Some titles could be earned by hitting a percentage of a quiz
- [ ] Orders
  - [ ] Make an order system to purchase coins
- [ ] Create another user profiles page
  - [ ] Create a page with another user details
  - [ ] Create a sysmtem to follow another user
  - [ ] Create a system of notifications to follow user posts
- [ ] Think about some leaderboard
  - [ ] Mark a game as played, make some integration to get hours played (if exists)
  - [ ] Leaderboard for GCStatus missions, titles, coins and level
- [ ] Create a ticket for support
  - [ ] The ticket can be used to get support
  - [ ] The ticket can be used to report some suspicious activity
- [ ] Create a staging API with github environments
- [ ] AWS
  - [ ] SNS
  - [ ] Lambda
- [ ] Create a chat between users
  - [ ] User can be able to chat another users
  - [ ] User can be able to create a group and chat them
  - [ ] User can be able to change the group name and avatar
  - [ ] User can be able to add and remove members (creator or admin)
  - [ ] User can be able to add admins on groups (owner only)
  - [ ] Chat should use realtime
- [ ] Create quizz that could reward with some coins and experience, maybe titles
- [ ] Award with coins and experience on comment, heart a game, or something else
  - [ ] Heart a game;
  - [ ] Make a comment on game details
  - [ ] Make a comment on blogs
  - [ ] Check how to prevent spam (award only once for each awardable)
- [ ] Create a method to receive all main data from API for HOME
  - [ ] Method should return notifications
  - [ ] Method should return home banners
  - [ ] Method should return 9 popular games
  - [ ] Method should return the next most awaited release (and should stay for one week as released)
  - [ ] Method should return 9 hot games
  - [ ] Method should return 9 most liked games
  - [ ] Method should return 9 upcoming games
- [ ] Coupons for coins purchase
- [ ] Admin system
  - [ ] ...
