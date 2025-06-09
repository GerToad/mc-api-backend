
---

## API Backend (Go) README  
```markdown
# Minecraft Hosting API (Go)

## Overview  
This service provides a RESTful API for provisioning and managing Minecraft server instances, user accounts, billing, and AI-driven features. It follows clean architecture principles and integrates with Terraform for infrastructure orchestration. :contentReference[oaicite:17]{index=17}

## Features  
- **Server Management**: Create, update, and delete Minecraft server instances. :contentReference[oaicite:18]{index=18}  
- **User Authentication**: JWT-based auth with secure middleware. :contentReference[oaicite:19]{index=19}  
- **Terraform Integration**: Triggers Terraform runs via a job queue for infra provisioning.   
- **AI Services Orchestration**: Endpoints to enable/disable DDoS protection and support chatbots.   

## Directory Structure  

api-backend/
├── cmd/
│ └── server/ # Main application entry point
├── internal/
│ ├── handlers/ # HTTP handlers
│ ├── services/ # Business logic & Terraform orchestration
│ └── middleware/ # Auth, logging, metrics
├── pkg/
│ └── models/ # Shared data models
├── configs/
│ └── config.yaml # Application configuration schema
├── Dockerfile
├── go.mod
└── README.md

- Uses a **Hexagonal Architecture** to separate external adapters (HTTP) from core logic.   
- Configuration is loaded from YAML and environment variables.   

## Prerequisites  
- Go 1.20+ installed. :contentReference[oaicite:24]{index=24}  
- Docker (for containerized deployments). :contentReference[oaicite:25]{index=25}  
- Access to a message queue (e.g., RabbitMQ, AWS SQS) for Terraform job orchestration.   

## Installation 
1. **Clone the Repository**  
   ```bash
   git clone https://github.com/your-org/api-backend.git
   cd api-backend
