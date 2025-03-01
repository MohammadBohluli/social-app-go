package db

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"math/rand/v2"

	"github.com/MohammadBohluli/social-app-go/internal/store"
)

var usernames = []string{
	"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hank", "Ivy", "Jack",
	"Karen", "Leo", "Mona", "Nathan", "Olivia", "Paul", "Quincy", "Rachel", "Steve", "Tracy",
	"Uma", "Victor", "Wendy", "Xander", "Yvonne", "Zane", "Abby", "Ben", "Catherine", "Daniel",
	"Emma", "Fred", "Gina", "Henry", "Isla", "Jake", "Kylie", "Liam", "Megan", "Noah",
	"Oscar", "Penny", "Quinn", "Riley", "Sam", "Tina", "Ulysses", "Vera", "Will", "Xenia",
	"Yasmine", "Zack", "Aaron", "Bella", "Chris", "Diana", "Ethan", "Fiona", "George", "Holly",
	"Ian", "Julia", "Kevin", "Laura", "Mike", "Natalie", "Owen", "Pam", "Quinton", "Rose",
	"Scott", "Tessa", "Umar", "Vanessa", "Wayne", "Xiomara", "Yuri", "Zelda", "Adrian", "Becky",
	"Caleb", "Delilah", "Elliot", "Faith", "Gordon", "Heather", "Isaac", "Jasmine", "Kyle", "Linda",
	"Mason", "Nina", "Oliver", "Paula", "Ron", "Sophia", "Trevor", "Ursula", "Vince", "Willa",
}

var titles = []string{
	"Mastering Go Basics",
	"Concurrency in Go",
	"Building APIs with Go",
	"Go vs. Python: A Comparison",
	"Understanding Goroutines",
	"Go for Backend Development",
	"Handling Errors in Go",
	"Go Modules Explained",
	"JWT Authentication in Go",
	"Testing in Go: Best Practices",
	"Go and Docker Integration",
	"Using gRPC with Go",
	"Web Scraping with Go",
	"Go Design Patterns",
	"Database Management in Go",
	"Logging in Go Applications",
	"Building a CLI in Go",
	"Go Performance Optimization",
	"Deploying Go Applications",
	"Security Best Practices in Go",
}

var contents = []string{
	"The future of AI and its impact on jobs",
	"How to optimize SQL queries for performance",
	"A deep dive into WebSockets in Go",
	"Understanding the SOLID principles in software development",
	"10 must-know Linux commands for developers",
	"How to build a REST API with NestJS",
	"Exploring GraphQL vs REST: Pros and Cons",
	"Understanding Docker and containerization",
	"The importance of clean code and best practices",
	"How to implement caching for better performance",
	"An introduction to event-driven architecture",
	"The role of message queues in microservices",
	"A beginner’s guide to CI/CD pipelines",
	"How to secure your API endpoints effectively",
	"Introduction to distributed databases",
	"Understanding the CAP theorem in database systems",
	"Why logging and monitoring are essential in production",
	"A step-by-step guide to writing unit tests",
	"How to implement OAuth2 authentication",
	"The basics of Web3 and blockchain development",
}

var tags = []string{
	"golang", "backend", "api", "web-development", "database",
	"microservices", "docker", "cloud", "testing", "security",
	"devops", "graphql", "authentication", "performance",
	"design-patterns", "linux", "caching", "logging", "monitoring", "CI/CD",
}

var comments = []string{
	"Great explanation! Thanks for sharing.",
	"I disagree with some points, but overall a solid post.",
	"Can you provide a code example for better understanding?",
	"This really helped me, appreciate it!",
	"What are the best practices for this approach?",
	"I tried this, but ran into an error. Any suggestions?",
	"Amazing content! Looking forward to more.",
	"Could you compare this with an alternative method?",
	"Thanks! This saved me a lot of time.",
	"I think there's a typo in the third paragraph.",
	"Can this be optimized further?",
	"Nice work! Any real-world use cases?",
	"Does this approach scale well?",
	"Would love to see a follow-up on this topic.",
	"This is outdated, any updates on recent changes?",
	"You made it so easy to understand!",
	"I implemented this and it works perfectly!",
	"Can you explain why this works better than other methods?",
	"I’ve seen mixed opinions on this. What’s your take?",
	"This deserves more attention! Well done.",
}

func Seed(s store.Storage, db *sql.DB) {
	ctx := context.Background()
	tx, _ := db.BeginTx(ctx, nil)

	users := generateUsers(10)
	for _, user := range users {
		if err := s.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("❌Error creating user seed: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(20, users)
	for _, post := range posts {

		if err := s.Posts.Create(ctx, post); err != nil {
			log.Println("❌Error creating post seed: ", err)
			return
		}
	}

	comments := generateComments(50, users, posts)
	for _, comment := range comments {
		if err := s.Comments.Create(ctx, *comment); err != nil {
			log.Println("❌Error creating comment seed: ", err)
			return
		}
	}

	log.Println("✅ Seed is completly")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123456",
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {

		user := users[rand.IntN(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.IntN(len(titles))],
			Content: contents[rand.IntN(len(contents))],
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}

	return cms
}
