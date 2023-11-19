## All The Highlights: Centralizing Readwise Highlights

**All The Highlights** streamlines the management and accessibility of book highlights from Readwise. It serves as a crucial bridge between collecting insights and integrating them into personal knowledge systems. By automating the fetching, processing, and displaying of highlights, this tool simplifies the task of managing reading insights, making them easily accessible and usable for further personal or professional development.

### Features:

- Fetches book data and highlights from the Readwise Export API.
- Filters out books without highlights.
- Generates an HTML representation of the books and their associated highlights.
- Allows easy copy & paste functionality for each highlight and book metadata.
- Automated deployment to GitHub Pages using GitHub Actions.

### How to Use:

1. **Set Up**: Clone this repository to your local machine or GitHub workspace.

2. **Configuration**: Ensure you have a Readwise token. You will need to set this as an environment variable named `READWISE_TOKEN`.

3. **Run Locally**: Execute the main program using the command `go run main.go`.

4. **Deployment**: The project is configured with a GitHub Actions workflow (`deploy.yml`) that automatically builds and deploys the generated HTML to GitHub Pages upon push to the `main` branch or on a daily schedule.

5. **Viewing the Output**: Once deployed, the highlights can be viewed on GitHub Pages at the configured URL.

### Dependencies:

- Go programming language (tested with version 1.16).
- Readwise Export API for fetching the book data and highlights.
- The [grr](https://github.com/klausbreyer/grr) library for rendering HTML templates.

### Contributing:

If you have suggestions, bug reports, or feature requests, please open an issue on the GitHub repository. Contributions via pull requests are also welcome!
