import { useEffect, useState } from "react";

const defaultState = { title: "", description: "" };

export default function ItemForm({
  onCreate,
  onUpdate,
  editingItem,
  loading,
  onCancel,
}) {
  const [form, setForm] = useState(defaultState);

  useEffect(() => {
    if (editingItem) {
      setForm({
        title: editingItem.title,
        description: editingItem.description || "",
      });
    } else {
      setForm(defaultState);
    }
  }, [editingItem]);

  function handleChange(event) {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    if (!form.title.trim()) {
      return;
    }
    const payload = {
      title: form.title.trim(),
      description: form.description.trim(),
    };
    let ok = false;
    if (editingItem) {
      ok = await onUpdate(editingItem.id, payload);
    } else {
      ok = await onCreate(payload);
    }
    if (ok) {
      setForm(defaultState);
      onCancel?.();
    }
  }

  return (
    <div className="card">
      <div className="card-header">
        <h2>{editingItem ? "✏️ Edit Item" : "➕ New Item"}</h2>
      </div>
      <form className="form" onSubmit={handleSubmit}>
        <label>
          <span>Title</span>
          <input
            type="text"
            name="title"
            value={form.title}
            onChange={handleChange}
            placeholder="Enter item title"
            required
          />
        </label>

        <label>
          <span>Description</span>
          <textarea
            name="description"
            value={form.description}
            onChange={handleChange}
            rows={4}
            placeholder="Add optional details..."
          />
        </label>

        <div className="form-actions">
          <button type="submit" className="primary" disabled={loading}>
            {loading ? "Saving..." : editingItem ? "Update Item" : "Create Item"}
          </button>
          {editingItem && (
            <button
              type="button"
              className="secondary"
              onClick={() => {
                setForm(defaultState);
                onCancel?.();
              }}
              disabled={loading}
            >
              Cancel
            </button>
          )}
        </div>
      </form>
    </div>
  );
}

