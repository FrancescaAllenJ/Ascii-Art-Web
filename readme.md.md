# ASCII Art Web

## 📜 Description

This project is a web-based version of the `ascii-art` CLI program written in Go.  
It allows users to enter text, choose from three different ASCII art banner styles, and view the result rendered on a webpage.

It uses standard Go libraries only and serves an HTML form with input and banner options, then displays the ASCII art using banner files.

---

## 👩‍💻 Author

- Francesca Allen (`ftafrial`)

---

## 🚀 Usage

### Run the server:

```bash
go run .

Then visit:
👉 http://localhost:8080

Instructions:
Type your desired text (e.g. Hello)

Select a banner style:

Standard

Shadow

Thinkertoy

Click “Generate ASCII”

View the ASCII art result on a separate page

Use the “Generate another” link to go back

🧠 Implementation Details
The web server uses html/template to render HTML pages.

ASCII art generation is handled by a custom package in asciiart/generate.go.

Banner files are stored in /banners/*.txt and must contain 855 lines (95 printable characters × 9 lines each).

Each character is mapped by:

go
Copy
Edit
start := int(rune - 32) * 9
Result is shown via /ascii-art POST response, rendered in <pre> tags for formatting.

📁 Project Structure
go
Copy
Edit
ascii-art-web/
├── main.go
├── asciiart/
│   └── generate.go
├── templates/
│   ├── index.html
│   └── result.html
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
├── go.mod
├── go.sum
└── README.md
🧪 Status
✅ Fully working
✅ Banner switching works
✅ Input/output handled correctly
✅ HTTP error handling included
✅ Terminal logs included for debugging
✅ Gitea version tagged v1.0-audit-ready

🔧 Something Wrong?
Submit an issue
Or even better:
Propose a change!

---

## 📜 Commit History

This project was developed using Git, with a complete commit history available in the [Gitea repository](https://learn.01founders.co/git/ftafrial/ascii-art-web).  
Each commit reflects a real step in the development process — from setup, through implementation, bug fixing, styling, and final audit preparation.

The `feature/ascii-integration` branch contains the full working version, and the `v1.0-audit-ready` tag marks the project completion point.
