# Goapi

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/andycai/goapi)

Goapi is a comprehensive web-based management platform designed for Unity
game development teams. It provides a suite of tools and services to streamline
the game development workflow, including version control management,
configuration management, and development tools integration.

## Features

- **Version Control Integration**
  - SVN and Git support
  - Browse repositories
  - Manage branches and commits
  - Code review functionality

- **Game Configuration Management**
  - Game configuration editor
  - Configuration version control
  - Multi-environment support
  - Configuration validation

- **Development Tools**
  - Luban configuration tool integration
  - Image asset management
  - File management system
  - Bug tracking system

- **Team Collaboration**
  - User management and authentication
  - Role-based access control
  - Team notes and documentation
  - Activity logging

- **CI/CD Integration**
  - Unity build automation
  - Build task management
  - Server configuration management
  - Deployment automation

## Project Structure

```
.
├── core/           # Core application framework
├── modules/        # Feature modules
│   ├── adminlog/   # Admin activity logging
│   ├── browse/     # Repository browser
│   ├── bugtracker/ # Bug tracking system
│   ├── citask/     # CI/CD task management
│   ├── filemanager/# File management
│   ├── gameconf/   # Game configuration
│   ├── git/        # Git integration
│   ├── luban/      # Luban tool integration
│   ├── svn/        # SVN integration
│   └── ...         # Other modules
├── models/         # Data models
├── utils/          # Utility functions
├── templates/      # HTML templates
├── public/         # Static assets
└── sql/           # Database schemas
```

## Technology Stack

- **Backend**
  - Go 1.22+
  - Fiber web framework
  - GORM database ORM
  - HTML template engine

- **Frontend**
  - HTML/CSS/JavaScript
  - Modern web components
  - Responsive design

## Getting Started

### Prerequisites

- Go 1.22 or higher
- MySQL/MariaDB
- Git
- SVN (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/andycai/unitool.git
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure the application:
   - Copy `conf.toml.example` to `conf.toml`
   - Update the configuration settings

4. Build the application:
   ```bash
   ./build.sh    # For Unix-like systems
   build.bat     # For Windows
   ```

### Running

1. Start the server:
   ```bash
   ./start.sh    # For Unix-like systems
   start.bat     # For Windows
   ```

2. Access the web interface at `http://localhost:8080`

### Configuration

The main configuration file `conf.toml` includes settings for:

- Server configuration
- Database connection
- Version control systems
- Authentication settings
- Build paths and options

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## Contributing

Contributions are welcome! Please feel free to submit pull requests.

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
