import { useState, useEffect } from "react";
import api from "../api/client";

export default function UserManagement({ currentUser }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadUsers();
  }, []);

  async function loadUsers() {
    setLoading(true);
    setError(null);
    try {
      const data = await api.fetchUsers();
      setUsers(data);
    } catch (err) {
      setError(err.response?.data?.error || "Failed to load users");
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(userId, username) {
    if (!window.confirm(`Delete user "${username}"?`)) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await api.deleteUser(userId);
      setUsers((prev) => prev.filter((u) => u.id !== userId));
    } catch (err) {
      setError(err.response?.data?.error || "Failed to delete user");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="card">
      <div className="card-header">
        <h2>👥 User Management</h2>
        <div style={{ display: "flex", gap: "0.75rem", alignItems: "center" }}>
          <span className="badge-count">{users.length}</span>
          <button
            type="button"
            className="secondary"
            onClick={loadUsers}
            disabled={loading}
          >
            {loading ? "Loading..." : "🔄 Refresh"}
          </button>
        </div>
      </div>

      {error && <div className="error-message">⚠️ {error}</div>}

      <div className="user-list">
        {users.length === 0 && !loading && (
          <div className="empty-state">
            <p className="muted">No users found</p>
          </div>
        )}
        {users.map((user) => (
          <div key={user.id} className="user-item">
            <div className="user-info">
              <strong>{user.username}</strong>
              <span className="badge">{user.role}</span>
              <span className="muted">
                Joined {new Date(user.created_at).toLocaleDateString()}
              </span>
            </div>
            {user.id !== currentUser.id ? (
              <button
                type="button"
                className="danger"
                onClick={() => handleDelete(user.id, user.username)}
                disabled={loading}
              >
                🗑️ Delete
              </button>
            ) : (
              <span className="badge-owner">You</span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
