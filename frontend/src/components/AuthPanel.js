import { useState } from "react";

const defaultForm = { username: "", password: "" };

export default function AuthPanel({ onLogin, onRegister, loading }) {
  const [mode, setMode] = useState("login");
  const [form, setForm] = useState(defaultForm);

  function handleChange(event) {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    if (!form.username || !form.password) {
      return;
    }
    const payload = {
      username: form.username.trim(),
      password: form.password,
    };
    let ok = false;
    if (mode === "login") {
      ok = await onLogin(payload);
    } else {
      ok = await onRegister(payload);
    }
    if (ok) {
      setForm(defaultForm);
      if (mode === "register") {
        setMode("login");
      }
    }
  }

  function toggleMode() {
    setMode((prev) => (prev === "login" ? "register" : "login"));
  }

  return (
    <div className="card auth-card">
      <div className="auth-header">
        <h2>{mode === "login" ? "Welcome Back" : "Create Account"}</h2>
        <p className="muted">
          {mode === "login"
            ? "Sign in to continue"
            : "Register to get started"}
        </p>
      </div>
      <form onSubmit={handleSubmit} className="form">
        <label>
          <span>Username</span>
          <input
            type="text"
            name="username"
            value={form.username}
            onChange={handleChange}
            autoComplete="username"
            placeholder="Enter your username"
            required
          />
        </label>

        <label>
          <span>Password</span>
          <input
            type="password"
            name="password"
            value={form.password}
            onChange={handleChange}
            autoComplete={mode === "login" ? "current-password" : "new-password"}
            placeholder="Enter your password"
            required
            minLength={6}
          />
        </label>

        <button type="submit" className="primary" disabled={loading}>
          {loading ? "Processing..." : mode === "login" ? "Sign In" : "Register"}
        </button>
      </form>

      <div className="auth-footer">
        <button type="button" className="link-button" onClick={toggleMode}>
          {mode === "login"
            ? "Need an account? Register"
            : "Already have an account? Sign in"}
        </button>
      </div>
    </div>
  );
}

