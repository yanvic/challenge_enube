// src/pages/Import.jsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";

export default function ImportPage() {
  const [file, setFile] = useState(null);
  const [progress, setProgress] = useState(0);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleUpload = async () => {
    if (!file) return alert("Selecione um arquivo");
    setLoading(true);
    setProgress(0);

    const formData = new FormData();
    formData.append("file", file);

    try {
      await api.post("/excel", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        onUploadProgress: (evt) => {
          const pct = Math.round((evt.loaded * 100) / evt.total);
          setProgress(pct);
        },
      });
      alert("Upload realizado, aguarde o processamento.");
      navigate("/home");
    } catch (err) {
      console.error(err);
      alert("Erro no import");
    } finally {
      setLoading(false);
      setProgress(0);
    }
  };

  return (
    <div style={styles.page}>
      <div style={styles.card}>
        <h3>Importador de Excel / CSV</h3>
        <p style={{ color: "#6b7280" }}>Selecione o arquivo (.xlsx, .csv) e clique em importar.</p>

        <input type="file" accept=".csv, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-excel" onChange={e => setFile(e.target.files[0])} />
        <div style={{ marginTop: 12, display: "flex", gap: 8 }}>
          <button onClick={handleUpload} disabled={loading} style={styles.primary}>{loading ? `Enviando ${progress}%` : "Importar"}</button>
          <button onClick={() => navigate("/home")} style={styles.secondary}>Voltar</button>
        </div>

        {loading && <div style={{ marginTop: 12 }}>
          <div style={{ height: 8, background: "#eee", borderRadius: 8 }}>
            <div style={{ width: `${progress}%`, height: "100%", background: "#2563eb", borderRadius: 8 }} />
          </div>
        </div>}
      </div>
    </div>
  );
}

const styles = {
  page: { maxWidth: 900, margin: "30px auto", padding: "0 16px" },
  card: { padding: 20, borderRadius: 12, background: "#fff", boxShadow: "0 8px 30px rgba(2,6,23,0.06)" },
  primary: { padding: "8px 14px", borderRadius: 8, border: "none", background: "#059669", color: "#fff", cursor: "pointer" },
  secondary: { padding: "8px 14px", borderRadius: 8, border: "1px solid #e6e9ee", background: "#fff", cursor: "pointer" },
};
