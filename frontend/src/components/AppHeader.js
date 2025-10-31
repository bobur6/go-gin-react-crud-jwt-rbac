export default function AppHeader({ user, onLogout }) {
  return (
    <header className="app-header">
      <div>
        <h1>PortalGo</h1>
        <p className="muted">Secure item management</p>
      </div>
      {user && (
        <div className="user-pill">
          <span>{user.username}</span>
          <span className="badge">{user.role}</span>
          <button type="button" onClick={onLogout} className="secondary">
            Sign Out
          </button>
        </div>
      )}
    </header>
  );
}

