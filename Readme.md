# Project Time Tracking System

[中文版](README_zh.md)

## Project Overview

The Project Time Tracking System is a web application developed in Go, designed to help companies or teams manage and track employee work hours. This system provides an intuitive interface for employees to submit their time entries, administrators to view statistical data, and supports exporting detailed Excel reports.

## Key Features

- Employee Time Submission: Employees can easily record daily work hours and project information
- Admin Dashboard: Administrators can view overall time statistics and project progress
- Project Management: Add, edit, and delete projects
- Employee Management: Add, edit, and delete employee accounts
- Data Export: Generate detailed Excel reports, including project and employee time statistics

## Tech Stack

- Backend: Go
- Database: SQLite
- Frontend: HTML, CSS (Tailwind CSS), JavaScript
- Dependency Management: Go Modules

## Installation Guide

1. Clone the repository:
   ```
   git clone https://github.com/VxNull/project-time-tracker.git
   ```

2. Navigate to the project directory:
   ```
   cd project-time-tracker
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Run the application:
   ```
   go run main.go
   ```

5. Access the application in your browser at `http://localhost:8080`

## Usage Instructions

### Administrators

1. Log in using the default admin account (username: admin, password: password123)
2. View overall time statistics in the dashboard
3. Manage projects and employee accounts
4. Export time reports

### Employees

1. Log in using the assigned account
2. Submit daily work hours
3. View personal time statistics

## Contribution Guidelines

We welcome all forms of contributions, including but not limited to:

- Reporting bugs
- Suggesting new features
- Improving documentation
- Submitting code fixes or new features

Please follow these steps:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Contact Information

If you have any questions or suggestions, please contact us through:

- Project Link: [https://github.com/VxNull/project-time-tracker](https://github.com/VxNull/project-time-tracker)
- Issue Tracker: [https://github.com/VxNull/project-time-tracker/issues](https://github.com/VxNull/project-time-tracker/issues)
