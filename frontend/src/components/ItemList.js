function formatDate(value) {
  if (!value) {
    return "";
  }
  try {
    return new Date(value).toLocaleString();
  } catch {
    return value;
  }
}

export default function ItemList({
  items,
  currentUser,
  loading,
  onEdit,
  onDelete,
}) {
  if (!items.length) {
    return (
      <div className="card">
        <div className="card-header">
          <h2>📋 Items</h2>
        </div>
        <div className="empty-state">
          <p className="muted">
            No items yet. Create your first item using the form.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="card">
      <div className="card-header">
        <h2>📋 Items</h2>
        <span className="badge-count">{items.length}</span>
      </div>
      <ul className="item-list">
        {items.map((item) => {
          const canEdit =
            currentUser?.role === "admin" ||
            currentUser?.username === item.owner;
          const canDelete = currentUser?.role === "admin";

          return (
            <li key={item.id} className="item-card">
              <div className="item-card__header">
                <h3>{item.title}</h3>
                <span className="badge-owner">{item.owner}</span>
              </div>
              {item.description && (
                <p className="item-card__description">{item.description}</p>
              )}
              <div className="item-card__meta">
                <span>📅 {formatDate(item.created_at)}</span>
              </div>
              <div className="item-card__actions">
                {canEdit && (
                  <button
                    type="button"
                    className="secondary"
                    onClick={() => onEdit(item)}
                    disabled={loading}
                  >
                    ✏️ Edit
                  </button>
                )}
                {canDelete && (
                  <button
                    type="button"
                    className="danger"
                    onClick={() => onDelete(item.id)}
                    disabled={loading}
                  >
                    🗑️ Delete
                  </button>
                )}
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

