ASCII Art Web
📜 Description

A web GUI for the ascii-art project written in Go.
Type text, choose a banner (standard, shadow, thinkertoy), and get ASCII art in your browser.

Uses only Go standard library packages.

Banners are embedded via go:embed (no runtime file reads).

Handles Windows/macOS/Linux newlines and multi-line input (with a blank separator line between rendered blocks).

👩‍💻 Author

Francesca Allen (ftafrial)

🧰 Requirements

Go 1.21+ (or compatible)

🚀 Quick Start
Run the server
go build
go run .


Open: http://localhost:8080

How to use

Type your text (supports multi-line).

Pick a banner: standard, shadow, or thinkertoy.

Submit to render the ASCII art on the result page.

Optional: CLI helper (great for audit checks)
go build -o ascii-check ./cmd/ascii-check
./ascii-check -banner standard -text "Hello\nThere"

🔌 HTTP Endpoints

GET / → main page (HTML form)

405 Method Not Allowed if not GET

404 Not Found if path ≠ / or templates missing

POST /ascii-art → render ASCII art and return result page

400 Bad Request for missing/invalid inputs (unknown banner, unsupported chars, malformed banners)

405 Method Not Allowed if not POST

404 Not Found if path ≠ /ascii-art or templates missing

500 Internal Server Error for unexpected failures

Quick probes (with server running):

curl -i -X POST http://localhost:8080/ | head -n 1          # 405
curl -i -X POST http://localhost:8080/ascii-art | head -n 1 # 400
curl -i -s -d "inputText=Hi&banner=bogus" http://localhost:8080/ascii-art | head -n 1 # 400
curl -i http://localhost:8080/not-found | head -n 1         # 404

🧠 Implementation Details

Glyph set: printable ASCII 32..126.

Banner format: each glyph has 9 lines (8 visible rows + 1 separator);
total lines = 95 × 9 + 1 = 856 (the last line is blank).

Indexing: for rune r ∈ [32..126], start := (r-32) * 9; rows are start + 0..7.

Multi-line input: normalize \r\n/\r → \n.

Render each non-empty input line as an 8-row block.

Insert one blank separator line between blocks.

Explicit empty input lines add another blank line.

Empty input → empty output.

Embedding: //go:embed asciiart/banners/*.txt compiles banners into the binary.
Files are loaded from the embedded FS, normalized to LF, validated at 856 lines, and cached.

✅ Allowed Packages (Std-lib only)

This project intentionally uses only the standard library (e.g., net/http, html/template, embed, strings, log, sync).

Check:

# Should print nothing (no external modules)
go list -m -f '{{if not .Main}}{{.Path}}{{end}}' all

🧪 Tests

Unit tests cover:

Single-line rendering shape (top row check)

Unsupported characters error

Multi-line rendering with a blank separator

Empty input → empty output

Run:

go test -v ./asciiart

✅ How this passes the audit
1) Only standard packages
# Should print nothing (no external modules)
go list -m -f '{{if not .Main}}{{.Path}}{{end}}' all

2) Exact ASCII output checks (CLI helper)

Build the small CLI:

go build -o ascii-check ./cmd/ascii-check


Compare helper (Git Bash / macOS Terminal / WSL). Paste once:

compare() {
  # usage: compare "<input text>" <banner> <audit-file>
  local input="$1" banner="$2" expfile="$3"
  local actual expected
  actual="$(./ascii-check -banner "$banner" -text "$input" | sed 's/\r$//')"
  # audit blocks use "$" at end of each line to show EOL — strip it:
  expected="$(sed -E 's/[[:space:]]*\$$//' "$expfile" | sed 's/\r$//')"
  if diff -u <(printf "%s" "$expected") <(printf "%s" "$actual") >/dev/null; then
    echo "OK  - $banner  <<$input>>"
  else
    echo "FAIL- $banner  <<$input>>"
    diff -u <(printf "%s" "$expected") <(printf "%s" "$actual") || true
  fi
}


Create a tests/ folder and paste the exact blocks from the audit (keep the trailing $ markers). Examples:

mkdir -p tests

cat > tests/hello.standard.txt <<'EOF'
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
EOF

cat > tests/HELLO.standard.txt <<'EOF'
 _    _   ______   _        _         ____   $
| |  | | |  ____| | |      | |       / __ \  $
| |__| | | |__    | |      | |      | |  | | $
|  __  | |  __|   | |      | |      | |  | | $
| |  | | | |____  | |____  | |____  | |__| | $
|_|  |_| |______| |______| |______|  \____/  $
                                             $
                                             $
EOF

cat > tests/hello_n_there.standard.txt <<'EOF'
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $
EOF


Compare:

compare "hello" standard tests/hello.standard.txt
compare "HELLO" standard tests/HELLO.standard.txt
compare "Hello\nThere" standard tests/hello_n_there.standard.txt


Repeat for the other audit strings by saving each block into tests/*.txt and calling compare "<input>" standard tests/file.txt.
For shadow and thinkertoy, save their versions likewise and change the -banner argument.

3) HTTP status codes (Functional requirements)

With the server running (go run .), run:

curl -i -X POST http://localhost:8080/ | head -n 1          # 405
curl -i -X POST http://localhost:8080/ascii-art | head -n 1 # 400
curl -i -s -d "inputText=Hi&banner=bogus" http://localhost:8080/ascii-art | head -n 1 # 400
curl -i http://localhost:8080/not-found | head -n 1         # 404

📁 Project Structure
ascii-art-web/
├─ main.go
├─ go.mod
├─ templates/
│  ├─ index.html
│  └─ result.html
├─ asciiart/
│  ├─ generate.go
│  ├─ generate_test.go
│  └─ banners/
│     ├─ standard.txt       # 856 lines (final blank line)
│     ├─ shadow.txt         # 856 lines (final blank line)
│     └─ thinkertoy.txt     # 856 lines (final blank line)
└─ cmd/
   └─ ascii-check/
      └─ main.go            # small CLI helper (optional)

🧷 Troubleshooting

Banner file validation: “got 855 lines, want 856” → add a final newline to the banner file.

Unknown route returns 200: ensure handlers check r.URL.Path for exact matches (/ and /ascii-art) and return 404 otherwise.

Windows line endings: inputs are normalized automatically; banners are normalized during embed load.

Module import: module path is learn.01founders.co/git/ftafrial/ascii-art-web.
Import the package as:
import "learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"

🔓 Open Source / Reuse

This module is self-contained (no third-party deps) and easy to reuse in other apps needing ASCII-art rendering via web or CLI.

🏷️ Tags / Branches

Working branch: feature/ascii-integration

Release tag: v1.0-audit-ready