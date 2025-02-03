# Public Music Playlist Viewer API
## Overview
The Public Music Playlist Viewer API is a real-time backend service designed for public spaces like cafes, restaurants, or retail stores. It provides a robust API for managing and broadcasting music playlist information, enabling client applications to display currently playing songs and upcoming tracks in real-time.
Features

- Real-time playlist state management
- WebSocket support for live updates
- RESTful API endpoints for playlist data
- Music service provider integration support
- Authentication system for admin operations
- Rate limiting and connection management
- Comprehensive logging system
- Health check endpoints
- Metrics collection for monitoring

## Technical Stack

- Go 1.22
- Standard library WebSocket implementation
- SQL database for persistent storage
- Redis for caching (optional)
- Music service integration support
- In-memory state management

## System Requirements

- Go 1.22 or higher
- Minimum 1GB RAM
- PostgreSQL 14+ or MySQL 8+
- Redis 6+ (optional)

## Architecture
The system implements a publisher-subscriber pattern:

- Music source updates are received via SSE
- Backend maintains the current state
- Changes are broadcast to clients
- State persistence ensures system reliability

## API Documentation

- Authentication: JWT
- Base URL: `/api/v1`

## Key Endpoints:

### Playlist Management:

- GET `/playlist/current` - Get current track
- GET `/playlist/queue` - Get upcoming tracks
- WS `/ws/playlist` - Real-time updates


### Admin Operations:

- POST `/admin/playlist` - Update playlist
- DELETE `/admin/track/{id}` - Remove track
- PUT `/admin/track/reorder` - Reorder tracks


### System Operations:

- GET `/health` - System health check
- GET `/metrics` - System metrics (protected)



## Setup Instructions

- Clone the repository
- Configure environment variables
- Set up the database
- Build and run the application

## Development

- Install Go 1.22
- Run `go mod tidy`
- Copy `.env.example` to `.env` and configure
- Start the development server: `go run main.go`

## Production Deployment

- Build the application
- Configure environment variables
- Set up a reverse proxy
- Configure SSL/TLS
- Start the service
- Monitor system health

## Security Considerations

- JWT-based authentication
- Rate limiting implementation
- CORS policy configuration
- Input validation and sanitization
- Secure WebSocket connections
- Database connection encryption

## Scaling Considerations

- Horizontal scaling support
- Database connection pooling
- Redis caching layer
- Load balancer configuration
- Realtime connection management

## Contributing

- Fork the repository
- Create a feature branch
- Commit your changes
- Push to the branch
- Create a Pull Request

## License
MIT License