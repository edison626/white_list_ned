package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	username       = "123"  //用户
	password       = "123"  //密码
	port           = "8812" //端口
	filePath       = ".url"
	filePathDomain = ".domain"
	defaultText    = "# Enter your content here"
)

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/editor", editorHandler)
	http.HandleFunc("/save", saveHandler)

	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form.Get("user")
		pass := r.Form.Get("pass")

		if user == username && pass == password {
			// 登录成功，设置登录状态Cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "loggedin",
				Value: "true",
				Path:  "/",
			})

			http.Redirect(w, r, "/editor", http.StatusSeeOther)
		} else {
			fmt.Fprint(w, "Invalid username or password.")
		}
		return
	}
	tmpl := `
	<!DOCTYPE html>
	<html>

	<head>
		<meta charset="UTF-8">
	<title>域名加白中心</title>
		<style>
			body,
			html {
				height: 100%;
				margin: 0;
				padding: 0;
				background-color: #d2d6de;
			}

			.container {
				width: 300px;
				height: 200px;
				position: absolute;
				top: 50%;
				left: 50%;
				transform: translate(-50%, -50%);
				background-color: white;
				border-radius: 2px;
				box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);
			}

			.title {
				text-align: center;
				padding-top: 30px;
				font-size: 18px;
				font-weight: bold;
			}

			.form-group {
				margin-top: 20px;
				margin-left: 15px;
				margin-right: 15px;
			}

			.form-group input {
				width: 100%;
				height: 30px;
				line-height: 30px;
				padding: 0 10px;
				border: 1px solid #ccc;
				border-radius: 2px;
				box-sizing: border-box;
			}

			.form-group button {
				width: 100%;
				height: 32px;
				line-height: 32px;
				margin-top: 15px;
				background-color: #209eff;
				color: white;
				border: none;
				border-radius: 2px;
				cursor: pointer;
			}

			.form-group button:hover {
				background-color: #1089ff;
			}
		</style>
	</head>

	<body>
		<div class="container">
			<div class="title">IP加白</div>
			<form method="post" action="/">
				<div class="form-group">
					<input type="text" name="user" placeholder="Username">
				</div>
				<div class="form-group">
					<input type="password" name="pass" placeholder="Password">
				</div>
				<div class="form-group">
					<button type="submit">Login</button>
				</div>
			</form>
		</div>
	</body>

	</html>
	`

	fmt.Fprint(w, tmpl)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("loggedin")
	if err != nil || cookie.Value != "true" {
		// 未登录，重定向到登录页面
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 从文件中读取内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Error reading file.", http.StatusInternalServerError)
		return
	}

	contentDomain, err := os.ReadFile(filePathDomain)
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Error reading file.", http.StatusInternalServerError)
		return
	}

	data := struct {
		Content       string
		ContentDomain string
		Success       bool
	}{
		Content:       string(content), // 将读取的内容赋值给结构体的字段
		ContentDomain: string(contentDomain),
		Success:       r.URL.Query().Get("success") == "true",
	}

	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>欢迎使用YQ</title>
			<style>
				.form-section {
					display: flex;
					justify-content: space-between;
					margin-bottom: 10px; /* Adds some space between the form sections */
				}
				.form-section div {
					display: flex;
					flex-direction: column;
					margin-right: 20px; /* Adds some space between the textareas */
				}
				body {
					display: flex;
					flex-direction: column;
					justify-content: center;
					align-items: center;
					height: 100vh;
				}
				textarea {
					width: 500px;
					height: 300px;
				}
				button {
					margin-top: 10px;
				}
				.success {
					display: none;
					color: green;
					margin-top: 10px;
				}
				.error {
        			color: red; /* This makes the text color red */
        			margin-top: 10px; /* Adds spacing above the error message */
        			display: none; /* Initially hides the error message */
				}
				.form-group input {
					width: 100%;
					height: 30px;
					line-height: 30px;
					padding: 0 10px;
					border: 1px;
					border-radius: 2px;
					box-sizing: border-box;
				}
			</style>
			<script>
			window.addEventListener('DOMContentLoaded', (event) => {
				const form = document.querySelector('form');
				const textarea = document.querySelector('textarea[name="content"]');
				const successMessage = document.querySelector('.success');
				const errorMessage = document.createElement('p');
				errorMessage.classList.add('error');
				form.parentNode.insertBefore(errorMessage, form);
	
				form.addEventListener('submit', (e) => {
					const lines = textarea.value.split(',');
					// const ipPattern = /^(?:\d{1,3}\.){3}\d{1,3}$/;
					const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
					for (let line of lines) {
						if (!ipPattern.test(line)) {
							e.preventDefault(); // Prevent form submission
							errorMessage.textContent = "错误 IP 地址格式: " + line;
							errorMessage.style.display = 'block';
							return; // Stop checking further if any line is invalid
						}
					}
					errorMessage.style.display = 'none'; // Hide error message if all lines are valid
				});
	
				if (successMessage) {
					successMessage.style.display = 'block';
					setTimeout(() => {
						successMessage.style.display = 'none';
					}, 2000);
				}
			});
		</script>
		</head>
		<body>
			<h2>输入加白域名</h2>
			{{ if .Success }}
				<p class="success">保存成功！</p>
			{{ end }}
			<form action="/save" method="post">
			<div class="form-section">
				<div>
					<label for="content">加白IP</label>
					<textarea name="content">{{ .Content }}</textarea>
				</div>
				<div>
					<label for="contentdomain">域名</label>
					<textarea name="contentdomain">{{ .ContentDomain }}</textarea>
				</div>
			</div>
			<button type="submit">保存启动</button>
		</form>
		</form>
		</body>
		</html>
	`

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error parsing template.", http.StatusInternalServerError)
		return
	}

	// err = t.Execute(w, dataDomain)
	// if err != nil {
	// 	log.Println("Error executing template:", err)
	// 	http.Error(w, "Error executing template.", http.StatusInternalServerError)
	// }

	err = t.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template.", http.StatusInternalServerError)
	}
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content := r.FormValue("content")
		fmt.Printf("\n加白IP:%s", content)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			log.Println("Error writing file:", err)
			http.Error(w, "Error writing file.", http.StatusInternalServerError)
			return
		}

		contentdomain := r.FormValue("contentdomain")
		fmt.Printf("\n域名: %s", contentdomain)
		err2 := os.WriteFile(filePathDomain, []byte(contentdomain), 0644)
		if err2 != nil {
			log.Println("Error writing file:", err2)
			http.Error(w, "Error writing file.", http.StatusInternalServerError)
			return
		}

		cmd := exec.Command("bash", ".dd.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Println("Error executing command:", err)
			http.Error(w, "Error executing command.", http.StatusInternalServerError)
			return
		}

		redirectURL := "/editor?success=true"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
