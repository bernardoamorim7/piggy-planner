# Piggy Planner

The self hosted and easy to use finance tracker.

Piggy Planner is a simple budgeting application that allows users to track their expenses and income. The application is built using Go, Templ, HTMX, Chart.js, SQLite and DaisyUI with TailwindCSS.

## Usage

Just download the latest release and run the executable. The application will start a web server on `localhost:8777` unless you specify a different address using an environment variable called `PIGGY_PORT`. The application will create a SQLite database in the current working directory called `piggy_planner.db`.

## Building from Source

To build the application from source, you will need:

- Go 1.22 or later
- Node.js 20 or later
- Templ

First, clone the repository:

```bash
git clone "https://github.com/bernardoamorim7/piggy-planner.git"
```

Install the required Node.js packages:

```bash
npm install
```

Install templ using the following command:

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

Now you can either use the Makefile 

```bash
make build
```

or build the application manually, generating the templ files, the CSS file and building the Go application:

```bash
templ generate && npx tailwindcss -i ./web/assets/css/input.css -o ./web/assets/css/tailwind.css --minify && go build -o piggy-planner main.go 
```

## Contributing

Pull requests are welcome, but this project is a personal project and I may not have the time to review and merge them. If you find a bug or have a feature request, please open an issue.


## License

This project is licensed under the [MIT License](LICENSE).