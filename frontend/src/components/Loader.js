export default function Loader({ visible }) {
  if (!visible) {
    return null;
  }

  return (
    <div className="loader-overlay">
      <div className="loader">Loading...</div>
    </div>
  );
}
