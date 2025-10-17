# Bidlancer

**Bidlancer** is a backend service for a freelancing platform that connects project owners with freelancers through a transparent and competitive bidding system.

---

## What Is Bidlancer?

Bidlancer is a web application that provides a platform for **freelancers** and **project owners** to collaborate efficiently and transparently.

### How It Works

1. **Project Creation** – A project owner creates a new post describing their project.  
2. **Project Discovery** – Freelancers use our **advanced search** feature to find projects that match their skills and interests.  
3. **Bidding** – Freelancers place **bids** on projects they want to work on.  
4. **Selection** – The project owner reviews the bids and selects a freelancer (or team) to collaborate with.  

---

## Key Features

- **Transparency First**  
  All authenticated users can view bid amounts and bidder profiles.  
  Both freelancers and project owners have access to detailed project histories, reviews, and comments.  

- **Team Collaboration**  
  Freelancers can form **teams** to work together on larger projects.  
  Our system ensures full clarity about **who** is part of a bidding team and **who** is responsible for what.  

- **Competitive Bidding System**  
  Real-time updates and open visibility into bids foster a fair and dynamic marketplace.  

---

## Tech Stack

| Layer | Technology |
|-------|-------------|
| **Language** | Go |
| **Database** | PostgreSQL |
| **Caching** | Redis |
| **Message Broker** | Kafka |
| **CDC (Change Data Capture)** | Debezium |
| **Search Engine** | Elasticsearch |
| **Containerization** | Docker |

---

## Architecture Notes

Bidlancer is designed with **high scalability** and **microservice-friendly** principles in mind.  
The project is still **under active development**.

---

## Roadmap

- [ ] Add real-time notification service via Kafka    
- [ ] Kubernetes
- [ ] Monitoring
- [ ] Email microservice
- [ ] Email notifications on projects
- [ ] Email notifications for when someone is out bidded
- [ ] AI integration (far fetched)

---

## Contributing

We welcome contributions!  
Please open an issue or submit a pull request if you'd like to improve Bidlancer.

---

## License

This project is licensed under the **MIT License** – see the [LICENSE](LICENSE) file for details.

