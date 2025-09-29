import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import api from "../api";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await api.post("/login", { email, password });
      localStorage.setItem("token", res.data.token);
      navigate("/home");
    } catch (err) {
      alert("Erro no login");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h2 style={styles.title}>Entrar</h2>
        <form onSubmit={handleSubmit} style={styles.form}>
          <input autoFocus type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} style={styles.input} />
          <input type="password" placeholder="Senha" value={password} onChange={e => setPassword(e.target.value)} style={styles.input} />
          <button type="submit" style={styles.button} disabled={loading}>{loading ? "Entrando..." : "Entrar"}</button>
        </form>
        <p style={styles.text}>
          Não tem conta? <Link to="/register" style={styles.link}>Cadastre-se</Link>
        </p>
      </div>
    </div>
  );
}

const styles = {
  container: { height: "calc(100vh - 64px)", display: "flex", justifyContent: "center", alignItems: "center", background: "#f6f8fa" },
  card: { width: 360, padding: "2rem", borderRadius: 12, background: "#fff", boxShadow: "0 6px 20px rgba(15, 23, 42, 0.06)", textAlign: "center" },
  title: { marginBottom: "1rem", color: "#111827" },
  form: { display: "flex", flexDirection: "column", gap: 12 },
  input: { padding: "0.9rem", borderRadius: 8, border: "1px solid #e6e9ee", fontSize: 14 },
  button: { padding: "0.9rem", borderRadius: 8, border: "none", background: "#2563eb", color: "#fff", fontWeight: 700, cursor: "pointer" },
  text: { marginTop: 12, color: "#6b7280" },
  link: { color: "#2563eb", fontWeight: 700, textDecoration: "none" },
};
