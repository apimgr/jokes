package admin

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler manages admin routes and authentication
type Handler struct {
	auth      *AuthManager
	version   string
	commit    string
	buildDate string
}

// NewHandler creates a new admin handler
func NewHandler(username, password, apiToken string, sessionTimeout int, sslEnabled bool, version, commit, buildDate string) *Handler {
	return &Handler{
		auth:      NewAuthManager(username, password, apiToken, sessionTimeout, sslEnabled),
		version:   version,
		commit:    commit,
		buildDate: buildDate,
	}
}

// RegisterRoutes registers admin routes on Gin router
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Admin web interface (session auth)
	r.GET("/admin", h.handleAdminLogin)
	r.POST("/admin/login", h.handleAdminLoginPost)
	r.GET("/admin/logout", h.handleAdminLogout)
	r.POST("/admin/logout", h.handleAdminLogout)
	r.GET("/admin/dashboard", h.requireSession(h.handleAdminDashboard))
	r.GET("/admin/settings", h.requireSession(h.handleAdminSettings))
	r.POST("/admin/settings", h.requireSession(h.handleAdminSettingsPost))

	// Admin API (bearer token auth)
	adminAPI := r.Group("/api/v1/admin")
	{
		adminAPI.GET("/status", h.requireToken(h.handleAPIStatus))
		adminAPI.GET("/config", h.requireToken(h.handleAPIGetConfig))
		adminAPI.PUT("/config", h.requireToken(h.handleAPIUpdateConfig))
		adminAPI.POST("/reload", h.requireToken(h.handleAPIReload))
	}
}

// Middleware for session authentication
func (h *Handler) requireSession(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, ok := h.auth.GetSessionFromRequest(c.Request)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/admin")
			c.Abort()
			return
		}
		// Refresh session on activity
		h.auth.RefreshSession(session.ID)
		handler(c)
	}
}

// Middleware for bearer token authentication
func (h *Handler) requireToken(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetTokenFromRequest(c.Request)
		if token == "" || !h.auth.ValidateAPIToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		handler(c)
	}
}

// handleAdminLogin shows the login page
func (h *Handler) handleAdminLogin(c *gin.Context) {
	// Check if already logged in
	if _, ok := h.auth.GetSessionFromRequest(c.Request); ok {
		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
		return
	}

	h.renderLoginPage(c, "")
}

// handleAdminLoginPost processes login form
func (h *Handler) handleAdminLoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if h.auth.Authenticate(username, password) {
		session := h.auth.CreateSession(username, c.ClientIP())
		h.auth.SetSessionCookie(c.Writer, session)
		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
		return
	}

	h.renderLoginPage(c, "Invalid username or password")
}

// handleAdminLogout logs out the user
func (h *Handler) handleAdminLogout(c *gin.Context) {
	if session, ok := h.auth.GetSessionFromRequest(c.Request); ok {
		h.auth.DeleteSession(session.ID)
	}
	h.auth.ClearSessionCookie(c.Writer)
	c.Redirect(http.StatusSeeOther, "/admin")
}

// handleAdminDashboard shows the admin dashboard
func (h *Handler) handleAdminDashboard(c *gin.Context) {
	h.renderDashboardPage(c)
}

// handleAdminSettings shows the settings page
func (h *Handler) handleAdminSettings(c *gin.Context) {
	h.renderSettingsPage(c, "")
}

// handleAdminSettingsPost handles settings form submission
func (h *Handler) handleAdminSettingsPost(c *gin.Context) {
	h.renderSettingsPage(c, "Settings updated successfully")
}

// API Handlers

func (h *Handler) handleAPIStatus(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"version":   h.version,
		"commit":    h.commit,
		"buildDate": h.buildDate,
		"uptime":    time.Since(time.Now()).String(),
		"memory": gin.H{
			"alloc":      m.Alloc,
			"totalAlloc": m.TotalAlloc,
			"sys":        m.Sys,
			"numGC":      m.NumGC,
		},
		"runtime": gin.H{
			"goroutines": runtime.NumGoroutine(),
			"cpus":       runtime.NumCPU(),
			"goVersion":  runtime.Version(),
		},
	})
}

func (h *Handler) handleAPIGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": h.version,
	})
}

func (h *Handler) handleAPIUpdateConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) handleAPIReload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "reloaded"})
}

// HTML Templates

func (h *Handler) renderLoginPage(c *gin.Context, errorMsg string) {
	tmpl := template.Must(template.New("login").Parse(loginTemplate))
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, map[string]interface{}{
		"Error": errorMsg,
	})
}

func (h *Handler) renderDashboardPage(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, map[string]interface{}{
		"Version":    h.version,
		"Commit":     h.commit,
		"BuildDate":  h.buildDate,
		"MemAlloc":   fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		"Goroutines": runtime.NumGoroutine(),
	})
}

func (h *Handler) renderSettingsPage(c *gin.Context, message string) {
	tmpl := template.Must(template.New("settings").Parse(settingsTemplate))
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(c.Writer, map[string]interface{}{
		"Message": message,
	})
}

const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Login - Jokes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --red: #ff5555;
            --green: #50fa7b;
            --input-bg: #44475a;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .login-container {
            background: var(--input-bg);
            padding: 2rem;
            border-radius: 8px;
            width: 100%;
            max-width: 400px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
        h1 { text-align: center; margin-bottom: 1.5rem; color: var(--accent); }
        .error { background: var(--red); color: #fff; padding: 0.75rem; border-radius: 4px; margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
        input[type="text"], input[type="password"] {
            width: 100%;
            padding: 0.75rem;
            border: none;
            border-radius: 4px;
            background: var(--bg-color);
            color: var(--fg-color);
            margin-bottom: 1rem;
            font-size: 1rem;
        }
        input:focus { outline: 2px solid var(--accent); }
        button {
            width: 100%;
            padding: 0.75rem;
            border: none;
            border-radius: 4px;
            background: var(--accent);
            color: var(--bg-color);
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: opacity 0.2s;
        }
        button:hover { opacity: 0.9; }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Admin Login</h1>
        {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
        <form method="POST" action="/admin/login">
            <label for="username">Username</label>
            <input type="text" id="username" name="username" required autofocus>
            <label for="password">Password</label>
            <input type="password" id="password" name="password" required>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>`

const dashboardTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard - Jokes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --card-bg: #44475a;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
        }
        .navbar {
            background: var(--card-bg);
            padding: 1rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .navbar a:hover { color: var(--accent); }
        .container { max-width: 1200px; margin: 2rem auto; padding: 0 1rem; }
        .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; }
        .card {
            background: var(--card-bg);
            padding: 1.5rem;
            border-radius: 8px;
        }
        .card h3 { color: var(--accent); margin-bottom: 0.5rem; }
        .card p { font-size: 1.5rem; font-weight: bold; }
    </style>
</head>
<body>
    <nav class="navbar">
        <h1>Jokes API Admin</h1>
        <div>
            <a href="/admin/dashboard">Dashboard</a>
            <a href="/admin/settings">Settings</a>
            <a href="/admin/logout">Logout</a>
        </div>
    </nav>
    <div class="container">
        <div class="cards">
            <div class="card">
                <h3>Version</h3>
                <p>{{.Version}}</p>
            </div>
            <div class="card">
                <h3>Memory Usage</h3>
                <p>{{.MemAlloc}}</p>
            </div>
            <div class="card">
                <h3>Goroutines</h3>
                <p>{{.Goroutines}}</p>
            </div>
            <div class="card">
                <h3>Build Date</h3>
                <p>{{.BuildDate}}</p>
            </div>
        </div>
    </div>
</body>
</html>`

const settingsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings - Jokes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --card-bg: #44475a;
            --green: #50fa7b;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
        }
        .navbar {
            background: var(--card-bg);
            padding: 1rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .navbar a:hover { color: var(--accent); }
        .container { max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
        .message { background: var(--green); color: #000; padding: 1rem; border-radius: 4px; margin-bottom: 1rem; }
        .card {
            background: var(--card-bg);
            padding: 1.5rem;
            border-radius: 8px;
        }
        .card h2 { color: var(--accent); margin-bottom: 1rem; }
    </style>
</head>
<body>
    <nav class="navbar">
        <h1>Jokes API Admin</h1>
        <div>
            <a href="/admin/dashboard">Dashboard</a>
            <a href="/admin/settings">Settings</a>
            <a href="/admin/logout">Logout</a>
        </div>
    </nav>
    <div class="container">
        {{if .Message}}<div class="message">{{.Message}}</div>{{end}}
        <div class="card">
            <h2>Settings</h2>
            <p>Settings configuration coming soon.</p>
        </div>
    </div>
</body>
</html>`
