import { useState } from "react";
import AppHeader from "./components/AppHeader";
import AuthPanel from "./components/AuthPanel";
import ItemForm from "./components/ItemForm";
import ItemList from "./components/ItemList";
import Loader from "./components/Loader";
import Notification from "./components/Notification";
import UserManagement from "./components/UserManagement";
import { useAppContext } from "./context/AppContext";
import "./App.css";

function App() {
  const {
    state: { user, items, loading, error, notification },
    actions: {
      login,
      register,
      logout,
      createItem,
      updateItem,
      deleteItem,
      setError,
      setNotification,
    },
  } = useAppContext();
  const [editingItem, setEditingItem] = useState(null);

  async function handleDelete(id) {
    const ok = await deleteItem(id);
    if (ok && editingItem && editingItem.id === id) {
      setEditingItem(null);
    }
  }

  function handleLogout() {
    logout();
    setEditingItem(null);
  }

  return (
    <div className={`app-container ${!user ? 'app-container--centered' : ''}`}>
      <AppHeader user={user} onLogout={handleLogout} />

      {(notification || error) && (
        <div className="notification-wrapper">
          <Notification
            type="success"
            message={notification}
            onClose={() => setNotification(null)}
          />
          <Notification
            type="error"
            message={error}
            onClose={() => setError(null)}
          />
        </div>
      )}

      <Loader visible={loading} />

      {!user ? (
        <AuthPanel onLogin={login} onRegister={register} loading={loading} />
      ) : (
        <div className="content-grid">
          {user.role === "admin" && (
            <div className="admin-section">
              <UserManagement currentUser={user} />
            </div>
          )}
          <ItemForm
            onCreate={createItem}
            onUpdate={updateItem}
            editingItem={editingItem}
            loading={loading}
            onCancel={() => setEditingItem(null)}
          />
          <ItemList
            items={items}
            currentUser={user}
            loading={loading}
            onEdit={(item) => setEditingItem(item)}
            onDelete={handleDelete}
          />
        </div>
      )}
    </div>
  );
}

export default App;
