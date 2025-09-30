import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import api from "../api";

export default function Register() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post("/register", { username, email, password });
      alert("Cadastro realizado com sucesso!");
      navigate("/login");
    } catch (err) {
      alert("Erro no cadastro");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h2 style={styles.title}>Cadastro</h2>
        <form onSubmit={handleSubmit} style={styles.form}>
          <input type="text" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} style={styles.input} />
          <input autoFocus type="email" placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} style={styles.input} />
          <input type="password" placeholder="Senha" value={password} onChange={e => setPassword(e.target.value)} style={styles.input} />
          <button type="submit" style={styles.button} disabled={loading}>{loading ? "Cadastrando..." : "Cadastrar"}</button>
        </form>
        <p style={styles.text}>
          Já tem conta? <Link to="/login" style={styles.link}>Fazer login</Link>
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
  button: { padding: "0.9rem", borderRadius: 8, border: "none", background: "#059669", color: "#fff", fontWeight: 700, cursor: "pointer" },
  text: { marginTop: 12, color: "#6b7280" },
  link: { color: "#2563eb", fontWeight: 700, textDecoration: "none" },
};
